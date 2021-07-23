## About

diamond 是一个 golang 开发的完全开源的 devops 自动化运维平台

- [x] 堡垒机(登录，鉴权，录屏)
- [x] 服务器(分组，权限)
- [x] 权限控制(前端，后端)
- [ ] 云资源(同步，操作)
- [ ] 作业平台(批量执行)
- [ ] 定时任务
- [ ] 域名(同步，操作)
- [ ] 监控

更多功能，持续开发中。。。

## Require

```bash
golang 1.16+
mysql 5.7+
redis 3.2+
```

## Installation

```bash
go build
```

## Settings

```bash
config.json
```

## Usage

```bash
Available Commands:
  api         start api server[开启 api 服务器]
  completion  generate the autocompletion script for the specified shell
  help        Help about any command
  migrate     auto migrate[运行自动迁移]
  seed        seed the database[创建admin账户]
  sshd        start sshd server[开启 sshd 服务器]
  syncperm    sync permissions info[同步权限信息]
```

## Shoulders

- [gin-gonic/gin](https://github.com/gin-gonic/gin)
- [gorm](https://gorm.io/)
- [spf13/cobra](https://github.com/spf13/cobra)
- [spf13/viper](https://github.com/spf13/viper)
- [gliderlabs/ssh](https://github.com/gliderlabs/ssh)
- [go-redis/redis](https://github.com/go-redis/redis)
- [dgrijalva/jwt-go](https://github.com/dgrijalva/jwt-go)
- [gobuffalo/nulls](https://github.com/gobuffalo/nulls)
- [google/uuid](https://github.com/google/uuid)
- [olekukonko/tablewriter](https://github.com/olekukonko/tablewriter)
- [pquerna/otp](https://github.com/pquerna/otp)
- [golang.org/x/term](https://golang.org/x/term)
- [robfig/cron](https://github.com/robfig/cron)
