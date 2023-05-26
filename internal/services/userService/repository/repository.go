package repository

import (
	"context"
	"github.com/rs/zerolog"
	"tacy/internal/services/userService"
	"tacy/internal/storage"
	"tacy/pkg/botlogger"
)

type Service struct {
	Storage storage.Storage
	Logger  zerolog.Logger
}

func (s Service) InputPhotoFromAdmin(ctx context.Context, photo []byte) error {
	err := s.Storage.InsertPhoto(ctx, photo)
	if err != nil {
		s.Logger.Warn().Err(err).Msg("Unable to input photo")
		return err
	}
	return nil
}

func (s Service) InputComplimentFromAdmin(ctx context.Context) {
	//TODO implement me
	panic("implement me")
}

func (s Service) InputThoughtsFromUser(ctx context.Context) {
	//TODO implement me
	panic("implement me")
}

func (s Service) OutputComplimentAndPhotoRandom(ctx context.Context) {
	//TODO implement me
	panic("implement me")
}

func (s Service) OutputAllPhotos(ctx context.Context) ([][]byte, error) {
	photos, err := s.Storage.GetAllPhotos(ctx)
	if err != nil {
		s.Logger.Warn().Err(err)
		return nil, err
	}
	return photos, nil
}

func (s Service) OutputAllCompliments(ctx context.Context) {
	//TODO implement me
	panic("implement me")
}

func InitNewUserService(s storage.Storage) userService.Service {
	return Service{
		Storage: s,
		Logger:  botlogger.GetLogger(),
	}
}
