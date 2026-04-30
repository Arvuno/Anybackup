# `foundation-cli mysql authorize`

## 命令摘要

| 项目 | 值 |
|---|---|
| CLI | `foundation-cli mysql authorize --tenant-id <tenant-id> --data '<json>'` |
| 方法 | `POST` |
| 路径 | `/resource_center/v1/databasealone/mysql/authorize` |
| 风险 | `写入` |

## 请求参数

### 公共参数

| 字段 | CLI 参数 | 类型 | 必填 | 说明 |
|---|---|---|---|---|
| tenantId | `--tenant-id` | string | 否 | 租户上下文；可由环境变量 FOUNDATION_TENANT_ID 提供 |
| endpoint | `--endpoint` | string | 是 | Foundation 接口地址 |
| ak | `--ak` | string | 是 | 访问密钥 |
| sk | `--sk` | string | 是 | 访问密钥密文 |
| targetVersion | `--target-version` | string | 否 | 默认值为 `8.0.9.0` |

### 命令参数

| 字段 | CLI 参数 | 类型 | 必填 | 说明 |
|---|---|---|---|---|
| data | `--data`, `-d` | string | 是（二选一） | 行内 JSON 请求体 |
| bodyFile | `--body-file` | string | 是（二选一） | JSON 请求体文件路径 |

## Body 参数（`--data`）

顶层请求体：`MySQLAuthorizeRequest`

| 字段 | 类型 | 必填 | 说明 |
|---|---|---|---|
| instanceName | string | 是 | MySQL 实例名称 |
| clientId | string | 是 | 客户端唯一标识 |
| username | string | 是 | 数据库用户名 |
| password | string | 是 | 数据库密码 |
| systemId | string | 是 | 生产系统 ID |
| resourceId | string | 是 | 保护对象 ID |
| osUserName | string | 是 | 实例所属 OS 用户名 |
| type | int | 是 | 应用类型（MySQL 场景固定为 `202`） |
| authType | int | 否 | 授权类型 |
| portType | int | 否 | 端口类型 |
| dbPort | int | 否 | 数据库端口 |
| dbSSL | bool | 否 | 是否启用 SSL |
| dbPath | string | 否 | 数据库路径 |
| dbSocket | string | 否 | Socket 路径 |
| loginPath | string | 否 | 登录路径 |
| hostName | string | 否 | 主机名 |
| sslCa | string | 否 | SSL CA 内容 |
| sslCaPath | string | 否 | SSL CA 路径 |
| sslCert | string | 否 | SSL 证书内容 |
| sslCipher | string | 否 | SSL 加密套件 |
| sslKey | string | 否 | SSL 私钥内容 |
| sslMode | string | 否 | SSL 模式 |
| dbConfig | string | 否 | 数据库配置文件路径 |
| keyRing | string | 否 | 密钥环配置路径 |

## 枚举列表

- `type`：`202=MySQL`。
- `dbSSL`：`false=关闭 SSL`，`true=开启 SSL`。
- `authType`、`portType`：当前文档无稳定枚举定义，请以接口实时返回和后端契约为准。

## 请求体关键校验

CLI 当前会重点校验以下字段：

- `instanceName`
- `clientId`
- `username`
- `password`
- `systemId`
- `resourceId`
- `osUserName`
- `type`

## 请求体示例

```json
{
  "instanceName": "mysql3306_U0HYTDM3RENXS4EJ",
  "authType": 0,
  "portType": 0,
  "dbPort": 3306,
  "dbSSL": false,
  "username": "root",
  "password": "dsadadad",
  "dbPath": "",
  "dbSocket": "/var/lib/mysql/mysql.sock",
  "loginPath": "",
  "hostName": "",
  "sslCa": "",
  "sslCaPath": "",
  "sslCert": "",
  "sslCipher": "",
  "sslKey": "",
  "sslMode": "",
  "type": 202,
  "systemId": "13caa58fd84bab1beb0b7c9d76b7efae",
  "clientId": "13caa58fd84bab1beb0b7c9d76b7efae",
  "dbConfig": "/",
  "keyRing": "/",
  "resourceId": "428e4c05bf87deef20a86e670f3eb9a9",
  "osUserName": "mysql"
}
```

## 返回案例

说明：
1. 当前仓库未附带该命令的真实接口成功响应。
2. 下例仅用于说明常见返回结构层级，字段名和取值以后端当前版本实际返回为准。
3. 建议后续补充 1 份真实成功响应和 1 份典型失败响应。

```json
{
  "error": null,
  "status": "success",
  "responseData": {
    "instanceName": "mysql3306_U0HYTDM3RENXS4EJ",
    "resourceId": "428e4c05bf87deef20a86e670f3eb9a9",
    "authorized": true
  }
}
```
## 示例

```bash
foundation-cli mysql authorize \
  --tenant-id <tenant-id> \
  --endpoint <endpoint> \
  --ak <ak> \
  --sk <sk> \
  --data '{"instanceName":"mysql3306_U0HYTDM3RENXS4EJ","authType":0,"portType":0,"dbPort":3306,"dbSSL":false,"username":"root","password":"dsadadad","dbPath":"","dbSocket":"/var/lib/mysql/mysql.sock","loginPath":"","hostName":"","sslCa":"","sslCaPath":"","sslCert":"","sslCipher":"","sslKey":"","sslMode":"","type":202,"systemId":"13caa58fd84bab1beb0b7c9d76b7efae","clientId":"13caa58fd84bab1beb0b7c9d76b7efae","dbConfig":"/","keyRing":"/","resourceId":"428e4c05bf87deef20a86e670f3eb9a9","osUserName":"mysql"}'
```


