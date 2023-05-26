package userService

import (
	"context"
)

type Service interface {
	InputPhotoFromAdmin(ctx context.Context, photo []byte) error
	InputComplimentFromAdmin(ctx context.Context)
	InputThoughtsFromUser(ctx context.Context)
	OutputComplimentAndPhotoRandom(ctx context.Context)
	OutputAllPhotos(ctx context.Context) ([][]byte, error)
	OutputAllCompliments(ctx context.Context)
}
