package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"
)

var (
	ErrInvalidKeySize = errors.New("encryption key must be 32 bytes for AES-256")
	ErrInvalidCiphertext = errors.New("invalid ciphertext format")
)

// EncryptCredentials 使用AES-256-GCM加密凭证
func EncryptCredentials(plaintext []byte, key []byte) ([]byte, error) {
	if len(key) != 32 {
		return nil, ErrInvalidKeySize
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	// 加密并附加认证标签
	ciphertext := gcm.Seal(nonce, nonce, plaintext, nil)
	return ciphertext, nil
}

// DecryptCredentials 使用AES-256-GCM解密凭证
func DecryptCredentials(ciphertext []byte, key []byte) ([]byte, error) {
	if len(key) != 32 {
		return nil, ErrInvalidKeySize
	}

	if len(ciphertext) < 12 { // GCM nonce size
		return nil, ErrInvalidCiphertext
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonceSize := gcm.NonceSize()
	nonce, encryptedData := ciphertext[:nonceSize], ciphertext[nonceSize:]

	plaintext, err := gcm.Open(nil, nonce, encryptedData, nil)
	if err != nil {
		return nil, err
	}

	return plaintext, nil
}

// EncodeBase64 将字节编码为Base64字符串
func EncodeBase64(data []byte) string {
	return base64.StdEncoding.EncodeToString(data)
}

// DecodeBase64 将Base64字符串解码为字节
func DecodeBase64(encoded string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(encoded)
}