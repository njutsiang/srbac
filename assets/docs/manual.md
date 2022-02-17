# 使用指南

## 通过 APISIX 访问 SRBAC

假设你已部署好 Redis：

Redis Host：127.0.0.1 <br>
Redis Port：6379 <br>
Redis Password：无

假设你已部署好 APISIX：

Admin API Host：http://127.0.0.1:9180 <br>
Admin API API-KEY：edd1c9f034335f136f87ad84b625c8f1

将 SRBAC 添加到 APISIX 的上游列表： <br>
更多关于 APISIX 的操作见 [官方文档](https://apisix.apache.org/zh/docs/apisix/admin-api/)

```shell
curl -X POST 'http://127.0.0.1:9180/apisix/admin/upstreams' \
-H 'X-API-KEY: edd1c9f034335f136f87ad84b625c8f1' \
-d '{
    "name": "SRBAC",
    "type": "roundrobin",
    "nodes": {
        "127.0.0.1:8000": 1
    },
    "retries": 2,
    "retry_timeout": 3,
    "timeout": {
        "connect": 3,
        "send": 3,
        "read": 3
    },
    "scheme": "http"
}'

# 添加后得到 upstream_id：394927496115520497
```

添加路由，使 SRBAC 的所有请求都经过 APISIX：

```shell
curl -X POST 'http://127.0.0.1:9180/apisix/admin/routes' \
-H 'X-API-KEY: edd1c9f034335f136f87ad84b625c8f1' \
-d '{
    "name":"SRBAC",
    "desc":"srbac-service",
    "host":"srbac.local.com",
    "uri":"/*",
    "methods":["GET","POST","PUT","DELETE","HEAD","OPTIONS"],
    "plugins":{
        "token-auth":{
            "header":"authorization",
            "query":"token",
            "cookie":"user_token",
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
            "host":"srbac-service",
            "headers":{
                "X-Gateway":"apisix",
                "X-Target-Service":"srbac-service"
            }
        },
        "response-rewrite":{
            "headers":{
                "X-Gateway":"apisix",
                "X-Target-Service":"srbac-service"
            }
        }
    },
    "upstream_id":"394927496115520497"
}'
```

添加静态文件路由，使静态文件不经过鉴权，如果静态文件使用了单独域名部署，就没有该步骤了：

```shell
curl -X POST 'http://127.0.0.1:9180/apisix/admin/routes' \
-H 'X-API-KEY: edd1c9f034335f136f87ad84b625c8f1' \
-d '{
    "name":"SRBAC",
    "desc":"srbac-service",
    "host":"srbac.local.com",
    "uri":"/assets/*",
    "methods":["GET","HEAD","OPTIONS"],
    "plugins":{
        "proxy-rewrite":{
            "host":"srbac-service",
            "headers":{
                "X-Gateway":"apisix",
                "X-Target-Service":"srbac-service"
            }
        },
        "response-rewrite":{
            "headers":{
                "X-Gateway":"apisix",
                "X-Target-Service":"srbac-service"
            }
        }
    },
    "upstream_id":"394927496115520497"
}'
```

## 通过 APISIX + SRBAC 访问你的任何服务

如同以上将 SRBAC 加入到 APISIX 一样，你的其他服务都能够以相同的方式加入到 APISIX，并且享受 SRBAC 提供的完全解耦的网关鉴权，这里将对其中的配置项加以说明。

1. 首先需要在 SRBAC 管理后台添加服务、接口节点、分配相关权限等等；
2. 然后再将服务作为上游添加到 APISIX；
3. 最后向 APISIX 添加路由规则，并且使用插件：token-auth、static-jwt-auth、rbac-access。

```shell
curl -X POST 'http://127.0.0.1:9180/apisix/admin/routes' \
-H 'X-API-KEY: edd1c9f034335f136f87ad84b625c8f1' \
-d '{
    // 项目名称
    "name":"订单中心",
    
    // 项目域名
    "host":"order.local.com",
    
    // 需要经过网关鉴权的路由，/* 表示所有路由
    "uri":"/*",
    
    // 允许的请求方式
    "methods":["GET","POST","PUT","DELETE","HEAD","OPTIONS"],
    
    // 以下配置插件
    "plugins":{
    
        // 启用用户身份识别插件
        // 如下配置表示依次尝试从 header、query、cookie 中获取用户 token
        // 并且用户的登录状态是保存在 Redis 中，能够通过这样的键值查询到：user:token:1f2cde20df5b02771b0f3c2746ad9deb
        // 从 Redis 中查询到的用户数据需要是至少包含 id 的 JSON 字符串，例如：{"id": 1,"name": "超级管理员",......}
        // 如果项目中保持用户登录状态的方式不是 token，而是 JWT，则需要改用插件：static-jwt-auth
        "token-auth":{
            "header":"authorization",
            "query":"token",
            "cookie":"user_token",
            "redis": {
                "host": "127.0.0.1",
                "port": 6379,
                "password": "",
                "database": 0
            }
        },

        // 启用用户身份识别插件
        // 该插件和 token-auth 插件只需要二选一
        // secret：JWT 加密的密钥
        // response_err：JWT 不存在时，是否返回 401 错误，false 表示不返回错误，继续执行 rbac-access 插件，可能会被允许访问一些允许匿名访问的接口
        "static-jwt-auth":{
            "header":"authorization",
            "query":"token",
            "secret":"123456",
            "response_err":false
        },
        
        // 启用 SRBAC 鉴权插件
        // 这里的 Redis 连接配置必须和 SRBAC 项目的 Redis 是同一个
        "rbac-access":{
            "redis": {
                "host": "127.0.0.1",
                "port": 6379,
                "password": "",
                "database": 0
            }
        },
        
        // 重写请求
        // host：项目服务标识，将以此识别目标服务，该值为配置在 SRBAC 管理后台的服务标识
        "proxy-rewrite":{
            "host":"order-service",
            "headers":{
                "X-Gateway":"apisix",
                "X-Target-Service":"order-service"
            }
        },
        
        // 重写响应
        "response-rewrite":{
            "headers":{
                "X-Gateway":"apisix",
                "X-Target-Service":"srbac-service"
            }
        }
    },
    
    // 上游 id，从调用 APISIX 添加上游的接口返回值中得到
    "upstream_id":"394927496115520497"
}'
```