package authService

import (
	"context"
	"github.com/mr-linch/go-tg"
)

type Service interface {
	GetRoleById(ctx context.Context, UserId tg.UserID) (int, error)
}
