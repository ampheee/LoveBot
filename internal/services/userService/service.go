package userService

import (
	"context"
)

type Service interface {
	InputPhotoFromAdmin(ctx context.Context, photo []byte) error
	InputComplimentFromAdmin(ctx context.Context, compliment string) error
	InputThoughtsFromUser(ctx context.Context)
	OutputComplimentAndPhotoByRandom(ctx context.Context)
	OutputAllPhotos(ctx context.Context) ([][]byte, error)
	OutputAllCompliments(ctx context.Context) ([]string, error)
}
