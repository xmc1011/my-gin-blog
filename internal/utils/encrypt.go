package utils

import "golang.org/x/crypto/bcrypt"

// 使用 bcrypt 对字符串进行加密生成一个哈希值
func BcryptHash(str string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(str), bcrypt.DefaultCost)
	return string(bytes), err
}

// 使用 bcrypt 对比 明文字符串 和 哈希值
func BcryptCheck(plain, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(plain))
	return err == nil
}
