package main

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/binary"
	"ft_otp/logger"
	"time"
)

func NewHOTP(key []byte) uint32 {
	counter := uint64(time.Now().Unix() / 30)
	counterInBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(counterInBytes, counter)
	
	hmac_sha1 := hmac.New(sha1.New, key)
	_, err := hmac_sha1.Write(counterInBytes)
	if err != nil {
		logger.LogError(err)
	}
	hmac_result := hmac_sha1.Sum(nil)
	
	offset :=  int(hmac_result[19] & 0xf)
	binCode := (int(hmac_result[offset]) & 0x7f) << 24 | (int(hmac_result[offset + 1]) & 0xff) << 16 | (int(hmac_result[offset + 2]) & 0xff) <<  8 | (int(hmac_result[offset +3 ]) & 0xff)
	
	return uint32(binCode) % 1000000
}
