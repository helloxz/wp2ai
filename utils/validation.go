package utils

import (
	"net/mail"
	"regexp"
)

// 接收一个字符串作为参数，然后验证它是否是一个有效的电子邮件地址。
func IsEmail(email string) bool {
	// 验证电子邮件地址是否有效
	_, err := mail.ParseAddress(email)
	if err != nil {
		return false // 验证失败，返回 "false"
	}
	return true
}

// 验证密码
func IsPassword(password string) bool {
	if len(password) < 6 {
		return false
	}

	// 定义正则表达式
	pattern := `^[A-Za-z0-9!@#$%^&\*\.]+$`
	regex := regexp.MustCompile(pattern)

	return regex.MatchString(password)
}
