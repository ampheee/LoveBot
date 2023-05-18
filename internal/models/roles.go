package models

const (
	AdminRole int = iota + 1
	UserRole
)

type UserRoleS struct {
	UserId   int64
	UserRole int8
}
