package repository

import (
	"context"

	"github.com/mystaline/paperid-test/internal/entity"
)

var _ UserRepositor = (*UserRepository)(nil)

type UserRepository struct{}

type UserRepositor interface {
	GetOneUserByID(ctx context.Context, db DBTX, id int64) (*entity.User, error)
}

func NewUserRepository() UserRepositor {
	return &UserRepository{}
}

func (r *UserRepository) GetOneUserByID(
	ctx context.Context,
	db DBTX,
	id int64,
) (*entity.User, error) {
	rawQuery := `SELECT id, full_name, created_at, updated_at FROM users WHERE id = $1 LIMIT 1`
	rows := db.QueryRow(ctx, rawQuery, id)

	user := entity.User{}
	err := rows.Scan(&user.ID, &user.FullName, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return &user, nil
}
