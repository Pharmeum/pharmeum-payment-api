package handlers

import (
	"fmt"

	"github.com/pkg/errors"
)

var (
	ErrInvalidToken      = errors.New("invalid token")
	ErrInvalidWalletKind = errors.New("invalid wallet kind, supported patient and doctor wallets")
)

func ErrResponse(code int, err error) []byte {
	return []byte(fmt.Sprintf(`{"code": %d, "error": "%s"}`, code, err.Error()))
}
