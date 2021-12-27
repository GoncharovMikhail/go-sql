package sql

import (
	"database/sql"
	db "github.com/GoncharovMikhail/go-sql/pkg/db/user"
	uuid "github.com/jackc/pgtype/ext/gofrs-uuid"
	"testing"
)

/*type pgContainer struct {
	testcontainers.Container
}
*/
//var pgC *pgContainer = nil

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
