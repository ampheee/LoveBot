package userService

import (
	"context"
)

type USRepository interface {
	GetComplimentByRandom(ctx context.Context) (string, error)
	GetPhotoByRandom(ctx context.Context) (Photo, error)
	GetAllPhotos(ctx context.Context) ([]Photo, error)
	GetAllCompliments(ctx context.Context) ([]Compliment, error)
	InsertPhotoFromAdmin(ctx context.Context, buffer []byte) error
	InsertThoughtsFromUser(ctx context.Context, thoughts string) error
	InsertComplimentFromAdmin(ctx context.Context, compliment string) error
	UpdateComplimentById(ctx context.Context, id int) error
	UpdateImageById(ctx context.Context, id int) error
	RefreshImages(ctx context.Context) error
	RefreshCompliments(ctx context.Context) error
}

type USUsecase interface {
	InsertPhotoFromAdmin(ctx context.Context, photo []byte) error
	InsertComplimentFromAdmin(ctx context.Context, compliment string) error
	InputThoughtsFromUser(ctx context.Context, thoughts string) error
	OutputComplimentAndPhotoByRandom(ctx context.Context) ([]byte, string, error)
	OutputAllPhotos(ctx context.Context) ([][]byte, error)
	OutputAllCompliments(ctx context.Context) ([]string, error)
}

type Compliment struct {
	Id         int
	Compliment string
}

type Photo struct {
	Id    int
	Photo []byte
}
