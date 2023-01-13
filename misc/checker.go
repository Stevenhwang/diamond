package misc

import "github.com/Stevenhwang/gommon/tools"

func Checker(ip string, hashedPassword string, password string) bool {
	if tools.CheckPassword(hashedPassword, password) {
		return true
	} else {
		Cache.Set(ip, []byte{1}) // 加入黑名单
		return false
	}
}
