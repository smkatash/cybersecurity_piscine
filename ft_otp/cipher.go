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

func getSecret() []byte {
	secret := make([]byte, 32)
	if _, err := rand.Read(secret); err != nil {
		logger.LogError(err)
	}
	if err := os.WriteFile(secretFile, secret, 0600); err != nil {
		logger.LogError(err)
	}
	return secret
}

func EncodeKey(inputKey []byte) bool {
	if !isValidKey(inputKey) {
        logger.LogErrorStr("Key must be 64 hexadecimal characters.")
        return false
    }

	keyInBytes, err := hex.DecodeString(string(inputKey))
	if err != nil {
        logger.LogError(err)
        return false
    }

    secretKey := getSecret()
    if secretKey == nil {
        logger.LogErrorStr("Failed to generate secret key.")
        return false
    }

    aesCipher, err := aes.NewCipher(secretKey)
    if err != nil {
        logger.LogError(err)
        return false
    }

    gcm, err := cipher.NewGCM(aesCipher)
    if err != nil {
        logger.LogError(err)
        return false
    }

    nonce := make([]byte, gcm.NonceSize())
    if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
        logger.LogError(err)
        return false
    }

    ciphertext := gcm.Seal(nil, nonce, keyInBytes, nil)

    if err := os.WriteFile("ft_otp.key", append(nonce, ciphertext...), 0644); err != nil {
        logger.LogError(err)
        return false
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