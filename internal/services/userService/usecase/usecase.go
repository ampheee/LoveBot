package usecase

import (
	"context"
	"fmt"
	"github.com/rs/zerolog"
	"tacy/internal/services/userService"
	"tacy/pkg/botlogger"
)

type Usecase struct {
	Repo   userService.USRepository
	logger zerolog.Logger
}

func (u *Usecase) OutputhAllThoughts(ctx context.Context) ([]string, error) {
	//TODO implement me
	panic("implement me")
}

func (u *Usecase) InsertPhoto(ctx context.Context, msg []byte, fromUser bool) error {
	return u.Repo.InsertPhoto(ctx, msg, fromUser)
}

func (u *Usecase) InsertComplimentFromAdmin(ctx context.Context, compliment string) error {
	return u.Repo.InsertComplimentFromAdmin(ctx, compliment)
}

func (u *Usecase) InputThoughtsFromUser(ctx context.Context, thougth string) error {
	return u.Repo.InsertThoughtsFromUser(ctx, thougth)
}

func (u *Usecase) OutputComplimentAndPhotoByRandom(ctx context.Context) ([]byte, string, error) {
	reservedPhoto, errPhoto := u.Repo.GetPhotoByRandom(ctx)
	reservedCompliment, errCompliment := u.Repo.GetComplimentByRandom(ctx)
	if errCompliment != nil {
		u.logger.Warn().Err(errCompliment).Msg("something went wrong while get compliment")
		return nil, "", errCompliment
	} else if errPhoto != nil {
		u.logger.Warn().Err(errPhoto).Msg("something went wrong while get photo")
		return nil, "", errPhoto
	}
	return reservedPhoto.Photo, reservedCompliment, nil
}

func (u *Usecase) OutputAllCompliments(ctx context.Context) ([]string, error) {
	getted, err := u.Repo.GetAllCompliments(ctx)
	if err != nil {
		return nil, err
	}
	var compliments []string
	for _, val := range getted {
		compliments = append(compliments, fmt.Sprintf("%d. %s", val.Id, val.Compliment))
	}
	return compliments, nil
}

func (u *Usecase) OutputAllPhotos(ctx context.Context) ([][]byte, error) {
	photos, err := u.Repo.GetAllPhotos(ctx)
	if err != nil {
		u.logger.Warn().Err(err)
		return nil, err
	}
	var bytePhotos [][]byte
	for _, photo := range photos {
		bytePhotos = append(bytePhotos, photo.Photo)
	}
	return bytePhotos, nil
}

func NewUSUsecase(repo userService.USRepository) userService.USUsecase {
	return &Usecase{
		Repo:   repo,
		logger: botlogger.GetLogger(),
	}
}
