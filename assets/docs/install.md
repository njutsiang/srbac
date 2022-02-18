# 安装部署

## 依赖服务

请自行安装以下服务，版本号以下面经过测试的版本为示例，其他版本号是否兼容请以各开源项目官方变更日志为准。

- lua-5.1
- openresty-1.19.9.1
- etcd-v3.5.1
- apisix-2.11.0
- apisix-dashboard-2.10.1
- redis-*
- mysql-5.7 | 8.*

## 部署 SRBAC

下载最新的 Release 版本 [https://github.com/njutsiang/srbac/releases](https://github.com/njutsiang/srbac/releases)

1. 下载解压

```shell
wget https://github.com/njutsiang/srbac/releases/download/[版本号]/srbac-bin.zip
unzip srbac-bin.zip
cd srbac-bin
```

2. 修改配置文件

```
vi ./config.yaml
```

```text
# 监听端口
listen:
  port: 8000

# 默认超级管理员
# 仅在系统初始化的时候生成超级管理员，后续可在 Web 管理后台修改超管账号和密码
super:
  username: admin
  password: 123456

# MySQL 连接配置
mysql:
  host: 127.0.0.1
  port: 3306
  user: root
  password: 123456
  db: srbac
  charset: utf8mb4
  
# Redis 连接配置
redis:
  host: 127.0.0.1
  port: 6379
  password:
  db: 0

# Cookie 加密的密钥
cookie:
  key: 19be30e8a9850e5cb5d1f7e3e47bde35
```

3. 导入 MySQL 表结构文件

```shell
mysql -h127.0.0.1 -uroot -p123456 -P3306 [你的库名] < mysql.sql
```

4. 启动服务

```shell
./main
```

假设你部署好的 SRBAC：

IP：127.0.0.1 <br>
Port：8000（为什么不是 80？因为 80 端口留给一会儿 APISIX 用）<br>
Domain：srbac.local.com

## 安装 SRBAC 的 APISIX 插件

```shell
# 1. 进入 SRBAC 项目主目录
cd ./srbac

# 2. 将 SRBAC 提供的 APISIX 插件拷贝到你部署好的 APISIX 的插件目录
cp ./apisix-plugins/* /your-apisix-path/apisix/plugins

# 3. 编辑 APISIX 配置文件
vi /your-apisix-path/conf/config.yaml
# 将 SRBAC 提供的插件添加到 plugins 末尾，就像这样：
plugins:
  - ......
  - token-auth
  - static-jwt-auth
  - rbac-access

# 4. 重启 APISIX
/your-apisix-path/bin/apisix stop
/your-apisix-path/bin/apisix start
```

更多使用指南，[见详细文档](https://github.com/njutsiang/srbac/blob/main/assets/docs/manual.md)
