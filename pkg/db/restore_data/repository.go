package restore_data

import (
	"context"
	"database/sql"
	"github.com/GoncharovMikhail/go-sql/errors"
	"github.com/GoncharovMikhail/go-sql/pkg/entity"
)

type SQLRestoreDataRepository interface {
	SaveInTx(context.Context, *entity.RestoreDataEntity, *sql.Tx) (*entity.RestoreDataEntity, errors.Errors)
	FindOneByUsernameInTx(context.Context, string, *sql.Tx) (*entity.RestoreDataEntity, bool, error)
}
