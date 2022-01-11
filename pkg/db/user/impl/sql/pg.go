package sql

import (
	"context"
	"database/sql"
	"github.com/GoncharovMikhail/go-sql/errors"
	dbConsts "github.com/GoncharovMikhail/go-sql/pkg/db/consts"
	"github.com/GoncharovMikhail/go-sql/pkg/db/user"
	"github.com/GoncharovMikhail/go-sql/pkg/entity"
	"github.com/Masterminds/squirrel"
)

type postgresUserRepository struct{}

func NewPostgresUserRepository() user.SQLUserRepository {
	return &postgresUserRepository{}
}

func (postgresUserRepository *postgresUserRepository) SaveInTx(
	ctx context.Context,
	entity *entity.UserDataEntity,
	tx *sql.Tx,
) (*entity.UserDataEntity, errors.Errors) {
	err := squirrel.
		Insert("\"user\"").
		Columns("username", "password").
		Values(entity.Username, entity.Password).
		Suffix(dbConsts.Suffix).
		PlaceholderFormat(squirrel.Dollar).
		RunWith(tx).
		ScanContext(ctx, &entity.Id, &entity.Username, &entity.Password)
	if err != nil {
		errTxRollback := tx.Rollback()
		if errTxRollback != nil {
			return nil,
				errors.NewErrors(
					errors.BuildSimpleErrMsg("err", err),
					err,
					errors.NewErrors(
						errors.BuildSimpleErrMsg("errTxRollback", errTxRollback),
						errTxRollback,
						nil,
					),
				)
		}
		return nil, errors.NewErrors(
			errors.BuildSimpleErrMsg("err", err),
			err,
			nil,
		)
	}
	err = tx.Commit()
	if err != nil {
		return nil, errors.NewErrors(
			errors.BuildSimpleErrMsg("err", err),
			err,
			nil,
		)
	}
	return entity, nil
}

func (postgresUserRepository *postgresUserRepository) FindOneByUsernameInTx(
	ctx context.Context,
	username string,
	tx *sql.Tx,
) (*entity.UserDataEntity, bool, errors.Errors) {
	var ue = &entity.UserDataEntity{}
	err := squirrel.
		Select("*").
		From("\"user\"").
		Where(squirrel.Eq{"username": username}).
		PlaceholderFormat(squirrel.Dollar).
		RunWith(tx).
		ScanContext(ctx, &ue.Id, &ue.Username, &ue.Password)
	if err != nil {
		return nil,
			false,
			errors.NewErrors(
				errors.BuildSimpleErrMsg("err", err),
				err,
				nil,
			)
	}
	return ue,
		true,
		nil
}
