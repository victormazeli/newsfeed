package models

import (
	"github.com/kamva/mgm/v3"
	"time"
)

type User struct {
	mgm.DefaultModel `bson:",inline"`
	Email            string    `json:"email"`
	Picture          string    `json:"picture"`
	FullName         string    `json:"full_name"`
	Password         string    `json:"password"`
	Otp              string    `json:"otp"`
	OtpExpireTime    time.Time `json:"otp_expire_time"`
	IsVerified       bool      `json:"is_verified"`
	IsOtpVerified    bool      `json:"is_otp_verified"`
	IsPasswordReset  bool      `json:"is_password_reset"`
	Topics           []Topic
	UpdatedAt        string `json:"updated_at" bson:"updated_at"`
	CreatedAt        string `json:"created_at" bson:"created_at"`
}

type Topic struct {
	Topic string `json:"topic"`
}
