package userService

import "context"

type Service interface {
	InputPhotoFromAdmin(ctx context.Context)
	InputComplimentFromAdmin(ctx context.Context)
	InputThoughtsFromUser(ctx context.Context)
	OutputComplimentAndPhotoRandom(ctx context.Context)
	OutputAllPhotos(ctx context.Context)
	OutputAllCompliments(ctx context.Context)
}
