package main

import  (
	"crypto/rand"
	"log"
	"os"
	"encoding/hex"
	"os/exec"
	"strings"
	"fmt"
)

func GenerateSecretKey() {
    keyBytes := make([]byte, 32)
    _, err := rand.Read(keyBytes)
    if err != nil {
		log.Fatal(err)
    }
	key := hex.EncodeToString(keyBytes)
	err = os.WriteFile(secretKey, []byte(key), 0600)
    if err != nil {
		log.Fatal(err)
    }
}

func GetSecretKey() string {
	key, err := os.ReadFile(secretKey)
	if err != nil {
		log.Fatal(err)
    }
	if len(key) < 16 {
		log.Fatal("stockholm: key must be at least 16 characters long.")
	}
	return string(key)
}

func EncryptFile(srcFile string) {
	destFile := srcFile
	if !strings.HasSuffix(srcFile, ".ft") {
		destFile = srcFile + ".ft"
	}

	cmd := exec.Command("openssl", "enc", "-aes-256-cbc", "-pbkdf2", "-salt", "-in", srcFile, "-out",  destFile, "-k", GetSecretKey())
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
	if (!silentMode) {
		fmt.Println("encrypted: " + srcFile)
	}
	if destFile != srcFile {
		err = os.Remove(srcFile)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func DecryptFile(srcFile string) {
	destFile := srcFile
	if strings.HasSuffix(srcFile, ".ft") {
		destFile = strings.TrimSuffix(destFile, ".ft")
		if strings.HasSuffix(destFile, ".ft") {
			return
		}
	}
	cmd := exec.Command("openssl", "enc", "-aes-256-cbc", "-d", "-pbkdf2", "-in", srcFile, "-out",  destFile, "-k", GetSecretKey())
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
	if (!silentMode) {
		fmt.Println("decrypted: " + destFile)
	}
	if destFile != srcFile {
		err = os.Remove(srcFile)
		if err != nil {
			log.Fatal(err)
		}
	}
}
