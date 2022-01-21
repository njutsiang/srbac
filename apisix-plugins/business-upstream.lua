local ngx = ngx
local core = require("apisix.core")
local redis = require("resty.redis")

-- 定义插件配置的数据结构
local schema = {
    type = "object",
    properties = {
        header = {
            type = "string",
        },
        query = {
            type = "string",
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

-- 定义插件
local _M = {
    version = 1.0,
    priority = 98,
    name = "business-upstream",
    schema = schema,
}

-- 校验插件配置数据
function _M.check_schema(conf)
    return core.schema.check(schema, conf)
end

-- 重写响应头
local function set_response_header(company_id)
    if not company_id then
        company_id = 0
    end
    core.response.set_header("X-Company-Id", tostring(company_id))
end

-- 重写响应头
local function set_response_header_result(upstream_id)
    core.response.set_header("X-Apisix-Upstream-Id", upstream_id)
end

-- 定义响应的错误
-- Lua 是强类型，1 和 "1" 是不相等的
local function get_response_error()
    core.response.set_header("Content-Type", "application/json")
    return 400, core.json.encode({
        code = 400,
        message = "没有找到 company_id"
    }, true)
end

-- 通过公司类型获取 upstream_id
local function get_upstream_id_by_level(company_level)
    local upstream_ids = {
        vip0 = "00000000000000001760",
        vip1 = "00000000000000001762",
        vip2 = "00000000000000001765"
    }
    local upstream_id = upstream_ids[company_level]
    if not upstream_id then
        company_level = "vip0"
    end
    core.response.set_header("X-Company-Level", company_level)
    return upstream_ids[company_level]
end

-- 通过 company_id 获取公司类型
local function get_company_level(conf, company_id)
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

    -- 从 Redis 读取公司信息
    local company_level, company_err = red:get("level:company_id:"..company_id)
    if company_level then
        if company_level == ngx.null then
            core.log.error("从 Redis 读取公司信息为 NULL")
            return ""
        end
    else
        core.log.error("从 Redis 读取公司信息失败: ", company_err)
        return ""
    end

    -- 将当前 Redis 连接放入连接池
    -- 10000，空闲连接保持 10s
    -- 100，连接池最大保持 100 个连接
    ok, err = red:set_keepalive(10000, 100)
    if not ok then
        core.log.error("设置 Redis 连接池失败: ", err)
    end

    company_level = tostring(company_level)
    if company_level == "userdata: NULL" then
        company_level = ""
    end

    return company_level
end

-- 清除用户自定义请求头
local function clear_header()
    ngx.req.clear_header('X-Company-Id')
    ngx.req.clear_header('X-Company-Level')
    ngx.req.clear_header('X-Apisix-Upstream-Id')
end

-- 在 access 阶段执行
function _M.access(conf, ctx)
    clear_header()

    -- 从 Header 或 Query 中获取 company_id
    local company_id = core.request.header(ctx, conf.header);
    if not company_id then
        local uri_args = core.request.get_uri_args(ctx) or {}
        company_id = uri_args[conf.query]
    end
    set_response_header(company_id)
    if not company_id then
        return get_response_error()
    end

    -- 指定上游 id
    -- 具体信息还可以在上游配置中详细配置，包括：
    -- 后端节点、负载均衡算法、健康检测
    local company_level = get_company_level(conf, company_id)
    ctx.upstream_id = get_upstream_id_by_level(company_level)
    set_response_header_result(ctx.upstream_id)
end

-- 返回插件
return _M