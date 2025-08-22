package auth

import (
    "crypto/aes"
    "crypto/cipher"
    "encoding/base64"
)

// Encrypt encrypts sensitive data using AES
func Encrypt(plaintext string, key []byte) (string, error) {
    block, err := aes.NewCipher(key)
    if err != nil { return "", err }

    ciphertext := make([]byte, len(plaintext))
    stream := cipher.NewCFBEncrypter(block, key[:block.BlockSize()])
    stream.XORKeyStream(ciphertext, []byte(plaintext))

    return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// Decrypt decrypts AES encrypted data
func Decrypt(ciphertext string, key []byte) (string, error) {
    data, _ := base64.StdEncoding.DecodeString(ciphertext)
    block, err := aes.NewCipher(key)
    if err != nil { return "", err }

    plaintext := make([]byte, len(data))
    stream := cipher.NewCFBDecrypter(block, key[:block.BlockSize()])
    stream.XORKeyStream(plaintext, data)

    return string(plaintext), nil
}
