package entity

import (
	"database/sql"
	pgUuidType "github.com/jackc/pgtype/ext/gofrs-uuid"
)

type RestoreDataEntity struct {
	UserId      pgUuidType.UUID `db:"user_id"`
	Email       string          `db:"email"`
	PhoneNumber sql.NullString  `db:"phone_number,omitempty"`
}
