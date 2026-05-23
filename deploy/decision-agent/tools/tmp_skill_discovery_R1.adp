{
  "toolbox": {
    "configs": [
      {
        "box_id": "f182d75a-7fd2-421a-a8e0-7064d75e39af",
        "box_name": "tmp_skill_discovery_R1",
        "box_desc": "临时工具箱 让LLM发现可用skill_id列表 排查R1用 稳定后删除",
        "box_svc_url": "http://agent-operator-integration:9000",
        "status": "published",
        "category_type": "other_category",
        "category_name": "未分类",
        "is_internal": false,
        "source": "custom",
        "tools": [
          {
            "tool_id": "ca4bc23d-e383-4d9e-9904-3a36d36e1680",
            "name": "list_skills",
            "description": "列出当前业务域下所有已注册的 skill，返回每条 skill 的 skill_id、name、description、status。使用场景：当你需要调用 builtin_skill_load(skill_id) 但不知道 skill_id 时，必须先调用本工具获取列表，再根据 description 匹配任务目标，选择合适的 skill 进行加载。**切勿自行臆测 skill_id；skill_id 必须是 UUID 格式。**",
            "status": "disabled",
            "metadata_type": "openapi",
            "metadata": {
              "version": "1c9d62ad-1e7d-4460-8755-daff2fa32ee4",
              "summary": "list_skills",
              "description": "列出当前业务域下所有已注册的 skill，返回每条 skill 的 skill_id、name、description、status。使用场景：当你需要调用 builtin_skill_load(skill_id) 但不知道 skill_id 时，必须先调用本工具获取列表，再根据 description 匹配任务目标，选择合适的 skill 进行加载。**切勿自行臆测 skill_id；skill_id 必须是 UUID 格式。**",
              "server_url": "http://agent-operator-integration:9000",
              "path": "/api/agent-operator-integration/v1/skills",
              "method": "GET",
              "create_time": 1776824809491349500,
              "update_time": 1776824809491349500,
              "create_user": "266c6a42-6131-4d62-8f39-853e7093701c",
              "update_user": "266c6a42-6131-4d62-8f39-853e7093701c",
              "api_spec": {
                "parameters": [
                  {
                    "name": "page",
                    "in": "query",
                    "description": "页码，从 1 开始",
                    "required": false,
                    "schema": {
                      "default": 1,
                      "type": "integer"
                    }
                  },
                  {
                    "name": "page_size",
                    "in": "query",
                    "description": "每页条数",
                    "required": false,
                    "schema": {
                      "default": 50,
                      "type": "integer"
                    }
                  },
                  {
                    "name": "all",
                    "in": "query",
                    "description": "是否跨所有业务域返回（false=仅当前业务域）",
                    "required": false,
                    "schema": {
                      "default": false,
                      "type": "boolean"
                    }
                  },
                  {
                    "name": "x-account-id",
                    "in": "header",
                    "description": "调用方账户 ID（平台自动注入）",
                    "required": false,
                    "schema": {
                      "type": "string"
                    }
                  },
                  {
                    "name": "x-account-type",
                    "in": "header",
                    "description": "调用方账户类型（平台自动注入）",
                    "required": false,
                    "schema": {
                      "type": "string"
                    }
                  },
                  {
                    "name": "x-business-domain",
                    "in": "header",
                    "description": "业务域 ID（平台自动注入）",
                    "required": false,
                    "schema": {
                      "type": "string"
                    }
                  }
                ],
                "request_body": {
                  "description": "",
                  "content": {},
                  "required": false
                },
                "responses": [
                  {
                    "status_code": "200",
                    "description": "skill 列表",
                    "content": {
                      "application/json": {
                        "schema": {
                          "properties": {
                            "data": {
                              "items": {
                                "properties": {
                                  "description": {
                                    "description": "skill 功能说明，用于判断是否适合当前任务",
                                    "type": "string"
                                  },
                                  "name": {
                                    "description": "skill 显示名",
                                    "type": "string"
                                  },
                                  "skill_id": {
                                    "description": "skill 唯一标识（UUID）",
                                    "type": "string"
                                  },
                                  "status": {
                                    "description": "published / offline",
                                    "type": "string"
                                  }
                                },
                                "type": "object"
                              },
                              "type": "array"
                            },
                            "page": {
                              "type": "integer"
                            },
                            "page_size": {
                              "type": "integer"
                            },
                            "total": {
                              "type": "integer"
                            }
                          },
                          "type": "object"
                        }
                      }
                    }
                  }
                ],
                "components": {
                  "schemas": {}
                },
                "callbacks": null,
                "security": null,
                "tags": null,
                "external_docs": null
              }
            },
            "use_rule": "",
            "global_parameters": {
              "name": "",
              "description": "",
              "required": false,
              "in": "",
              "type": "",
              "value": null
            },
            "create_time": 1776824809491562000,
            "update_time": 1776915259189479700,
            "create_user": "266c6a42-6131-4d62-8f39-853e7093701c",
            "update_user": "266c6a42-6131-4d62-8f39-853e7093701c",
            "extend_info": null,
            "resource_object": "tool",
            "source_id": "1c9d62ad-1e7d-4460-8755-daff2fa32ee4",
            "source_type": "openapi",
            "script_type": "",
            "code": "",
            "dependencies": [],
            "dependencies_url": ""
          },
          {
            "tool_id": "51382ef3-b35b-44a6-8a53-c670cbf53f10",
            "name": "list_skills_v2",
            "description": "列出当前业务域下所有已注册的 skill，返回每条 skill 的 skill_id、name、description、status。使用场景：当你需要调用 builtin_skill_load(skill_id) 但不知道 skill_id 时，必须先调用本工具获取列表，再根据 description 匹配任务目标，选择合适的 skill 进行加载。**切勿自行臆测 skill_id；skill_id 必须是 UUID 格式。**",
            "status": "enabled",
            "metadata_type": "openapi",
            "metadata": {
              "version": "6e7048b0-b98e-46a9-b31f-b6564ae31604",
              "summary": "list_skills_v2",
              "description": "列出当前业务域下所有已注册的 skill，返回每条 skill 的 skill_id、name、description、status。使用场景：当你需要调用 builtin_skill_load(skill_id) 但不知道 skill_id 时，必须先调用本工具获取列表，再根据 description 匹配任务目标，选择合适的 skill 进行加载。**切勿自行臆测 skill_id；skill_id 必须是 UUID 格式。**",
              "server_url": "http://agent-operator-integration:9000",
              "path": "/api/agent-operator-integration/v1/skills",
              "method": "GET",
              "create_time": 1776827592637636400,
              "update_time": 1776827592637636400,
              "create_user": "266c6a42-6131-4d62-8f39-853e7093701c",
              "update_user": "266c6a42-6131-4d62-8f39-853e7093701c",
              "api_spec": {
                "parameters": [
                  {
                    "name": "page",
                    "in": "query",
                    "description": "页码，从 1 开始",
                    "required": false,
                    "schema": {
                      "default": 1,
                      "type": "integer"
                    }
                  },
                  {
                    "name": "page_size",
                    "in": "query",
                    "description": "每页条数",
                    "required": false,
                    "schema": {
                      "default": 50,
                      "type": "integer"
                    }
                  },
                  {
                    "name": "all",
                    "in": "query",
                    "description": "是否跨所有业务域返回（false=仅当前业务域）",
                    "required": false,
                    "schema": {
                      "default": false,
                      "type": "boolean"
                    }
                  },
                  {
                    "name": "X-Authorization",
                    "in": "header",
                    "description": "OAuth Bearer token（格式：Bearer <token>）",
                    "required": false,
                    "schema": {
                      "type": "string"
                    }
                  },
                  {
                    "name": "x-account-id",
                    "in": "header",
                    "description": "调用方账户 ID（平台自动注入）",
                    "required": false,
                    "schema": {
                      "type": "string"
                    }
                  },
                  {
                    "name": "x-account-type",
                    "in": "header",
                    "description": "调用方账户类型（平台自动注入）",
                    "required": false,
                    "schema": {
                      "type": "string"
                    }
                  },
                  {
                    "name": "x-business-domain",
                    "in": "header",
                    "description": "业务域 ID（平台自动注入）",
                    "required": false,
                    "schema": {
                      "type": "string"
                    }
                  }
                ],
                "request_body": {
                  "description": "",
                  "content": {},
                  "required": false
                },
                "responses": [
                  {
                    "status_code": "200",
                    "description": "skill 列表",
                    "content": {
                      "application/json": {
                        "schema": {
                          "properties": {
                            "data": {
                              "items": {
                                "properties": {
                                  "description": {
                                    "description": "skill 功能说明，用于判断是否适合当前任务",
                                    "type": "string"
                                  },
                                  "name": {
                                    "description": "skill 显示名",
                                    "type": "string"
                                  },
                                  "skill_id": {
                                    "description": "skill 唯一标识（UUID）",
                                    "type": "string"
                                  },
                                  "status": {
                                    "description": "published / offline",
                                    "type": "string"
                                  }
                                },
                                "type": "object"
                              },
                              "type": "array"
                            },
                            "page": {
                              "type": "integer"
                            },
                            "page_size": {
                              "type": "integer"
                            },
                            "total": {
                              "type": "integer"
                            }
                          },
                          "type": "object"
                        }
                      }
                    }
                  }
                ],
                "components": {
                  "schemas": {}
                },
                "callbacks": null,
                "security": null,
                "tags": null,
                "external_docs": null
              }
            },
            "use_rule": "",
            "global_parameters": {
              "name": "",
              "description": "",
              "required": false,
              "in": "",
              "type": "",
              "value": null
            },
            "create_time": 1776827592637820700,
            "update_time": 1776827607071534300,
            "create_user": "266c6a42-6131-4d62-8f39-853e7093701c",
            "update_user": "266c6a42-6131-4d62-8f39-853e7093701c",
            "extend_info": null,
            "resource_object": "tool",
            "source_id": "6e7048b0-b98e-46a9-b31f-b6564ae31604",
            "source_type": "openapi",
            "script_type": "",
            "code": "",
            "dependencies": [],
            "dependencies_url": ""
          }
        ],
        "create_time": 1776824784296692500,
        "update_time": 1776824816779447600,
        "create_user": "266c6a42-6131-4d62-8f39-853e7093701c",
        "update_user": "266c6a42-6131-4d62-8f39-853e7093701c",
        "metadata_type": "openapi"
      }
    ]
  }
}