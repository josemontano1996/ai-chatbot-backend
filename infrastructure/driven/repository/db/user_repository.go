package sqlcrepo

import (
	"context"
	"fmt"

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
