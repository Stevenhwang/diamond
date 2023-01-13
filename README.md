## About

diamond 是一个 golang 开发的完全开源的 devops 自动化运维平台

- [x] 堡垒机(支持 web 终端和 ssh 客户端，支持密码和公钥连接)
- [x] 服务器(分配管理)
- [x] 权限控制(前端菜单，后端接口，ACL)
- [ ] 云资源(同步，操作)
- [ ] 作业平台(批量执行)
- [ ] 定时任务
- [ ] 域名(同步，操作)

持续开发中。。。

## Require

```bash
golang 1.18+
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
  seed        seed admin account[创建admin账户]
  sshd        start sshd server[开启sshd服务]
```
