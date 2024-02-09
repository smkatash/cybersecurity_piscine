package main

import (
	"fmt"
	"os"
	"strings"

	"ft_otp/logger"
	"github.com/skip2/go-qrcode"
)

func ReadKeyFromFile(file string) []byte {
	key, err := os.ReadFile(file)   
	if err != nil {
		logger.LogError(err)
	}
	return key
}

func GenerateQR(seed []byte) bool {
	result := "otpauth://totp/ktashbae:ft_otp_user?secret=" + string(seed) + "&issuer=ktashbae"
	qrCode, err := qrcode.New(result, qrcode.Medium)
	if err != nil {
		logger.LogError(err)
	}

	file, err := os.Create("qrcode.png")
	if err != nil {
		logger.LogError(err)
	}
	defer file.Close()

	if err := qrCode.Write(256, file); err != nil {
		logger.LogError(err)
	}
	return true
}

func ft_otp(args []string) {
	if args[1] == "-g" && strings.HasSuffix(args[2], ".hex") {
		if EncodeKey(ReadKeyFromFile(args[2])) {
			fmt.Println("Key was successfully saved in ft_otp.key.")
		}
	} else if args[1] == "-g" {
		logger.LogErrorStr("key must be 64 hexadecimal characters.")
		os.Exit(1)
	} else if args[1] == "-k" && strings.HasSuffix(args[2], ".key") {
		if key := DecodeKey(ReadKeyFromFile(args[2])); key != nil {
			fmt.Println(NewHOTP(key))
		} else {
			logger.LogErrorStr("key was compromised.")
		}
	} else {
		logger.LogUsage()
	}
}

func main() {
	if len(os.Args) != 3 {
		logger.LogUsage()
	}
	ft_otp(os.Args)
}
