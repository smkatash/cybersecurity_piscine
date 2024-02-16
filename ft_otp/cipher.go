package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"ft_otp/logger"
	"io"
	"os"
	"regexp"
	"encoding/hex"
)

var secretFile = "secret.key"

func isValidKey(input []byte) bool {
	if len(input) != 64 {
		return false
	}
	str := string(input)
	hexPattern := "^[0-9a-fA-F]+$"
	regex := regexp.MustCompile(hexPattern)
	return regex.MatchString(str)
}

func CreateSecret() []byte {
	secret := make([]byte, 32)
	if _, err := rand.Read(secret); err != nil {
		logger.LogError(err)
	}
	secretKey := hex.EncodeToString(secret)
	if err := os.WriteFile(secretFile, []byte(secretKey), 0600); err != nil {
		logger.LogError(err)
	}
	return []byte(secretKey)
}

func EncodeKey(inputKey []byte) bool {
	if !isValidKey(inputKey) {
		logger.LogErrorStr("key must be 64 hexadecimal characters.")
		os.Exit(1)
	}

	secretKey := CreateSecret()
	if secretKey == nil {
		return false
	}

	ciph, err := aes.NewCipher(secretKey)
	if err != nil {
		logger.LogError(err)
	}

	gmc, err := cipher.NewGCM(ciph)
	if err != nil {
		logger.LogError(err)
	}
	
	nonce := make([]byte, gmc.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		logger.LogError(err)
	}
	
	if err = os.WriteFile("ft_otp.key", gmc.Seal(nonce, nonce, inputKey, nil), 0777); err != nil {
		logger.LogError(err)
	}
	if GenerateQR(secretKey) {
		fmt.Println("Seed QR code saved as qrcode.png")
	}
	return true
}

func DecodeKey(inputKey []byte) []byte {
	secretKey := ReadKeyFromFile(secretFile)

	ciph, err := aes.NewCipher(secretKey)
	if err != nil {
		logger.LogError(err)
	}

	gmc, err := cipher.NewGCM(ciph)
	if err != nil {
		logger.LogError(err)
	}

	nonceSize := gmc.NonceSize()
	if nonceSize > len(inputKey) {
		return nil
	}

	nonce, inputKey := inputKey[:nonceSize], inputKey[nonceSize:]
	key, err := gmc.Open(nil, nonce, inputKey, nil)
	if err != nil {
		logger.LogError(err)
	}
	return key
}