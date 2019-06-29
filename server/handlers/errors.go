package handlers

import (
	"fmt"

	"github.com/pkg/errors"
)

var (
	ErrInvalidToken = errors.New("invalid token")
)

func ErrResponse(code int, err error) []byte {
	return []byte(fmt.Sprintf(`{"code": %d, "error": "%s"}`, code, err.Error()))
}
