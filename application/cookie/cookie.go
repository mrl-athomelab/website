package cookie

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type Cookie struct {
	CreatedAt time.Time
	TTL       time.Duration
	Payload   interface{}
}

func (c *Cookie) Set(ctx *gin.Context, name, key string) error {
	bytes, err := json.Marshal(c)
	if err != nil {
		return err
	}
	encoded := base64.StdEncoding.EncodeToString(bytes)
	validator := computeHmac256(encoded, key)
	cookie := &http.Cookie{}
	cookie.Name = name
	cookie.MaxAge = int(c.TTL.Seconds())
	cookie.Value = fmt.Sprintf("%s.%s", validator, encoded)
	http.SetCookie(ctx.Writer, cookie)
	return nil
}

func computeHmac256(message string, secret string) string {
	key := []byte(secret)
	h := hmac.New(sha256.New, key)
	h.Write([]byte(message))
	return hex.EncodeToString(h.Sum(nil))
}

func Get(ctx *gin.Context, name, key string) (*Cookie, error) {
	cookie, err := ctx.Cookie(name)
	if err != nil {
		return nil, err
	}
	parts := strings.Split(cookie, ".")
	if len(parts) != 2 {
		return nil, fmt.Errorf("Invalid parts length, got %d want 2", len(parts))
	}
	if computeHmac256(parts[1], key) != parts[0] {
		return nil, fmt.Errorf("Invalid checksum, got %s", parts[0])
	}
	bytes, err := base64.StdEncoding.DecodeString(parts[1])
	if err != nil {
		return nil, err
	}
	result := &Cookie{}
	err = json.Unmarshal(bytes, result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (c *Cookie) IsExpired() bool {
	return c.CreatedAt.Add(c.TTL).Before(time.Now())
}
