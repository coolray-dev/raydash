// GENERATED BY THE COMMAND ABOVE; DO NOT EDIT
// This file was generated by swaggo/swag

package docs

import (
	"bytes"
	"encoding/json"
	"strings"

	"github.com/alecthomas/template"
	"github.com/swaggo/swag"
)

var doc = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{.Description}}",
        "title": "{{.Title}}",
        "contact": {},
        "license": {
            "name": "GPLv3",
            "url": "https://www.gnu.org/licenses/gpl-3.0.html"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/users": {
            "get": {
                "description": "Simply list out all users",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Users"
                ],
                "summary": "All Users",
                "operationId": "users.Index",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Group ID",
                        "name": "gid",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "Node Token",
                        "name": "Authorization",
                        "in": "header"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/users.indexResponse"
                        }
                    },
                    "403": {
                        "description": "Forbidden",
                        "schema": {
                            "$ref": "#/definitions/handler.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/handler.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/users/{username}": {
            "get": {
                "description": "Return user according to username in url",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Users"
                ],
                "summary": "Show the required user",
                "operationId": "users.Show",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Username",
                        "name": "username",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Node Token",
                        "name": "Authorization",
                        "in": "header"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/users.showResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/handler.ErrorResponse"
                        }
                    }
                }
            },
            "delete": {
                "description": "Delete a user according to username",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Users"
                ],
                "summary": "Delete a user",
                "operationId": "users.Destroy",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Username",
                        "name": "username",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Node Token",
                        "name": "Authorization",
                        "in": "header"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/users.destroyResponse"
                        }
                    },
                    "403": {
                        "description": "Forbidden",
                        "schema": {
                            "$ref": "#/definitions/handler.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/handler.ErrorResponse"
                        }
                    }
                }
            },
            "patch": {
                "description": "Update the user provided",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Users"
                ],
                "summary": "Update user",
                "operationId": "users.Update",
                "parameters": [
                    {
                        "description": "User Object",
                        "name": "user",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.User"
                        }
                    },
                    {
                        "type": "string",
                        "description": "Username",
                        "name": "username",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Node Token",
                        "name": "Authorization",
                        "in": "header"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/users.updateResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/handler.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/users/{username}/groups": {
            "get": {
                "description": "Return a list of groups of a user",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Users"
                ],
                "summary": "List all groups",
                "operationId": "users.Groups",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Username",
                        "name": "username",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Node Token",
                        "name": "Authorization",
                        "in": "header"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/users.groupsResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/handler.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/users/{username}/nodes": {
            "get": {
                "description": "Return a list of nodes of a user",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Users"
                ],
                "summary": "List all nodes",
                "operationId": "users.Nodes",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Username",
                        "name": "username",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Node Token",
                        "name": "Authorization",
                        "in": "header"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/users.nodesResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/handler.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/users/{username}/services": {
            "get": {
                "description": "Return a list of services of a user",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Users"
                ],
                "summary": "List all services",
                "operationId": "users.Services",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Username",
                        "name": "username",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Node Token",
                        "name": "Authorization",
                        "in": "header"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/users.servicesResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/handler.ErrorResponse"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "handler.ErrorResponse": {
            "type": "object",
            "properties": {
                "error": {
                    "type": "string"
                }
            }
        },
        "models.Group": {
            "type": "object",
            "properties": {
                "created_at": {
                    "type": "string"
                },
                "description": {
                    "type": "string"
                },
                "groupname": {
                    "type": "string"
                },
                "id": {
                    "type": "integer"
                },
                "updated_at": {
                    "type": "string"
                }
            }
        },
        "models.Node": {
            "type": "object",
            "properties": {
                "access_token": {
                    "type": "string"
                },
                "created_at": {
                    "type": "string"
                },
                "current_traffic": {
                    "type": "integer"
                },
                "description": {
                    "type": "string"
                },
                "hasMultiPort": {
                    "type": "boolean"
                },
                "hasUDP": {
                    "type": "boolean"
                },
                "host": {
                    "description": "The Host to access v2ray",
                    "type": "string"
                },
                "id": {
                    "type": "integer"
                },
                "max_traffic": {
                    "type": "integer"
                },
                "name": {
                    "type": "string"
                },
                "ports": {
                    "type": "string"
                },
                "settings": {
                    "type": "object",
                    "$ref": "#/definitions/models.settings"
                },
                "updated_at": {
                    "type": "string"
                }
            }
        },
        "models.Service": {
            "type": "object",
            "properties": {
                "alterid": {
                    "type": "integer"
                },
                "created_at": {
                    "type": "string"
                },
                "description": {
                    "type": "string"
                },
                "email": {
                    "type": "string"
                },
                "host": {
                    "type": "string"
                },
                "id": {
                    "type": "integer"
                },
                "name": {
                    "type": "string"
                },
                "nid": {
                    "type": "integer"
                },
                "port": {
                    "type": "integer"
                },
                "protocol": {
                    "type": "string"
                },
                "security": {
                    "type": "string"
                },
                "uid": {
                    "type": "integer"
                },
                "updated_at": {
                    "type": "string"
                },
                "uuid": {
                    "type": "string"
                }
            }
        },
        "models.ShadowsocksSetting": {
            "type": "object"
        },
        "models.User": {
            "type": "object",
            "properties": {
                "created_at": {
                    "type": "string"
                },
                "current_traffic": {
                    "type": "integer"
                },
                "email": {
                    "type": "string"
                },
                "id": {
                    "type": "integer"
                },
                "max_traffic": {
                    "type": "integer"
                },
                "subscription_token": {
                    "type": "string"
                },
                "updated_at": {
                    "type": "string"
                },
                "username": {
                    "type": "string"
                },
                "uuid": {
                    "type": "string"
                }
            }
        },
        "models.VmessSetting": {
            "type": "object",
            "properties": {
                "protocol": {
                    "type": "string"
                }
            }
        },
        "models.settings": {
            "type": "object",
            "properties": {
                "listen": {
                    "type": "string"
                },
                "port": {
                    "type": "integer"
                },
                "shadowsocksSettings": {
                    "type": "object",
                    "$ref": "#/definitions/models.ShadowsocksSetting"
                },
                "vmessSettings": {
                    "type": "object",
                    "$ref": "#/definitions/models.VmessSetting"
                }
            }
        },
        "users.destroyResponse": {
            "type": "object",
            "properties": {
                "user": {
                    "type": "string"
                }
            }
        },
        "users.groupsResponse": {
            "type": "object",
            "properties": {
                "groups": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/models.Group"
                    }
                }
            }
        },
        "users.indexResponse": {
            "type": "object",
            "properties": {
                "total": {
                    "type": "integer"
                },
                "users": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/models.User"
                    }
                }
            }
        },
        "users.nodesResponse": {
            "type": "object",
            "properties": {
                "nodes": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/models.Node"
                    }
                }
            }
        },
        "users.servicesResponse": {
            "type": "object",
            "properties": {
                "services": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/models.Service"
                    }
                }
            }
        },
        "users.showResponse": {
            "type": "object",
            "properties": {
                "user": {
                    "type": "object",
                    "$ref": "#/definitions/models.User"
                }
            }
        },
        "users.updateResponse": {
            "type": "object",
            "properties": {
                "user": {
                    "type": "object",
                    "$ref": "#/definitions/models.User"
                }
            }
        }
    }
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
	Version:     "1.0",
	Host:        "localhost",
	BasePath:    "/v1",
	Schemes:     []string{},
	Title:       "RayDash API",
	Description: "A Swagger UI For RayDash API",
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
	swag.Register(swag.Name, &s{})
}
