package models

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID
	Phone        string
	BirthdayDate string
	DefaultTag   string
	Token        Token
	SmsCode      *Code
}

type SignUpByPhoneDTO struct {
	Phone        string
	BirthdayDate string
	DefaultTag   string
}

func (d *SignUpByPhoneDTO) Validate() error {
	return validation.ValidateStruct(
		d,
		validation.Field(&d.Phone, validation.Required),
		validation.Field(&d.BirthdayDate, validation.Required),
		validation.Field(&d.DefaultTag, validation.Required),
	)
}

type TempUser struct {
	ID          uuid.UUID
	Phone       string
	Code        string
	IsConfirmed bool
}
