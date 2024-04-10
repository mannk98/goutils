package utils

import (
	"crypto/rand"
	"crypto/sha1"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/rs/xid"
)

/* return Prestring + _uuidv4 */
func GenerateIdv4(preString string) (string, error) {
	newUUID, err := uuid.NewRandom()
	if err != nil {
		return "", err
	}
	uuidStr := newUUID.String()
	return preString + uuidStr, err
}

/*Get n random bytes*/
func GenerateRandomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	// Note that err == nil only if we read len(b) bytes.
	if err != nil {
		return nil, err
	}

	return b, nil
}

func GenerateID() string {
	return xid.New().String()
}

/* Get random string with n character */
func GenerateRandomString(s int) (string, error) {
	b, err := GenerateRandomBytes(s)
	return base64.URLEncoding.EncodeToString(b), err
}

/* creat sha1 of data*/
func GenerateSha1Bytes(data []byte) string {
	h := sha1.New()
	h.Write([]byte(data))
	return hex.EncodeToString(h.Sum(nil))
}

func GenerateSha1String(data string) string {
	h := sha1.New()
	h.Write([]byte(data))
	return hex.EncodeToString(h.Sum(nil))
}

func GenerateTokenSha1(key int) string {
	if key == 0 {
		key = 1998
	}
	nowtimestam := time.Now().Unix() + int64(key)
	fmt.Print(nowtimestam)
	return GenerateSha1String(Int64ToString(nowtimestam))
}

/* compare token (15s limit between create and check) (ex: http post and response confrim)*/
func TokenSha1IsMatch(key int, token string) bool {
	if key == 0 {
		key = 1998
	}
	nowtimestam := time.Now().Unix() + int64(key)
	for i := -15; i < 15; i++ {
		if token == GenerateSha1String(strconv.FormatInt((nowtimestam+int64(i)), 10)) {
			fmt.Println(strconv.FormatInt((nowtimestam + int64(i)), 10))
			fmt.Println(token)
			return true
		}
	}
	return false
}
