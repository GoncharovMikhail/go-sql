package merged

import (
	"context"
	"database/sql"
	"github.com/GoncharovMikhail/go-sql/pkg/entity"
)

type UserOps interface {
	SaveOrUpdate(ctx context.Context, ue *entity.UserEntity) (*entity.UserEntity, error)
}

type postgresUserOps struct {
	db *sql.DB
}

func NewPostgresUserOps(db *sql.DB) UserOps {
	return &postgresUserOps{
		db: db,
	}
}

func (postgresUserOps *postgresUserOps) SaveOrUpdate(ctx context.Context, ue *entity.UserEntity) (*entity.UserEntity, error) {
	_, err := postgresUserOps.db.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelDefault,
		ReadOnly:  false,
	})
	if err != nil {
		return nil, err
	}
	return nil, nil
}
