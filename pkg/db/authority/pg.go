package authority

import (
	"context"
	"database/sql"
	"github.com/GoncharovMikhail/go-sql/errors"
	"github.com/GoncharovMikhail/go-sql/pkg/entity"
	"github.com/Masterminds/squirrel"
)

func SaveInTx(ctx context.Context, ae *entity.AuthorityEntity, tx *sql.Tx) (*entity.AuthorityEntity, errors.Errors) {
	err := squirrel.
		Insert("authority").
		Columns("name").
		Values(ae.Name).
		PlaceholderFormat(squirrel.Dollar).
		RunWith(tx).
		ScanContext(ctx, &ae.Id, &ae.Name)
	if err != nil {
		errTxRollback := tx.Rollback()
		if err != nil {
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
	return ae, nil
}

func FindAllByUsernameInTx(ctx context.Context, username string, tx *sql.Tx) ([]*entity.AuthorityEntity, bool) {
	var authorityNames []string
	err := squirrel.
		Select("name").
		From("authority").
		Join("user_authority USING (authority_id)").
		Join("user USING (user_id)").
		Where(squirrel.Eq{"username": username}).
		PlaceholderFormat(squirrel.Dollar).
		RunWith(tx).
		ScanContext(ctx, authorityNames)
	if err != nil {
		return nil, false
	}
	var res []*entity.AuthorityEntity
	for _, name := range authorityNames {
		_ = append(res, &entity.AuthorityEntity{Name: name})
	}
	return res, true
}
