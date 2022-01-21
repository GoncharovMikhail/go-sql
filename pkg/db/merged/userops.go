package merged

import (
	"context"
	"database/sql"
	restoreDataRepositoryPkg "github.com/GoncharovMikhail/go-sql/pkg/db/restore_data"
	pdRestoreDataRepositoryPkg "github.com/GoncharovMikhail/go-sql/pkg/db/restore_data/impl/sql"
	userRepositoryPkg "github.com/GoncharovMikhail/go-sql/pkg/db/user"
	pgUserRepositoryPkg "github.com/GoncharovMikhail/go-sql/pkg/db/user/impl/sql"
	userStatusRepositoryPkg "github.com/GoncharovMikhail/go-sql/pkg/db/user_status"
	pgUserStatusRepositoryPkg "github.com/GoncharovMikhail/go-sql/pkg/db/user_status/impl/sql"
	"github.com/GoncharovMikhail/go-sql/pkg/entity"
)

type UserOps interface {
	SaveOrUpdate(ctx context.Context, ue *entity.UserEntity) (*entity.UserEntity, error)
}

type postgresUserOps struct {
	db  *sql.DB
	ur  userRepositoryPkg.SQLUserRepository
	rdr restoreDataRepositoryPkg.SQLRestoreDataRepository
	usr userStatusRepositoryPkg.SQLUserStatusRepository
}

func NewPostgresUserOps(db *sql.DB,
	ur userRepositoryPkg.SQLUserRepository,
	rdr restoreDataRepositoryPkg.SQLRestoreDataRepository,
	usr userStatusRepositoryPkg.SQLUserStatusRepository,
) *postgresUserOps {
	return &postgresUserOps{
		db:  db,
		ur:  ur,
		rdr: rdr,
		usr: usr,
	}
}

// NewPostgresUserOpsWithDefaults returns default implementation
// @Prefer
func NewPostgresUserOpsWithDefaults(db *sql.DB) UserOps {
	return NewPostgresUserOps(
		db,
		pgUserRepositoryPkg.NewPostgresUserDataRepository(),
		pdRestoreDataRepositoryPkg.NewPostgresRestoreDataRepository(db),
		pgUserStatusRepositoryPkg.NewPostgresUserStatusRepository(db),
	)
}

func (postgresUserOps *postgresUserOps) SaveOrUpdate(ctx context.Context, ue *entity.UserEntity) (*entity.UserEntity, error) {
	tx, err := postgresUserOps.db.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelDefault,
		ReadOnly:  false,
	})
	if err != nil {
		return nil, err
	}
	return nil, nil
}
