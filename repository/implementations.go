package repository

import (
	"context"
	"database/sql"
	"fmt"
)

func (r *Repository) InsertUser(ctx context.Context, input User) (output InsertUserOutput, err error) {
	stmt, err := r.Db.PrepareContext(ctx, "INSERT INTO users(phone_number, full_name, password, password_salt) VALUES($1,$2,$3,$4) RETURNING id")
	if err != nil {
		return
	}

	err = stmt.QueryRowContext(ctx, input.PhoneNumber, input.FullName, input.Password, input.PasswordSalt).Scan(&output.ID)
	if err != nil {
		return
	}

	return
}

func (r *Repository) GetUserByPhoneNumber(ctx context.Context, input GetUserByPhoneNumberInput) (output User, err error) {
	stmt, err := r.Db.PrepareContext(ctx, "SELECT id, phone_number, full_name, password, password_salt FROM users WHERE phone_number = $1")
	if err != nil {
		return
	}

	err = stmt.QueryRowContext(ctx, input.PhoneNumber).Scan(
		&output.ID,
		&output.PhoneNumber,
		&output.FullName,
		&output.Password,
		&output.PasswordSalt,
	)

	if err != nil {
		fmt.Println(sql.ErrNoRows)
		if err != sql.ErrNoRows {
			return
		}
	}

	return
}

func (r *Repository) GetUserByID(ctx context.Context, input GetUserByIDInput) (output UserInfo, err error) {
	fmt.Println(input.ID.String())
	stmt, err := r.Db.PrepareContext(ctx, "SELECT phone_number, full_name FROM users WHERE id = $1")
	if err != nil {
		return
	}

	err = stmt.QueryRowContext(ctx, input.ID.String()).Scan(
		&output.PhoneNumber,
		&output.FullName,
	)

	if err != nil {
		return
	}

	return
}

func (r *Repository) GetUserByFullName(ctx context.Context, input GetUserByFullNameInput) (output UserInfo, err error) {
	stmt, err := r.Db.PrepareContext(ctx, "SELECT phone_number, full_name FROM users WHERE LOWER(full_name) = LOWER($1)")
	if err != nil {
		return
	}

	err = stmt.QueryRowContext(ctx, input.FullName).Scan(
		&output.PhoneNumber,
		&output.FullName,
	)

	if err != nil {
		return
	}

	return
}

func (r *Repository) UpdateUser(ctx context.Context, input UpdateUserInput) (err error) {
	stmt, err := r.Db.PrepareContext(ctx, "UPDATE users SET phone_number = $1, full_name = $2 WHERE id = $3")
	if err != nil {
		return
	}

	_, err = stmt.ExecContext(ctx, input.PhoneNumber, input.FullName, input.ID)
	if err != nil {
		return
	}
	return
}

func (r *Repository) UpdateLastLoginAndSuccessfullyLogin(ctx context.Context, input UpdateLastLoginAndSuccessfullyLoginInput) (err error) {
	stmt, err := r.Db.PrepareContext(ctx, "UPDATE users SET successfully_login = (successfully_login + 1), last_login = $1 WHERE id = $2")
	if err != nil {
		return
	}

	_, err = stmt.ExecContext(ctx, input.LastLogin, input.ID)
	if err != nil {
		return
	}
	return
}
