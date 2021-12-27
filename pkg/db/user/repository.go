package user

import (
	"context"
)

type UserRepository interface {
	Save(context.Context, *UserEntity) (*UserEntity, error)
	FindOneByUsername(ctx context.Context, username string) (retUser *UserEntity, retErr error)
}
