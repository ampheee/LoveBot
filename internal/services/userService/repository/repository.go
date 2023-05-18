package repository

import (
	"github.com/rs/zerolog"
	"tacy/internal/storage"
	"tacy/pkg/botlogger"
)

type Service struct {
	Storage storage.Storage
	Logger  zerolog.Logger
}

func (s Service) InputPhotoFromUser() {
	//TODO implement me
	panic("implement me")
}

func (s Service) SendComplimentByRandom() {
	//TODO implement me
	panic("implement me")
}

func (s Service) InputThoughtsFromUser() {
	//TODO implement me
	panic("implement me")
}

func (s Service) SendPhotoByRandom() {
	//TODO implement me
	panic("implement me")
}

func (s Service) SendComplimentAndPhotoRandom() {
	//TODO implement me
	panic("implement me")
}

func InitNewService(s storage.Storage) Service.UserService {
	return Service{
		Storage: s,
		Logger:  botlogger.GetLogger(),
	}
}
