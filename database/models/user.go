package models

import (
	"github.com/kamva/mgm/v3"
	"time"
)

type User struct {
	mgm.DefaultModel `bson:",inline"`
	Email            string    `json:"email" bson:"email"`
	Picture          string    `json:"picture" bson:"picture"`
	FullName         string    `json:"full_name" bson:"full_name"`
	Password         string    `json:"password" bson:"password"`
	PhoneNumber      string    `json:"phone_number" bson:"phone_number"`
	Otp              string    `json:"otp" bson:"otp"`
	OtpExpireTime    time.Time `json:"otp_expire_time" bson:"otp_expire_time"`
	IsVerified       bool      `json:"is_verified" bson:"is_verified"`
	IsOtpVerified    bool      `json:"is_otp_verified" bson:"is_otp_verified"`
	IsPasswordReset  bool      `json:"is_password_reset" bson:"is_password_reset"`
	Topics           []*string `json:"topics" bson:"topics"`
	CreatedAt        time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt        time.Time `json:"updated_at" bson:"updated_at"`
}

//type Topic struct {
//	Topic string `json:"topic"`
//}
