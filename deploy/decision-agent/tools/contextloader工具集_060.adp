{
  "toolbox": {
    "configs": [
      {
        "box_id": "e521d454-4a0b-4dc9-8a28-d0986de1cef9",
        "box_name": "contextloader工具集_060",
        "box_desc": "contextloader工具集",
        "box_svc_url": "http://agent-retrieval:30779",
        "status": "published",
        "category_type": "data_query",
        "category_name": "数据查询",
        "is_internal": false,
        "source": "custom",
        "tools": [
          {
            "tool_id": "05275bb1-46e2-4727-9c6f-97d9ea0af94b",
            "name": "get_logic_properties_values",
            "description": "根据 query 生成 dynamic_params，批量查询指定对象的逻辑属性值。",
            "status": "enabled",
            "metadata_type": "openapi",
            "metadata": {
              "version": "35fdd885-a378-4d1a-a4bb-2af441de55a5",
              "summary": "get_logic_properties_values",
              "description": "根据 query 生成 dynamic_params，批量查询指定对象的逻辑属性值。",
              "server_url": "http://agent-retrieval:30779",
              "path": "/api/agent-retrieval/in/v1/kn/logic-property-resolver",
              "method": "POST",
              "create_time": 1776908488240137500,
              "update_time": 1776908488240137500,
              "create_user": "266c6a42-6131-4d62-8f39-853e7093701c",
              "update_user": "266c6a42-6131-4d62-8f39-853e7093701c",
              "api_spec": {
                "parameters": [
                  {
                    "name": "x-account-id",
                    "in": "header",
                    "description": "账户ID，用于内部服务调用时传递账户信息",
                    "required": false,
                    "schema": {
                      "type": "string"
                    }
                  },
                  {
                    "name": "x-account-type",
                    "in": "header",
                    "description": "账户类型：user(用户), app(应用), anonymous(匿名)",
                    "required": false,
                    "schema": {
                      "enum": [
                        "user",
                        "app",
                        "anonymous"
                      ],
                      "type": "string"
                    }
                  },
                  {
                    "name": "response_format",
                    "in": "query",
                    "description": "响应格式：json 或 toon，默认 json",
                    "required": false,
                    "schema": {
                      "default": "json",
                      "enum": [
                        "json",
                        "toon"
                      ],
                      "type": "string"
                    }
                  }
                ],
                "request_body": {
                  "description": "",
                  "content": {
                    "application/json": {
                      "examples": {
                        "示例": {
                          "value": {
                            "_instance_identities": [
                              {
                                "company_id": "company_000001"
                              }
                            ],
                            "kn_id": "kn_medical",
                            "ot_id": "company",
                            "properties": [
                              "approved_drug_count",
                              "business_health_score"
                            ],
                            "query": "最近一年这些药企的药品上市数量和健康度"
                          }
                        }
                      },
                      "schema": {
                        "$ref": "#/components/schemas/ResolveLogicPropertiesRequest"
                      }
                    }
                  },
                  "required": false
                },
                "responses": [
                  {
                    "status_code": "200",
                    "description": "ok",
                    "content": {
                      "application/json": {
                        "schema": {
                          "$ref": "#/components/schemas/ResolveLogicPropertiesResponse"
                        }
                      }
                    }
                  },
                  {
                    "status_code": "400",
                    "description": "bad request",
                    "content": {
                      "application/json": {
                        "schema": {
                          "$ref": "#/components/schemas/Error"
                        }
                      }
                    }
                  },
                  {
                    "status_code": "500",
                    "description": "internal error",
                    "content": {
                      "application/json": {
                        "schema": {
                          "$ref": "#/components/schemas/Error"
                        }
                      }
                    }
                  }
                ],
                "components": {
                  "schemas": {
                    "ResolveOptions": {
                      "description": "【可选配置】控制接口行为的高级选项\n",
                      "properties": {
                        "return_debug": {
                          "type": "boolean",
                          "description": "是否返回 debug（dynamic_params、warnings 等）。默认 false"
                        }
                      },
                      "type": "object"
                    },
                    "ResolveLogicPropertiesResponse": {
                      "description": "成功返回 datas；缺参时返回 error_code、missing（含 hint）",
                      "oneOf": [
                        {
                          "$ref": "#/components/schemas/ObjectPropertiesValuesResponse"
                        },
                        {
                          "$ref": "#/components/schemas/MissingParamsError"
                        }
                      ]
                    },
                    "ObjectPropertiesValuesResponse": {
                      "properties": {
                        "debug": {
                          "$ref": "#/components/schemas/ResolveDebugInfo"
                        },
                        "datas": {
                          "type": "array",
                          "description": "与 _instance_identities 顺序对齐，每项含主键和请求的 properties",
                          "items": {
                            "type": "object"
                          }
                        }
                      },
                      "type": "object",
                      "required": [
                        "datas"
                      ]
                    },
                    "ResolveDebugInfo": {
                      "type": "object",
                      "properties": {
                        "trace_id": {
                          "type": "string"
                        },
                        "warnings": {
                          "type": "array",
                          "items": {
                            "type": "string"
                          }
                        },
                        "dynamic_params": {
                          "type": "object"
                        },
                        "now_ms": {
                          "type": "integer"
                        }
                      }
                    },
                    "MissingParamsError": {
                      "type": "object",
                      "properties": {
                        "error_code": {
                          "type": "string"
                        },
                        "message": {
                          "type": "string"
                        },
                        "missing": {
                          "items": {
                            "type": "object",
                            "properties": {
                              "params": {
                                "items": {
                                  "properties": {
                                    "hint": {
                                      "type": "string"
                                    },
                                    "name": {
                                      "type": "string"
                                    }
                                  },
                                  "type": "object"
                                },
                                "type": "array"
                              },
                              "property": {
                                "type": "string"
                              }
                            }
                          },
                          "type": "array"
                        }
                      }
                    },
                    "Error": {
                      "type": "object",
                      "properties": {
                        "message": {
                          "type": "string"
                        },
                        "error_code": {
                          "type": "string"
                        }
                      }
                    },
                    "ResolveLogicPropertiesRequest": {
                      "properties": {
                        "ot_id": {
                          "description": "对象类ID。例 company、drug",
                          "type": "string"
                        },
                        "properties": {
                          "description": "逻辑属性名列表（metric/operator）。自动生成 dynamic_params 并查询。",
                          "items": {
                            "type": "string"
                          },
                          "type": "array"
                        },
                        "query": {
                          "description": "用户查询，需含时间（如\"最近一年\"）、统计维度、业务上下文，用于生成 dynamic_params",
                          "type": "string"
                        },
                        "_instance_identities": {
                          "items": {
                            "type": "object"
                          },
                          "type": "array",
                          "description": "对象实例标识数组。**必须从上游提取，不可臆造。** 流程：先调 query_object_instance 或 query_instance_subgraph → 从每个对象的 _instance_identity 字段取值 → 按原顺序组成数组传入。"
                        },
                        "additional_context": {
                          "description": "可选。补充上下文，如 timezone、instant、step、对象属性等，帮助生成 dynamic_params。",
                          "type": "string"
                        },
                        "kn_id": {
                          "type": "string",
                          "description": "知识网络ID。例 kn_medical"
                        },
                        "options": {
                          "$ref": "#/components/schemas/ResolveOptions"
                        }
                      },
                      "type": "object",
                      "required": [
                        "kn_id",
                        "ot_id",
                        "query",
                        "_instance_identities",
                        "properties"
                      ]
                    }
                  }
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
            "create_time": 1776908488245462300,
            "update_time": 1776908488245462300,
            "create_user": "266c6a42-6131-4d62-8f39-853e7093701c",
            "update_user": "266c6a42-6131-4d62-8f39-853e7093701c",
            "extend_info": null,
            "resource_object": "tool",
            "source_id": "35fdd885-a378-4d1a-a4bb-2af441de55a5",
            "source_type": "openapi",
            "script_type": "",
            "code": "",
            "dependencies": [],
            "dependencies_url": ""
          },
          {
            "tool_id": "2fd071fa-a696-4fee-91e1-5b2dc190e88b",
            "name": "kn_schema_search",
            "description": "基于用户查询意图，返回业务知识网络中相关的概念信息",
            "status": "enabled",
            "metadata_type": "openapi",
            "metadata": {
              "version": "9c5adc7c-6591-4948-8976-14a5ff99f287",
              "summary": "kn_schema_search",
              "description": "基于用户查询意图，返回业务知识网络中相关的概念信息",
              "server_url": "http://agent-retrieval:30779",
              "path": "/api/agent-retrieval/in/v1/kn/semantic-search",
              "method": "POST",
              "create_time": 1776908488240137500,
              "update_time": 1776908488240137500,
              "create_user": "266c6a42-6131-4d62-8f39-853e7093701c",
              "update_user": "266c6a42-6131-4d62-8f39-853e7093701c",
              "api_spec": {
                "parameters": [
                  {
                    "name": "x-account-id",
                    "in": "header",
                    "description": "账户ID，用于内部服务调用时传递账户信息",
                    "required": false,
                    "schema": {
                      "type": "string"
                    }
                  },
                  {
                    "name": "x-account-type",
                    "in": "header",
                    "description": "账户类型：user(用户), app(应用), anonymous(匿名)",
                    "required": false,
                    "schema": {
                      "enum": [
                        "user",
                        "app",
                        "anonymous"
                      ],
                      "type": "string"
                    }
                  },
                  {
                    "name": "response_format",
                    "in": "query",
                    "description": "响应格式：json 或 toon，默认 json",
                    "required": false,
                    "schema": {
                      "default": "json",
                      "enum": [
                        "json",
                        "toon"
                      ],
                      "type": "string"
                    }
                  }
                ],
                "request_body": {
                  "description": "",
                  "content": {
                    "application/json": {
                      "schema": {
                        "$ref": "#/components/schemas/SemanticSearchRequest"
                      }
                    }
                  },
                  "required": false
                },
                "responses": [
                  {
                    "status_code": "200",
                    "description": "成功返回相关概念信息",
                    "content": {
                      "application/json": {
                        "schema": {
                          "$ref": "#/components/schemas/SemanticSearchResponse"
                        }
                      }
                    }
                  },
                  {
                    "status_code": "400",
                    "description": "参数错误",
                    "content": {
                      "application/json": {
                        "schema": {
                          "$ref": "#/components/schemas/ErrorResponse"
                        }
                      }
                    }
                  },
                  {
                    "status_code": "500",
                    "description": "服务器内部错误",
                    "content": {
                      "application/json": {
                        "schema": {
                          "$ref": "#/components/schemas/ErrorResponse"
                        }
                      }
                    }
                  }
                ],
                "components": {
                  "schemas": {
                    "ResourceInfo": {
                      "description": "数据来源信息",
                      "properties": {
                        "name": {
                          "type": "string",
                          "description": "视图名称"
                        },
                        "type": {
                          "type": "string",
                          "description": "数据来源类型"
                        },
                        "id": {
                          "description": "数据视图id",
                          "type": "string"
                        }
                      },
                      "type": "object"
                    },
                    "Concept": {
                      "type": "object",
                      "properties": {
                        "concept_name": {
                          "description": "概念类名称",
                          "type": "string"
                        },
                        "concept_type": {
                          "type": "string",
                          "description": "概念类型",
                          "enum": [
                            "object_type",
                            "relation_type",
                            "action_type"
                          ]
                        },
                        "concept_detail": {
                          "description": "概念类详情，根据concept_type返回不同结构：\n- 当concept_type为\"object_type\"时，返回ObjectTypeDetail结构，包含对象类的完整信息\n- 当concept_type为\"relation_type\"时，返回RelationTypeDetail结构，包含关系类的完整信息\n- 当concept_type为\"action_type\"时，返回ActionTypeDetail结构，包含行动类的完整信息\n",
                          "oneOf": [
                            {
                              "$ref": "#/components/schemas/ObjectTypeDetail"
                            },
                            {
                              "$ref": "#/components/schemas/RelationTypeDetail"
                            },
                            {
                              "$ref": "#/components/schemas/ActionTypeDetail"
                            }
                          ]
                        },
                        "concept_id": {
                          "description": "概念类ID",
                          "type": "string"
                        }
                      }
                    },
                    "RelationTypeDetail": {
                      "type": "object",
                      "description": "关系类概念详情",
                      "properties": {
                        "_score": {
                          "description": "分数",
                          "type": "number",
                          "format": "float"
                        },
                        "type": {
                          "description": "关系类型",
                          "type": "string"
                        },
                        "tags": {
                          "description": "标签",
                          "items": {
                            "type": "string"
                          },
                          "type": "array"
                        },
                        "comment": {
                          "description": "备注",
                          "type": "string"
                        },
                        "target_object_type_id": {
                          "type": "string",
                          "description": "目标对象类ID"
                        },
                        "source_object_type_id": {
                          "description": "起点对象类ID",
                          "type": "string"
                        },
                        "name": {
                          "type": "string",
                          "description": "关系类名称"
                        },
                        "id": {
                          "description": "关系类id",
                          "type": "string"
                        },
                        "module_type": {
                          "type": "string",
                          "description": "模块类型"
                        }
                      }
                    },
                    "DataProperty": {
                      "properties": {
                        "type": {
                          "description": "属性数据类型",
                          "type": "string"
                        },
                        "comment": {
                          "type": "string",
                          "description": "备注"
                        },
                        "condition_operations": {
                          "description": "该数据属性支持的查询条件操作符列表。\n",
                          "items": {
                            "type": "string",
                            "enum": [
                              "==",
                              "!=",
                              ">",
                              "<",
                              ">=",
                              "<=",
                              "in",
                              "not_in",
                              "like",
                              "not_like",
                              "range",
                              "out_range",
                              "exist",
                              "not_exist",
                              "regex",
                              "match",
                              "knn"
                            ]
                          },
                          "type": "array"
                        },
                        "display_name": {
                          "type": "string",
                          "description": "属性显示名称"
                        },
                        "mapped_field": {
                          "description": "视图字段信息"
                        },
                        "name": {
                          "description": "属性名称",
                          "type": "string"
                        }
                      },
                      "type": "object",
                      "description": "数据属性结构定义"
                    },
                    "ObjectTypeDetail": {
                      "description": "对象类概念详情",
                      "properties": {
                        "_score": {
                          "type": "number",
                          "format": "float",
                          "description": "分数"
                        },
                        "tags": {
                          "type": "array",
                          "description": "标签",
                          "items": {
                            "type": "string"
                          }
                        },
                        "data_properties": {
                          "description": "数据属性",
                          "items": {
                            "$ref": "#/components/schemas/DataProperty"
                          },
                          "type": "array"
                        },
                        "name": {
                          "description": "对象名称",
                          "type": "string"
                        },
                        "comment": {
                          "description": "备注",
                          "type": "string"
                        },
                        "primary_keys": {
                          "description": "主键字段",
                          "items": {
                            "type": "string"
                          },
                          "type": "array"
                        },
                        "logic_properties": {
                          "items": {
                            "type": "object"
                          },
                          "type": "array",
                          "description": "逻辑属性"
                        },
                        "module_type": {
                          "type": "string",
                          "description": "模块类型"
                        },
                        "id": {
                          "description": "对象id",
                          "type": "string"
                        },
                        "data_source": {
                          "$ref": "#/components/schemas/ResourceInfo"
                        }
                      },
                      "type": "object"
                    },
                    "SemanticSearchRequest": {
                      "type": "object",
                      "required": [
                        "query",
                        "kn_id"
                      ],
                      "properties": {
                        "kn_id": {
                          "type": "string",
                          "description": "业务知识网络ID"
                        },
                        "max_concepts": {
                          "description": "最大返回概念数量",
                          "default": 10,
                          "type": "integer"
                        },
                        "query": {
                          "description": "用户自然语言查询",
                          "type": "string"
                        },
                        "rerank_action": {
                          "type": "string",
                          "description": "重排动作",
                          "default": "default",
                          "enum": [
                            "default",
                            "vector",
                            "llm"
                          ]
                        },
                        "search_scope": {
                          "$ref": "#/components/schemas/SearchScope"
                        }
                      }
                    },
                    "SemanticSearchResponse": {
                      "properties": {
                        "concepts": {
                          "type": "array",
                          "items": {
                            "$ref": "#/components/schemas/Concept"
                          }
                        }
                      },
                      "type": "object"
                    },
                    "SearchScope": {
                      "properties": {
                        "include_action_types": {
                          "description": "是否包含行作类",
                          "type": "boolean"
                        },
                        "include_object_types": {
                          "type": "boolean",
                          "description": "是否包含对象类"
                        },
                        "include_relation_types": {
                          "type": "boolean",
                          "description": "是否包含关系类"
                        },
                        "concept_groups": {
                          "description": "限定的概念分组",
                          "items": {
                            "type": "string"
                          },
                          "type": "array"
                        }
                      },
                      "type": "object",
                      "description": "【可选】搜索域配置\n"
                    },
                    "ErrorResponse": {
                      "type": "object",
                      "properties": {
                        "detail": {
                          "description": "错误详情",
                          "type": "object"
                        },
                        "link": {
                          "type": "string",
                          "description": "错误链接"
                        },
                        "solution": {
                          "description": "解决方案",
                          "type": "string"
                        },
                        "code": {
                          "type": "string",
                          "description": "错误码"
                        },
                        "description": {
                          "type": "string",
                          "description": "错误描述"
                        }
                      }
                    },
                    "ActionTypeDetail": {
                      "type": "object",
                      "description": "行动类概念详情",
                      "properties": {
                        "module_type": {
                          "type": "string",
                          "description": "模块类型"
                        },
                        "name": {
                          "type": "string",
                          "description": "行动类名称"
                        },
                        "object_type_id": {
                          "description": "行动类所绑定的对象类ID",
                          "type": "string"
                        },
                        "tags": {
                          "items": {
                            "type": "string"
                          },
                          "type": "array",
                          "description": "标签"
                        },
                        "_score": {
                          "format": "float",
                          "description": "分数",
                          "type": "number"
                        },
                        "comment": {
                          "type": "string",
                          "description": "备注"
                        },
                        "id": {
                          "type": "string",
                          "description": "行动类ID"
                        }
                      }
                    }
                  }
                },
                "callbacks": null,
                "security": null,
                "tags": [
                  "SemanticSearch"
                ],
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
            "create_time": 1776908488245462300,
            "update_time": 1776908488245462300,
            "create_user": "266c6a42-6131-4d62-8f39-853e7093701c",
            "update_user": "266c6a42-6131-4d62-8f39-853e7093701c",
            "extend_info": null,
            "resource_object": "tool",
            "source_id": "9c5adc7c-6591-4948-8976-14a5ff99f287",
            "source_type": "openapi",
            "script_type": "",
            "code": "",
            "dependencies": [],
            "dependencies_url": ""
          },
          {
            "tool_id": "52b35175-cee3-41ea-91c0-1d70e8371f9c",
            "name": "kn_search",
            "description": "基于知识网络的智能检索工具，支持传入完整的问题或一个或多个关键词，能够检索问题或关键词的属性信息和上下文信息。\r\n支持概念召回、语义实例召回、多轮对话等功能。\r\n",
            "status": "enabled",
            "metadata_type": "openapi",
            "metadata": {
              "version": "1ef67af5-ef70-4859-b1f9-816141a08984",
              "summary": "kn_search",
              "description": "基于知识网络的智能检索工具，支持传入完整的问题或一个或多个关键词，能够检索问题或关键词的属性信息和上下文信息。\n支持概念召回、语义实例召回、多轮对话等功能。\n",
              "server_url": "http://agent-retrieval:30779",
              "path": "/api/agent-retrieval/in/v1/kn/kn_search",
              "method": "POST",
              "create_time": 1776908488240137500,
              "update_time": 1776908488240137500,
              "create_user": "266c6a42-6131-4d62-8f39-853e7093701c",
              "update_user": "266c6a42-6131-4d62-8f39-853e7093701c",
              "api_spec": {
                "parameters": [
                  {
                    "name": "x-account-id",
                    "in": "header",
                    "description": "账户ID，用于内部服务调用时传递账户信息",
                    "required": false,
                    "schema": {
                      "type": "string"
                    }
                  },
                  {
                    "name": "x-account-type",
                    "in": "header",
                    "description": "账户类型：user(用户), app(应用), anonymous(匿名)",
                    "required": false,
                    "schema": {
                      "enum": [
                        "user",
                        "app",
                        "anonymous"
                      ],
                      "type": "string"
                    }
                  },
                  {
                    "name": "response_format",
                    "in": "query",
                    "description": "响应格式：json 或 toon，默认 json",
                    "required": false,
                    "schema": {
                      "default": "json",
                      "enum": [
                        "json",
                        "toon"
                      ],
                      "type": "string"
                    }
                  }
                ],
                "request_body": {
                  "description": "kn_search 请求体",
                  "content": {
                    "application/json": {
                      "schema": {
                        "$ref": "#/components/schemas/KnSearchRequest"
                      }
                    }
                  },
                  "required": false
                },
                "responses": [
                  {
                    "status_code": "500",
                    "description": "服务器内部错误",
                    "content": {
                      "application/json": {
                        "schema": {
                          "$ref": "#/components/schemas/ErrorResponse"
                        }
                      }
                    }
                  },
                  {
                    "status_code": "200",
                    "description": "成功返回检索结果",
                    "content": {
                      "application/json": {
                        "schema": {
                          "$ref": "#/components/schemas/KnSearchResponse"
                        }
                      }
                    }
                  },
                  {
                    "status_code": "400",
                    "description": "参数错误",
                    "content": {
                      "application/json": {
                        "schema": {
                          "$ref": "#/components/schemas/ErrorResponse"
                        }
                      }
                    }
                  }
                ],
                "components": {
                  "schemas": {
                    "Node": {
                      "description": "节点数据，至少包含 object_type_id、<object_type_id>_name、unique_identities",
                      "properties": {
                        "object_type_id": {
                          "type": "string"
                        },
                        "unique_identities": {
                          "type": "object",
                          "description": "对象的唯一标识信息"
                        }
                      },
                      "type": "object"
                    },
                    "ErrorResponse": {
                      "properties": {
                        "message": {
                          "description": "错误详情",
                          "type": "string"
                        },
                        "error": {
                          "description": "错误信息",
                          "type": "string"
                        }
                      },
                      "type": "object"
                    },
                    "KnSearchRequest": {
                      "type": "object",
                      "required": [
                        "query",
                        "kn_id"
                      ],
                      "properties": {
                        "query": {
                          "type": "string",
                          "description": "用户查询问题或关键词，多个关键词之间用空格隔开"
                        },
                        "retrieval_config": {
                          "description": "召回配置参数，用于控制不同类型的召回场景（概念召回、语义实例召回、属性过滤）。如果不提供，将使用系统默认配置。",
                          "properties": {
                            "concept_retrieval": {
                              "$ref": "#/components/schemas/ConceptRetrievalConfig"
                            }
                          },
                          "type": "object"
                        },
                        "enable_rerank": {
                          "description": "是否启用重排序。如果为true，则启用重排序。",
                          "default": true,
                          "type": "boolean"
                        },
                        "kn_id": {
                          "description": "指定的知识网络ID，必须传递",
                          "type": "string"
                        },
                        "only_schema": {
                          "default": false,
                          "type": "boolean",
                          "description": "是否只召回概念（schema），不召回语义实例。如果为True，则只返回object_types、relation_types和action_types，不返回nodes。"
                        }
                      }
                    },
                    "RelationType": {
                      "required": [
                        "concept_id",
                        "concept_name",
                        "source_object_type_id",
                        "target_object_type_id"
                      ],
                      "properties": {
                        "concept_type": {
                          "description": "概念类型: relation_type",
                          "type": "string"
                        },
                        "source_object_type_id": {
                          "type": "string",
                          "description": "源对象类型ID"
                        },
                        "target_object_type_id": {
                          "type": "string",
                          "description": "目标对象类型ID"
                        },
                        "concept_id": {
                          "type": "string",
                          "description": "概念ID"
                        },
                        "concept_name": {
                          "type": "string",
                          "description": "概念名称"
                        }
                      },
                      "type": "object"
                    },
                    "ActionType": {
                      "type": "object",
                      "description": "操作类型信息。精简模式（schema_brief=True）下仅包含：id, name, action_type, object_type_id, object_type_name, comment, tags, kn_id",
                      "properties": {
                        "action_type": {
                          "type": "string",
                          "description": "操作类型（如：add, modify等）"
                        },
                        "comment": {
                          "type": "string",
                          "description": "注释说明"
                        },
                        "id": {
                          "type": "string",
                          "description": "操作类型ID"
                        },
                        "kn_id": {
                          "description": "知识网络ID",
                          "type": "string"
                        },
                        "name": {
                          "description": "操作类型名称",
                          "type": "string"
                        },
                        "object_type_id": {
                          "description": "对象类型ID",
                          "type": "string"
                        },
                        "object_type_name": {
                          "type": "string",
                          "description": "对象类型名称"
                        },
                        "tags": {
                          "type": "array",
                          "description": "标签列表",
                          "items": {
                            "type": "string"
                          }
                        }
                      }
                    },
                    "KnSearchResponse": {
                      "description": "检索结果，返回object_types/relation_types/action_types，并返回语义实例nodes/message。\n多轮时由concept_retrieval.return_union控制 nodes 的并集/增量。\n",
                      "properties": {
                        "object_types": {
                          "items": {
                            "$ref": "#/components/schemas/ObjectType"
                          },
                          "type": "array",
                          "description": "对象类型列表（概念召回时返回）。\n当schema_brief=True时，仅包含：concept_id, concept_name, comment, data_properties（仅name和display_name）, logic_properties（仅name和display_name）, sample_data（当include_sample_data=True时）。\n当schema_brief=False时，包含完整字段（包括primary_keys, display_key, sample_data等）\n"
                        },
                        "relation_types": {
                          "description": "关系类型列表（概念召回时返回）。\n精简模式和完整模式均包含：concept_id, concept_name, source_object_type_id, target_object_type_id\n",
                          "items": {
                            "$ref": "#/components/schemas/RelationType"
                          },
                          "type": "array"
                        },
                        "action_types": {
                          "items": {
                            "$ref": "#/components/schemas/ActionType"
                          },
                          "type": "array",
                          "description": "操作类型列表（概念召回时返回）。\n当schema_brief=True时，每个action_type仅包含以下字段：id, name, action_type, object_type_id, object_type_name, comment, tags, kn_id\n"
                        },
                        "message": {
                          "type": "string",
                          "description": "提示信息（例如未召回到实例数据时返回原因说明）"
                        },
                        "nodes": {
                          "description": "语义实例召回结果（当不提供conditions且召回到实例时返回），与条件召回节点风格对齐的扁平列表。\n每个节点至少包含 object_type_id、<object_type_id>_name、unique_identities\n",
                          "items": {
                            "$ref": "#/components/schemas/Node"
                          },
                          "type": "array"
                        }
                      },
                      "type": "object"
                    },
                    "ConceptRetrievalConfig": {
                      "description": "概念召回/概念流程配置参数（原最外层参数已收敛到此处）",
                      "properties": {
                        "include_sample_data": {
                          "default": false,
                          "type": "boolean",
                          "description": "是否获取对象类型的样例数据。True会为每个召回对象类型获取一条样例数据。"
                        },
                        "schema_brief": {
                          "type": "boolean",
                          "description": "概念召回时是否返回精简schema。True仅返回必要字段（概念ID/名称/关系source&target），不返回大字段。",
                          "default": true
                        },
                        "top_k": {
                          "type": "integer",
                          "description": "概念召回返回最相关关系类型数量（对象类型会随关系类型自动过滤）。",
                          "default": 10
                        }
                      },
                      "type": "object"
                    },
                    "DataProperty": {
                      "properties": {
                        "display_name": {
                          "type": "string",
                          "description": "属性显示名称"
                        },
                        "name": {
                          "type": "string",
                          "description": "属性名称"
                        },
                        "comment": {
                          "type": "string",
                          "description": "属性描述（非精简模式）"
                        }
                      },
                      "type": "object"
                    },
                    "LogicProperty": {
                      "type": "object",
                      "properties": {
                        "name": {
                          "type": "string",
                          "description": "属性名称"
                        },
                        "display_name": {
                          "description": "属性显示名称",
                          "type": "string"
                        }
                      }
                    },
                    "ObjectType": {
                      "properties": {
                        "concept_type": {
                          "type": "string",
                          "description": "概念类型: object_type"
                        },
                        "concept_id": {
                          "type": "string",
                          "description": "概念ID"
                        },
                        "sample_data": {
                          "type": "object",
                          "description": "样例数据（当include_sample_data=True时返回，无论schema_brief是否为True）"
                        },
                        "data_properties": {
                          "items": {
                            "$ref": "#/components/schemas/DataProperty"
                          },
                          "type": "array",
                          "description": "对象属性列表。精简模式下仅包含name和display_name字段（数量不截断）"
                        },
                        "logic_properties": {
                          "type": "array",
                          "description": "逻辑属性列表（指标等）。精简模式下仅包含name和display_name字段（数量不截断）",
                          "items": {
                            "$ref": "#/components/schemas/LogicProperty"
                          }
                        },
                        "comment": {
                          "type": "string",
                          "description": "概念描述"
                        },
                        "display_key": {
                          "type": "string",
                          "description": "显示字段名（用于获取instance_name）。仅当schema_brief=False时返回"
                        },
                        "primary_keys": {
                          "type": "array",
                          "description": "主键字段列表（支持多个主键）。仅当schema_brief=False时返回",
                          "items": {
                            "type": "string"
                          }
                        },
                        "concept_name": {
                          "description": "概念名称",
                          "type": "string"
                        }
                      },
                      "type": "object",
                      "required": [
                        "concept_id",
                        "concept_name"
                      ]
                    }
                  }
                },
                "callbacks": null,
                "security": null,
                "tags": [
                  "kn-search"
                ],
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
            "create_time": 1776908488245462300,
            "update_time": 1776908488245462300,
            "create_user": "266c6a42-6131-4d62-8f39-853e7093701c",
            "update_user": "266c6a42-6131-4d62-8f39-853e7093701c",
            "extend_info": null,
            "resource_object": "tool",
            "source_id": "1ef67af5-ef70-4859-b1f9-816141a08984",
            "source_type": "openapi",
            "script_type": "",
            "code": "",
            "dependencies": [],
            "dependencies_url": ""
          },
          {
            "tool_id": "5467284c-f25a-4665-8f05-a320dd6e1ec3",
            "name": "get_kn_index_build_status",
            "description": "查询最新50个构建任务的整体状态（按创建时间倒排）。如果所有任务都已完成则返回completed，如果有任务正在运行则返回running",
            "status": "enabled",
            "metadata_type": "openapi",
            "metadata": {
              "version": "572c393c-a788-4276-98fc-e463a13b315a",
              "summary": "get_kn_index_build_status",
              "description": "查询最新50个构建任务的整体状态（按创建时间倒排）。如果所有任务都已完成则返回completed，如果有任务正在运行则返回running",
              "server_url": "http://agent-retrieval:30779",
              "path": "/api/agent-retrieval/in/v1/kn/full_ontology_building_status",
              "method": "GET",
              "create_time": 1776908488240137500,
              "update_time": 1776908488240137500,
              "create_user": "266c6a42-6131-4d62-8f39-853e7093701c",
              "update_user": "266c6a42-6131-4d62-8f39-853e7093701c",
              "api_spec": {
                "parameters": [
                  {
                    "name": "x-account-id",
                    "in": "header",
                    "description": "账户ID，用于内部服务调用时传递账户信息",
                    "required": false,
                    "schema": {
                      "type": "string"
                    }
                  },
                  {
                    "name": "x-account-type",
                    "in": "header",
                    "description": "账户类型：user(用户), app(应用), anonymous(匿名)",
                    "required": false,
                    "schema": {
                      "enum": [
                        "user",
                        "app",
                        "anonymous"
                      ],
                      "type": "string"
                    }
                  },
                  {
                    "name": "kn_id",
                    "in": "query",
                    "description": "业务知识网络ID",
                    "required": true,
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
                    "description": "成功返回构建状态",
                    "content": {
                      "application/json": {
                        "example": {
                          "kn_id": "d5levlh818p1vl2slp60",
                          "state": "completed",
                          "state_detail": "All latest 50 jobs are completed"
                        },
                        "schema": {
                          "$ref": "#/components/schemas/BuildStatusSimpleResponse"
                        }
                      }
                    }
                  },
                  {
                    "status_code": "400",
                    "description": "参数错误",
                    "content": {
                      "application/json": {
                        "schema": {
                          "$ref": "#/components/schemas/ErrorResponse"
                        }
                      }
                    }
                  },
                  {
                    "status_code": "401",
                    "description": "未授权",
                    "content": {
                      "application/json": {
                        "schema": {
                          "$ref": "#/components/schemas/ErrorResponse"
                        }
                      }
                    }
                  },
                  {
                    "status_code": "500",
                    "description": "服务器内部错误",
                    "content": {
                      "application/json": {
                        "schema": {
                          "$ref": "#/components/schemas/ErrorResponse"
                        }
                      }
                    }
                  }
                ],
                "components": {
                  "schemas": {
                    "BuildStatusSimpleResponse": {
                      "properties": {
                        "state": {
                          "enum": [
                            "running",
                            "completed"
                          ],
                          "type": "string",
                          "description": "构建状态（running表示有任务正在运行，completed表示所有任务都已完成）"
                        },
                        "state_detail": {
                          "type": "string",
                          "description": "状态详情"
                        },
                        "kn_id": {
                          "type": "string",
                          "description": "业务知识网络ID"
                        }
                      },
                      "type": "object",
                      "description": "构建状态响应",
                      "required": [
                        "kn_id",
                        "state",
                        "state_detail"
                      ]
                    },
                    "ErrorResponse": {
                      "type": "object",
                      "properties": {
                        "description": {
                          "type": "string",
                          "description": "错误描述"
                        },
                        "detail": {
                          "description": "错误详情",
                          "type": "object"
                        },
                        "link": {
                          "description": "错误链接",
                          "type": "string"
                        },
                        "solution": {
                          "type": "string",
                          "description": "解决方案"
                        },
                        "code": {
                          "type": "string",
                          "description": "错误码"
                        }
                      }
                    }
                  }
                },
                "callbacks": null,
                "security": null,
                "tags": [
                  "OntologyJob"
                ],
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
            "create_time": 1776908488245462300,
            "update_time": 1776908488245462300,
            "create_user": "266c6a42-6131-4d62-8f39-853e7093701c",
            "update_user": "266c6a42-6131-4d62-8f39-853e7093701c",
            "extend_info": null,
            "resource_object": "tool",
            "source_id": "572c393c-a788-4276-98fc-e463a13b315a",
            "source_type": "openapi",
            "script_type": "",
            "code": "",
            "dependencies": [],
            "dependencies_url": ""
          },
          {
            "tool_id": "3c78c9da-0bd5-48b4-b23b-960972d4d2af",
            "name": "create_kn_index_build_job",
            "description": "创建一个全量构建业务知识网络的任务",
            "status": "enabled",
            "metadata_type": "openapi",
            "metadata": {
              "version": "32207cc2-79fe-4650-b5e2-dd96ee93fc1d",
              "summary": "create_kn_index_build_job",
              "description": "创建一个全量构建业务知识网络的任务",
              "server_url": "http://agent-retrieval:30779",
              "path": "/api/agent-retrieval/in/v1/kn/full_build_ontology",
              "method": "POST",
              "create_time": 1776908488240137500,
              "update_time": 1776908488240137500,
              "create_user": "266c6a42-6131-4d62-8f39-853e7093701c",
              "update_user": "266c6a42-6131-4d62-8f39-853e7093701c",
              "api_spec": {
                "parameters": [
                  {
                    "name": "x-account-id",
                    "in": "header",
                    "description": "账户ID，用于内部服务调用时传递账户信息",
                    "required": false,
                    "schema": {
                      "type": "string"
                    }
                  },
                  {
                    "name": "x-account-type",
                    "in": "header",
                    "description": "账户类型：user(用户), app(应用), anonymous(匿名)",
                    "required": false,
                    "schema": {
                      "enum": [
                        "user",
                        "app",
                        "anonymous"
                      ],
                      "type": "string"
                    }
                  }
                ],
                "request_body": {
                  "description": "",
                  "content": {
                    "application/json": {
                      "example": {
                        "kn_id": "kn_1234567890",
                        "name": "全量构建任务"
                      },
                      "schema": {
                        "$ref": "#/components/schemas/CreateJobRequest"
                      }
                    }
                  },
                  "required": false
                },
                "responses": [
                  {
                    "status_code": "201",
                    "description": "创建成功",
                    "content": {
                      "application/json": {
                        "example": {
                          "id": "job_1234567890"
                        },
                        "schema": {
                          "$ref": "#/components/schemas/CreateJobResponse"
                        }
                      }
                    }
                  },
                  {
                    "status_code": "400",
                    "description": "参数错误",
                    "content": {
                      "application/json": {
                        "schema": {
                          "$ref": "#/components/schemas/ErrorResponse"
                        }
                      }
                    }
                  },
                  {
                    "status_code": "401",
                    "description": "未授权",
                    "content": {
                      "application/json": {
                        "schema": {
                          "$ref": "#/components/schemas/ErrorResponse"
                        }
                      }
                    }
                  },
                  {
                    "status_code": "500",
                    "description": "服务器内部错误",
                    "content": {
                      "application/json": {
                        "schema": {
                          "$ref": "#/components/schemas/ErrorResponse"
                        }
                      }
                    }
                  }
                ],
                "components": {
                  "schemas": {
                    "CreateJobResponse": {
                      "type": "object",
                      "properties": {
                        "id": {
                          "description": "任务ID",
                          "type": "string"
                        }
                      }
                    },
                    "ErrorResponse": {
                      "type": "object",
                      "properties": {
                        "code": {
                          "type": "string",
                          "description": "错误码"
                        },
                        "description": {
                          "type": "string",
                          "description": "错误描述"
                        },
                        "detail": {
                          "type": "object",
                          "description": "错误详情"
                        },
                        "link": {
                          "type": "string",
                          "description": "错误链接"
                        },
                        "solution": {
                          "description": "解决方案",
                          "type": "string"
                        }
                      }
                    },
                    "CreateJobRequest": {
                      "type": "object",
                      "required": [
                        "kn_id",
                        "name"
                      ],
                      "properties": {
                        "kn_id": {
                          "type": "string",
                          "description": "业务知识网络ID"
                        },
                        "name": {
                          "type": "string",
                          "description": "任务名称"
                        }
                      }
                    }
                  }
                },
                "callbacks": null,
                "security": null,
                "tags": [
                  "OntologyJob"
                ],
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
            "create_time": 1776908488245462300,
            "update_time": 1776908488245462300,
            "create_user": "266c6a42-6131-4d62-8f39-853e7093701c",
            "update_user": "266c6a42-6131-4d62-8f39-853e7093701c",
            "extend_info": null,
            "resource_object": "tool",
            "source_id": "32207cc2-79fe-4650-b5e2-dd96ee93fc1d",
            "source_type": "openapi",
            "script_type": "",
            "code": "",
            "dependencies": [],
            "dependencies_url": ""
          },
          {
            "tool_id": "8542b7c2-f82a-4c1e-ab8c-83e2e73ccfdf",
            "name": "query_instance_subgraph",
            "description": "基于预定义的关系路径查询知识图谱中的对象子图。支持多条路径查询，每条路径返回独立子图。对象以map形式返回，支持过滤条件和排序。query_type需设为\"relation_path\"。\r\n",
            "status": "enabled",
            "metadata_type": "openapi",
            "metadata": {
              "version": "5b876cf6-b31d-46b7-a35b-a673131028fc",
              "summary": "query_instance_subgraph",
              "description": "基于预定义的关系路径查询知识图谱中的对象子图。支持多条路径查询，每条路径返回独立子图。对象以map形式返回，支持过滤条件和排序。query_type需设为\"relation_path\"。\n",
              "server_url": "http://agent-retrieval:30779",
              "path": "/api/agent-retrieval/in/v1/kn/query_instance_subgraph",
              "method": "POST",
              "create_time": 1776908488240137500,
              "update_time": 1776908488240137500,
              "create_user": "266c6a42-6131-4d62-8f39-853e7093701c",
              "update_user": "266c6a42-6131-4d62-8f39-853e7093701c",
              "api_spec": {
                "parameters": [
                  {
                    "name": "x-account-id",
                    "in": "header",
                    "description": "账户ID，用于内部服务调用时传递账户信息",
                    "required": false,
                    "schema": {
                      "type": "string"
                    }
                  },
                  {
                    "name": "x-account-type",
                    "in": "header",
                    "description": "账户类型：user(用户), app(应用), anonymous(匿名)",
                    "required": false,
                    "schema": {
                      "enum": [
                        "user",
                        "app",
                        "anonymous"
                      ],
                      "type": "string"
                    }
                  },
                  {
                    "name": "kn_id",
                    "in": "query",
                    "description": "业务知识网络ID",
                    "required": true,
                    "schema": {
                      "type": "string"
                    }
                  },
                  {
                    "name": "include_logic_params",
                    "in": "query",
                    "description": "包含逻辑属性的计算参数，默认false，返回结果不包含逻辑属性的字段和值",
                    "required": false,
                    "schema": {
                      "type": "boolean"
                    }
                  },
                  {
                    "name": "response_format",
                    "in": "query",
                    "description": "响应格式：json 或 toon，默认 json",
                    "required": false,
                    "schema": {
                      "default": "json",
                      "enum": [
                        "json",
                        "toon"
                      ],
                      "type": "string"
                    }
                  }
                ],
                "request_body": {
                  "description": "子图查询请求体",
                  "content": {
                    "application/json": {
                      "schema": {
                        "$ref": "#/components/schemas/SubGraphQueryBaseOnTypePath"
                      }
                    }
                  },
                  "required": false
                },
                "responses": [
                  {
                    "status_code": "200",
                    "description": "对象子图查询响应体",
                    "content": {
                      "application/json": {
                        "schema": {
                          "$ref": "#/components/schemas/PathEntries"
                        }
                      }
                    }
                  }
                ],
                "components": {
                  "schemas": {
                    "Sort": {
                      "required": [
                        "field",
                        "direction"
                      ],
                      "properties": {
                        "direction": {
                          "enum": [
                            "desc",
                            "asc"
                          ],
                          "type": "string",
                          "description": "排序方向"
                        },
                        "field": {
                          "type": "string",
                          "description": "排序字段"
                        }
                      },
                      "type": "object",
                      "description": "排序字段"
                    },
                    "Condition": {
                      "description": "过滤条件结构，用于构建对象实例的查询筛选条件。\n\n**重要规则：**\n- `value_from` 和 `value` 必须同时使用，不能单独使用\n- `value_from` 当前仅支持 \"const\"（常量值）\n- 当使用 `value_from: \"const\"` 时，必须同时提供 `value` 字段\n",
                      "required": [
                        "operation"
                      ],
                      "properties": {
                        "operation": {
                          "type": "string",
                          "description": "查询条件操作符。**注意：** 虽然这里列出了所有可能的操作符，但每个对象类实际支持的操作符列表以对象类定义中的 `condition_operations` 字段为准。",
                          "enum": [
                            "and",
                            "or",
                            "==",
                            "!=",
                            ">",
                            ">=",
                            "<",
                            "<=",
                            "in",
                            "not_in",
                            "like",
                            "not_like",
                            "exist",
                            "not_exist",
                            "match"
                          ]
                        },
                        "sub_conditions": {
                          "description": "子过滤条件数组，用于逻辑操作符(and/or)的组合查询",
                          "items": {
                            "$ref": "#/components/schemas/Condition"
                          },
                          "type": "array"
                        },
                        "value": {
                          "description": "字段值，格式根据操作符类型而定：\n- 比较操作符: 单个值\n- 范围查询: [min, max]数组\n- 集合操作: 值数组\n- 向量搜索: 特定格式数组\n\n**必须与 `value_from: \"const\"` 同时使用**\n",
                          "oneOf": [
                            {
                              "type": "string"
                            },
                            {
                              "type": "number"
                            },
                            {
                              "type": "boolean"
                            },
                            {
                              "type": "array",
                              "items": {
                                "oneOf": [
                                  {
                                    "type": "string"
                                  },
                                  {
                                    "type": "number"
                                  },
                                  {
                                    "type": "boolean"
                                  }
                                ]
                              }
                            }
                          ]
                        },
                        "value_from": {
                          "type": "string",
                          "description": "字段值来源。\n\n**重要：** 当前仅支持 \"const\"（常量值），且必须与 `value` 字段同时使用\n",
                          "enum": [
                            "const"
                          ]
                        },
                        "field": {
                          "type": "string",
                          "description": "字段名称，也即对象类的属性名称"
                        }
                      },
                      "type": "object"
                    },
                    "RelationPath": {
                      "required": [
                        "relations",
                        "length"
                      ],
                      "properties": {
                        "relations": {
                          "items": {
                            "$ref": "#/components/schemas/Relation"
                          },
                          "type": "array",
                          "description": "路径的边集合，沿着路径顺序出现的边"
                        },
                        "length": {
                          "description": "当前路径的长度",
                          "type": "integer"
                        }
                      },
                      "type": "object",
                      "description": "对象的关系路径"
                    },
                    "ObjectTypeOnPath": {
                      "description": "路径中的对象类信息",
                      "required": [
                        "id",
                        "condition",
                        "limit"
                      ],
                      "properties": {
                        "limit": {
                          "type": "integer",
                          "description": "对象类获取对象数量的限制"
                        },
                        "sort": {
                          "items": {
                            "$ref": "#/components/schemas/Sort"
                          },
                          "type": "array",
                          "description": "对当前对象类的排序字段"
                        },
                        "condition": {
                          "$ref": "#/components/schemas/Condition"
                        },
                        "id": {
                          "description": "对象类id",
                          "type": "string"
                        }
                      },
                      "type": "object"
                    },
                    "ObjectSubGraphResponse": {
                      "type": "object",
                      "description": "对象子图",
                      "required": [
                        "objects",
                        "relation_paths",
                        "total_count",
                        "search_after"
                      ],
                      "properties": {
                        "objects": {
                          "description": "子图中的对象map，格式为：\n{\n  \"对象ID1\": {ObjectInfoInSubgraph对象1},\n  \"对象ID2\": {ObjectInfoInSubgraph对象2}\n}\n其中key是ObjectInfoInSubgraph中的id属性，value是完整的ObjectInfoInSubgraph对象。\n动态数据字段，其值可以是基本类型、MetricProperty或OperatorProperty\n",
                          "type": "object"
                        },
                        "relation_paths": {
                          "items": {
                            "$ref": "#/components/schemas/RelationPath"
                          },
                          "type": "array",
                          "description": "对象的关系路径集合"
                        },
                        "search_after": {
                          "items": {},
                          "type": "array",
                          "description": "表示返回的最后一个起点类对象的排序值，获取这个用于下一次 search_after 分页"
                        },
                        "total_count": {
                          "type": "integer",
                          "description": "起点对象类的总条数"
                        }
                      }
                    },
                    "TypeEdge": {
                      "type": "object",
                      "description": "路径中的边信息。**方向和顺序极其重要**！通过关系类id确定边，通过路径的起点对象类id和终点对象类id来确定当前路径的方向为正向还是反向，与关系类的起终点一致为正向，相反则为反向。每个TypeEdge必须与路径中的前后对象类型严格对应，这直接影响查询结果的正确性。",
                      "required": [
                        "relation_type_id",
                        "source_object_type_id",
                        "target_object_type_id"
                      ],
                      "properties": {
                        "source_object_type_id": {
                          "type": "string",
                          "description": "路径的起点对象类id"
                        },
                        "target_object_type_id": {
                          "description": "路径的终点对象类id",
                          "type": "string"
                        },
                        "relation_type_id": {
                          "type": "string",
                          "description": "关系类id"
                        }
                      }
                    },
                    "RelationTypePath": {
                      "type": "object",
                      "description": "基于路径获取对象子图。**这是查询的核心结构**！用于定义完整的关系路径查询模板，包括路径中的所有对象类型和关系类型。object_types和relation_types数组的顺序**必须严格对应**，共同构成一个完整的关系路径。",
                      "required": [
                        "relation_types",
                        "object_types"
                      ],
                      "properties": {
                        "limit": {
                          "type": "integer",
                          "description": "当前路径返回的路径数量的限制。"
                        },
                        "object_types": {
                          "type": "array",
                          "description": "路径中的对象类集合，**顺序必须严格**与路径中节点出现顺序保持一致。对于n跳路径，object_types数组长度应为n+1，且必须按照source_object_type → 中间节点 → target_object_type的顺序排列。如果某个节点没有过滤条件或者排序或者限制数量，也必须保留其id字段以确保顺序正确。",
                          "items": {
                            "$ref": "#/components/schemas/ObjectTypeOnPath"
                          }
                        },
                        "relation_types": {
                          "description": "路径的边集合，**顺序必须严格**按照路径中关系出现的顺序排列。对于n跳路径，relation_types数组长度应为n，且必须与object_types数组中的对象类型严格对应：第i个relation_type的source_object_type_id必须等于object_types数组中第i个对象的id，target_object_type_id必须等于object_types数组中第i+1个对象的id。",
                          "items": {
                            "$ref": "#/components/schemas/TypeEdge"
                          },
                          "type": "array"
                        }
                      }
                    },
                    "PathEntries": {
                      "description": "路径子图返回体",
                      "required": [
                        "entries"
                      ],
                      "properties": {
                        "entries": {
                          "description": "路径子图",
                          "items": {
                            "$ref": "#/components/schemas/ObjectSubGraphResponse"
                          },
                          "type": "array"
                        }
                      },
                      "type": "object"
                    },
                    "Relation": {
                      "description": "一度关系（边）",
                      "required": [
                        "relation_type_id",
                        "relation_type_name",
                        "source_object_id",
                        "target_object_id"
                      ],
                      "properties": {
                        "relation_type_name": {
                          "description": "关系类名称",
                          "type": "string"
                        },
                        "source_object_id": {
                          "description": "起点对象id",
                          "type": "string"
                        },
                        "target_object_id": {
                          "description": "终点对象id",
                          "type": "string"
                        },
                        "relation_type_id": {
                          "description": "关系类id",
                          "type": "string"
                        }
                      },
                      "type": "object"
                    },
                    "SubGraphQueryBaseOnTypePath": {
                      "required": [
                        "relation_type_paths"
                      ],
                      "properties": {
                        "relation_type_paths": {
                          "items": {
                            "$ref": "#/components/schemas/RelationTypePath"
                          },
                          "type": "array",
                          "description": "关系类路径集合,数组中可以包含多条不同的关系路径，系统会同时查询并返回所有路径的结果。每条路径必须符合严格的顺序和方向要求。"
                        }
                      },
                      "type": "object",
                      "description": "查询请求的顶层结构。用于基于关系类路径查询对象子图。relation_type_paths数组中可以包含多条不同的关系路径，系统会同时查询并返回所有路径的结果。每条路径必须符合严格的顺序和方向要求。"
                    }
                  }
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
            "create_time": 1776908488245462300,
            "update_time": 1776908488245462300,
            "create_user": "266c6a42-6131-4d62-8f39-853e7093701c",
            "update_user": "266c6a42-6131-4d62-8f39-853e7093701c",
            "extend_info": null,
            "resource_object": "tool",
            "source_id": "5b876cf6-b31d-46b7-a35b-a673131028fc",
            "source_type": "openapi",
            "script_type": "",
            "code": "",
            "dependencies": [],
            "dependencies_url": ""
          },
          {
            "tool_id": "f46fa5df-f371-447f-8451-2d2f34cc78e9",
            "name": "query_object_instance",
            "description": "根据单个对象类查询对象实例，该接口基于业务知识网络语义检索接口返回的对象类定义，查询具体的对象实例数据。",
            "status": "enabled",
            "metadata_type": "openapi",
            "metadata": {
              "version": "150c8354-5e7a-41d2-aee6-92b02469f11d",
              "summary": "query_object_instance",
              "description": "根据单个对象类查询对象实例，该接口基于业务知识网络语义检索接口返回的对象类定义，查询具体的对象实例数据。",
              "server_url": "http://agent-retrieval:30779",
              "path": "/api/agent-retrieval/in/v1/kn/query_object_instance",
              "method": "POST",
              "create_time": 1776908488240137500,
              "update_time": 1776908488240137500,
              "create_user": "266c6a42-6131-4d62-8f39-853e7093701c",
              "update_user": "266c6a42-6131-4d62-8f39-853e7093701c",
              "api_spec": {
                "parameters": [
                  {
                    "name": "x-account-id",
                    "in": "header",
                    "description": "账户ID，用于内部服务调用时传递账户信息",
                    "required": false,
                    "schema": {
                      "type": "string"
                    }
                  },
                  {
                    "name": "x-account-type",
                    "in": "header",
                    "description": "账户类型：user(用户), app(应用), anonymous(匿名)",
                    "required": false,
                    "schema": {
                      "enum": [
                        "user",
                        "app",
                        "anonymous"
                      ],
                      "type": "string"
                    }
                  },
                  {
                    "name": "kn_id",
                    "in": "query",
                    "description": "业务知识网络ID",
                    "required": true,
                    "schema": {
                      "type": "string"
                    }
                  },
                  {
                    "name": "ot_id",
                    "in": "query",
                    "description": "对象类ID",
                    "required": true,
                    "schema": {
                      "type": "string"
                    }
                  },
                  {
                    "name": "include_logic_params",
                    "in": "query",
                    "description": "包含逻辑属性的计算参数，默认false，返回结果不包含逻辑属性的字段和值",
                    "required": false,
                    "schema": {
                      "type": "boolean"
                    }
                  },
                  {
                    "name": "response_format",
                    "in": "query",
                    "description": "响应格式：json 或 toon，默认 json",
                    "required": false,
                    "schema": {
                      "default": "json",
                      "enum": [
                        "json",
                        "toon"
                      ],
                      "type": "string"
                    }
                  }
                ],
                "request_body": {
                  "description": "",
                  "content": {
                    "application/json": {
                      "schema": {
                        "$ref": "#/components/schemas/FirstQueryWithSearchAfter"
                      }
                    }
                  },
                  "required": false
                },
                "responses": [
                  {
                    "status_code": "200",
                    "description": "ok",
                    "content": {
                      "application/json": {
                        "schema": {
                          "$ref": "#/components/schemas/ObjectDataResponse"
                        }
                      }
                    }
                  }
                ],
                "components": {
                  "schemas": {
                    "LogicSource": {
                      "required": [
                        "type",
                        "id"
                      ],
                      "properties": {
                        "name": {
                          "type": "string",
                          "description": "名称。查看详情时返回。"
                        },
                        "type": {
                          "enum": [
                            "metric",
                            "operator"
                          ],
                          "type": "string",
                          "description": "数据来源类型"
                        },
                        "id": {
                          "type": "string",
                          "description": "数据来源ID"
                        }
                      },
                      "type": "object",
                      "description": "数据来源"
                    },
                    "VectorConfig": {
                      "description": "向量索引的配置",
                      "required": [
                        "dimension"
                      ],
                      "properties": {
                        "dimension": {
                          "description": "向量维度",
                          "type": "integer"
                        }
                      },
                      "type": "object"
                    },
                    "FirstQueryWithSearchAfter": {
                      "type": "object",
                      "description": "分页查询的第一次查询请求",
                      "properties": {
                        "sort": {
                          "description": "排序字段，默认使用 @timestamp排序，排序方向为 desc",
                          "items": {
                            "$ref": "#/components/schemas/Sort"
                          },
                          "type": "array"
                        },
                        "condition": {
                          "$ref": "#/components/schemas/Condition"
                        },
                        "limit": {
                          "type": "integer",
                          "description": "返回的数量，默认值 10。范围 1-100",
                          "default": 10
                        },
                        "need_total": {
                          "type": "boolean",
                          "description": "是否需要总数，默认false"
                        },
                        "properties": {
                          "items": {
                            "type": "string"
                          },
                          "type": "array",
                          "description": "指定返回的对象属性字段列表，默认返回所有属性。"
                        }
                      }
                    },
                    "FulltextConfig": {
                      "description": "全文索引的配置",
                      "required": [
                        "analyzer",
                        "field_keyword"
                      ],
                      "properties": {
                        "field_keyword": {
                          "description": "是否保留原始字符串，保留原始字符串可用于精确匹配。默认是false",
                          "type": "boolean"
                        },
                        "analyzer": {
                          "description": "分词器",
                          "enum": [
                            "standard",
                            "ik_max_word"
                          ],
                          "type": "string"
                        }
                      },
                      "type": "object"
                    },
                    "ViewField": {
                      "properties": {
                        "name": {
                          "type": "string",
                          "description": "字段名称"
                        },
                        "type": {
                          "type": "string",
                          "description": "视图字段类型，查看时有此字段"
                        },
                        "display_name": {
                          "type": "string",
                          "description": "字段显示名.查看时有此字段"
                        }
                      },
                      "type": "object",
                      "description": "视图字段信息",
                      "required": [
                        "name"
                      ]
                    },
                    "Condition": {
                      "properties": {
                        "operation": {
                          "enum": [
                            "and",
                            "or",
                            "==",
                            "!=",
                            ">",
                            ">=",
                            "<",
                            "<=",
                            "in",
                            "not_in",
                            "like",
                            "not_like",
                            "exist",
                            "not_exist",
                            "match"
                          ],
                          "type": "string",
                          "description": "查询条件操作符。\n**注意：** 虽然这里列出了所有可能的操作符，但每个对象类实际支持的操作符列表以对象类定义中的 `condition_operations` 字段为准。\n"
                        },
                        "sub_conditions": {
                          "items": {
                            "$ref": "#/components/schemas/Condition"
                          },
                          "type": "array",
                          "description": "子过滤条件数组，用于逻辑操作符(and/or)的组合查询"
                        },
                        "value": {
                          "description": "字段值，格式根据操作符类型而定：\n- 比较操作符: 单个值\n- 范围查询: [min, max]数组\n- 集合操作: 值数组\n- 向量搜索: 特定格式数组\n\n**必须与 `value_from: \"const\"` 同时使用**\n"
                        },
                        "value_from": {
                          "enum": [
                            "const"
                          ],
                          "type": "string",
                          "description": "字段值来源。\n\n**重要：** 当前仅支持 \"const\"（常量值），且必须与 `value` 字段同时使用\n"
                        },
                        "field": {
                          "description": "字段名称，也即对象类的属性名称",
                          "type": "string"
                        }
                      },
                      "type": "object",
                      "description": "过滤条件结构，用于构建对象实例的查询筛选条件。\n\n**重要规则：**\n- `value_from` 和 `value` 必须同时使用，不能单独使用\n- `value_from` 当前仅支持 \"const\"（常量值）\n- 当使用 `value_from: \"const\"` 时，必须同时提供 `value` 字段\n",
                      "required": [
                        "operation"
                      ]
                    },
                    "Sort": {
                      "description": "排序字段",
                      "required": [
                        "field",
                        "direction"
                      ],
                      "properties": {
                        "direction": {
                          "enum": [
                            "desc",
                            "asc"
                          ],
                          "type": "string",
                          "description": "排序方向"
                        },
                        "field": {
                          "type": "string",
                          "description": "排序字段"
                        }
                      },
                      "type": "object"
                    },
                    "DataSource": {
                      "required": [
                        "type",
                        "id"
                      ],
                      "properties": {
                        "type": {
                          "enum": [
                            "data_view"
                          ],
                          "type": "string",
                          "description": "数据来源类型为数据视图"
                        },
                        "id": {
                          "description": "数据视图ID",
                          "type": "string"
                        },
                        "name": {
                          "type": "string",
                          "description": "名称。查看详情时返回。"
                        }
                      },
                      "type": "object",
                      "description": "数据来源"
                    },
                    "ConceptGroup": {
                      "type": "object",
                      "description": "概念分组",
                      "required": [
                        "id",
                        "name"
                      ],
                      "properties": {
                        "name": {
                          "description": "概念分组名称",
                          "type": "string"
                        },
                        "id": {
                          "description": "概念分组ID",
                          "type": "string"
                        }
                      }
                    },
                    "ObjectTypeDetail": {
                      "properties": {
                        "branch": {
                          "type": "string",
                          "description": "分支ID"
                        },
                        "creator": {
                          "type": "string",
                          "description": "创建人ID"
                        },
                        "logic_properties": {
                          "description": "逻辑属性",
                          "items": {
                            "$ref": "#/components/schemas/LogicProperty"
                          },
                          "type": "array"
                        },
                        "updater": {
                          "type": "string",
                          "description": "最近一次修改人"
                        },
                        "name": {
                          "type": "string",
                          "description": "对象类名称"
                        },
                        "update_time": {
                          "type": "integer",
                          "format": "int64",
                          "description": "最近一次更新时间"
                        },
                        "comment": {
                          "description": "备注（可以为空）",
                          "type": "string"
                        },
                        "data_source": {
                          "$ref": "#/components/schemas/DataSource"
                        },
                        "create_time": {
                          "description": "创建时间",
                          "type": "integer",
                          "format": "int64"
                        },
                        "icon": {
                          "type": "string",
                          "description": "图标"
                        },
                        "detail": {
                          "type": "string",
                          "description": "说明书。按需返回，若指定了include_detail=true，则返回，否则不返回。列表查询时不返回此字段"
                        },
                        "tags": {
                          "items": {
                            "type": "string"
                          },
                          "type": "array",
                          "description": "标签。 （可以为空）"
                        },
                        "id": {
                          "description": "对象类ID",
                          "type": "string"
                        },
                        "display_key": {
                          "type": "string",
                          "description": "对象实例的显示属性"
                        },
                        "primary_keys": {
                          "description": "主键",
                          "items": {
                            "type": "string"
                          },
                          "type": "array"
                        },
                        "data_properties": {
                          "description": "数据属性",
                          "items": {
                            "$ref": "#/components/schemas/DataProperty"
                          },
                          "type": "array"
                        },
                        "kn_id": {
                          "type": "string",
                          "description": "业务知识网络id"
                        },
                        "color": {
                          "description": "颜色",
                          "type": "string"
                        },
                        "module_type": {
                          "type": "string",
                          "description": "模块类型",
                          "enum": [
                            "object_type"
                          ]
                        },
                        "concept_groups": {
                          "type": "array",
                          "description": "概念分组id",
                          "items": {
                            "$ref": "#/components/schemas/ConceptGroup"
                          }
                        }
                      },
                      "type": "object",
                      "description": "对象类信息"
                    },
                    "Parameter4Metric": {
                      "description": "逻辑参数",
                      "required": [
                        "name",
                        "value_from",
                        "operation"
                      ],
                      "properties": {
                        "operation": {
                          "type": "string",
                          "description": "操作符。映射指标模型的属性时，此字段必须",
                          "enum": [
                            "in",
                            "=",
                            "!=",
                            ">",
                            ">=",
                            "<",
                            "<="
                          ]
                        },
                        "value": {
                          "type": "string",
                          "description": "参数值。value_from=property时，填入的是对象类的数据属性名称；value_from=input时，不设置此字段"
                        },
                        "value_from": {
                          "description": "值来源",
                          "enum": [
                            "property",
                            "input"
                          ],
                          "type": "string"
                        },
                        "name": {
                          "type": "string",
                          "description": "参数名称"
                        }
                      },
                      "type": "object"
                    },
                    "Parameter4Operator": {
                      "required": [
                        "name",
                        "value_from"
                      ],
                      "properties": {
                        "source": {
                          "description": "参数来源",
                          "type": "string"
                        },
                        "type": {
                          "description": "参数类型",
                          "type": "string"
                        },
                        "value": {
                          "type": "string",
                          "description": "参数值。value_from=property时，填入的是对象类的数据属性名称；value_from=input时，不设置此字段"
                        },
                        "value_from": {
                          "type": "string",
                          "description": "值来源",
                          "enum": [
                            "property",
                            "input"
                          ]
                        },
                        "name": {
                          "description": "参数名称",
                          "type": "string"
                        }
                      },
                      "type": "object",
                      "description": "逻辑参数"
                    },
                    "ObjectDataResponse": {
                      "properties": {
                        "datas": {
                          "type": "array",
                          "description": "对象实例数据。动态数据字段，其值可以是基本类型、MetricProperty或OperatorProperty",
                          "items": {
                            "type": "object"
                          }
                        },
                        "object_type": {
                          "$ref": "#/components/schemas/ObjectTypeDetail"
                        },
                        "search_after": {
                          "type": "array",
                          "description": "表示返回的最后一个文档的排序值，获取这个用于下一次 search_after 分页",
                          "items": {}
                        },
                        "total_count": {
                          "type": "integer",
                          "description": "总条数"
                        }
                      },
                      "type": "object",
                      "description": "节点（对象类）信息",
                      "required": [
                        "groups",
                        "type",
                        "datas",
                        "search_after"
                      ]
                    },
                    "DataProperty": {
                      "required": [
                        "name",
                        "display_name",
                        "type",
                        "comment",
                        "mapped_field",
                        "index",
                        "fulltext_config",
                        "vector_config"
                      ],
                      "properties": {
                        "type": {
                          "description": "属性数据类型。除了视图的字段类型之外，还有 metric、objective、event、trace、log、operator",
                          "type": "string"
                        },
                        "vector_config": {
                          "$ref": "#/components/schemas/VectorConfig"
                        },
                        "comment": {
                          "description": "属性描述",
                          "type": "string"
                        },
                        "display_name": {
                          "type": "string",
                          "description": "属性显示名"
                        },
                        "fulltext_config": {
                          "$ref": "#/components/schemas/FulltextConfig"
                        },
                        "index": {
                          "description": "是否开启索引，默认是true",
                          "type": "boolean"
                        },
                        "mapped_field": {
                          "$ref": "#/components/schemas/ViewField"
                        },
                        "name": {
                          "description": "属性名称。只能包含小写英文字母、数字、下划线（_）、连字符（-），且不能以下划线和连字符开头",
                          "type": "string"
                        }
                      },
                      "type": "object",
                      "description": "数据属性"
                    },
                    "Parameter": {
                      "oneOf": [
                        {
                          "$ref": "#/components/schemas/Parameter4Operator"
                        },
                        {
                          "$ref": "#/components/schemas/Parameter4Metric"
                        }
                      ],
                      "type": "object",
                      "description": "逻辑/指标参数"
                    },
                    "LogicProperty": {
                      "required": [
                        "name",
                        "data_source",
                        "parameters"
                      ],
                      "properties": {
                        "parameters": {
                          "description": "逻辑所需的参数",
                          "items": {
                            "$ref": "#/components/schemas/Parameter"
                          },
                          "type": "array"
                        },
                        "type": {
                          "description": "属性数据类型。除了视图的字段类型之外，还有 metric、objective、event、trace、log、operator",
                          "type": "string"
                        },
                        "comment": {
                          "description": "属性描述",
                          "type": "string"
                        },
                        "data_source": {
                          "$ref": "#/components/schemas/LogicSource"
                        },
                        "display_name": {
                          "type": "string",
                          "description": "属性显示名"
                        },
                        "index": {
                          "description": "是否开启索引，默认是true",
                          "type": "boolean"
                        },
                        "name": {
                          "description": "属性名称。只能包含小写英文字母、数字、下划线（_）、连字符（-），且不能以下划线和连字符开头",
                          "type": "string"
                        }
                      },
                      "type": "object",
                      "description": "逻辑属性"
                    }
                  }
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
            "create_time": 1776908488245462300,
            "update_time": 1776908488245462300,
            "create_user": "266c6a42-6131-4d62-8f39-853e7093701c",
            "update_user": "266c6a42-6131-4d62-8f39-853e7093701c",
            "extend_info": null,
            "resource_object": "tool",
            "source_id": "150c8354-5e7a-41d2-aee6-92b02469f11d",
            "source_type": "openapi",
            "script_type": "",
            "code": "",
            "dependencies": [],
            "dependencies_url": ""
          },
          {
            "tool_id": "00929c5f-3375-4ddc-9fb0-c48a24707f39",
            "name": "get_action_info",
            "description": "根据对象实例标识召回关联行动，返回 _dynamic_tools。",
            "status": "enabled",
            "metadata_type": "openapi",
            "metadata": {
              "version": "d25c679f-d451-4253-b8ca-f862a529368c",
              "summary": "get_action_info",
              "description": "根据对象实例标识召回关联行动，返回 _dynamic_tools。",
              "server_url": "http://agent-retrieval:30779",
              "path": "/api/agent-retrieval/in/v1/kn/get_action_info",
              "method": "POST",
              "create_time": 1776908488240137500,
              "update_time": 1776908488240137500,
              "create_user": "266c6a42-6131-4d62-8f39-853e7093701c",
              "update_user": "266c6a42-6131-4d62-8f39-853e7093701c",
              "api_spec": {
                "parameters": [
                  {
                    "name": "x-account-id",
                    "in": "header",
                    "description": "账户ID，用于内部服务调用时传递账户信息",
                    "required": false,
                    "schema": {
                      "type": "string"
                    }
                  },
                  {
                    "name": "x-account-type",
                    "in": "header",
                    "description": "账户类型：user(用户), app(应用), anonymous(匿名)",
                    "required": false,
                    "schema": {
                      "enum": [
                        "user",
                        "app",
                        "anonymous"
                      ],
                      "type": "string"
                    }
                  }
                ],
                "request_body": {
                  "description": "",
                  "content": {
                    "application/json": {
                      "examples": {
                        "multi_instance_example": {
                          "summary": "多对象实例示例",
                          "value": {
                            "_instance_identities": [
                              {
                                "disease_id": "disease_000001"
                              },
                              {
                                "disease_id": "disease_000002"
                              }
                            ],
                            "at_id": "generate_treatment_plan",
                            "kn_id": "kn_medical"
                          }
                        },
                        "single_instance_example": {
                          "summary": "单对象实例示例",
                          "value": {
                            "_instance_identities": [
                              {
                                "disease_id": "disease_000001"
                              }
                            ],
                            "at_id": "generate_treatment_plan",
                            "kn_id": "kn_medical"
                          }
                        }
                      },
                      "schema": {
                        "$ref": "#/components/schemas/ActionRecallRequest"
                      }
                    }
                  },
                  "required": false
                },
                "responses": [
                  {
                    "status_code": "200",
                    "description": "成功返回动态工具列表",
                    "content": {
                      "application/json": {
                        "schema": {
                          "$ref": "#/components/schemas/ActionRecallResponse"
                        }
                      }
                    }
                  },
                  {
                    "status_code": "400",
                    "description": "请求参数错误",
                    "content": {
                      "application/json": {
                        "examples": {
                          "invalid_request": {
                            "value": {
                              "code": "INVALID_REQUEST",
                              "description": "_instance_identities 格式错误"
                            }
                          }
                        },
                        "schema": {
                          "$ref": "#/components/schemas/ErrorResponse"
                        }
                      }
                    }
                  },
                  {
                    "status_code": "500",
                    "description": "服务器内部错误",
                    "content": {
                      "application/json": {
                        "schema": {
                          "$ref": "#/components/schemas/ErrorResponse"
                        }
                      }
                    }
                  },
                  {
                    "status_code": "502",
                    "description": "上游服务不可用",
                    "content": {
                      "application/json": {
                        "schema": {
                          "$ref": "#/components/schemas/ErrorResponse"
                        }
                      }
                    }
                  }
                ],
                "components": {
                  "schemas": {
                    "ActionRecallRequest": {
                      "required": [
                        "kn_id",
                        "at_id"
                      ],
                      "properties": {
                        "at_id": {
                          "description": "行动类ID（从 Schema 获取）",
                          "type": "string"
                        },
                        "kn_id": {
                          "type": "string",
                          "description": "知识网络ID"
                        },
                        "_instance_identities": {
                          "description": "对象实例标识列表（可选）。每个元素为主键键值对，必须从 query_object_instance 或 query_instance_subgraph 返回的 _instance_identity 字段提取，不可臆造。",
                          "items": {
                            "type": "object"
                          },
                          "type": "array"
                        }
                      },
                      "type": "object"
                    },
                    "ActionRecallResponse": {
                      "required": [
                        "_dynamic_tools"
                      ],
                      "properties": {
                        "headers": {
                          "type": "object"
                        },
                        "_dynamic_tools": {
                          "items": {
                            "properties": {
                              "api_url": {
                                "type": "string"
                              },
                              "description": {
                                "type": "string"
                              },
                              "fixed_params": {
                                "type": "object"
                              },
                              "name": {
                                "type": "string"
                              },
                              "parameters": {
                                "type": "object"
                              }
                            },
                            "type": "object"
                          },
                          "type": "array",
                          "description": "Function Call 格式的工具列表"
                        }
                      },
                      "type": "object"
                    },
                    "ErrorResponse": {
                      "type": "object",
                      "properties": {
                        "code": {
                          "type": "string"
                        },
                        "description": {
                          "type": "string"
                        }
                      }
                    }
                  }
                },
                "callbacks": null,
                "security": null,
                "tags": [
                  "action-recall"
                ],
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
            "create_time": 1776908488245462300,
            "update_time": 1776908488245462300,
            "create_user": "266c6a42-6131-4d62-8f39-853e7093701c",
            "update_user": "266c6a42-6131-4d62-8f39-853e7093701c",
            "extend_info": null,
            "resource_object": "tool",
            "source_id": "d25c679f-d451-4253-b8ca-f862a529368c",
            "source_type": "openapi",
            "script_type": "",
            "code": "",
            "dependencies": [],
            "dependencies_url": ""
          },
          {
            "tool_id": "74498bbd-2bdf-4db3-a4da-a90133c7fb65",
            "name": "find_skills",
            "description": "基于业务上下文召回 Skill 候选列表。\r\n\r\n**召回模式选择规则（自动判断，无需显式指定）：**\r\n\r\n| 请求参数 | 召回模式 |\r\n|---------|---------|\r\n| 仅 kn_id | 网络级（Mode 1） |\r\n| kn_id + object_type_id | 对象类级（Mode 2） |\r\n| kn_id + object_type_id + instance_identities | 实例级（Mode 3） |\r\n\r\n**网络级特殊规则：** 仅当 skills ObjectType 与任意其他 ObjectType 都不存在 RelationType 时才返回结果，\r\n否则返回空列表（需要更精确的 object_type_id 才能召回）。\r\n\r\n**skill_query 说明：** 传入时会对 skills 实例的 name/description 字段追加文本过滤条件（支持 knn/match/like，\r\n优先使用已构建的向量/索引能力），若 BKN 中 skills ObjectType 不存在或元数据获取失败则返回 502。\r\n",
            "status": "enabled",
            "metadata_type": "openapi",
            "metadata": {
              "version": "d856956c-2c45-4c1d-8428-9db906ac0655",
              "summary": "find_skills",
              "description": "基于业务上下文召回 Skill 候选列表。\n\n**召回模式选择规则（自动判断，无需显式指定）：**\n\n| 请求参数 | 召回模式 |\n|---------|---------|\n| 仅 kn_id | 网络级（Mode 1） |\n| kn_id + object_type_id | 对象类级（Mode 2） |\n| kn_id + object_type_id + instance_identities | 实例级（Mode 3） |\n\n**网络级特殊规则：** 仅当 skills ObjectType 与任意其他 ObjectType 都不存在 RelationType 时才返回结果，\n否则返回空列表（需要更精确的 object_type_id 才能召回）。\n\n**skill_query 说明：** 传入时会对 skills 实例的 name/description 字段追加文本过滤条件（支持 knn/match/like，\n优先使用已构建的向量/索引能力），若 BKN 中 skills ObjectType 不存在或元数据获取失败则返回 502。\n",
              "server_url": "http://agent-retrieval:30779",
              "path": "/api/agent-retrieval/in/v1/kn/find_skills",
              "method": "POST",
              "create_time": 1776908488240137500,
              "update_time": 1776908488240137500,
              "create_user": "266c6a42-6131-4d62-8f39-853e7093701c",
              "update_user": "266c6a42-6131-4d62-8f39-853e7093701c",
              "api_spec": {
                "parameters": [
                  {
                    "name": "x-account-id",
                    "in": "header",
                    "description": "账户 ID，用于内部服务调用时传递账户信息",
                    "required": false,
                    "schema": {
                      "type": "string"
                    }
                  },
                  {
                    "name": "x-account-type",
                    "in": "header",
                    "description": "账户类型：user（用户）、app（应用）、anonymous（匿名）",
                    "required": false,
                    "schema": {
                      "enum": [
                        "user",
                        "app",
                        "anonymous"
                      ],
                      "type": "string"
                    }
                  },
                  {
                    "name": "response_format",
                    "in": "query",
                    "description": "响应格式：json 或 toon，默认 json",
                    "required": false,
                    "schema": {
                      "default": "json",
                      "enum": [
                        "json",
                        "toon"
                      ],
                      "type": "string"
                    }
                  }
                ],
                "request_body": {
                  "description": "",
                  "content": {
                    "application/json": {
                      "examples": {
                        "instance_level": {
                          "summary": "实例级召回（kn_id + object_type_id + instance_identities）",
                          "value": {
                            "instance_identities": [
                              {
                                "contract_id": "C-2024-001"
                              }
                            ],
                            "kn_id": "kn_legal",
                            "object_type_id": "contract",
                            "top_k": 10
                          }
                        },
                        "network_level": {
                          "summary": "网络级召回（仅 kn_id）",
                          "value": {
                            "kn_id": "kn_legal",
                            "top_k": 10
                          }
                        },
                        "object_type_level": {
                          "summary": "对象类级召回（kn_id + object_type_id）",
                          "value": {
                            "kn_id": "kn_legal",
                            "object_type_id": "contract",
                            "top_k": 10
                          }
                        },
                        "with_skill_query": {
                          "summary": "带语义过滤的实例级召回",
                          "value": {
                            "instance_identities": [
                              {
                                "contract_id": "C-2024-001"
                              }
                            ],
                            "kn_id": "kn_legal",
                            "object_type_id": "contract",
                            "skill_query": "合同审查",
                            "top_k": 5
                          }
                        }
                      },
                      "schema": {
                        "$ref": "#/components/schemas/FindSkillsRequest"
                      }
                    }
                  },
                  "required": false
                },
                "responses": [
                  {
                    "status_code": "200",
                    "description": "成功返回 Skill 候选列表（可能为空列表）",
                    "content": {
                      "application/json": {
                        "examples": {
                          "empty_result_network_scope": {
                            "summary": "无匹配 Skill（网络级但 skills 有关联）",
                            "value": {
                              "entries": [],
                              "message": "当前知识网络中的 Skill 不是全局生效，网络级范围未返回结果。请补充 object_type_id；如已定位到实例，再补充 instance_identities 后重试。"
                            }
                          },
                          "empty_result_no_binding": {
                            "summary": "无匹配 Skill（对象类未配置 Skill 绑定）",
                            "value": {
                              "entries": [],
                              "message": "当前对象类未配置 Skill 绑定关系，无法在该范围内召回 Skill。请确认该对象类是否已绑定 Skill。"
                            }
                          },
                          "success_with_results": {
                            "summary": "返回 Skill 候选",
                            "value": {
                              "entries": [
                                {
                                  "description": "对合同条款进行全面审查，识别风险点",
                                  "name": "合同审查",
                                  "skill_id": "skill_contract_review"
                                },
                                {
                                  "description": "提取合同中的关键条款和义务",
                                  "name": "关键条款提取",
                                  "skill_id": "skill_clause_extract"
                                }
                              ]
                            }
                          }
                        },
                        "schema": {
                          "$ref": "#/components/schemas/FindSkillsResponse"
                        }
                      }
                    }
                  },
                  {
                    "status_code": "400",
                    "description": "请求参数错误",
                    "content": {
                      "application/json": {
                        "examples": {
                          "instance_without_object_type": {
                            "value": {
                              "code": "INVALID_REQUEST",
                              "description": "instance_identities 不为空时 object_type_id 也必须提供"
                            }
                          },
                          "missing_kn_id": {
                            "value": {
                              "code": "INVALID_REQUEST",
                              "description": "kn_id 为必填字段"
                            }
                          }
                        },
                        "schema": {
                          "$ref": "#/components/schemas/ErrorResponse"
                        }
                      }
                    }
                  },
                  {
                    "status_code": "502",
                    "description": "上游服务不可用（BKN 或 ontology-query 异常）",
                    "content": {
                      "application/json": {
                        "examples": {
                          "skills_ot_not_found": {
                            "value": {
                              "code": "BAD_GATEWAY",
                              "description": "skill_query requires skills ObjectType (id=skills) but none found in kn_id=kn_legal"
                            }
                          }
                        },
                        "schema": {
                          "$ref": "#/components/schemas/ErrorResponse"
                        }
                      }
                    }
                  }
                ],
                "components": {
                  "schemas": {
                    "FindSkillsRequest": {
                      "properties": {
                        "instance_identities": {
                          "items": {
                            "type": "object"
                          },
                          "type": "array",
                          "description": "对象实例标识列表。每个元素为主键键值对，必须从 query_object_instance 或\nquery_instance_subgraph 返回的 _instance_identity 字段提取，不可臆造。\n传入时 object_type_id 也必须提供，否则返回 400。\n"
                        },
                        "kn_id": {
                          "type": "string",
                          "description": "知识网络 ID"
                        },
                        "object_type_id": {
                          "type": "string",
                          "description": "业务对象类型 ID（从 kn_search 或 kn_schema_search 返回的 concept_id 获取）。\n不传时为网络级召回；传入时为对象类级或实例级召回。\n"
                        },
                        "skill_query": {
                          "type": "string",
                          "description": "可选的 Skill 语义过滤词，对 skills 实例的 name/description 字段追加文本过滤。\nBKN 已构建向量时使用 knn，已构建全文索引时使用 match，否则使用 like。\n若 skills ObjectType 元数据获取失败则返回 502。\n"
                        },
                        "top_k": {
                          "type": "integer",
                          "description": "最多返回的 Skill 数量，默认 10，最大 20",
                          "default": 10
                        }
                      },
                      "type": "object",
                      "required": [
                        "kn_id",
                        "object_type_id"
                      ]
                    },
                    "FindSkillsResponse": {
                      "type": "object",
                      "required": [
                        "entries"
                      ],
                      "properties": {
                        "entries": {
                          "description": "Skill 候选列表，按命中层级优先级（实例级 > 对象类级 > 网络级）和相关性分数降序排列。\n同优先级同分数时按 skill_id 字典序排列，保证顺序稳定。\n无匹配时返回空数组。\n",
                          "items": {
                            "$ref": "#/components/schemas/SkillItem"
                          },
                          "type": "array"
                        },
                        "message": {
                          "type": "string",
                          "description": "空结果说明信息。仅当 entries 为空且接口返回 200 时出现。\n用于向调用方（Agent）解释当前为什么没有结果，以及下一步建议。\nentries 非空时不返回此字段。支持多语言（由请求 X-Language 头决定）。\n"
                        }
                      }
                    },
                    "SkillItem": {
                      "required": [
                        "skill_id",
                        "name"
                      ],
                      "properties": {
                        "name": {
                          "type": "string",
                          "description": "Skill 名称"
                        },
                        "skill_id": {
                          "description": "Skill 唯一标识，由 execution-factory 写入 BKN skills ObjectType 时确定",
                          "type": "string"
                        },
                        "description": {
                          "type": "string",
                          "description": "Skill 功能描述（可选，BKN 中无此属性时不返回）"
                        }
                      },
                      "type": "object"
                    },
                    "ErrorResponse": {
                      "type": "object",
                      "properties": {
                        "code": {
                          "description": "错误码",
                          "type": "string"
                        },
                        "description": {
                          "type": "string",
                          "description": "错误详情"
                        }
                      }
                    }
                  }
                },
                "callbacks": null,
                "security": null,
                "tags": [
                  "skill-recall"
                ],
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
            "create_time": 1776908488245462300,
            "update_time": 1776908488245462300,
            "create_user": "266c6a42-6131-4d62-8f39-853e7093701c",
            "update_user": "266c6a42-6131-4d62-8f39-853e7093701c",
            "extend_info": null,
            "resource_object": "tool",
            "source_id": "d856956c-2c45-4c1d-8428-9db906ac0655",
            "source_type": "openapi",
            "script_type": "",
            "code": "",
            "dependencies": [],
            "dependencies_url": ""
          }
        ],
        "create_time": 1776908488238766300,
        "update_time": 1776908495056935700,
        "create_user": "266c6a42-6131-4d62-8f39-853e7093701c",
        "update_user": "266c6a42-6131-4d62-8f39-853e7093701c",
        "metadata_type": "openapi"
      }
    ]
  }
}