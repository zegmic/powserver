package challange_test

import (
	"github.com/zegmic/powserver/pkg/challange"
	"strings"
	"testing"
)

func TestChallengeCreate(t *testing.T) {
	h := challange.Hash{}
	challenge := h.Create()
	if len(strings.Split(challenge, ".")) != 3 {
		t.Fail()
	}
}
