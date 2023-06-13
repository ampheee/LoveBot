package authService

import (
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/mr-linch/go-tg"
)

type Service interface {
	GetRoleById(ctx context.Context, UserId tg.UserID) (int, error)
	InsertNewUser(ctx context.Context, id tg.UserID, conn *pgxpool.Conn) error
}
