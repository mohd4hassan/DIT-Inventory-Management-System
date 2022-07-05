package context

import (
	"context"

	"IMS/models"
)

type privateKey string

const (
	userKey privateKey = "user"
)


func WithUser(ctx context.Context, user *models.User) context.Context {
	return context.WithValue(ctx, "user", user)
}

func User(ctx context.Context) *models.User {
	if temp := ctx.Value("user"); temp != nil {
		if user, ok := temp.(*models.User); ok {
			return user
		}
	}
	return nil
}
