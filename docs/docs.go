// Package docs Code generated by swaggo/swag. DO NOT EDIT
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {},
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/api/v1/balance": {
            "get": {
                "description": "Показывает баланс пользователя на счёте по юзернейму",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "users",
                    "wallets"
                ],
                "summary": "Shows wallet balance",
                "parameters": [
                    {
                        "type": "string",
                        "description": "JWT token",
                        "name": "Authorization",
                        "in": "header",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK"
                    },
                    "400": {
                        "description": "Bad Request"
                    }
                }
            }
        },
        "/api/v1/exchange": {
            "post": {
                "description": "Позволяет обменять валюту на другую, курс можно узнать в /api/v1/exchange/rates",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "exchange"
                ],
                "summary": "Exchanger endpoint",
                "parameters": [
                    {
                        "type": "string",
                        "description": "JWT token",
                        "name": "Authorization",
                        "in": "header",
                        "required": true
                    },
                    {
                        "description": "Exchange query in json format",
                        "name": "amount",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/storages.Exchanger"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK"
                    },
                    "400": {
                        "description": "Bad Request"
                    }
                }
            }
        },
        "/api/v1/exchange/rates": {
            "get": {
                "description": "Позволяет узнать актуальный курс по отношению к доллару",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "exchange"
                ],
                "summary": "Exchanger endpoint",
                "parameters": [
                    {
                        "type": "string",
                        "description": "JWT token",
                        "name": "Authorization",
                        "in": "header",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK"
                    },
                    "400": {
                        "description": "Bad Request"
                    }
                }
            }
        },
        "/api/v1/login": {
            "post": {
                "description": "Авторизация пользователя",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "users"
                ],
                "summary": "Authorize  user",
                "parameters": [
                    {
                        "description": "Данные пользователя",
                        "name": "input",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/storages.User"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK"
                    },
                    "400": {
                        "description": "Bad Request"
                    }
                }
            }
        },
        "/api/v1/register": {
            "post": {
                "description": "Регистрация нового пользователя",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "users"
                ],
                "summary": "Register user",
                "parameters": [
                    {
                        "description": "user info",
                        "name": "input",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/storages.User"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created"
                    },
                    "400": {
                        "description": "Bad Request"
                    }
                }
            }
        },
        "/api/v1/wallet/deposit": {
            "post": {
                "description": "Пополнить баланс пользователя",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "wallets",
                    "users"
                ],
                "summary": "Deposit balance",
                "parameters": [
                    {
                        "type": "string",
                        "description": "JWT token",
                        "name": "Authorization",
                        "in": "header",
                        "required": true
                    },
                    {
                        "description": "Deposit query in json format",
                        "name": "amount",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/storages.Deposit"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK"
                    },
                    "400": {
                        "description": "Bad Request"
                    }
                }
            }
        },
        "/api/v1/wallet/withdraw": {
            "post": {
                "description": "Снять средства со счёта пользователя",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "users",
                    "wallets"
                ],
                "summary": "Withdraw amount",
                "parameters": [
                    {
                        "type": "string",
                        "description": "JWT token",
                        "name": "Authorization",
                        "in": "header",
                        "required": true
                    },
                    {
                        "description": "Withdraw query in json format",
                        "name": "amount",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/storages.Withdraw"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK"
                    },
                    "400": {
                        "description": "Bad Request"
                    }
                }
            }
        }
    },
    "definitions": {
        "storages.Deposit": {
            "type": "object",
            "properties": {
                "amount": {
                    "type": "number"
                },
                "currency": {
                    "type": "string"
                }
            }
        },
        "storages.Exchanger": {
            "type": "object",
            "properties": {
                "amount": {
                    "type": "number"
                },
                "from_currency": {
                    "type": "string"
                },
                "to_currency": {
                    "type": "string"
                }
            }
        },
        "storages.User": {
            "type": "object",
            "properties": {
                "email": {
                    "type": "string"
                },
                "password": {
                    "type": "string"
                },
                "username": {
                    "type": "string"
                }
            }
        },
        "storages.Withdraw": {
            "type": "object",
            "properties": {
                "amount": {
                    "type": "number"
                },
                "currency": {
                    "type": "string"
                }
            }
        }
    },
    "securityDefinitions": {
        "ApiKeyAuth": {
            "type": "apiKey",
            "name": "Authorization",
            "in": "header"
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "localhost:8080",
	BasePath:         "/",
	Schemes:          []string{},
	Title:            "Currency wallet",
	Description:      "API service",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
