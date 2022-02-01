// Package docs GENERATED BY THE COMMAND ABOVE; DO NOT EDIT
// This file was generated by swaggo/swag
package docs

import (
	"bytes"
	"encoding/json"
	"strings"
	"text/template"

	"github.com/swaggo/swag"
)

var doc = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {
            "name": "Yuanji",
            "url": "https://gimo.me",
            "email": "self@gimo.me"
        },
        "license": {
            "name": "MIT License",
            "url": "https://github.com/masakichi/echobin/blob/main/LICENSE"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/delete": {
            "delete": {
                "consumes": [
                    "application/json",
                    "multipart/form-data",
                    "application/x-www-form-urlencoded"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "HTTP methods"
                ],
                "summary": "The request's query parameters.",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/main.otherMethodResponse"
                        }
                    }
                }
            }
        },
        "/deny": {
            "get": {
                "produces": [
                    "text/plain"
                ],
                "tags": [
                    "Response formats"
                ],
                "summary": "Returns page denied by robots.txt rules.",
                "responses": {
                    "200": {
                        "description": "Denied message"
                    }
                }
            }
        },
        "/get": {
            "get": {
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "HTTP methods"
                ],
                "summary": "The request's query parameters.",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/main.getMethodResponse"
                        }
                    }
                }
            }
        },
        "/headers": {
            "get": {
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Request inspection"
                ],
                "summary": "Return the incoming request's HTTP headers.",
                "responses": {
                    "200": {
                        "description": "The request’s headers.",
                        "schema": {
                            "$ref": "#/definitions/main.requestHeadersResponse"
                        }
                    }
                }
            }
        },
        "/html": {
            "get": {
                "produces": [
                    "text/html"
                ],
                "tags": [
                    "Response formats"
                ],
                "summary": "Returns a simple HTML document.",
                "responses": {
                    "200": {
                        "description": "An HTML page."
                    }
                }
            }
        },
        "/ip": {
            "get": {
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Request inspection"
                ],
                "summary": "Returns the requester's IP Address.",
                "responses": {
                    "200": {
                        "description": "The Requester’s IP Address.",
                        "schema": {
                            "$ref": "#/definitions/main.requestIPResponse"
                        }
                    }
                }
            }
        },
        "/json": {
            "get": {
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Response formats"
                ],
                "summary": "Returns a simple JSON document.",
                "responses": {
                    "200": {
                        "description": "An JSON document."
                    }
                }
            }
        },
        "/patch": {
            "patch": {
                "consumes": [
                    "application/json",
                    "multipart/form-data",
                    "application/x-www-form-urlencoded"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "HTTP methods"
                ],
                "summary": "The request's query parameters.",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/main.otherMethodResponse"
                        }
                    }
                }
            }
        },
        "/post": {
            "post": {
                "consumes": [
                    "application/json",
                    "multipart/form-data",
                    "application/x-www-form-urlencoded"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "HTTP methods"
                ],
                "summary": "The request's query parameters.",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/main.otherMethodResponse"
                        }
                    }
                }
            }
        },
        "/put": {
            "put": {
                "consumes": [
                    "application/json",
                    "multipart/form-data",
                    "application/x-www-form-urlencoded"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "HTTP methods"
                ],
                "summary": "The request's query parameters.",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/main.otherMethodResponse"
                        }
                    }
                }
            }
        },
        "/robots.txt": {
            "get": {
                "produces": [
                    "text/plain"
                ],
                "tags": [
                    "Response formats"
                ],
                "summary": "Returns some robots.txt rules.",
                "responses": {
                    "200": {
                        "description": "Robots file"
                    }
                }
            }
        },
        "/status/{codes}": {
            "get": {
                "produces": [
                    "text/plain"
                ],
                "tags": [
                    "Status codes"
                ],
                "summary": "Return status code or random status code if more than one are given",
                "parameters": [
                    {
                        "type": "string",
                        "description": "codes",
                        "name": "codes",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "100": {
                        "description": "Informational responses"
                    },
                    "200": {
                        "description": "Success"
                    },
                    "300": {
                        "description": "Redirection"
                    },
                    "400": {
                        "description": "Client Errors"
                    },
                    "500": {
                        "description": "Server Errors"
                    }
                }
            },
            "put": {
                "produces": [
                    "text/plain"
                ],
                "tags": [
                    "Status codes"
                ],
                "summary": "Return status code or random status code if more than one are given",
                "parameters": [
                    {
                        "type": "string",
                        "description": "codes",
                        "name": "codes",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "100": {
                        "description": "Informational responses"
                    },
                    "200": {
                        "description": "Success"
                    },
                    "300": {
                        "description": "Redirection"
                    },
                    "400": {
                        "description": "Client Errors"
                    },
                    "500": {
                        "description": "Server Errors"
                    }
                }
            },
            "post": {
                "produces": [
                    "text/plain"
                ],
                "tags": [
                    "Status codes"
                ],
                "summary": "Return status code or random status code if more than one are given",
                "parameters": [
                    {
                        "type": "string",
                        "description": "codes",
                        "name": "codes",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "100": {
                        "description": "Informational responses"
                    },
                    "200": {
                        "description": "Success"
                    },
                    "300": {
                        "description": "Redirection"
                    },
                    "400": {
                        "description": "Client Errors"
                    },
                    "500": {
                        "description": "Server Errors"
                    }
                }
            },
            "delete": {
                "produces": [
                    "text/plain"
                ],
                "tags": [
                    "Status codes"
                ],
                "summary": "Return status code or random status code if more than one are given",
                "parameters": [
                    {
                        "type": "string",
                        "description": "codes",
                        "name": "codes",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "100": {
                        "description": "Informational responses"
                    },
                    "200": {
                        "description": "Success"
                    },
                    "300": {
                        "description": "Redirection"
                    },
                    "400": {
                        "description": "Client Errors"
                    },
                    "500": {
                        "description": "Server Errors"
                    }
                }
            },
            "patch": {
                "produces": [
                    "text/plain"
                ],
                "tags": [
                    "Status codes"
                ],
                "summary": "Return status code or random status code if more than one are given",
                "parameters": [
                    {
                        "type": "string",
                        "description": "codes",
                        "name": "codes",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "100": {
                        "description": "Informational responses"
                    },
                    "200": {
                        "description": "Success"
                    },
                    "300": {
                        "description": "Redirection"
                    },
                    "400": {
                        "description": "Client Errors"
                    },
                    "500": {
                        "description": "Server Errors"
                    }
                }
            }
        },
        "/user-agent": {
            "get": {
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Request inspection"
                ],
                "summary": "Return the incoming requests's User-Agent header.",
                "responses": {
                    "200": {
                        "description": "The request’s User-Agent header.",
                        "schema": {
                            "$ref": "#/definitions/main.requestUserAgentResponse"
                        }
                    }
                }
            }
        },
        "/xml": {
            "get": {
                "produces": [
                    "text/xml"
                ],
                "tags": [
                    "Response formats"
                ],
                "summary": "Returns a simple XML document.",
                "responses": {
                    "200": {
                        "description": "An XML document."
                    }
                }
            }
        }
    },
    "definitions": {
        "main.getMethodResponse": {
            "type": "object",
            "properties": {
                "args": {
                    "type": "object",
                    "additionalProperties": true
                },
                "headers": {
                    "type": "object",
                    "additionalProperties": {
                        "type": "string"
                    }
                },
                "origin": {
                    "type": "string"
                },
                "url": {
                    "type": "string"
                }
            }
        },
        "main.otherMethodResponse": {
            "type": "object",
            "properties": {
                "args": {
                    "type": "object",
                    "additionalProperties": true
                },
                "data": {
                    "type": "string"
                },
                "files": {
                    "type": "object",
                    "additionalProperties": true
                },
                "form": {},
                "headers": {
                    "type": "object",
                    "additionalProperties": {
                        "type": "string"
                    }
                },
                "json": {
                    "type": "object",
                    "additionalProperties": true
                },
                "origin": {
                    "type": "string"
                },
                "url": {
                    "type": "string"
                }
            }
        },
        "main.requestHeadersResponse": {
            "type": "object",
            "properties": {
                "headers": {
                    "type": "object",
                    "additionalProperties": {
                        "type": "string"
                    }
                }
            }
        },
        "main.requestIPResponse": {
            "type": "object",
            "properties": {
                "origin": {
                    "type": "string"
                }
            }
        },
        "main.requestUserAgentResponse": {
            "type": "object",
            "properties": {
                "user-agent": {
                    "type": "string"
                }
            }
        }
    },
    "tags": [
        {
            "description": "Testing different HTTP verbs",
            "name": "HTTP methods"
        },
        {
            "description": "Generates responses with given status code",
            "name": "Status codes"
        },
        {
            "description": "Inspect the request data",
            "name": "Request inspection"
        },
        {
            "description": "Returns responses in different data formats",
            "name": "Response formats"
        }
    ]
}`

type swaggerInfo struct {
	Version     string
	Host        string
	BasePath    string
	Schemes     []string
	Title       string
	Description string
}

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = swaggerInfo{
	Version:     "0.1",
	Host:        "",
	BasePath:    "",
	Schemes:     []string{},
	Title:       "echobin API",
	Description: "A simple HTTP Request & Response Service.",
}

type s struct{}

func (s *s) ReadDoc() string {
	sInfo := SwaggerInfo
	sInfo.Description = strings.Replace(sInfo.Description, "\n", "\\n", -1)

	t, err := template.New("swagger_info").Funcs(template.FuncMap{
		"marshal": func(v interface{}) string {
			a, _ := json.Marshal(v)
			return string(a)
		},
		"escape": func(v interface{}) string {
			// escape tabs
			str := strings.Replace(v.(string), "\t", "\\t", -1)
			// replace " with \", and if that results in \\", replace that with \\\"
			str = strings.Replace(str, "\"", "\\\"", -1)
			return strings.Replace(str, "\\\\\"", "\\\\\\\"", -1)
		},
	}).Parse(doc)
	if err != nil {
		return doc
	}

	var tpl bytes.Buffer
	if err := t.Execute(&tpl, sInfo); err != nil {
		return doc
	}

	return tpl.String()
}

func init() {
	swag.Register("swagger", &s{})
}
