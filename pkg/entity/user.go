package entity

type UserEntity struct {
	*UserDataEntity
	*RestoreDataEntity
	*UserStatusEntity
}
