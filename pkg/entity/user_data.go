package entity

import (
	pgUuidType "github.com/jackc/pgtype/ext/gofrs-uuid"
)

type UserDataEntity struct {
	Id       pgUuidType.UUID `db:"id"`
	Username string          `db:"username"`
	Password string          `db:"password"`
}
