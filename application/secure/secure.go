package secure

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"
)

func MD5Hash(text string) string {
	hasher := md5.New()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}

func computeHmac256(message string, secret string) string {
	key := []byte(secret)
	h := hmac.New(sha256.New, key)
	h.Write([]byte(message))
	return hex.EncodeToString(h.Sum(nil))
}

func GenerateToken(message, key string) string {
	accessToken := fmt.Sprintf("%d_%s", time.Now().Unix(), message)
	return accessToken + "_" + computeHmac256(accessToken, key)
}
