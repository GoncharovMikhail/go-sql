package sql

import (
	"context"
	"database/sql"
	"github.com/Masterminds/squirrel"
	"sql/pkg/db/user"
)

type PostgresUserRepository struct {
	Db *sql.DB
}

var sb = squirrel.StatementBuilderType{}

const (
	saveUserQuery        = `INSERT INTO "user"(username, password) VALUES (?, ?) RETURNING *`
	saveRestoreDataQuery = `INSERT INTO restore_data(user_id, email, phone_number) VALUES (?, ?, ?) RETURNING *`
)

func (s *PostgresUserRepository) Save(ctx context.Context, entity *user.UserEntity) (*user.UserEntity, error) {
	tx, err := s.Db.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelDefault,
		ReadOnly:  false,
	})
	if err != nil {
		return nil, err
	}
	rowUser := tx.QueryRowContext(ctx, saveUserQuery, entity.Username, entity.Password)
	errIdScan := rowUser.Scan(&entity.Id)
	if errIdScan != nil {
		errTxRollback := tx.Rollback()
		if errTxRollback != nil {
			return nil, errTxRollback
		}
		return nil, err
	}
	if entity.RestoreData == nil {
		return entity, nil
	}
	if rowSaveRestoreData := tx.QueryRowContext(ctx, saveRestoreDataQuery, entity.Id, entity.Email, entity.PhoneNumber); rowSaveRestoreData.Err() != nil {
		return nil, rowSaveRestoreData.Err()
	}
	err = tx.Commit()
	if err != nil {
		return nil, err
	}
	return entity, nil
}

func (s *PostgresUserRepository) FindOneByUsername(ctx context.Context, username string) (retUser *user.UserEntity, retErr error) {
	query := sb.
		Select("*").
		From("\"user\"").
		Where(map[string]interface{}{"username": username}).
		Limit(1)
	rows, err := query.QueryContext(ctx)
	if err != nil {
		//return status.Error(codes.Internal, retErr.Error()), nil
	}
	defer func() {
		cerr := rows.Close()
		if retErr == nil && cerr != nil {
			//retErr = status.Error(codes.Internal, cerr.Error())
		}
	}()
	for rows.Next() {
		err := rows.Scan(
			&retUser.Id,
			&retUser.Username,
			&retUser.Password,
		)
		if err != nil {
			retErr = err
		}
	}
	return retUser, retErr
}
