# SRBAC

Service-And-Role-Based Access Control 基于服务和角色的访问控制

## 依赖服务

- MySQL：持久化存储服务、权限、角色、用户数据
- Redis：MySQL 数据的高速缓存
- etcd：存储 apisix 的动态配置数据
- apisix：高性能鉴权网关

## 安装部署

[见详细文档](http://)

## SRBAC 核心

- 服务
- 权限
- 角色
- 用户
- 访问控制

## SRBAC 亮点

- 适用于多服务集群
- 适用于微服务网关
- 实现对于服务和接口的访问控制
- 实现跨服务的角色权限控制
- 权限控制和业务代码完全解耦
- 安全、可靠、高效

## SRBAC 权限节点类型

- 菜单权限节点：是否有权限显示某些菜单或按钮
- 接口权限节点：是否有权限请求某些接口
- 数据权限节点：是否同一个接口不同角色权限的用户请求获取到不同的数据，或执行不同的行为

## 网关鉴权

- 当网关判断到接口允许匿名访问时，直接将请求代理到该服务
- 当网关判断到用户没有登录时，直接响应：401 未登录
- 当网关判断到用户没有接口权限时，直接响应：403 没有权限
- 当网关判断到用户有接口权限时，在请求头中携带该用户拥有的该服务的数据权限，再将请求代理到该服务

## 系统组成

- 服务管理
  - 每个服务必须有一个唯一标识
  - 服务的增删改查
- 用户管理
  - 每个用户必须有一个唯一 id
  - 用户的增删改查
- 角色管理
  - 角色的增删改查
  - 配置一个角色拥有哪些服务的权限
- 权限节点管理
  - 基于服务管理权限节点，即所有权限节点是挂在具体的服务下面的
  - 菜单权限节点的增删改查
  - 接口权限节点的增删改查
  - 数据权限节点的增删改查
- 角色权限分配
  - 给角色分配具体服务下面的菜单权限节点
  - 给角色分配具体服务下面的接口权限节点
  - 给角色分配具体服务下面的数据权限节点
- 用户角色权限分配
  - 给用户分配角色
  - 给用户分配菜单权限节点
  - 给用户分配接口权限节点
  - 给用户分配数据权限节点
  - 用户的最终权限是所分配的角色和直接分配的权限节点的并集

## APISIX 插件

- business-upstream
  - 企业动态路由插件，适用于 SaaS 企业平台，根据企业的付费等级动态路由，达到不同企业等级之间物理隔离的目的
- rbac-access
  - 动态鉴权的核心插件，实现基于 SRBAC 模型的动态鉴权
- static-jwt-auth
  - 相对于 APISIX 官方插件 jwt-auth 而言更轻量级的插件，它无需添加 consumer 既可以实现 auth 相关功能
- token-auth
  - 相对于 APISIX 官方插件 key-auth 而言更轻量级、扩展性更强的插件，依托 Redis 可以完成用户相关的更多的操作：可以动态维护用户状态、用户状态自动过期、APISIX 多节点集群数据共享

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
