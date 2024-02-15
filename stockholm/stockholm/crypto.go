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

var secretKey string = "secret.key"

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
	fmt.Println("Getting " + string(key))
	if err != nil {
		log.Fatal(err)
    }
	if len(key) < 16 {
		log.Fatal("Key must be at least 16 characters long.")
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
		log.Fatal("Error encrypting file:", err)
	}
	if (!silentMode) {
		fmt.Println("encrypted: " + srcFile)
	}
	if destFile != srcFile {
		err = os.Remove(srcFile)
		if err != nil {
			log.Fatal("Error removing file: ", err)
		}
	}
}

func DecryptFile(srcFile string) {
	destFile := srcFile
	if !strings.HasSuffix(srcFile, ".ft") {
		destFile = strings.TrimSuffix(srcFile, ".ft")
		if strings.HasSuffix(destFile, ".ft") {
			return
		}
	}
	fmt.Println(destFile)

	cmd := exec.Command("openssl", "enc", "-aes-256-cbc", "-d", "-in", srcFile, "-out",  destFile, "-k", GetSecretKey())
	err := cmd.Run()
	if err != nil {
		log.Fatal("Error encrypting file:", err)
	}
	if (!silentMode) {
		fmt.Println("decrypted: " + destFile)
	}
	if destFile != srcFile {
		err = os.Remove(srcFile)
		if err != nil {
			log.Fatal("Error removing file: ", err)
		}
	}
}
