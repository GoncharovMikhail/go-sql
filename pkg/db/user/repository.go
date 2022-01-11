package user

import (
	"context"
	"database/sql"
	"github.com/GoncharovMikhail/go-sql/errors"
	"github.com/GoncharovMikhail/go-sql/pkg/entity"
)

type SQLUserRepository interface {
	SaveInTx(context.Context, *entity.UserDataEntity, *sql.Tx) (*entity.UserDataEntity, errors.Errors)
	FindOneByUsernameInTx(context.Context, string, *sql.Tx) (*entity.UserDataEntity, bool, errors.Errors)
}
