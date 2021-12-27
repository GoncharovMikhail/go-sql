package user

type UserSaveRequest struct {
	Username    string  `json:"username"`
	Password    string  `json:"password"`
	Email       *string `json:"email,omitempty"`
	PhoneNumber *string `json:"phone_number,omitempty"`
}
