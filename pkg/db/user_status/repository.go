package user_status

import (
	"context"
	"database/sql"
	"github.com/GoncharovMikhail/go-sql/errors"
	"github.com/GoncharovMikhail/go-sql/pkg/entity"
)

type SQLUserStatusRepository interface {
	SaveOrUpdateInTx(context.Context, *entity.UserStatusEntity, *sql.Tx) (*entity.UserStatusEntity, errors.Errors)
	FindOneByUsernameInTx(context.Context, string, *sql.Tx) (*entity.UserStatusEntity, bool, errors.Errors)
}
