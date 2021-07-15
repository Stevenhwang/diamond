package utils

import (
	"crypto/md5"
	"encoding/hex"
	"math/rand"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// FindValInSlice 查询val是否在 string slice 中
func FindValInSlice(slice []string, val string) bool {
	for _, v := range slice {
		if v == val {
			return true
		}
	}
	return false
}

// GeneratePassword 生成bycrypt加密密码
func GeneratePassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

// CheckPassword 检验bycrypt加密密码
func CheckPassword(hashedPassword string, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}

// MD5 生成md5字符串
func MD5(v string) string {
	d := []byte(v)
	m := md5.New()
	m.Write(d)
	return hex.EncodeToString(m.Sum(nil))
}

// GenerateRandCode 生成随机码
// candidate 候选字符串
// length  随机码长度
func GenerateRandCode(candidate string, length int) string {
	var letters = []rune(candidate)
	b := make([]rune, length)
	rand.Seed(time.Now().UnixNano())
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

// RemoveDupInSlice 去掉 string slice 中的重复元素
func RemoveDupInSlice(slice []string) []string {
	result := []string{}
	tempMap := map[string]byte{} // 存放不重复主键
	for _, e := range slice {
		l := len(tempMap)
		tempMap[e] = 0
		if len(tempMap) != l { // 加入map后，map长度变化，则元素不重复
			result = append(result, e)
		}
	}
	return result
}
