package restore_data

import (
	"context"
	"database/sql"
	"github.com/GoncharovMikhail/go-sql/errors"
	dbConsts "github.com/GoncharovMikhail/go-sql/pkg/db/consts"
	"github.com/GoncharovMikhail/go-sql/pkg/entity"
	"github.com/Masterminds/squirrel"
)

func SaveInTx(ctx context.Context, rde *entity.RestoreDataEntity, tx *sql.Tx) (*entity.RestoreDataEntity, errors.Errors) {
	err := squirrel.
		Insert("restore_data").
		Columns("user_id", "email", "phone_number").
		Values(rde.UserId, rde.Email, rde.PhoneNumber).
		Suffix(dbConsts.Suffix).
		PlaceholderFormat(squirrel.Dollar).
		RunWith(tx).
		ScanContext(ctx, &rde.UserId, &rde.Email, &rde.PhoneNumber)
	if err != nil {
		errTxRollback := tx.Rollback()
		if errTxRollback != nil {
			return nil,
				errors.NewErrors(
					errors.BuildSimpleErrMsg("err", err),
					err,
					errors.NewErrors(
						errors.BuildSimpleErrMsg("errTxRollback", err),
						errTxRollback,
						nil,
					),
				)
		}
		return nil,
			errors.NewErrors(
				errors.BuildSimpleErrMsg("err", err),
				err,
				nil,
			)

	}
	err = tx.Commit()
	if err != nil {
		return nil,
			errors.NewErrors(
				errors.BuildSimpleErrMsg("err", err),
				err,
				nil,
			)

	}
	return rde, nil
}

func FindOneByUsernameInTx(ctx context.Context, username string, tx *sql.Tx) (*entity.RestoreDataEntity, bool, error) {
	var rde = &entity.RestoreDataEntity{}
	err := squirrel.
		Select("user_id", "email", "phone_number").
		From("restore_data").
		Join("user USING (user_id)").
		Where(squirrel.Eq{"username": username}).
		PlaceholderFormat(squirrel.Dollar).
		RunWith(tx).
		ScanContext(ctx, &rde.UserId, &rde.Email, rde.PhoneNumber)
	if err != nil {
		errTxRollback := tx.Rollback()
		if errTxRollback != nil {
			return nil, false, errTxRollback
		}
		return nil, false, err
	}
	return rde, true, nil
}
