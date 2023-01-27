## About

diamond 是一个 golang 开发的完全开源的 devops 自动化运维平台

- [x] 堡垒机(支持 web 终端和 ssh 客户端，支持密码和公钥连接)
- [x] 权限控制(前端菜单，后端接口，服务器分配，ACL 权限)
- [x] 任务平台(linux command，可使用 ansible 或其他工具辅助)
- [x] 定时任务(支持秒级的 crontab 表达式)
- [x] 支持 ssh 隧道，便于 navicate 连接使用
- [x] 内置 socks5 代理，便于跟服务器内网通讯(比如使用 ftp 客户端等)

持续开发中。。。

## Require

```bash
Linux
golang 1.19+
mysql 5.8+
redis 3.2+
```

## Installation

```bash
前端：cd frontend && npm install && npm run build
后端: go build
编译出的二进制使用embed将前端打包的dist嵌入，所以可以单独部署
```

## Settings

```bash
config.yml
```

## Usage

```bash
Available Commands:
  app         start app server[开启app服务]
  seed        seed user account[创建用户账户]
  socks5      start socks5 server[开启socks5服务]
  sshd        start sshd server[开启sshd服务]
```
