package entity

import (
	"database/sql"
	pgUuidType "github.com/jackc/pgtype/ext/gofrs-uuid"
)

type UserStatusEntity struct {
	UserId                pgUuidType.UUID `db:"user_id"`
	AccountNonExpired     sql.NullBool    `db:"account_non_expired"`
	AccountNonLocked      sql.NullBool    `db:"account_non_locked"`
	CredentialsNonExpired sql.NullBool    `db:"credentials_non_expired"`
	Enabled               sql.NullBool    `db:"enabled"`
}
