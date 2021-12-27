package sql

import (
	"database/sql"
	uuid "github.com/jackc/pgtype/ext/gofrs-uuid"
	"github.com/testcontainers/testcontainers-go"
	db "sql/pkg/db/user"
	"testing"
)

type pgContainer struct {
	testcontainers.Container
}

var pgC *pgContainer = nil

func init() {

}

func TestSaveUser(t *testing.T) {
	var user *db.UserEntity = &db.UserEntity{
		Id:       uuid.UUID{},
		Username: "test",
		Password: "test",
		RestoreData: &db.RestoreData{
			//UserId:      "",
			Email:       "test",
			PhoneNumber: sql.NullString{},
		},
	}
	print(user)
}
