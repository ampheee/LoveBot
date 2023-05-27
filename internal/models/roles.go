package models

const (
	AdminRole int = iota + 1
	UserRole
	UnknownRole
)

type UserRoleS struct {
	UserId   int64
	UserRole int8
}
