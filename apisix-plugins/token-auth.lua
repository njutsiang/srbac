local ngx = ngx
local core = require("apisix.core")
local redis = require("resty.redis")

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
        "header", "query", "redis"
    },
}

-- 定义数据结构
-- 在定义 Consumer 时，使用该插件时的数据结构
local consumer_schema = {
    type = "object",
    properties = {
        token = {
            type = "string",
        },
    },
    required = {"key"},
}

-- 定义插件
-- type = "auth"，表示该插件是认证类型的插件，可以在定义 Consumer 时使用
local _M = {
    version = 1.0,
    priority = 99,
    type = "auth",
    name = "token-auth",
    schema = schema,
    consumer_schema = consumer_schema,
}

-- 检验数据结构
function _M.check_schema(conf, schema_type)
    if schema_type == core.schema.TYPE_CONSUMER then
        return core.schema.check(consumer_schema, conf)
    else
        return core.schema.check(schema, conf)
    end
end

-- 重写请求头和响应头
local function set_header(ctx, user_id, user_json)
    core.request.set_header(ctx, "X-User-Id", user_id)
    core.request.set_header(ctx, "X-User-Data", user_json)
    core.response.set_header("X-User-Id", user_id)
    core.response.set_header("X-User-Data", user_json)
end

-- 在 rewrite 阶段执行
function _M.rewrite(conf, ctx)

    -- 从请求头获取 token，请求头的 Key 不区分大小写
    local token = core.request.header(ctx, conf.header)

    -- 从 Uri Query 参数中获取 token，Uri Query 参数的 Key 是要区分大小写的
    if not (token and token ~= "") then
        local uri_args = core.request.get_uri_args(ctx) or {}
        token = uri_args[conf.query]
    end

    -- 如果 token 存在，就从 Redis 中去读取更多数据
    -- 如果 token 不存在，获取到为 nil
    if not (token and token ~= "") then
        set_header(ctx, 0, "NULL")
        return nil
    end

    -- 实例化 Redis，超时时间为 1s
    local red = redis:new()
    red:set_timeouts(1000, 1000, 1000)

    -- Redis 建立连接
    local ok, err = red:connect(conf.redis.host, conf.redis.port)
    if not ok then
        core.log.error("Redis 连接失败: ", err)
        set_header(ctx, 0, "NULL")
        return nil
    end

    -- Redis 验证密码
    if conf.redis.password and conf.redis.password ~= "" then
        ok, err = red:auth(conf.redis.password)
        if not ok then
            core.log.error("Redis 密码错误: ", err)
            set_header(ctx, 0, "NULL")
            return nil
        end
    end

    -- 切换 Redis 数据库
    ok, err = red:select(conf.redis.database)
    if not ok then
        core.log.error("Redis 切换数据库失败: ", err)
        set_header(ctx, 0, "NULL")
        return nil
    end

    -- 从 Redis 读取用户信息
    local user_json, user_err = red:get("user:token:"..token)
    if not user_json then
        core.log.error("从 Redis 读取用户失败: ", user_err)
        set_header(ctx, 0, "NULL")
        return nil
    end

    if user_json == ngx.null then
        core.log.error("从 Redis 读取到用户为 NULL")
        set_header(ctx, 0, "NULL")
        return nil
    end

    local user_data = core.json.decode(user_json)
    if user_data then
        local user_id = tostring(user_data.id)
        ctx.user_id = user_id
        ctx.consumer_name = user_id
        ctx.consumer_ver = 1
        ctx.consumer = {
            id = user_id,
            username = user_id,
            consumer_name = user_id,
            auth_conf = {
                key = token
            },
            plugins = {}
        }
        ctx.consumer.plugins["key-auth"] = {
            key = token
        }
        set_header(ctx, user_id, user_json)
    else
        core.log.error("从 Redis 读取到用户不是 JSON")
        set_header(ctx, 0, "NO JSON")
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