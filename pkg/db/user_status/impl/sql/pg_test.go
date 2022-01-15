package sql

import (
	"github.com/GoncharovMikhail/go-sql/const/test"
	"github.com/GoncharovMikhail/go-sql/pkg/db/user_status"
	"github.com/GoncharovMikhail/go-sql/pkg/db/util"
	"github.com/GoncharovMikhail/go-sql/pkg/entity"
	"github.com/gofrs/uuid"
	"log"
	"testing"
)

//lateInit
var (
	entityToSave *entity.UserDataEntity
	uuidUsername uuid.UUID
	repository   user_status.SQLUserStatusRepository
)

func init() {
	var err error
	uuidUsername, err = uuid.NewV1()
	if err != nil {
		log.Panic(err)
	}
	entityToSave = &entity.UserDataEntity{
		Username: uuidUsername.String(),
		Password: uuidUsername.String(),
	}
	defer util.MustCloseDb(test.DB)
	repository = NewPostgresUserStatusRepository()
}

func TestPostgresUserStatusRepository_FindOneByUsernameInTx(t *testing.T) {

}
