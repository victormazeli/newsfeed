// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package model

type Article struct {
	ID          *string   `json:"id"`
	Creator     []*string `json:"creator"`
	Title       *string   `json:"title"`
	Description *string   `json:"description"`
	ImageURL    *string   `json:"image_url"`
	Link        *string   `json:"link"`
	SourceID    *string   `json:"source_id"`
	PubDate     *string   `json:"pubDate"`
	Content     *string   `json:"content"`
	Category    []*string `json:"category"`
	Likes       []*string `json:"likes"`
	IsLiked     *bool     `json:"isLiked"`
}

type Category struct {
	ID   *string `json:"id"`
	Name *string `json:"name"`
}

type CompleteRegistration struct {
	Topics []*string `json:"topics"`
}

type CreateUser struct {
	Email       string `json:"email"`
	FullName    string `json:"full_name"`
	PhoneNumber string `json:"phone_number"`
	Password    string `json:"password"`
}

type ErrorResponse struct {
	Status  *string    `json:"status"`
	Results *NewsError `json:"results"`
}

type ForgotPassword struct {
	Email string `json:"email"`
}

type GenericResponse struct {
	Message string `json:"message"`
}

type GoogleAuth struct {
	AccessToken string `json:"access_token"`
}

type GoogleAuthModel struct {
	Sub           *string `json:"sub"`
	Name          *string `json:"name"`
	GivenName     *string `json:"given_name"`
	FamilyName    *string `json:"family_name"`
	Email         *string `json:"email"`
	Picture       *string `json:"picture"`
	EmailVerified *bool   `json:"email_verified"`
	Locale        *string `json:"locale"`
}

type Login struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

type Logout struct {
	Token string `json:"token"`
}

type NewsError struct {
	Message *string `json:"message"`
	Code    *string `json:"code"`
}

type NewsQuery struct {
	Source   *string `json:"source"`
	Category *string `json:"category"`
	PageSize int     `json:"pageSize"`
	Page     int     `json:"page"`
}

type ResetPassword struct {
	Email       string `json:"email"`
	NewPassword string `json:"newPassword"`
}

type Response struct {
	Status       *string    `json:"status"`
	TotalResults *int       `json:"totalResults"`
	Results      []*Article `json:"results"`
}

type Source struct {
	ID       *string   `json:"id"`
	Name     *string   `json:"name"`
	URL      *string   `json:"url"`
	Category []*string `json:"category"`
	Icon     *string   `json:"icon"`
}

type SourceLogo struct {
	Name   *string `json:"name"`
	Domain *string `json:"domain"`
	Icon   *string `json:"icon"`
}

type SourceResponse struct {
	Status  *string   `json:"status"`
	Results []*Source `json:"results"`
}

type User struct {
	ID              string    `json:"_id"`
	Email           *string   `json:"email"`
	Picture         *string   `json:"picture"`
	FullName        *string   `json:"full_name"`
	Topics          []*string `json:"topics"`
	IsVerified      *bool     `json:"is_verified"`
	IsOtpVerified   *bool     `json:"is_otp_verified"`
	IsPasswordReset *bool     `json:"is_password_reset"`
	CreatedAt       string    `json:"created_at"`
	UpdatedAt       string    `json:"updated_at"`
}

type VerifyOtp struct {
	Otp string `json:"otp"`
}

type PromptContent struct {
	Content *string `json:"content"`
}

type PromptResponse struct {
	Result *string `json:"result"`
}
