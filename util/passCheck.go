package util

import (
	"github.com/kataras/golog"
	"golang.org/x/crypto/bcrypt"
)

var passwordToken = "MIICWwIBAAKBgQCraRaunSw1bMXeGL908snY6mVbWzp2nlIqfo6UKYd"

// Hash 密码hash
func Hash(password string) string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(passwordToken+password), bcrypt.DefaultCost)
	if err != nil {
		golog.Fatalf("pkg.password.password %s", err)
		return ""
	}

	return string(bytes)
}

// Verify 密码hash验证
func Verify(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(passwordToken+password))
	return err == nil
}
