package user

import (
	"context"
	"database/sql"
	"github.com/GoncharovMikhail/go-sql/errors"
	dbConsts "github.com/GoncharovMikhail/go-sql/pkg/db/sql/consts"
	"github.com/GoncharovMikhail/go-sql/pkg/db/sql/entity_information"
	"github.com/GoncharovMikhail/go-sql/pkg/db/sql/util"
	"github.com/GoncharovMikhail/go-sql/pkg/entity"
	"github.com/Masterminds/squirrel"
	"github.com/gofrs/uuid"
)

const (
	TableName    = "\"user\""
	IdColumnName = "id"
	Username     = "username"
	Password     = "password"
)

func SaveOrUpdateInTx(ctx context.Context, ude *entity.UserDataEntity, tx *sql.Tx) (*entity.UserDataEntity, errors.Errors, *sql.Tx) {
	var isNew bool
	var errorz errors.Errors
	isNew, errorz, tx = entity_information.IsNew(ctx, TableName, IdColumnName, ude.Id, tx)
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
		From(TableName).
		Where(squirrel.Eq{Username: userUsername}).
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
		Insert(TableName).
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
	columnNames = append(columnNames, Username, Password)
	columnValues := make([]interface{}, 0)
	columnValues = append(columnValues, ude.Username, ude.Password)
	if ude.Id.UUID != uuid.Nil {
		columnNames = append(columnNames, IdColumnName)
		columnValues = append(columnValues, ude.Id)
	}
	return columnNames,
		columnValues
}

func update(ctx context.Context, ude *entity.UserDataEntity, tx *sql.Tx) (*entity.UserDataEntity, errors.Errors, *sql.Tx) {
	err := squirrel.
		Update(TableName).
		Set(Username, ude.Username).
		Set(Password, ude.Password).
		Where(squirrel.Eq{IdColumnName: ude.Id.UUID}).
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
