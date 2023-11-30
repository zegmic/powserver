package challange

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"math/big"
	"strconv"
	"strings"
	"time"
)

type Hash struct {
}

func (e *Hash) Create() string {
	now := time.Now().UnixMilli()
	nowStr := strconv.FormatInt(now+1000*30, 10)
	nowBase := base64.StdEncoding.EncodeToString([]byte(nowStr))

	n, _ := rand.Int(rand.Reader, big.NewInt(10000000000))
	numBase := base64.StdEncoding.EncodeToString([]byte(n.String()))

	challengeToSign := fmt.Sprintf("%s.%s", numBase, nowBase)
	signature := sign(challengeToSign)
	signBase := base64.StdEncoding.EncodeToString([]byte(signature))

	return fmt.Sprintf("%s.%s.%s", numBase, nowBase, signBase)
}

func (e *Hash) Verify(challenge string, solution int) bool {
	parts := strings.Split(challenge, ".")
	if len(parts) != 3 {
		return false
	}

	if !validSignature(challenge) {
		return false
	}

	if !validTs(challenge) {
		return false
	}

	puzzle := fmt.Sprintf("%d.%s", solution, challenge)
	dig := sha256.New()
	dig.Write([]byte(puzzle))
	sum := dig.Sum(nil)

	return verify(sum)
}

func verify(sum []byte) bool {
	return strings.HasPrefix(hex.EncodeToString(sum[:]), "00000")
}

func validTs(challenge string) bool {
	parts := strings.Split(challenge, ".")
	ts, err := base64.StdEncoding.DecodeString(parts[1])
	if err != nil {
		return false
	}

	val, err := strconv.ParseInt(string(ts), 10, 64)
	if err != nil {
		return false
	}

	return val-time.Now().UnixMilli() > 0
}

func validSignature(challenge string) bool {
	parts := strings.Split(challenge, ".")
	origSign, err := base64.StdEncoding.DecodeString(parts[2])
	if err != nil {
		return false
	}
	signature := sign(fmt.Sprintf("%s.%s", parts[0], parts[1]))

	return signature == string(origSign)
}

func sign(val string) string {
	secretHash := md5.New()
	// in a real system it should be e.g. env var
	secretHash.Write([]byte("s3cr37"))
	key := secretHash.Sum(nil)

	sig := hmac.New(sha256.New, key)
	sig.Write([]byte(val))
	hexStr := hex.EncodeToString(sig.Sum(nil))

	return hexStr
}