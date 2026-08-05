package request

import (
	"errors"
)

type UserRegisterRequest struct {
	Email           string `json:"email" validate:"required,email,max=255"`
	Password        string `json:"password" validate:"required,min=8,max=72"`
	ConfirmPassword string `json:"confirm_password" validate:"required,min=8,max=72"`
	FirstName       string `json:"first_name" validate:"required,max=100"`
	LastName        string `json:"last_name" validate:"required,max=100"`
	Phone           string `json:"phone" validate:"required,max=30"`
}

func (req *UserRegisterRequest) Validate() error {
	if req.Password != req.ConfirmPassword {
		return errors.New("passwords do not match")
	}
	return nil
}

//func (req *UserRegisterRequest) NewUser() *model.User {
//	return &model.User{
//		Email:     req.Email,
//		FirstName: req.FirstName,
//		LastName:  req.LastName,
//		Phone:     req.Phone,
//	}
//}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email,max=255"`
	Password string `json:"password" validate:"required,min=8,max=72"`
}
