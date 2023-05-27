package storage

import (
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/mr-linch/go-tg"
)

type Storage interface {
	GetRoleById(ctx context.Context, id tg.UserID) (int, error)
	GetComplimentById(ctx context.Context) (string, error)
	GetPhotoById(ctx context.Context) ([]byte, error)
	GetAllPhotos(ctx context.Context) ([][]byte, error)
	GetAllCompliments(ctx context.Context) ([]Compliment, error)
	InsertPhoto(ctx context.Context, buffer []byte) error
	InsertThoughts(ctx context.Context, thoughts string) error
	InsertCompliment(ctx context.Context, compliment string) error
	InsertNewUser(ctx context.Context, id tg.UserID, conn *pgxpool.Conn) error
}

type Compliment struct {
	Id         int
	Compliment string
}
