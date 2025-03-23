// Package docs Code generated by swaggo/swag. DO NOT EDIT
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {
            "name": "API Support",
            "email": "ericphlpp@proton.me"
        },
        "license": {
            "name": "MIT",
            "url": "https://github.com/AlgoHive-Coding-Puzzles/BeeAPI/blob/main/LICENSE"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/apikey": {
            "get": {
                "security": [
                    {
                        "Bearer": []
                    }
                ],
                "description": "Returns the current API key",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "API Key"
                ],
                "summary": "Check API key",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    }
                }
            }
        },
        "/name": {
            "get": {
                "description": "Returns the name of the server",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "App"
                ],
                "summary": "Get server name",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    }
                }
            }
        },
        "/ping": {
            "get": {
                "description": "Returns a pong response to check if the API is alive",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "App"
                ],
                "summary": "Health check endpoint",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    }
                }
            }
        },
        "/puzzle": {
            "get": {
                "description": "Returns details about a specific puzzle",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Puzzles"
                ],
                "summary": "Get puzzle details",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Theme name",
                        "name": "theme",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Puzzle Id",
                        "name": "puzzle",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/models.PuzzleResponse"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    }
                }
            },
            "delete": {
                "security": [
                    {
                        "Bearer": []
                    }
                ],
                "description": "Deletes a puzzle from a theme",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Puzzles"
                ],
                "summary": "Delete a puzzle",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Theme name",
                        "name": "theme",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Puzzle Id",
                        "name": "puzzle",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    }
                }
            }
        },
        "/puzzle/check/first": {
            "get": {
                "description": "Checks if the first solution matches the provided value",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Puzzles"
                ],
                "summary": "Check first solution",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Theme name",
                        "name": "theme",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Puzzle Id",
                        "name": "puzzle",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Unique ID for generation",
                        "name": "unique_id",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Solution to check",
                        "name": "solution",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "boolean"
                            }
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    }
                }
            }
        },
        "/puzzle/check/second": {
            "get": {
                "description": "Checks if the second solution matches the provided value",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Puzzles"
                ],
                "summary": "Check second solution",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Theme name",
                        "name": "theme",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Puzzle Id",
                        "name": "puzzle",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Unique ID for generation",
                        "name": "unique_id",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Solution to check",
                        "name": "solution",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "boolean"
                            }
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    }
                }
            }
        },
        "/puzzle/generate/input": {
            "get": {
                "description": "Generates puzzle input for a given puzzle",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Puzzles"
                ],
                "summary": "Generate puzzle input",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Theme name",
                        "name": "theme",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Puzzle Id",
                        "name": "puzzle",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Unique ID for generation",
                        "name": "unique_id",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    }
                }
            }
        },
        "/puzzle/upload": {
            "post": {
                "security": [
                    {
                        "Bearer": []
                    }
                ],
                "description": "Uploads a new puzzle to a theme",
                "consumes": [
                    "multipart/form-data"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Puzzles"
                ],
                "summary": "Upload a puzzle",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Theme name",
                        "name": "theme",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "file",
                        "description": "Puzzle file (.alghive)",
                        "name": "file",
                        "in": "formData",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    }
                }
            }
        },
        "/puzzles": {
            "get": {
                "description": "Returns all puzzles for a specific theme",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Puzzles"
                ],
                "summary": "Get puzzles for a theme",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Theme name",
                        "name": "theme",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/models.PuzzleResponse"
                            }
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    }
                }
            }
        },
        "/puzzles/ids": {
            "get": {
                "description": "Returns IDs of all puzzles for a specific theme",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Puzzles"
                ],
                "summary": "Get puzzle IDs",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Theme name",
                        "name": "theme",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "type": "string"
                            }
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    }
                }
            }
        },
        "/puzzles/names": {
            "get": {
                "description": "Returns names of all puzzles for a specific theme",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Puzzles"
                ],
                "summary": "Get puzzle names",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Theme name",
                        "name": "theme",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "type": "string"
                            }
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    }
                }
            }
        },
        "/theme": {
            "get": {
                "description": "Returns details of a specific theme by name",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Themes"
                ],
                "summary": "Get a specific theme",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Theme name",
                        "name": "name",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/models.ThemeResponse"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    }
                }
            },
            "post": {
                "security": [
                    {
                        "Bearer": []
                    }
                ],
                "description": "Creates a new theme with the given name",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Themes"
                ],
                "summary": "Create a new theme",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Theme name",
                        "name": "name",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    }
                }
            },
            "delete": {
                "security": [
                    {
                        "Bearer": []
                    }
                ],
                "description": "Deletes a theme with the given name",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Themes"
                ],
                "summary": "Delete a theme",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Theme name",
                        "name": "name",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    }
                }
            }
        },
        "/theme/reload": {
            "post": {
                "security": [
                    {
                        "Bearer": []
                    }
                ],
                "description": "Reloads all themes and puzzles",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Themes"
                ],
                "summary": "Reload themes",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    },
                    "429": {
                        "description": "Too Many Requests",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    }
                }
            }
        },
        "/themes": {
            "get": {
                "description": "Returns a list of all available themes",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Themes"
                ],
                "summary": "Get all themes",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/models.ThemeResponse"
                            }
                        }
                    }
                }
            }
        },
        "/themes/names": {
            "get": {
                "description": "Returns a list of theme names",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Themes"
                ],
                "summary": "Get theme names",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "type": "string"
                            }
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "models.PuzzleResponse": {
            "type": "object",
            "properties": {
                "author": {
                    "type": "string"
                },
                "cipher": {
                    "type": "string"
                },
                "compressedSize": {
                    "type": "integer"
                },
                "createdAt": {
                    "type": "string"
                },
                "difficulty": {
                    "type": "string"
                },
                "id": {
                    "type": "string"
                },
                "language": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                },
                "obscure": {
                    "type": "string"
                },
                "uncompressedSize": {
                    "type": "integer"
                },
                "updatedAt": {
                    "type": "string"
                }
            }
        },
        "models.ThemeResponse": {
            "type": "object",
            "properties": {
                "enigmes_count": {
                    "type": "integer"
                },
                "name": {
                    "type": "string"
                },
                "puzzles": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/models.PuzzleResponse"
                    }
                },
                "size": {
                    "type": "integer"
                }
            }
        }
    },
    "securityDefinitions": {
        "Bearer": {
            "type": "apiKey",
            "name": "Authorization",
            "in": "Bearer"
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "localhost:5000",
	BasePath:         "/",
	Schemes:          []string{},
	Title:            "BeeAPI Go",
	Description:      "API server for AlgoHive puzzles",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
