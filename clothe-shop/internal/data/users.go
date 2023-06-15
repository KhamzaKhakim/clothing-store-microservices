package data

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

type User struct {
	ID        int64    `json:"id"`
	Name      string   `json:"name"`
	Money     int64    `json:"money"`
	Email     string   `json:"email"`
	Password  password `json:"-"`
	Activated bool     `json:"activated"`
	Version   int      `json:"-"`
}

var AnonymousUser = &User{}

func (u *User) IsAnonymous() bool {
	return u == AnonymousUser
}

type password struct {
	plaintext *string
	hash      []byte
}

var (
	ErrDuplicateEmail = errors.New("duplicate email")
	ErrNotEnoughMoney = errors.New("not enough money")
)

type UserModel struct {
	DB *sql.DB
}

func (m UserModel) UpdateMoney(user *User, money int64) error {
	newMoney := user.Money - money
	if newMoney < 0 {
		return ErrNotEnoughMoney
	}

	query := `
UPDATE users
SET money = $1, version = version + 1
WHERE id = $2
RETURNING version`
	args := []any{
		newMoney,
		user.ID,
	}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	err := m.DB.QueryRowContext(ctx, query, args...).Scan(&user.Version)
	if err != nil {
		switch {
		case err.Error() == `pq: duplicate key value violates unique constraint "users_email_key"`:
			return ErrDuplicateEmail
		case errors.Is(err, sql.ErrNoRows):
			return ErrEditConflict
		default:
			return err
		}
	}
	return nil
}

func (m UserModel) Delete(id int64) error {
	if id < 1 {
		return ErrRecordNotFound
	}
	query := `
				DELETE FROM users
				WHERE id = $1`
	result, err := m.DB.Exec(query, id)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return ErrRecordNotFound
	}
	return nil
}
