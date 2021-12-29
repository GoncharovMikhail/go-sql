package sql

import (
	"context"
	"database/sql"
	"github.com/GoncharovMikhail/go-sql/pkg/db/user"
	"github.com/Masterminds/squirrel"
)

type PostgresUserRepository struct {
	Db *sql.DB
}

func (repository *PostgresUserRepository) Save(ctx context.Context, entity *user.UserEntity) (*user.UserEntity, error) {
	tx, err := repository.Db.BeginTx(
		ctx,
		&sql.TxOptions{
			Isolation: sql.LevelDefault,
			ReadOnly:  false,
		},
	)
	if err != nil {
		return nil, err
	}
	err = squirrel.
		Insert("\"user\"").
		Columns("username", "password").
		Values(entity.Username, entity.Password).
		Suffix("RETURNING *").
		PlaceholderFormat(squirrel.Dollar).
		RunWith(repository.Db).
		ScanContext(ctx, &entity.Id, &entity.Username, &entity.Password)
	if err != nil {
		errTxRollback := tx.Rollback()
		if errTxRollback != nil {
			return nil, errTxRollback
		}
		return nil, err
	}
	if entity.RestoreData == nil {
		err := tx.Commit()
		if err != nil {
			return nil, err
		}
		return entity, nil
	}
	err = squirrel.
		Insert("restore_data").
		Columns("user_id", "email", "phone_number").
		Values(entity.Id, entity.Email, &entity.PhoneNumber).
		Suffix("RETURNING *").
		PlaceholderFormat(squirrel.Dollar).
		RunWith(repository.Db).
		ScanContext(ctx, &entity.UserId, &entity.Email, &entity.PhoneNumber)
	if err != nil {
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
		PlaceholderFormat(squirrel.Dollar).
		RunWith(repository.Db).
		ToSql()
	if err != nil {
		return nil, err
	}
	row := repository.Db.QueryRowContext(ctx, query, args)
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
