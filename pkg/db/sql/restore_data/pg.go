package restore_data

import (
	"context"
	"database/sql"
	"github.com/GoncharovMikhail/go-sql/errors"
	dbConsts "github.com/GoncharovMikhail/go-sql/pkg/db/sql/consts"
	"github.com/GoncharovMikhail/go-sql/pkg/db/sql/entity_information"
	"github.com/GoncharovMikhail/go-sql/pkg/db/sql/util"
	"github.com/GoncharovMikhail/go-sql/pkg/entity"
	"github.com/Masterminds/squirrel"
)

const (
	tableName    = "restore_data"
	idColumnName = "user_id"
	email        = "email"
	phoneNumber  = "email"
)

func SaveOrUpdateInTx(ctx context.Context, rde *entity.RestoreDataEntity, tx *sql.Tx) (*entity.RestoreDataEntity, errors.Errors, *sql.Tx) {
	var isNew bool
	var errorz errors.Errors
	isNew, errorz, tx = entity_information.IsNew(ctx, tableName, idColumnName, rde.UserId, tx)
	if errorz != nil {
		errorz, tx = util.TxRollbackErrorHandle(errorz.Get(), tx)
		return nil,
			errorz,
			tx
	}
	if !isNew {
		return save(ctx, rde, tx)
	}
	return update(ctx, rde, tx)

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

func save(ctx context.Context, rde *entity.RestoreDataEntity, tx *sql.Tx) (*entity.RestoreDataEntity, errors.Errors, *sql.Tx) {
	err := squirrel.
		Insert(tableName).
		Columns(idColumnName, email, phoneNumber).
		Values(rde.UserId, rde.Email, rde.PhoneNumber).
		Suffix(dbConsts.Suffix).
		PlaceholderFormat(squirrel.Dollar).
		RunWith(tx).
		ScanContext(ctx, &rde.UserId, &rde.Email, &rde.PhoneNumber)
	if err != nil {
		var errorz errors.Errors
		errorz, tx = util.TxRollbackErrorHandle(err, tx)
		return nil,
			errorz,
			tx
	}
	return rde,
		nil,
		tx
}

func update(ctx context.Context, rde *entity.RestoreDataEntity, tx *sql.Tx) (*entity.RestoreDataEntity, errors.Errors, *sql.Tx) {
	err := squirrel.
		Update(tableName).
		Set(idColumnName, rde.UserId).
		Set(email, rde.Email).
		Set(phoneNumber, rde.PhoneNumber).
		Where(squirrel.Eq{idColumnName: rde.UserId}).
		Suffix(dbConsts.Suffix).
		PlaceholderFormat(squirrel.Dollar).
		RunWith(tx).
		ScanContext(ctx, &rde.UserId, &rde.Email, &rde.PhoneNumber)
	if err != nil {
		var errorz errors.Errors
		errorz, tx = util.TxRollbackErrorHandle(err, tx)
		return nil,
			errorz,
			tx
	}
	return rde,
		nil,
		tx
}
