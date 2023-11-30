package challange_test

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/zegmic/powserver/pkg/challange"
	"strings"
	"testing"
	"time"
)

func TestChallengeCreate(t *testing.T) {
	h := challange.Hash{}
	challenge := h.Create()
	if len(strings.Split(challenge, ".")) != 3 {
		t.Fail()
	}
}

func TestVerifyMissingParts(t *testing.T) {
	h := challange.Hash{}
	valid := h.Verify("rr.aa", 0)
	if valid {
		t.Fail()
	}
}

func TestVerifyInvalidSignature(t *testing.T) {
	h := challange.Hash{}
	valid := h.Verify("rr.aa.ee", 0)
	if valid {
		t.Fail()
	}
}

func TestVerifyTimestampDrifted(t *testing.T) {
	h := challange.Hash{}
	token := h.Create()
	h.SetNow(time.Now().UnixMilli() - 1000*30)
	valid := h.Verify(token, 0)
	if valid {
		t.Fail()
	}
}

func TestVerifyValidSolution(t *testing.T) {
	h := challange.Hash{}
	token := h.Create()
	nonce := solve(token)
	valid := h.Verify(token, nonce)
	if !valid {
		t.Fail()
	}
}

func solve(challenge string) int {
	for i := 0; i < 1000000000; i++ {
		sol := fmt.Sprintf("%d.%s", i, challenge)
		dig := sha256.New()
		dig.Write([]byte(sol))
		sum := dig.Sum(nil)

		if solved(hex.EncodeToString(sum[:])) {
			return i
		}
	}

	return -1
}

func solved(sum string) bool {
	return strings.HasPrefix(sum, "00000")
}
