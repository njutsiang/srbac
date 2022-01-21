local ngx = ngx
local core = require("apisix.core")
local jwt = require("resty.jwt")

-- 定义数据结构
-- 在定义 Route、Service 时，使用该插件时的数据结构
local schema = {
    type = "object",
    properties = {
        header = {
            type = "string"
        },
        query = {
            type = "string"
        },
        secret = {
            type = "string"
        },
        base64_encode = {
            type = "boolean",
            default = false
        },
        response_err = {
            type = "boolean",
            default = false
        }
    },
    required = {
        "header", "query", "secret"
    },
}

-- 定义插件
-- type = "auth"，表示该插件是认证类型的插件，可以在定义 Consumer 时使用
local _M = {
    version = 1.0,
    priority = 99,
    type = "auth",
    name = "static-jwt-auth",
    schema = schema,
}

-- 重写请求头和响应头
local function set_jwt_header(ctx, token)
    core.request.set_header(ctx, "X-JWT", token)
    core.response.set_header("X-JWT", token)
end

-- 重写请求头和响应头
local function set_jwt_failed_header(ctx, token)
    core.request.set_header(ctx, "X-JWT-Failed", token)
    core.response.set_header("X-JWT-Failed", token)
end

-- 重写请求头和响应头
local function set_user_header(ctx, jwt_token, payload)
    -- 处理数据
    local user_id = payload.user_id
    if not user_id then
        user_id = payload.key
    end
    if not user_id then
        user_id = payload.sub
    end
    if not user_id then
        user_id = 0
    end
    user_id = tostring(user_id)
    local user_data = core.json.encode(payload, true)

    -- 实现 auth 插件
    ctx.user_id = user_id
    ctx.consumer_name = user_id
    ctx.consumer_ver = 1
    ctx.consumer = {
        id = user_id,
        username = user_id,
        consumer_name = user_id,
        auth_conf = {
            key = jwt_token
        },
        plugins = {}
    }
    ctx.consumer.plugins["key-auth"] = {
        key = jwt_token
    }

    -- 写入到请求头、响应头
    core.request.set_header(ctx, "X-User-Id", user_id)
    core.request.set_header(ctx, "X-User-Data", user_data)
    core.response.set_header("X-User-Id", user_id)
    core.response.set_header("X-User-Data", user_data)
end

-- 获取错误响应
local function get_response_err(conf, message)
    core.response.set_header("Content-Type", "application/json")
    if conf.response_err then
        return 401, {
            err_code = 401,
            err_message = message
        }
    else
        return nil
    end
end

-- 清除用户自定义请求头
local function clear_header()
    ngx.req.clear_header('X-JWT')
    ngx.req.clear_header('X-JWT-Failed')
    ngx.req.clear_header('X-User-Id')
    ngx.req.clear_header('X-User-Data')
end

-- 在 rewrite 阶段执行
function _M.rewrite(conf, ctx)
    clear_header()

    -- 从请求头获取 JWT Token
    local jwt_token = core.request.header(ctx, conf.header)

    -- 从 Query 参数中获取 JWT Token
    if not (jwt_token and jwt_token ~= "") then
        local uri_args = core.request.get_uri_args(ctx) or {}
        jwt_token = uri_args[conf.query]
        if not (jwt_token and jwt_token ~= "") then
            set_jwt_header(ctx, "NULL")
            return get_response_err(conf, "JWT 不存在")
        end
    end

    -- 处理 JWT Token
    local prefix = string.sub(jwt_token, 1, 7)
    if prefix == "Bearer " or prefix == "bearer " then
        jwt_token = string.sub(jwt_token, 8)
    end
    set_jwt_header(ctx, jwt_token)

    -- 解析 JWT Token
    local jwt_obj = jwt:load_jwt(jwt_token)
    if not jwt_obj.valid then
        set_jwt_failed_header(ctx, "valid: "..jwt_obj.reason)
        return get_response_err(conf, jwt_obj.reason)
    end

    -- 校验 JWT Token 的签名
    jwt_obj = jwt:verify_jwt_obj(conf.secret, jwt_obj)
    if not jwt_obj.verified then
        set_jwt_failed_header(ctx, "verified: "..jwt_obj.reason)
        return get_response_err(conf, jwt_obj.reason)
    end

    -- 写入到请求头、响应头
    set_user_header(ctx, jwt_token, jwt_obj.payload)
end

return _M