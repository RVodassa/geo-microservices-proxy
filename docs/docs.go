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
            "email": "support@geo.com"
        },
        "license": {
            "name": "Apache 2.0"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/api/address/geocode": {
            "post": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "description": "Выполняет геокодирование на основе переданных параметров",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "geo-service"
                ],
                "summary": "Геокодирование",
                "parameters": [
                    {
                        "description": "Параметры геокодирования",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/entity.GeocodeRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Успешное геокодирование",
                        "schema": {
                            "$ref": "#/definitions/controller.Response"
                        }
                    },
                    "400": {
                        "description": "Ошибка клиента",
                        "schema": {
                            "$ref": "#/definitions/controller.ErrorResponse"
                        }
                    },
                    "403": {
                        "description": "Unauthorized",
                        "schema": {
                            "$ref": "#/definitions/controller.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Ошибка на сервере",
                        "schema": {
                            "$ref": "#/definitions/controller.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/api/address/search": {
            "post": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "description": "Выполняет поиск данных на основе переданных параметров",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "geo-service"
                ],
                "summary": "Поиск данных",
                "parameters": [
                    {
                        "description": "Параметры поиска",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/entity.SearchRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Успешный поиск",
                        "schema": {
                            "$ref": "#/definitions/controller.Response"
                        }
                    },
                    "400": {
                        "description": "Ошибка клиента",
                        "schema": {
                            "$ref": "#/definitions/controller.ErrorResponse"
                        }
                    },
                    "403": {
                        "description": "Unauthorized",
                        "schema": {
                            "$ref": "#/definitions/controller.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Ошибка на сервере",
                        "schema": {
                            "$ref": "#/definitions/controller.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/api/auth/login": {
            "post": {
                "description": "Проверяет учетные данные и возвращает токен доступа",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "auth"
                ],
                "summary": "Аутентификация пользователя",
                "parameters": [
                    {
                        "description": "Данные для входа",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/entity.LoginRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Успешная аутентификация",
                        "schema": {
                            "$ref": "#/definitions/controller.Response"
                        }
                    },
                    "400": {
                        "description": "Ошибка клиента",
                        "schema": {
                            "$ref": "#/definitions/controller.ErrorResponse"
                        }
                    },
                    "401": {
                        "description": "Неверные учетные данные",
                        "schema": {
                            "$ref": "#/definitions/controller.ErrorResponse"
                        }
                    },
                    "403": {
                        "description": "Unauthorized",
                        "schema": {
                            "$ref": "#/definitions/controller.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Ошибка на сервере",
                        "schema": {
                            "$ref": "#/definitions/controller.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/api/auth/register": {
            "post": {
                "description": "Сохраняет данные для авторизации в бд",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "auth"
                ],
                "summary": "Регистрация пользователя",
                "parameters": [
                    {
                        "description": "Данные для регистрации",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/entity.RegisterRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Успешная регистрация",
                        "schema": {
                            "$ref": "#/definitions/controller.Response"
                        }
                    },
                    "400": {
                        "description": "Ошибка клиента",
                        "schema": {
                            "$ref": "#/definitions/controller.ErrorResponse"
                        }
                    },
                    "403": {
                        "description": "Unauthorized",
                        "schema": {
                            "$ref": "#/definitions/controller.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Ошибка на сервере",
                        "schema": {
                            "$ref": "#/definitions/controller.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/api/user/list": {
            "post": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "description": "Возвращает список пользователей на основе переданных параметров",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "users"
                ],
                "summary": "Получить список пользователей",
                "parameters": [
                    {
                        "description": "Параметры запроса",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/entity.ListRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Успешный запрос",
                        "schema": {
                            "$ref": "#/definitions/controller.Response"
                        }
                    },
                    "400": {
                        "description": "Ошибка клиента",
                        "schema": {
                            "$ref": "#/definitions/controller.ErrorResponse"
                        }
                    },
                    "403": {
                        "description": "Unauthorized",
                        "schema": {
                            "$ref": "#/definitions/controller.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Ошибка на сервере",
                        "schema": {
                            "$ref": "#/definitions/controller.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/api/user/profile": {
            "post": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "description": "Возвращает профиль пользователя на основе переданных параметров",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "profile"
                ],
                "summary": "Получить профиль пользователя",
                "parameters": [
                    {
                        "description": "Параметры запроса",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/entity.ProfileRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Успешный запрос",
                        "schema": {
                            "$ref": "#/definitions/controller.Response"
                        }
                    },
                    "400": {
                        "description": "Ошибка клиента",
                        "schema": {
                            "$ref": "#/definitions/controller.ErrorResponse"
                        }
                    },
                    "403": {
                        "description": "Unauthorized",
                        "schema": {
                            "$ref": "#/definitions/controller.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Ошибка на сервере",
                        "schema": {
                            "$ref": "#/definitions/controller.ErrorResponse"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "controller.ErrorResponse": {
            "type": "object",
            "properties": {
                "message": {
                    "type": "string"
                }
            }
        },
        "controller.Response": {
            "type": "object",
            "properties": {
                "code": {
                    "description": "Сообщение для пользователя",
                    "type": "integer"
                },
                "data": {
                    "description": "Данные ответа"
                }
            }
        },
        "entity.GeocodeRequest": {
            "type": "object",
            "properties": {
                "lat": {
                    "type": "string"
                },
                "lng": {
                    "type": "string"
                }
            }
        },
        "entity.ListRequest": {
            "type": "object",
            "properties": {
                "limit": {
                    "type": "integer"
                },
                "offset": {
                    "type": "integer"
                }
            }
        },
        "entity.LoginRequest": {
            "type": "object",
            "properties": {
                "login": {
                    "type": "string"
                },
                "password": {
                    "type": "string"
                }
            }
        },
        "entity.ProfileRequest": {
            "type": "object",
            "properties": {
                "id": {
                    "type": "integer"
                }
            }
        },
        "entity.RegisterRequest": {
            "type": "object",
            "properties": {
                "login": {
                    "type": "string"
                },
                "password": {
                    "type": "string"
                }
            }
        },
        "entity.SearchRequest": {
            "type": "object",
            "properties": {
                "query": {
                    "type": "string"
                }
            }
        }
    },
    "securityDefinitions": {
        "BearerAuth": {
            "description": "Type \"Bearer\" followed by a space and the JWT token.",
            "type": "apiKey",
            "name": "Authorization",
            "in": "header"
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "",
	BasePath:         "",
	Schemes:          []string{},
	Title:            "Geo Microservices API",
	Description:      "API для работы с геоданными",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
