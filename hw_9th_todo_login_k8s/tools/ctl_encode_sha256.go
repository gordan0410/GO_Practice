package tools

import (
	"crypto/sha256"
	"encoding/hex"
)

// sha256加密
func Encode_password(password string) (string, error) {
	h := sha256.New224()
	_, err := h.Write([]byte(password))
	if err != nil {
		return "", err
	}
	result := hex.EncodeToString(h.Sum(nil))
	return result, nil
}
