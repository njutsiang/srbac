# SRBAC

Service-And-Role-Based Access Control 基于服务和角色的访问控制

## 依赖服务

- MySQL：持久化存储服务、权限、角色、用户数据
- Redis：MySQL 数据的高速缓存
- Etcd：存储 APISIX 的动态配置数据
- APISIX：高性能动态网关

## 安装部署

[见详细文档](https://github.com/njutsiang/srbac/blob/main/assets/docs/install.md)

## SRBAC 核心

- 网关
- 服务
- 用户
- 角色
- 权限
- 访问控制

## SRBAC 亮点

- 适用于多服务集群
- 适用于微服务网关
- 实现对于服务和接口的访问控制
- 实现跨服务的角色权限控制
- 实现访问控制和业务代码完全解耦，业务代码只需关注业务，无需为权限控制重复造轮子

## SRBAC 权限节点类型

- 菜单权限节点：是否有权限显示某些菜单或按钮
- 接口权限节点：是否有权限请求某些接口
- 数据权限节点：是否同一个接口不同角色权限的用户请求获取到不同的数据，或执行不同的行为

## 网关鉴权

- 当网关判断到接口允许匿名访问时，直接将请求代理到该服务
- 当网关判断到用户没有登录时，直接响应：401 未登录
- 当网关判断到用户没有接口权限时，直接响应：403 没有权限
- 当网关判断到用户有接口权限时，在请求头中携带该用户拥有的该服务的数据权限，再将请求代理到该服务

## 内容管理

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

- rbac-access
  - 动态鉴权的核心插件，实现基于 SRBAC 模型的动态鉴权
- static-jwt-auth
  - 相对于 APISIX 官方插件 jwt-auth 而言更轻量级的插件，它无需添加 consumer 既可以实现 auth 相关功能
- token-auth
  - 相对于 APISIX 官方插件 key-auth 而言更轻量级、扩展性更强的插件，依托 Redis 可以完成用户相关的更多的操作：可以动态维护用户状态、用户状态自动过期、APISIX 多节点集群数据共享
- business-upstream
  - 企业动态路由插件，适用于 SaaS 企业平台，根据企业的付费等级动态路由，达到不同企业等级之间物理隔离的目的

## 使用指南

[见详细文档](https://github.com/njutsiang/srbac/blob/main/assets/docs/manual.md)

## 界面截图

<img src="https://github.com/njutsiang/srbac/raw/main/assets/img/screely-1645003820980.png">
<img src="https://github.com/njutsiang/srbac/raw/main/assets/img/screely-1645004130502.png">
<img src="https://github.com/njutsiang/srbac/raw/main/assets/img/screely-1645004069102.png">
