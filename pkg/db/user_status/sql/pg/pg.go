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
)

const (
	tableName             = "user_status"
	idColumnName          = "user_id"
	accountNonExpired     = "account_non_expired"
	accountNonLocked      = "account_non_locked"
	credentialsNonExpired = "credentials_non_expired"
	enabled               = "enabled"
)

func SaveOrUpdateInTx(ctx context.Context, use *entity.UserStatusEntity, tx *sql.Tx) (*entity.UserStatusEntity, errors.Errors, *sql.Tx) {
	var isNew bool
	var errorz errors.Errors
	isNew, errorz, tx = entity_information.IsNew(ctx, tableName, idColumnName, use.UserId, tx)
	if errorz != nil {
		return nil, errorz, tx
	}
	if isNew {
		return save(ctx, use, tx)
	}
	return update(ctx, use, tx)
}

func FindOneByUsernameInTx(ctx context.Context, username string, tx *sql.Tx) (*entity.UserStatusEntity, bool, errors.Errors) {
	var use = &entity.UserStatusEntity{}
	err := squirrel.
		Select("user_id", "account_non_expired", "account_non_locked", "credentials_non_expired", "enabled").
		From("user_status").
		Join("user USING (user_id)").
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
		Insert(tableName).
		Columns(idColumnName, accountNonExpired, accountNonLocked, credentialsNonExpired, enabled).
		Values(use.UserId, use.AccountNonExpired, use.AccountNonLocked, use.CredentialsNonExpired, use.Enabled).
		Suffix(dbConsts.Suffix).
		PlaceholderFormat(squirrel.Dollar).
		RunWith(tx).
		ScanContext(ctx, &use.UserId, &use.AccountNonExpired, &use.AccountNonLocked, &use.CredentialsNonExpired, &use.Enabled)
	var errorz errors.Errors
	errorz, tx = util.TxRollbackErrorHandle(errSave, tx)
	if errorz != nil {
		return nil, errorz, tx
	}
	return use, nil, tx
}

func update(ctx context.Context, use *entity.UserStatusEntity, tx *sql.Tx) (*entity.UserStatusEntity, errors.Errors, *sql.Tx) {
	errUpdate := squirrel.
		Update(tableName).
		Set(accountNonExpired, use.AccountNonExpired).
		Set(accountNonLocked, use.AccountNonLocked).
		Set(credentialsNonExpired, use.CredentialsNonExpired).
		Set(enabled, use.Enabled).
		Where(squirrel.Eq{idColumnName: use.UserId}).
		ScanContext(ctx, &use.UserId, &use.AccountNonExpired, &use.AccountNonLocked, &use.CredentialsNonExpired, &use.Enabled)
	var errorz errors.Errors
	errorz, tx = util.TxRollbackErrorHandle(errUpdate, tx)
	if errorz != nil {
		return nil, errorz, tx
	}
	return use, nil, tx
}
