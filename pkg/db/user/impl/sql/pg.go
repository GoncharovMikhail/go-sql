package sql

import (
	"context"
	"database/sql"
	"github.com/GoncharovMikhail/go-sql/pkg/db/user"
	"github.com/Masterminds/squirrel"
	"strings"
)

type PostgresUserRepository struct {
	Db *sql.DB
}

const (
	saveUserQuery        = `INSERT INTO "user"(username, password) VALUES ($1, $2) RETURNING *`
	saveRestoreDataQuery = `INSERT INTO restore_data(user_id, email, phone_number) VALUES ($1, $2, $3) RETURNING *`
)

func (repository *PostgresUserRepository) Save(ctx context.Context, entity *user.UserEntity) (*user.UserEntity, error) {
	tx, err := repository.Db.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelDefault,
		ReadOnly:  false,
	})
	if err != nil {
		return nil, err
	}
	rowUser := tx.QueryRowContext(ctx, saveUserQuery, entity.Username, entity.Password)
	errIdScan := rowUser.Scan(
		&entity.Id,
		&entity.Username,
		&entity.Password,
	)
	if errIdScan != nil {
		errTxRollback := tx.Rollback()
		if errTxRollback != nil {
			return nil, errTxRollback
		}
		return nil, errIdScan
	}
	if entity.RestoreData == nil {
		err := tx.Commit()
		if err != nil {
			return nil, err
		}
		return entity, nil
	}
	rowSaveRestoreData := tx.QueryRowContext(ctx, saveRestoreDataQuery, entity.Id, entity.Email, entity.PhoneNumber)
	if rowSaveRestoreData.Err() != nil {
		return nil, rowSaveRestoreData.Err()
	}
	errRestoreDataScan := rowSaveRestoreData.Scan(
		&entity.RestoreData.UserId,
		&entity.RestoreData.Email,
		&entity.RestoreData.PhoneNumber,
	)
	if errRestoreDataScan != nil {
		return nil, err
	}
	err = tx.Commit()
	if err != nil {
		return nil, err
	}
	return entity, nil
}

func (repository *PostgresUserRepository) FindOneByUsername(ctx context.Context, username string) (*user.UserEntity, error) {
	query, args, err := squirrel.
		Select("*").
		From("\"user\"").
		Where(map[string]interface{}{"username": username}).
		ToSql()
	if err != nil {
		return nil, err
	}
	replace := strings.Replace(query, "?", "$1", -1)
	row := repository.Db.QueryRowContext(ctx, replace, args[0].(string))
	if row.Err() != nil {
		return nil, row.Err()
	}
	var retUser user.UserEntity
	err = row.Scan(
		&retUser.Id,
		&retUser.Username,
		&retUser.Password,
	)
	if err != nil {
		return nil, err
	}
	return &retUser, nil
}
