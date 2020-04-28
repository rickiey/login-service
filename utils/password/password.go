package password

import (
	"golang.org/x/crypto/bcrypt"
)

const (
	//$2a$(2 chars work cost)$(22 chars salt )(31 chars hash)
	DefaultSaltLen = 29
)

// 生成密码
// Get returns the encrypted string for specified utils
func Get(str string) (string, error) {
	//TODO:Perhaps than sha256 performance slightly low;But it is more safe 20170609

	hash, err := bcrypt.GenerateFromPassword([]byte(str), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}
// 对比加密后的内容和密码
// Compare is to compare encryptStr and provided utils
func Compare(encryptStr, str string) bool {
	// Comparing the password with the hash
	return nil == bcrypt.CompareHashAndPassword([]byte(encryptStr), []byte(str))
}

// GetSalt returns salt from EncryptPassword
func GetSalt(encryptStr string) string {
	return string(encryptStr[0 : DefaultSaltLen+1])
}
