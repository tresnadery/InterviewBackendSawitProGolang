// This file contains types that are used in the repository layer.
package repository

import (
	"github.com/google/uuid"
	"time"
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
	PhoneNumber       string `validate:"required,min=10,max=13,indonesian_phone_number"`
	FullName          string `validate:"required,min=3,max=60"`
	LastLogin         *time.Time
	SuccessfullyLogin int
}

type UserSecret struct {
	Password     string `validate:"required,min=6,max=64,valid_password"`
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

type UpdateUserInput struct {
	ID          uuid.UUID
	PhoneNumber string `validate:"required,min=10,max=13,indonesian_phone_number"`
	FullName    string `validate:"required,min=3,max=60"`
}

type UpdateLastLoginAndSuccessfullyLoginInput struct {
	ID        uuid.UUID
	LastLogin *time.Time
}
