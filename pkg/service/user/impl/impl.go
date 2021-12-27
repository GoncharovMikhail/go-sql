package impl

import (
	"context"
	"database/sql"
	model "sql/model/user"
	db "sql/pkg/db/user"
)

type UserServiceImpl struct {
	Ur db.UserRepository
}

func (u *UserServiceImpl) Save(request *model.UserSaveRequest) (*db.UserEntity, error) {
	var restoreData db.RestoreData
	if request.Email != nil {
		restoreData = db.RestoreData{
			Email: *request.Email,
		}
		if request.PhoneNumber != nil {
			restoreData.PhoneNumber = sql.NullString{
				String: *request.PhoneNumber,
			}
		}
	}
	entity := &db.UserEntity{
		Username:    request.Username,
		Password:    request.Password,
		RestoreData: &restoreData,
	}
	save, err := u.Ur.Save(context.Background(), entity)
	if err != nil {
		return nil, err
	}
	return save, nil
}

func (u *UserServiceImpl) FindOneByUsername(username string) (*db.UserEntity, bool) {
	byUsername, err := u.Ur.FindOneByUsername(context.Background(), username)
	if err != nil {
		return nil, false
	}
	return byUsername, true
}
