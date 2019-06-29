package utils

import (
	"context"

	"github.com/go-chi/jwtauth"
)

const userIDKey = "id"

//UserID provided from request context, takes ID from JWT claims
func UserID(ctx context.Context) uint64 {
	_, claims, _ := jwtauth.FromContext(ctx)
	userID, ok := claims[userIDKey]
	if !ok {
		return 0
	}

	return userID.(uint64)
}
