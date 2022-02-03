package pg

import (
	"context"
	"database/sql"
	"github.com/GoncharovMikhail/go-sql/errors"
	dbConsts "github.com/GoncharovMikhail/go-sql/pkg/db/consts"
	"github.com/GoncharovMikhail/go-sql/pkg/db/entity_information"
	"github.com/GoncharovMikhail/go-sql/pkg/db/util"
	"github.com/GoncharovMikhail/go-sql/pkg/entity"
	"github.com/Masterminds/squirrel"
	"github.com/gofrs/uuid"
)

const (
	tableName    = "\"user\""
	idColumnName = "id"
	username     = "username"
	password     = "password"
)

func SaveOrUpdateInTx(ctx context.Context, ude *entity.UserDataEntity, tx *sql.Tx) (*entity.UserDataEntity, errors.Errors, *sql.Tx) {
	var isNew bool
	var errorz errors.Errors
	isNew, errorz, tx = entity_information.IsNew(ctx, tableName, idColumnName, ude.Id, tx)
	if errorz != nil {
		errorz, tx = util.TxRollbackErrorHandle(errorz.Get(), tx)
		return nil,
			errorz,
			tx
	}
	if isNew {
		return save(ctx, ude, tx)
	}
	return update(ctx, ude, tx)
}

func FindOneByUsernameInTx(ctx context.Context, userUsername string, tx *sql.Tx) (*entity.UserDataEntity, bool, errors.Errors, *sql.Tx) {
	var ue = &entity.UserDataEntity{}
	err := squirrel.
		Select("*").
		From(tableName).
		Where(squirrel.Eq{username: userUsername}).
		PlaceholderFormat(squirrel.Dollar).
		RunWith(tx).
		ScanContext(ctx, &ue.Id, &ue.Username, &ue.Password)
	if err != nil {
		var errorz errors.Errors
		errorz, tx = util.TxRollbackErrorHandle(err, tx)
		return nil,
			false,
			errorz,
			tx
	}
	return ue,
		true,
		nil,
		tx
}

func save(ctx context.Context, ude *entity.UserDataEntity, tx *sql.Tx) (*entity.UserDataEntity, errors.Errors, *sql.Tx) {
	columnNames, columnValues := getColumnNamesAndColumnValues(ude)
	err := squirrel.
		Insert(tableName).
		Columns(columnNames...).
		Values(columnValues...).
		Suffix(dbConsts.Suffix).
		PlaceholderFormat(squirrel.Dollar).
		RunWith(tx).
		ScanContext(ctx, &ude.Id, &ude.Username, &ude.Password)
	if err != nil {
		var errorz errors.Errors
		errorz, tx = util.TxRollbackErrorHandle(err, tx)
		return nil,
			errorz,
			tx
	}
	return ude,
		nil,
		tx
}

func getColumnNamesAndColumnValues(ude *entity.UserDataEntity) ([]string, []interface{}) {
	columnNames := make([]string, 0)
	columnNames = append(columnNames, username, password)
	columnValues := make([]interface{}, 0)
	columnValues = append(columnValues, ude.Username, ude.Password)
	if ude.Id.UUID != uuid.Nil {
		columnNames = append(columnNames, idColumnName)
		columnValues = append(columnValues, ude.Id)
	}
	return columnNames,
		columnValues
}

func update(ctx context.Context, ude *entity.UserDataEntity, tx *sql.Tx) (*entity.UserDataEntity, errors.Errors, *sql.Tx) {
	err := squirrel.
		Update(tableName).
		Set(username, ude.Username).
		Set(password, ude.Password).
		Where(squirrel.Eq{idColumnName: ude.Id.UUID}).
		ScanContext(ctx, &ude.Id, &ude.Username, &ude.Password)
	if err != nil {
		var errorz errors.Errors
		errorz, tx = util.TxRollbackErrorHandle(err, tx)
		return nil,
			errorz,
			tx
	}
	return ude,
		nil,
		tx
}
