package user

import (
	"database/sql"
	uuid "github.com/jackc/pgtype/ext/gofrs-uuid"
)

type UserEntity struct {
	Id       uuid.UUID `db:"id"`
	Username string    `db:"username"`
	Password string    `db:"password"`
	*RestoreData
}

type RestoreData struct {
	UserId      string         `db:"user_id"`
	Email       string         `db:"email"`
	PhoneNumber sql.NullString `db:"phone_number,omitempty"`
}
