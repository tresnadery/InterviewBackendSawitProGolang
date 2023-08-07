// This file contains types that are used in the repository layer.
package repository

import (
	"github.com/google/uuid"
)

type GetTestByIdInput struct {
	Id string
}

type GetTestByIdOutput struct {
	Name string
}

type User struct {
	ID uuid.UUID
	UserInfo
	UserSecret
}

type UserInfo struct {
	PhoneNumber string
	FullName    string
}

type UserSecret struct {
	Password     string
	PasswordSalt string
}

type InsertUserOutput struct {
	ID uuid.UUID
}

type GetUserByPhoneNumberInput struct {
	PhoneNumber string
}

type GetUserByIDInput struct {
	ID uuid.UUID
}

type GetUserByFullNameInput struct {
	FullName string
}

type UpdateUserByID struct {
	ID          uuid.UUID
	PhoneNumber string
	FullName    string
}
