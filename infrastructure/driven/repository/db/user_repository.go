package sqlcrepo

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	sqlc "github.com/josemontano1996/ai-chatbot-backend/infrastructure/driven/repository/db/output"
	"github.com/josemontano1996/ai-chatbot-backend/internal/entities"
	outrepo "github.com/josemontano1996/ai-chatbot-backend/internal/ports/out/repositories"
)

type UserRepository struct {
	*sqlc.Queries
	db *pgxpool.Pool
}

func NewUserRepository(conn *pgxpool.Pool) outrepo.UserRepository {
	return &UserRepository{
		db:      conn,
		Queries: sqlc.New(conn),
	}
}

func (r *UserRepository) CreateUser(ctx context.Context, email string, hashedPassword string) (*entities.User, error) {
	params := sqlc.CreateUserParams{
		Email:    email,
		Password: hashedPassword,
	}

	createdUser, err := r.Queries.CreateUser(ctx, params)

	if err != nil {
		return nil, err
	}
	return entities.NewUserEntity(createdUser.ID.String(), createdUser.Email)
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (userEntity *entities.User, hashedPassword string, err error) {
	user, err := r.Queries.FindByEmail(ctx, email)

	if err != nil {
		return nil, "", fmt.Errorf("error finding user by email: %w", err)
	}

	userEntity, err = entities.NewUserEntity(user.ID.String(), user.Email)
	hashedPassword = user.Password
	return
}

func (r *UserRepository) FindById(ctx context.Context, id uuid.UUID) (userEntity *entities.User, hashedPassword string, err error) {
	user, err := r.Queries.FindById(ctx, id)

	if err != nil {
		return nil, "", fmt.Errorf("error finding user by id: %w", err)
	}

	userEntity, err = entities.NewUserEntity(user.ID.String(), user.Email)
	hashedPassword = user.Password
	return
}

func (r *UserRepository) UpdateUser(ctx context.Context, id uuid.UUID, new_email, new_hashed_password string) (*entities.User, error) {
	params := sqlc.UpdateUserParams{
		ID:          id,
		NewEmail:    new_email,
		NewPassword: new_hashed_password,
	}
	updatedUser, err := r.Queries.UpdateUser(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("error updating user: %w", err)
	}

	return entities.NewUserEntity(updatedUser.ID.String(), updatedUser.Email)
}

func (r *UserRepository) DeleteUser(ctx context.Context, id uuid.UUID) error {
	return r.Queries.DeleteUser(ctx, id)
}
