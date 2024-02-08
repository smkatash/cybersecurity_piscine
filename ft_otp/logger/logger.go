package logger

import (
	"log"
)

func LogUsage() {
	log.Fatal("usage: \n\t\t\t ./ft_otp -g [key.hex] \n\t\t\t ./ft_otp -k [ft_otp.key] \n\t\t\t ./ft_otp -qr")
}


func LogError(err error) {
	log.Fatalf("ft_otp: %s", err)
}

func LogErrorStr(err string) {
	log.Fatalf("ft_otp: error: %s", err)
}
