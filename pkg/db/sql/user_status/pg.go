package user_status

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
	TableName             = "user_status"
	IdColumnName          = "user_id"
	AccountNonExpired     = "account_non_expired"
	AccountNonLocked      = "account_non_locked"
	CredentialsNonExpired = "credentials_non_expired"
	Enabled               = "enabled"
)

func SaveOrUpdateInTx(ctx context.Context, use *entity.UserStatusEntity, tx *sql.Tx) (*entity.UserStatusEntity, errors.Errors, *sql.Tx) {
	var isNew bool
	var errorz errors.Errors
	isNew, errorz, tx = entity_information.IsNew(ctx, TableName, IdColumnName, use.UserId, tx)
	if errorz != nil {
		errorz, tx = util.TxRollbackErrorHandle(errorz.Get(), tx)
		return nil,
			errorz,
			tx
	}
	if isNew {
		return save(ctx, use, tx)
	}
	return update(ctx, use, tx)
}

func FindOneByUsernameInTx(ctx context.Context, username string, tx *sql.Tx) (*entity.UserStatusEntity, bool, errors.Errors) {
	var use = &entity.UserStatusEntity{}
	err := squirrel.
		Select(IdColumnName, AccountNonExpired, AccountNonLocked, CredentialsNonExpired, Enabled).
		From(TableName).
		Join("\"user\" as u on u.id = "+IdColumnName).
		Where(squirrel.Eq{"username": username}).
		RunWith(tx).
		ScanContext(ctx, &use.AccountNonExpired, &use.AccountNonLocked, &use.CredentialsNonExpired, &use.Enabled)
	if err != nil {
		return nil, false, nil
	}
	return use, true, nil
}

func save(ctx context.Context, use *entity.UserStatusEntity, tx *sql.Tx) (*entity.UserStatusEntity, errors.Errors, *sql.Tx) {
	errSave := squirrel.
		Insert(TableName).
		Columns(IdColumnName, AccountNonExpired, AccountNonLocked, CredentialsNonExpired, Enabled).
		Values(use.UserId, use.AccountNonExpired, use.AccountNonLocked, use.CredentialsNonExpired, use.Enabled).
		Suffix(dbConsts.Suffix).
		PlaceholderFormat(squirrel.Dollar).
		RunWith(tx).
		ScanContext(ctx, &use.UserId, &use.AccountNonExpired, &use.AccountNonLocked, &use.CredentialsNonExpired, &use.Enabled)
	if errSave != nil {
		var errorz errors.Errors
		errorz, tx = util.TxRollbackErrorHandle(errSave, tx)
		return nil,
			errorz,
			tx
	}
	return use,
		nil,
		tx
}

func update(ctx context.Context, use *entity.UserStatusEntity, tx *sql.Tx) (*entity.UserStatusEntity, errors.Errors, *sql.Tx) {
	errUpdate := squirrel.
		Update(TableName).
		Set(AccountNonExpired, use.AccountNonExpired).
		Set(AccountNonLocked, use.AccountNonLocked).
		Set(CredentialsNonExpired, use.CredentialsNonExpired).
		Set(Enabled, use.Enabled).
		Where(squirrel.Eq{IdColumnName: use.UserId}).
		ScanContext(ctx, &use.UserId, &use.AccountNonExpired, &use.AccountNonLocked, &use.CredentialsNonExpired, &use.Enabled)
	if errUpdate != nil {
		var errorz errors.Errors
		errorz, tx = util.TxRollbackErrorHandle(errUpdate, tx)
		return nil, errorz, tx
	}
	return use, nil, tx
}
