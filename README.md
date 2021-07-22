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
