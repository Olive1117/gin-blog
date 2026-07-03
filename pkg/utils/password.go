package utils

import "golang.org/x/crypto/bcrypt"

// HashPassword 对明文密码进行 bcrypt 哈希
func HashPassword(password string) (string, error) {
	// bcrypt.DefaultCost = 10，可以根据服务器性能调整（通常 10-12）
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// CheckPassword 验证明文密码是否与哈希匹配
func CheckPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
