package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson"
	"log"
	"net/http"
	"newsfeedbackend/config"
	"newsfeedbackend/database/models"
	"newsfeedbackend/graph/model"
	"newsfeedbackend/redis"
	"newsfeedbackend/utils"
	"time"
)

var httpClient = resty.New()

type Handler struct{}

// NewUser register new user
func (h Handler) NewUser(input model.CreateUser, env *config.Env) (*models.User, error) {
	user := &models.User{}

	e := mgm.Coll(user).First(bson.M{"email": input.Email}, user)

	if e != nil {
		var topics []models.Topic
		topicsArray := input.Topics
		for i := 0; i < len(topicsArray); i++ {
			topics = append(topics, models.Topic{
				Topic: *topicsArray[i],
			})
		}
		hashPassword, e := utils.GenerateFromPassword(input.Password)
		if e != nil {
			errMessage := errors.New("an error occurred")
			return nil, errMessage
		}
		otp := utils.GenerateOTP()
		expireTime := time.Now().Add(10 * time.Minute)
		newUser := models.User{
			Email:         input.Email,
			FullName:      input.FullName,
			IsVerified:    false,
			Otp:           otp,
			OtpExpireTime: expireTime,
			Password:      hashPassword,
			Picture:       input.Picture,
			Topics:        topics,
		}
		err := mgm.Coll(&models.User{}).Create(&newUser)

		if err != nil {
			errMessage := errors.New("an error occurred")
			return nil, errMessage
		}

		subject := "Registration"
		body := utils.GenerateOTPEmailTemplate(otp)
		res, _ := utils.SendEmail(env, input.Email, body, subject)
		log.Print("email sent successfully", res)

		return &newUser, nil
	}
	errMessage := errors.New("user already exist")
	return nil, errMessage

}

// Login login
func (h Handler) Login(input model.Login, env *config.Env) (*model.LoginResponse, error) {
	user := &models.User{}

	err := mgm.Coll(user).First(bson.M{"email": input.Email}, user)

	if err != nil {
		errMessage := errors.New("invalid credentials")
		return nil, errMessage
	}

	//if user.IsVerified == false {
	//	errMessage := errors.New("account not verified")
	//	return nil, errMessage
	//}

	passwordMatch, e := utils.ComparePasswordAndHash(input.Password, user.Password)

	if e != nil {
		errMessage := errors.New("an error occurred")
		return nil, errMessage
	}

	if passwordMatch == false {
		errMessage := errors.New("invalid credentials")
		return nil, errMessage
	}
	token := utils.GenerateToken(user.ID.Hex(), env.JwtKey)

	var topics []*model.Topic
	topicsArray := user.Topics
	for i := 0; i < len(topicsArray); i++ {
		topics = append(topics, &model.Topic{
			Topic: topicsArray[i].Topic,
		})
	}

	response := &model.LoginResponse{
		User: &model.User{
			Email:           user.Email,
			FullName:        user.FullName,
			ID:              user.ID.Hex(),
			IsPasswordReset: user.IsPasswordReset,
			IsOtpVerified:   user.IsOtpVerified,
			IsVerified:      user.IsVerified,
			Topics:          topics,
			UpdatedAt:       user.UpdatedAt,
			CreatedAt:       user.CreatedAt,
			Picture:         user.Picture,
		},
		Token: token,
	}

	return response, nil

}

// ForgotPassword user forgot password
func (h Handler) ForgotPassword(input model.ForgotPassword) (*model.GenericResponse, error) {
	user := &models.User{}

	err := mgm.Coll(user).First(bson.M{"email": input.Email}, user)

	if err != nil {
		errMessage := errors.New("user not found")
		return nil, errMessage
	}
	otp := utils.GenerateOTP()
	expireTime := time.Now().Add(10 * time.Minute)

	er := mgm.Coll(user).Update(&models.User{Otp: otp, OtpExpireTime: expireTime, IsPasswordReset: true})

	if er != nil {
		errMessage := errors.New("an error occurred")
		return nil, errMessage
	}

	response := &model.GenericResponse{
		Message: "password reset initiated",
	}
	return response, nil

}

func (h Handler) ResetPassword(input model.ResetPassword) (*model.GenericResponse, error) {
	user := &models.User{}

	err := mgm.Coll(user).First(bson.M{"email": input.Email}, user)

	if err != nil {
		errMessage := errors.New("password reset failed")
		return nil, errMessage
	}
	hasPassword, e := utils.GenerateFromPassword(input.NewPassword)

	if e != nil {
		errMessage := errors.New("an error occurred")
		return nil, errMessage
	}
	er := mgm.Coll(user).Update(&models.User{IsPasswordReset: false, Password: hasPassword})

	if er != nil {
		errMessage := errors.New("an error occurred")
		return nil, errMessage
	}

	response := &model.GenericResponse{
		Message: "password reset successful",
	}

	return response, nil

}

func (h Handler) VerifyEmail(input model.VerifyOtp) (*model.GenericResponse, error) {
	user := &models.User{}

	err := mgm.Coll(user).First(bson.M{"otp": input.Otp}, user)

	if err != nil {
		errMessage := errors.New("invalid token")
		return nil, errMessage
	}

	if user.IsVerified == true {
		errMessage := errors.New("account already verified")
		return nil, errMessage
	}

	// check if otp expired
	if user.OtpExpireTime.After(time.Now()) {
		errMessage := errors.New("otp expired")
		return nil, errMessage
	}

	er := mgm.Coll(user).Update(&models.User{IsVerified: true, Otp: ""})

	if er != nil {
		errMessage := errors.New("an error occurred")
		return nil, errMessage
	}

	response := &model.GenericResponse{
		Message: "email verification successful",
	}

	return response, nil
}

func (h Handler) VerifyResetOtp(input model.VerifyOtp) (*model.GenericResponse, error) {
	user := &models.User{}

	err := mgm.Coll(user).First(bson.M{"otp": input.Otp}, user)

	if err != nil {
		errMessage := errors.New("invalid token")
		return nil, errMessage
	}

	if user.IsPasswordReset == false {
		errMessage := errors.New("reset password not initiated")
		return nil, errMessage
	}

	// check if otp expired
	if user.OtpExpireTime.After(time.Now()) {
		errMessage := errors.New("otp expired")
		return nil, errMessage
	}

	er := mgm.Coll(user).Update(&models.User{IsVerified: true, Otp: ""})

	if er != nil {
		errMessage := errors.New("an error occurred")
		return nil, errMessage
	}

	response := &model.GenericResponse{
		Message: "otp verification successful",
	}

	return response, nil
}

func (h Handler) GetUser(userId string) *models.User {
	user := &models.User{}

	err := mgm.Coll(user).FindByID(userId, user)

	if err != nil {
		return nil
	}
	return user
}

func (h Handler) GetUserByEmail(email string) *models.User {
	user := &models.User{}

	err := mgm.Coll(user).First(bson.M{"email": email}, user)

	if err != nil {
		return nil
	}
	return user
}

func (h Handler) GetUserById(id interface{}) *models.User {
	user := &models.User{}

	err := mgm.Coll(user).FindByID(id, user)

	if err != nil {
		return nil
	}
	return user
}

func (h Handler) FetchNews(query string, env *config.Env, ctx context.Context) ([]*model.Article, error) {
	// check redis if data present
	result := redis.NewsCacheService{}.GetNews(ctx, query)

	if result != nil {
		return result, nil
	} else {
		res, err := httpClient.R().SetQueryString(fmt.Sprintf("q=%s&apiKey=%s", query, env.NewsApiKey)).SetHeader("Accept", "application/json").Get(env.NewsApiBaseUrl + "/everything")
		var response model.Response
		if err != nil {
			log.Print(err)
			return nil, err
		}
		if res.StatusCode() == http.StatusOK {
			err := json.Unmarshal(res.Body(), &response)
			if err != nil {
				return nil, err
			}
			var articleResponse []*model.Article
			for _, v := range response.Articles {
				article := model.Article{
					URL:         v.URL,
					Description: v.Description,
					Title:       v.URL,
					Author:      v.Author,
					Content:     v.Content,
					PublishedAt: v.PublishedAt,
					URLToImage:  v.URLToImage,
					Source:      v.Source,
				}
				articleResponse = append(articleResponse, &article)
			}
			redis.NewsCacheService{}.SetNews(ctx, query, articleResponse)
			return articleResponse, nil
		} else {
			err := errors.New("unable to fetch news")
			return nil, err
		}

	}

}
func (h Handler) NewsFeed(env *config.Env, id interface{}, ctx context.Context) ([]*model.Article, error) {
	user := &models.User{}

	err := mgm.Coll(user).FindByID(id, user)

	if err != nil {
		er := errors.New("user not found")
		return nil, er
	}

	var articleResponse []*model.Article

	for _, v := range user.Topics {
		article, er := h.FetchNews(v.Topic, env, ctx)
		if er != nil {
			break
		}
		articleResponse = append(articleResponse, article...)

	}
	if articleResponse != nil {
		return articleResponse, nil
	}
	er := errors.New("no data found")

	return nil, er

}

//func (h Handler) AskChatGPT() {
//	c := openai.NewClient(env.OpenAiKey)
//	//ctx := context.Background()
//
//	req := openai.CompletionRequest{
//		Model:            openai.GPT3TextDavinci003,
//		MaxTokens:        100,
//		Temperature:      0.3,
//		TopP:             1.0,
//		FrequencyPenalty: 0.0,
//		PresencePenalty:  0.0,
//		Prompt:           fmt.Sprintf("Translate from %s to %s: %s", input.LanguageFrom, input.LanguageTo, input.Text),
//	}
//	resp, err := c.CreateCompletion(ctx, req)
//	if err != nil {
//		fmt.Printf("Completion error: %v\n", err)
//		err := errors.New("operation failed")
//		return nil, err
//	}
//	translatedText := resp.Choices[0].Text
//	cleanedText := strings.ReplaceAll(translatedText, "\n", "")
//	response := &model.TranslationResponse{
//		TranslatedText: &cleanedText,
//	}
//	//fmt.Println(resp.Choices[0].Text)
//
//	return response, nil
//}
