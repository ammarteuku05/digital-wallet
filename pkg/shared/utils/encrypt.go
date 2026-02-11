package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"digital-wallet/pkg/response"
	"encoding/base64"
)

var bytess = []byte{35, 46, 57, 24, 85, 35, 24, 74, 87, 35, 88, 98, 66, 32, 14, 05}

func Encode(b []byte) string {
	return base64.StdEncoding.EncodeToString(b)
}

func Encrypt(text, MySecret string) (string, error) {
	block, err := aes.NewCipher([]byte(MySecret))
	if err != nil {
		return "", response.Wrap(err, "cannot init new chiper")
	}
	plainText := []byte(text)
	cfb := cipher.NewCFBEncrypter(block, bytess)
	cipherText := make([]byte, len(plainText))
	cfb.XORKeyStream(cipherText, plainText)
	return Encode(cipherText), nil
}
