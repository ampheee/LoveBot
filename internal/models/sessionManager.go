package models

type SessionStep int8

const (
	SessionStepInit = iota

	// admin steps
	SessionStepAdminMenuHandler

	SessionStepEnterAdminMenuHandler
	SessionStepInsertNewPhotoMenuHandler
	SessionStepInsertNewComplimentMenuHandler
	SessionStepGetAllPhotosMenuHandler
	SessionStepGetAllComplimentsHandler

	// user steps
	SessionStepUserMenuHandler
	SessionStepInsertSomeThoughts
	SessionStepGetCompliment
)

type Session struct {
	Step SessionStep
}
