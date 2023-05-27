package repository

import (
	"context"
	"errors"
	"github.com/rs/zerolog"
	"strconv"
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

func (s Service) InputComplimentFromAdmin(ctx context.Context, compliment string) error {
	err := s.Storage.InsertCompliment(ctx, compliment)
	if err != nil {
		s.Logger.Warn().Err(err)
		return err
	}
	return nil
}

func (s Service) InputThoughtsFromUser(ctx context.Context, thoughts string) error {
	err := s.Storage.InsertThoughts(ctx, thoughts)
	if err != nil {
		return err
	}
	return nil
}

func (s Service) OutputComplimentAndPhotoByRandom(ctx context.Context) (photo []byte, compliment string, err error) {
	photo, errPhoto := s.Storage.GetPhoto(ctx)
	compliment, errCompliment := s.Storage.GetCompliment(ctx)
	err = errors.Join(errPhoto, errCompliment)
	if err != nil {
		return nil, "", err
	}
	return photo, compliment, nil
}

func (s Service) OutputAllPhotos(ctx context.Context) ([][]byte, error) {
	photos, err := s.Storage.GetAllPhotos(ctx)
	if err != nil {
		s.Logger.Warn().Err(err)
		return nil, err
	}
	return photos, nil
}

func (s Service) OutputAllCompliments(ctx context.Context) ([]string, error) {
	compliments, err := s.Storage.GetAllCompliments(ctx)
	if err != nil {
		s.Logger.Warn().Err(err)
		return nil, err
	}
	var resCompliments []string
	for _, compliment := range compliments {
		resCompliments = append(resCompliments, strconv.Itoa(compliment.Id)+". "+compliment.Compliment)
	}
	return resCompliments, nil
}

func InitNewUserService(s storage.Storage) userService.Service {
	return Service{
		Storage: s,
		Logger:  botlogger.GetLogger(),
	}
}
