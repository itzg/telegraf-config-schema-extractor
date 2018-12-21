package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/influxdata/toml"
	"github.com/influxdata/toml/ast"
	"log"
	"os"
	"os/exec"
	"regexp"
)

func main() {
	if len(os.Args) != 2 {
		log.Fatal("Expected only a single telegraf input plugin name")
	}

	plugin := os.Args[1]

	cmd := exec.Command("telegraf", "--input-filter", plugin, "config")
	configOutBytes, err := cmd.Output()
	if err != nil {
		log.Fatal("Failed to run telegraf config command", err)
	}

	configSection := extractPluginConfigSection(plugin, configOutBytes)

	groups := processConfigGroups(configSection)
	groupsJson, err := json.Marshal(groups)
	if err != nil {
		log.Fatal("Failed to marshal groups", err)
	}

	fmt.Println(string(groupsJson))
}

type ConfigParam struct {
	Name string
	Type string
}

type ConfigGroup struct {
	Description string
	Params      []*ConfigParam
}

func processConfigGroups(configSection []byte) []*ConfigGroup {
	var groups []*ConfigGroup

	groupDescExp, err := regexp.Compile(`## (.+)`)
	if err != nil {
		log.Fatal("Failed to compile groupDesc regexp", err)
	}
	paramStartExp, err := regexp.Compile(`# ([[:word:]]+\s*=.+|\[.+?\])$`)
	if err != nil {
		log.Fatal("Failed to compile paramStartExp regexp", err)
	}
	paramContinueExp, err := regexp.Compile(`# (.+)`)
	if err != nil {
		log.Fatal("Failed to compile paramContinueExp regexp", err)
	}

	var currentGroup *ConfigGroup
	var groupDesc *bytes.Buffer
	var param *bytes.Buffer
	scanner := bufio.NewScanner(bytes.NewReader(configSection))
	for scanner.Scan() {
		line := scanner.Bytes()

		if captures := groupDescExp.FindSubmatch(line); captures != nil {
			if currentGroup != nil {
				if param != nil {
					// append previous param to previous group
					appendParamToGroup(param, currentGroup)
					param = nil
				}
				groups = append(groups, currentGroup)
				currentGroup = nil
			}

			if groupDesc == nil {
				groupDesc = new(bytes.Buffer)
				currentGroup = new(ConfigGroup)
			} else {
				groupDesc.WriteByte(' ')
			}
			groupDesc.Write(captures[1])

		} else if captures := paramStartExp.FindSubmatch(line); captures != nil {
			if groupDesc != nil {
				currentGroup.Description = groupDesc.String()
				groupDesc = nil
			}

			if param != nil {
				// append previous param to current group
				appendParamToGroup(param, currentGroup)
			}
			param = new(bytes.Buffer)
			param.Write(captures[1])
		} else if param != nil {
			if captures := paramContinueExp.FindSubmatch(line); captures != nil {
				param.WriteByte('\n')
				param.Write(captures[1])
			} else if currentGroup != nil {
				// append previous param to current group
				appendParamToGroup(param, currentGroup)
				param = nil
			}
		}
	}

	if currentGroup != nil {
		if param != nil {
			// append previous param to previous group
			appendParamToGroup(param, currentGroup)
		}
		groups = append(groups, currentGroup)
	}

	return groups
}

func appendParamToGroup(paramContent *bytes.Buffer, group *ConfigGroup) {
	table, err := toml.Parse(paramContent.Bytes())
	if err != nil {
		log.Println("Failed to parse param chunk", paramContent.String(), err)
		return
	}

	if len(table.Fields) != 1 {
		log.Println("Only expected one field to be parsed", paramContent.String())
		return
	}

	for key, entry := range table.Fields {
		name := key
		var paramType string

		switch typedEntry := entry.(type) {

		case *ast.KeyValue:
			switch typedEntry.Value.(type) {
			case *ast.String:
				paramType = "string"

			case *ast.Boolean:
				paramType = "boolean"

			default:
				fmt.Println("Unknown type given", paramContent.String())
				return
			}

		case *ast.Table:
			paramType = "map"
			name = deepestTableName(typedEntry)
		}

		param := &ConfigParam{Name: name, Type: paramType}
		group.Params = append(group.Params, param)

		return
	}
}

func deepestTableName(table *ast.Table) string {
	// iterate enough to grab first field of the table
	for _, entry := range table.Fields {

		if innerTable, ok := entry.(*ast.Table); ok {
			// contains a table
			return deepestTableName(innerTable)
		} else {
			// reached deepest table
			return table.Name
		}
	}

	return table.Name
}

func extractPluginConfigSection(plugin string, configOutBytes []byte) []byte {
	startExp, err := regexp.Compile(`^\[\[inputs\.` + plugin + `\]\]`)
	if err != nil {
		log.Fatal("Failed to compile section-start regexp", err)
	}
	endExp, err := regexp.Compile(`^\S.*`)
	if err != nil {
		log.Fatal("Failed to compile section-end regexp", err)
	}

	section := new(bytes.Buffer)
	sectionStarted := false
	lineScanner := bufio.NewScanner(bytes.NewReader(configOutBytes))
	for lineScanner.Scan() {
		line := lineScanner.Bytes()
		if !sectionStarted && startExp.Match(line) {
			sectionStarted = true
		} else if sectionStarted && endExp.Match(line) {
			sectionStarted = false
		} else if sectionStarted {
			section.Write(line)
			section.WriteByte('\n')
		}
	}

	return section.Bytes()
}
