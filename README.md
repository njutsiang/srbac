# SRBAC

Service-And-Role-Based Access Control 基于服务和角色的访问控制

## RBAC 核心：

- 权限
- 角色
- 用户
- 访问控制

## SRBAC 核心：

- 服务
- 权限
- 角色
- 用户
- 访问控制

## SRBAC 亮点：

- 适用于多服务集群
- 适用于微服务网关
- 实现对于服务和接口的访问控制
- 实现跨服务的角色权限控制
- 权限控制和业务代码完全解耦
- 安全、可靠、高效

## APISIX 使用示例

使用 business-upstream

```
curl -X POST 'http://127.0.0.1:9180/apisix/admin/routes' \
-H 'X-API-KEY: edd1c9f034335f136f87ad84b625c8f1' \
-d '{
    "name":"api",
    "desc":"企业动态路由",
    "uri":"/api/*",
    "methods":["GET","POST","PUT","DELETE","HEAD","OPTIONS"],
    "plugins":{
        "business-upstream":{
            "header": "X-Company-Id",
            "query": "company_id",
            "redis": {
                "host": "127.0.0.1",
                "port": 6379,
                "password": "",
                "database": 0
            }
        },
        "proxy-rewrite":{
            "regex_uri":["^/api/(.*)","/$1"],
            "host":"api",
            "headers":{
                "X-Gateway":"apisix",
                "X-Target-Server-Name":"api"
            }
        }
    },
    "upstream": {
        "type": "roundrobin",
        "nodes": {
            "127.0.0.1:12008": 1
        }
    }
}'
```

使用 static-jwt-auth

```
curl -X POST 'http://127.0.0.1:9180/apisix/admin/routes' \
-H 'X-API-KEY: edd1c9f034335f136f87ad84b625c8f1' \
-d '{
    "name":"user/v2",
    "uri":"/user/v2/*",
    "methods":["GET","POST","PUT","DELETE","HEAD","OPTIONS"],
    "plugins":{
        "static-jwt-auth":{
            "header":"authorization",
            "query":"token",
            "secret":"123456",
            "base64_encode":false,
            "response_err":true
        },
        "proxy-rewrite":{
            "regex_uri":["^/user/(.*)","/$1"],
            "host":"user-v2",
            "headers":{
                "X-Gateway":"apisix",
                "X-Target-Server-Name":"user-v2"
            }
        },
        "response-rewrite":{
            "headers":{
                "X-Gateway":"apisix",
                "X-Target-Server-Name":"user-v2"
            }
        }
    },
    "upstream": {
        "type": "roundrobin",
        "nodes": {
            "127.0.0.1:12008": 1
        }
    }
}'
```

使用 token-auth + rbac-access

```
curl -X POST 'http://127.0.0.1:9180/apisix/admin/routes' \
-H 'X-API-KEY: edd1c9f034335f136f87ad84b625c8f1' \
-d '{
    "name":"user/v1",
    "desc":"用户中心",
    "uri":"/user/v1/*",
    "methods":["GET","POST","PUT","DELETE","HEAD","OPTIONS"],
    "plugins":{
        "token-auth":{
            "header":"authorization",
            "query":"token",
            "redis": {
                "host": "127.0.0.1",
                "port": 6379,
                "password": "",
                "database": 0
            }
        },
        "rbac-access":{
            "redis": {
                "host": "127.0.0.1",
                "port": 6379,
                "password": "",
                "database": 0
            }
        },
        "proxy-rewrite":{
            "regex_uri":["^/user/(.*)","/$1"],
            "host":"user-v1",
            "headers":{
                "X-Gateway":"apisix",
                "X-Target-Server-Name":"user-v1"
            }
        },
        "response-rewrite":{
            "headers":{
                "X-Gateway":"apisix",
                "X-Target-Server-Name":"user-v1"
            }
        }
    },
    "upstream": {
        "type": "roundrobin",
        "nodes": {
            "127.0.0.1:12008": 1
        }
    }
}'
```

使用 token-auth + limit-count

```
curl -X POST 'http://127.0.0.1:9180/apisix/admin/routes' \
-H 'X-API-KEY: edd1c9f034335f136f87ad84b625c8f1' \
-d '{
    "name":"user/v3",
    "uri":"/user/v3/*",
    "methods":["GET","POST","PUT","DELETE","HEAD","OPTIONS"],
    "plugins":{
        "token-auth":{
            "header":"authorization",
            "query":"token",
            "redis": {
                "host": "127.0.0.1",
                "port": 6379,
                "password": "",
                "database": 0
            }
        },
        "proxy-rewrite":{
            "regex_uri":["^/user/(.*)","/$1"],
            "host":"user-v3",
            "headers":{
                "X-Gateway":"apisix",
                "X-Target-Server-Name":"user-v3"
            }
        },
        "limit-count":{
            "count":100,
            "time_window":10,
            "key_type":"var",
            "key":"consumer_name",
            "policy":"local",
            "allow_degradation":true,
            "show_limit_quota_header":true
        }
    },
    "upstream": {
        "type": "roundrobin",
        "nodes": {
            "127.0.0.1:12008": 1
        }
    }
}'
```
