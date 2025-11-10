package postgres

import (
	"context"
	"errors"
	"fmt"

	"ai-service/internal/service/user"
)

var ErrUserNotFound = errors.New("user not found")

type UserRepository struct {
	db *DB
}

func NewUserRepository(db *DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(ctx context.Context, user *user.User) (string, error) {
	query := `
		INSERT INTO users (user_id, phone, password)
		VALUES ($1, $2, NOW(), NOW())
		RETURNING user_id`
	err := r.db.Pool.QueryRow(ctx, query, user.ID, user.Phone, user.Password).
		Scan(&user.ID)
	if err != nil {
		return "", fmt.Errorf("error creating user: %w", err)
	}
	return user.ID, nil
}

func (r *UserRepository) GetByID(ctx context.Context, id string) (*user.User, error) {
	query := `SELECT user_id, phone, password FROM users WHERE user_id = $1`
	row := r.db.Pool.QueryRow(ctx, query, id)
	u := new(user.User)
	err := row.Scan(&u.ID, &u.ID, &u.Phone, &u.Password)
	if err != nil {
		return nil, fmt.Errorf("get user: %w", err)
	}
	return u, nil
}

func (r *UserRepository) GetByPhone(ctx context.Context, phone string) (*user.User, error) {
	q := `SELECT user_id, phone, password FROM users WHERE phone = $1`
	row := r.db.Pool.QueryRow(ctx, q, phone)
	u := &user.User{}
	if err := row.Scan(&u.ID, &u.Phone, &u.Password); err != nil {
		return nil, ErrUserNotFound
	}
	return u, nil
}

func (r *UserRepository) List(ctx context.Context, limit, offset int) ([]*user.User, error) {
	query := `SELECT user_id, phone, password FROM users LIMIT $1 OFFSET $2`
	rows, err := r.db.Pool.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("list users: %w", err)
	}
	defer rows.Close()

	var users []*user.User
	for rows.Next() {
		u := new(user.User)
		if err := rows.Scan(&u.ID, &u.ID, &u.Phone, &u.Password); err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, nil
}

func (r *UserRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM users WHERE user_id = $1`
	_, err := r.db.Pool.Exec(ctx, query, id)
	return err
}
