// This file contains the interfaces for the repository layer.
// The repository layer is responsible for interacting with the database.
// For testing purpose we will generate mock implementations of these
// interfaces using mockgen. See the Makefile for more information.
package repository

import "context"

type RepositoryInterface interface {
	InsertUser(ctx context.Context, input User) (output InsertUserOutput, err error)
	GetUserByPhoneNumber(ctx context.Context, input GetUserByPhoneNumberInput) (output User, err error)
	GetUserByID(ctx context.Context, input GetUserByIDInput) (output UserInfo, err error)
	GetUserByFullName(ctx context.Context, input GetUserByFullNameInput) (output UserInfo, err error)
	UpdateUserByID(ctx context.Context, input User) (err error)
}
