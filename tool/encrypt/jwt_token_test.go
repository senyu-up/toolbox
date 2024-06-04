package encrypt

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCreateToken(t *testing.T) {
	token, err := CreateToken("1", 10, "2")
	if err != nil {
		assert.Fail(t, err.Error())
		return
	}
	fmt.Println(token)
}

func TestParseToken(t *testing.T) {
	token, err := CreateToken("1ASAD", 10, "3")
	if err != nil {
		assert.Fail(t, err.Error())
		return
	}
	fmt.Println(token)
	claims, err := ParseToken(token, "3")
	if err != nil {
		assert.Fail(t, err.Error())
		return
	}
	fmt.Println(claims.ExpiresAt, claims.Data)
}
