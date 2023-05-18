package storage

import (
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/mr-linch/go-tg"
)

type Storage interface {
	GetRoleById(ctx context.Context, id tg.UserID) (int, error)
	GetCompliment(ctx context.Context) (string, error)
	GetPhoto(ctx context.Context) ([]byte, error)
	InsertPhoto(ctx context.Context, buffer []byte) error
	InsertThoughts(ctx context.Context, thoughts string) error
	InsertCompliment(ctx context.Context, compliment string) error
	InsertNewUser(ctx context.Context, id tg.UserID, conn *pgxpool.Conn) error
}
