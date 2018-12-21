This tool wraps [telegraf](https://www.influxdata.com/time-series-platform/telegraf/)'s `config`
operation extracting the configuration schema of a requested input plugin.

# Usage

You can pass the name of a telegraf input plugin

```
./telegraf-config-schema-extractor <input-plugin>
```

The full list of names is available by using `telegraf --input-list`.

# Building

Clone this repo outside of your `$GOPATH` and with Go 1.11 or greater installed, run

```bash
go build
```

The `telegraf-config-schema-extractor` executable will be produced in the current directory.

# Example

The following example extracts the schema of the [http_response](https://github.com/influxdata/telegraf/tree/master/plugins/inputs/http_response) plugin:

```bash
./telegraf-config-schema-extractor http_response
```

which would result in the following json (reformatted here for clarity):

```json
{
  "http_response": [
    {
      "Description": "Server address (default http://localhost)",
      "Params": {
        "address": {
          "Name": "address",
          "Type": "string",
          "Examples": [
            "http://localhost"
          ]
        }
      }
    },
    {
      "Description": "Set http_proxy (telegraf uses the system wide proxy settings if it's is not set)",
      "Params": {
        "http_proxy": {
          "Name": "http_proxy",
          "Type": "string",
          "Examples": [
            "http://localhost:8888"
          ]
        }
      }
    },
    {
      "Description": "Set response_timeout (default 5 seconds)",
      "Params": {
        "response_timeout": {
          "Name": "response_timeout",
          "Type": "string",
          "Examples": [
            "5s"
          ]
        }
      }
    },
    {
      "Description": "HTTP Request Method",
      "Params": {
        "method": {
          "Name": "method",
          "Type": "string",
          "Examples": [
            "GET"
          ]
        }
      }
    },
    {
      "Description": "Whether to follow redirects from the server (defaults to false)",
      "Params": {
        "follow_redirects": {
          "Name": "follow_redirects",
          "Type": "boolean",
          "Examples": [
            "false"
          ]
        }
      }
    },
    {
      "Description": "Optional HTTP Request Body",
      "Params": {
        "body": {
          "Name": "body",
          "Type": "string",
          "Examples": [
            "{'fake':'data'}\n"
          ]
        }
      }
    },
    {
      "Description": "Optional substring or regex match in body of the response",
      "Params": {
        "response_string_match": {
          "Name": "response_string_match",
          "Type": "string",
          "Examples": [
            "\"service_status\": \"up\"",
            "ok",
            "\".*_status\".?:.?\"up\""
          ]
        }
      }
    },
    {
      "Description": "Optional TLS Config",
      "Params": {
        "tls_ca": {
          "Name": "tls_ca",
          "Type": "string",
          "Examples": [
            "/etc/telegraf/ca.pem"
          ]
        },
        "tls_cert": {
          "Name": "tls_cert",
          "Type": "string",
          "Examples": [
            "/etc/telegraf/cert.pem"
          ]
        },
        "tls_key": {
          "Name": "tls_key",
          "Type": "string",
          "Examples": [
            "/etc/telegraf/key.pem"
          ]
        }
      }
    },
    {
      "Description": "Use TLS but skip chain \u0026 host verification",
      "Params": {
        "insecure_skip_verify": {
          "Name": "insecure_skip_verify",
          "Type": "boolean",
          "Examples": [
            "false"
          ]
        }
      }
    },
    {
      "Description": "HTTP Request Headers (all values must be strings)",
      "Params": {
        "headers": {
          "Name": "headers",
          "Type": "map"
        }
      }
    }
  ]
}
```