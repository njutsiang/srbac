local ngx = ngx
local ipairs = ipairs
local core = require("apisix.core")
local redis = require("resty.redis")

-- 定义插件配置的数据结构
local schema = {
    type = "object",
    properties = {
        redis = {
            type = "object",
            properties = {
                host = {
                    type = "string"
                },
                port = {
                    type = "integer"
                },
                password = {
                    type = "string"
                },
                database = {
                    type = "integer"
                }
            }
        }
    },
    required = {
        "redis"
    },
}

-- 定义插件
local _M = {
    version = 1.0,
    priority = 97,
    name = "rbac-access",
    schema = schema,
}

-- 校验插件配置数据
function _M.check_schema(conf)
    return core.schema.check(schema, conf)
end

-- 没有登录
local function get_response_401()
    core.response.set_header("Content-Type", "application/json")
    return 401, core.json.encode({
        code = 401,
        message = "没有登录"
    }, true)
end

-- 没有权限
local function get_response_403()
    core.response.set_header("Content-Type", "application/json")
    return 403, core.json.encode({
        code = 403,
        message = "没有权限"
    }, true)
end

-- 接口不存在
local function get_response_404()
    core.response.set_header("Content-Type", "application/json")
    return 404, core.json.encode({
        code = 404,
        message = "接口不存在"
    }, true)
end

-- 获取请求的 uri
local function get_uri(ctx)
    local uri = ngx.var.uri
    if ctx.var.upstream_uri and ctx.var.upstream_uri ~= "" then
        uri = ctx.var.upstream_uri
    end
    local i = core.string.find(uri, "?")
    if i then
        uri = string.sub(uri, 1, i - 1)
    end
    return uri
end

-- 判断指定角色对指定接口是否有权限
local function is_role_has_permission(red, role, service, method, uri)
    local key = "auth:role:" .. role .. ":service:" .. service .. ":apis"
    if red:sismember(key, method .. uri) == 1 then
        return true
    end
    return red:sismember(key, "*" .. uri) == 1
end

-- 判断指定用户对指定接口是否有权限
local function is_user_has_permission(red, user_id, service, method, uri)
    local key = "auth:user:" .. user_id .. ":service:" .. service .. ":apis"
    if red:sismember(key, method .. uri) == 1 then
        return true
    end
    return red:sismember(key, "*" .. uri) == 1
end

-- 获取用户的角色
local function get_roles(red, user_id)
    local key = "auth:user:" .. user_id .. ":roles"
    local roles = red:smembers(key);
    if roles then
        return roles
    else
        return {}
    end
end

-- 没有权限
local function is_has_permission(ctx, red, roles)
    local user_id = ctx.user_id
    local service = ctx.var.upstream_host
    local method = ngx.req.get_method()
    local uri = get_uri(ctx)
    for _, role in ipairs(roles) do
        if is_role_has_permission(red, role, service, method, uri) then
            return true
        end
    end
    if is_user_has_permission(red, user_id, service, method, uri) then
        return true
    end
    return false
end

-- 获取角色拥有的数据权限节点
local function get_auth_items_by_role(red, role, service)
    local key = "auth:role:" .. role .. ":service:" .. service .. ":items"
    local auth_items = red:smembers(key);
    if auth_items then
        return auth_items
    else
        return {}
    end
end

-- 获取用户拥有的数据权限节点
local function get_auth_items_by_user_id(red, user_id, service)
    local key = "auth:user:" .. user_id .. ":service:" .. service .. ":items"
    local auth_items = red:smembers(key);
    if auth_items then
        return auth_items
    else
        return {}
    end
end

-- 获取角色和用户拥有的数据权限节点
local function get_auth_items(red, user_id, roles, service)
    local auth_items = {}
    local _auth_items = {}
    for _, role in ipairs(roles) do
        _auth_items = get_auth_items_by_role(red, role, service)
        for _, auth_item in ipairs(_auth_items) do
            auth_items[auth_item] = 1
        end
    end
    _auth_items = get_auth_items_by_user_id(red, user_id, service)
    for _, auth_item in ipairs(_auth_items) do
        auth_items[auth_item] = 1
    end
    local result_items = {}
    for auth_item, _ in pairs(auth_items) do
        table.insert(result_items, auth_item)
    end
    return result_items
end

-- 清除用户自定义请求头
local function clear_header()
    ngx.req.clear_header('X-User-Auth-Items')
end

-- 在 access 阶段执行
function _M.access(conf, ctx)
    clear_header()

    -- 实例化 Redis，超时时间为 1s
    local red = redis:new()
    red:set_timeouts(1000, 1000, 1000)

    -- Redis 建立连接
    local ok, err = red:connect(conf.redis.host, conf.redis.port)
    if not ok then
        core.log.error("Redis 连接失败: ", err)
        return ""
    end

    -- Redis 验证密码
    if conf.redis.password and conf.redis.password ~= "" then
        ok, err = red:auth(conf.redis.password)
        if not ok then
            core.log.error("Redis 密码错误: ", err)
            return ""
        end
    end

    -- 切换 Redis 数据库
    ok, err = red:select(conf.redis.database)
    if not ok then
        core.log.error("Redis 切换数据库失败: ", err)
        return ""
    end

    -- 判断接口是否存在
    local key = "auth:service:" .. ctx.var.upstream_host .. ":apis"
    local uri = get_uri(ctx)
    local data = red:hget(key, ngx.req.get_method() .. uri)
    if not (type(data) == "string" and (data == "1" or data == "0")) then
        data = red:hget(key, "*" .. uri)
        if not (type(data) == "string" and (data == "1" or data == "0")) then
            return get_response_404()
        end
    end

    -- 如果接口是需要鉴权的
    if data == "1" then

        -- 判断用户是否已登录
        if not ctx.user_id then
            return get_response_401()
        end

        -- 查询用户拥有的角色
        local roles = get_roles(red, ctx.user_id)

        -- 判断用户是否有权限
        if not is_has_permission(ctx, red, roles) then
            return get_response_403()
        end

        -- 将用户拥有的数据权限节点写入到请求头
        local service = ctx.var.upstream_host
        local auth_items = get_auth_items(red, ctx.user_id, roles, service)
        core.response.set_header("X-User-Auth-Items", core.json.encode(auth_items, true))
    end

    -- 将当前 Redis 连接放入连接池
    -- 10000，空闲连接保持 10s
    -- 100，连接池最大保持 100 个连接
    ok, err = red:set_keepalive(10000, 100)
    if not ok then
        core.log.error("设置 Redis 连接池失败: ", err)
    end
end

-- 返回插件
return _M