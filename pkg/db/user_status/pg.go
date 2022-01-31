package user_status

import (
	"context"
	"database/sql"
	"github.com/GoncharovMikhail/go-sql/errors"
	dbConsts "github.com/GoncharovMikhail/go-sql/pkg/db/consts"
	"github.com/GoncharovMikhail/go-sql/pkg/entity"
	"github.com/Masterminds/squirrel"
)

func SaveOrUpdateInTx(ctx context.Context, use *entity.UserStatusEntity, tx *sql.Tx) (*entity.UserStatusEntity, errors.Errors) {
	err := squirrel.
		Insert("user_status").
		Columns("user_id", "account_non_expired", "account_non_locked", "credentials_non_expired", "enabled").
		Values(use.UserId, use.AccountNonLocked, use.AccountNonLocked, use.CredentialsNonExpired, use.Enabled).
		Suffix(dbConsts.Suffix).
		PlaceholderFormat(squirrel.Dollar).
		RunWith(tx).
		ScanContext(ctx, &use.UserId, &use.AccountNonExpired, &use.AccountNonLocked, &use.CredentialsNonExpired, &use.Enabled)
	if err != nil {
		return nil, errors.NewErrors(
			errors.BuildSimpleErrMsg("err", err),
			err,
			nil,
		)
	}
	return use, nil
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
