package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/kamva/mgm/v3"
	"github.com/sashabaranov/go-openai"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"net/http"
	"net/url"
	"newsfeedbackend/config"
	"newsfeedbackend/database/models"
	"newsfeedbackend/graph/model"
	"newsfeedbackend/redis"
	"newsfeedbackend/utils"
	"strings"
	"time"
)

var httpClient = resty.New()

type Handler struct{}

// NewUser register new user
func (h Handler) NewUser(input model.CreateUser, env *config.Env) (*models.User, error) {
	user := &models.User{}

	e := mgm.Coll(user).First(bson.M{"email": input.Email}, user)

	if e != nil {
		//var topics []models.Topic
		//topicsArray := input.Topics
		//for i := 0; i < len(topicsArray); i++ {
		//	topics = append(topics, models.Topic{
		//		Topic: *topicsArray[i],
		//	})
		//}
		hashPassword, e := utils.GenerateFromPassword(input.Password)
		if e != nil {
			errMessage := errors.New("an error occurred")
			return nil, errMessage
		}
		otp := utils.GenerateOTP()
		currentTime := getCurrentTime()
		expireTime := currentTime.Add(10 * time.Minute)
		log.Printf("time :%v", expireTime)
		newUser := &models.User{
			Email:         input.Email,
			FullName:      input.FullName,
			IsVerified:    false,
			Otp:           otp,
			OtpExpireTime: expireTime,
			Topics:        []*string{},
			Password:      hashPassword,
			CreatedAt:     getCurrentTime(),
			UpdatedAt:     getCurrentTime(),
		}
		err := mgm.Coll(&models.User{}).Create(newUser)

		if err != nil {
			log.Printf("err :%v", err)
			errMessage := errors.New("an error occurred")
			return nil, errMessage
		}

		subject := "Registration"
		body := utils.GenerateOTPEmailTemplate(otp)
		res, _ := utils.SendEmail(env, input.Email, body, subject)
		log.Print("email sent successfully", res)

		return newUser, nil
	}
	errMessage := errors.New("user already exist")
	return nil, errMessage

}

// complete registration

func (h Handler) CompleteRegistration(input model.CompleteRegistration, id interface{}, ctx context.Context) (*model.GenericResponse, error) {
	var updateData []string
	topicsArray := input.Topics
	for i := 0; i < len(topicsArray); i++ {
		updateData = append(updateData, *topicsArray[i])
	}

	userID, err := convertToString(id)
	if err != nil {
		return nil, err
	}

	userIDString, err := convertToObjectId(userID)

	if err != nil {
		return nil, err
	}
	opts := options.Update().SetUpsert(false)
	filter := bson.D{{"_id", userIDString}}
	data := bson.D{
		{"$set", bson.D{
			{"topics", updateData},
		}},
	}

	_, err = mgm.Coll(&models.User{}).UpdateOne(ctx, filter, data, opts)

	if err != nil {
		log.Printf("db error : %v", err)
	}

	response := &model.GenericResponse{
		Message: "Registration complete",
	}

	return response, nil
}

// Login login
func (h Handler) Login(input model.Login, env *config.Env, ctx context.Context) (*model.LoginResponse, error) {
	user := &models.User{}

	err := mgm.Coll(user).First(bson.M{"email": input.Email}, user)

	if err != nil && err == mongo.ErrNoDocuments {
		errMessage := errors.New("invalid credentials")
		return nil, errMessage
	}

	if user.IsVerified == false {
		errMessage := errors.New("account not verified")
		return nil, errMessage
	}

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
	//
	//var topics []*model.T
	//topicsArray := user.Topics
	//for i := 0; i < len(topicsArray); i++ {
	//	topics = append(topics, &model.Topic{
	//		Topic: topicsArray[i].Topic,
	//	})
	//}

	var appTokens []string

	appTokens = append(appTokens, token)

	redis.NewsCacheService{}.SetAppToken(ctx, "appToken", appTokens)

	response := &model.LoginResponse{
		Token: token,
	}

	return response, nil

}

// ForgotPassword user forgot password
func (h Handler) ForgotPassword(input model.ForgotPassword, env *config.Env, ctx context.Context) (*model.GenericResponse, error) {
	user := &models.User{}

	err := mgm.Coll(user).First(bson.M{"email": input.Email}, user)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			errMessage := errors.New("user not found")
			return nil, errMessage
		}
		return nil, err
	}
	otp := utils.GenerateOTP()
	currentTime := getCurrentTime()
	expireTime := currentTime.Add(10 * time.Minute)
	opts := options.Update().SetUpsert(false)
	filter := bson.D{{"email", input.Email}}
	update := bson.D{
		{"$set", bson.D{
			{"otp", otp},
			{"otp_expire_time", expireTime},
			{"is_password_reset", true},
		}},
	}

	_, er := mgm.Coll(user).UpdateOne(ctx, filter, update, opts)

	if er != nil {
		errMessage := errors.New("an error occurred")
		return nil, errMessage
	}

	subject := "Password Reset"
	body := utils.GenerateForgetPasswordEmailTemplate(otp)
	res, _ := utils.SendEmail(env, input.Email, body, subject)
	log.Print("email sent successfully", res)

	response := &model.GenericResponse{
		Message: "password reset initiated",
	}
	return response, nil

}

func (h Handler) ResendOtp(email string, ctx context.Context, env *config.Env) (*model.GenericResponse, error) {
	user := &models.User{}

	err := mgm.Coll(user).First(bson.M{"email": email}, user)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			errMessage := errors.New("user not found")
			return nil, errMessage
		}
		return nil, err
	}

	otp := utils.GenerateOTP()
	currentTime := getCurrentTime()
	expireTime := currentTime.Add(10 * time.Minute)
	opts := options.Update().SetUpsert(false)
	filter := bson.D{{"email", email}}
	update := bson.D{
		{"$set", bson.D{
			{"otp", otp},
			{"otp_expire_time", expireTime},
		}},
	}

	_, er := mgm.Coll(user).UpdateOne(ctx, filter, update, opts)

	if er != nil {
		errMessage := errors.New("an error occurred")
		return nil, errMessage
	}

	if user.IsPasswordReset == true {
		subject := "Password Reset"
		body := utils.GenerateForgetPasswordEmailTemplate(otp)
		res, _ := utils.SendEmail(env, email, body, subject)
		log.Print("email sent successfully", res)

	} else {
		subject := "Registration"
		body := utils.GenerateOTPEmailTemplate(otp)
		res, _ := utils.SendEmail(env, email, body, subject)
		log.Print("email sent successfully", res)

	}

	response := &model.GenericResponse{
		Message: "otp resent",
	}
	return response, nil

}

func (h Handler) ResetPassword(input model.ResetPassword, ctx context.Context) (*model.GenericResponse, error) {
	user := &models.User{}

	err := mgm.Coll(user).First(bson.M{"email": input.Email}, user)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			errMessage := errors.New("user not found")
			return nil, errMessage
		}
		return nil, err
	}
	hasPassword, e := utils.GenerateFromPassword(input.NewPassword)

	if e != nil {
		errMessage := errors.New("an error occurred")
		return nil, errMessage
	}
	opts := options.Update().SetUpsert(false)
	filter := bson.D{{"email", input.Email}}
	update := bson.D{
		{"$set", bson.D{
			{"password", hasPassword},
			{"is_password_reset", false},
		}},
	}

	_, er := mgm.Coll(user).UpdateOne(ctx, filter, update, opts)

	if er != nil {
		errMessage := errors.New("an error occurred")
		return nil, errMessage
	}

	response := &model.GenericResponse{
		Message: "password reset successful",
	}

	return response, nil

}

func (h Handler) VerifyEmail(input model.VerifyOtp, ctx context.Context) (*model.GenericResponse, error) {
	user := &models.User{}

	err := mgm.Coll(user).First(bson.M{"otp": input.Otp}, user)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			errMessage := errors.New("invalid token")
			return nil, errMessage
		}
		return nil, err
	}

	if user.IsVerified == true {
		errMessage := errors.New("account already verified")
		return nil, errMessage
	}

	// check if otp expired
	currentTime := getCurrentTime()
	if user.OtpExpireTime.After(currentTime) == false {
		errMessage := errors.New("otp expired")
		return nil, errMessage
	}

	opts := options.Update().SetUpsert(false)
	filter := bson.D{{"email", user.Email}}
	update := bson.D{
		{"$set", bson.D{
			{"is_verified", true},
			{"otp", ""},
			{"otp_expire_time", nil},
		}},
	}

	_, er := mgm.Coll(&models.User{}).UpdateOne(ctx, filter, update, opts)

	if er != nil {
		errMessage := errors.New("an error occurred")
		return nil, errMessage
	}

	response := &model.GenericResponse{
		Message: "email verification successful",
	}

	return response, nil
}

func (h Handler) VerifyResetOtp(input model.VerifyOtp, ctx context.Context) (*model.GenericResponse, error) {
	user := &models.User{}

	err := mgm.Coll(user).First(bson.M{"otp": input.Otp}, user)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			errMessage := errors.New("invalid token")
			return nil, errMessage
		}
		return nil, err
	}

	if user.IsPasswordReset == false {
		errMessage := errors.New("reset password not initiated")
		return nil, errMessage
	}

	// check if otp expired
	currentTime := getCurrentTime()
	if user.OtpExpireTime.After(currentTime) == false {
		errMessage := errors.New("otp expired")
		return nil, errMessage
	}

	opts := options.Update().SetUpsert(false)
	filter := bson.D{{"email", user.Email}}
	update := bson.D{
		{"$set", bson.D{
			{"is_password_reset", false},
			{"otp", ""},
			{"otp_expire_time", nil},
		}},
	}

	_, er := mgm.Coll(&models.User{}).UpdateOne(ctx, filter, update, opts)

	if er != nil {
		errMessage := errors.New("an error occurred")
		return nil, errMessage
	}

	response := &model.GenericResponse{
		Message: "otp verification successful",
	}

	return response, nil
}

func (h Handler) GoogleLogin(input model.GoogleAuth, env *config.Env, ctx context.Context) (*model.LoginResponse, error) {
	var payload model.GoogleAuthModel
	url := "https://www.googleapis.com/oauth2/v3/userinfo"

	res, err := httpClient.R().SetHeader("Accept", "application/json").SetAuthToken(input.AccessToken).Get(url)

	if err != nil {
		return nil, err
	}
	if res.StatusCode() == http.StatusOK {
		err := json.Unmarshal(res.Body(), &payload)
		if err != nil {
			return nil, err
		}

		user := &models.User{}

		er := mgm.Coll(user).First(bson.M{"email": payload.Email}, user)

		if er != nil {
			// create user
			newUser := models.User{
				Email:      *payload.Email,
				FullName:   *payload.Name,
				IsVerified: true,
				Password:   "",
			}
			err := mgm.Coll(&models.User{}).Create(&newUser)

			if err != nil {
				errMessage := errors.New("an error occurred")
				return nil, errMessage
			}
			token := utils.GenerateToken(user.ID.Hex(), env.JwtKey)

			var appTokens []string

			appTokens = append(appTokens, token)

			redis.NewsCacheService{}.SetAppToken(ctx, "appToken", appTokens)

			response := &model.LoginResponse{
				Token: token,
			}
			return response, nil
		}

		token := utils.GenerateToken(user.ID.Hex(), env.JwtKey)

		var appTokens []string

		appTokens = append(appTokens, token)

		redis.NewsCacheService{}.SetAppToken(ctx, "appToken", appTokens)

		response := &model.LoginResponse{
			Token: token,
		}

		return response, nil

	} else {
		err := errors.New("operation failed")
		return nil, err
	}

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

func (h Handler) NewsFeed(id interface{}, query model.NewsQuery, ctx context.Context) ([]*model.Article, error) {
	user := &models.User{}
	var news []models.News
	var categories []string
	var filter bson.M
	skip := (*query.Page - 1) * *query.PageSize

	err := mgm.Coll(user).FindByID(id, user)
	if err != nil {
		return nil, errors.New("user not found")
	}

	if query.Category != nil {
		categories = append(categories, *query.Category)
	} else {
		for _, topic := range user.Topics {
			if topic != nil {
				categories = append(categories, *topic)
			}
		}
	}

	if len(categories) == 0 {
		//log.Print("log ...")
		categories = append(categories, "top")
	}

	if query.Source == nil {
		filter = bson.M{
			"category": bson.M{"$in": categories},
		}
	} else {
		//log.Print("log 1 ...")
		//log.Printf(" categories :%v", categories)
		//log.Printf("source : %v", *query.Source)
		filter = bson.M{
			"category":  bson.M{"$in": categories},
			"source_id": query.Source,
		}
	}

	cursor, err := mgm.Coll(&models.News{}).Find(ctx, filter, options.Find().SetSkip(int64(skip)).SetLimit(int64(*query.PageSize)))
	if err != nil {
		return nil, err
	}

	if er := cursor.All(ctx, &news); err != nil {
		return nil, er
	}

	articleResponse := make([]*model.Article, 0, len(news))
	for _, newsData := range news {
		article := &model.Article{
			ID:          newsData.ID,
			Title:       newsData.Title,
			SourceID:    newsData.SourceID,
			Link:        newsData.Link,
			PubDate:     newsData.PubDate,
			Creator:     newsData.Creator,
			Category:    newsData.Category,
			ImageURL:    newsData.ImageURL,
			Description: newsData.Description,
			Content:     newsData.Content,
		}

		if article != nil {
			articleResponse = append(articleResponse, article)
		}

	}

	return articleResponse, nil
}
func (h Handler) AskChatGPT(input model.PromptContent, env *config.Env, ctx context.Context) (*model.PromptResponse, error) {
	c := openai.NewClient(env.OpenAiKey)
	prompt := "analyze this text in detail " + *input.Content

	var finalResponse string
	for {
		req := openai.CompletionRequest{
			Model:            openai.GPT3TextDavinci003,
			MaxTokens:        100,
			Temperature:      0.5,
			TopP:             1.0,
			FrequencyPenalty: 0.0,
			PresencePenalty:  0.0,
			Prompt:           prompt,
		}
		resp, err := c.CreateCompletion(ctx, req)
		if err != nil {
			log.Printf("Completion error: %v\n", err)
			return nil, err
		}

		if len(resp.Choices) > 0 {
			translatedText := resp.Choices[0].Text
			cleanedText := strings.ReplaceAll(translatedText, "\n", "")
			finalResponse += cleanedText
			// Check if the response is complete
			if resp.Choices[0].FinishReason == "stop" {
				break // End the conversation if the response is complete
			} else {
				prompt = finalResponse // Update the prompt with the accumulated response
				time.Sleep(5 * time.Second)
			}
		} else {
			log.Print("Empty response received from the API.")
		}
	}

	response := &model.PromptResponse{
		Result: &finalResponse,
	}

	return response, nil
}

func (h Handler) SeedSources(ctx context.Context, env *config.Env) ([]*model.Source, error) {
	var response model.SourceResponse
	var errorResponse model.ErrorResponse
	var sourceResult []*model.Source
	var bulkOps []mongo.WriteModel

	res, err := httpClient.R().
		SetQueryString(fmt.Sprintf("apiKey=%v&language=%s", env.NewsApiKey, "en")).
		SetHeader("Accept", "application/json").
		Get(env.NewsApiBaseUrl + "/sources")

	if err != nil {
		log.Printf("Error fetching sources: %v", err)
		return nil, err
	}

	if res.StatusCode() == http.StatusOK {
		if er := json.Unmarshal(res.Body(), &response); er != nil {
			log.Printf("Error decoding sources: %v", er)
			return nil, er
		}

		_, dbErr := mgm.Coll(&models.Source{}).DeleteMany(ctx, bson.M{})
		if dbErr != nil {
			log.Printf("Failed to clear collection: %v", dbErr)
			return nil, dbErr
		}

		for _, result := range response.Results {
			foundLogo, err := findLogo(*result.Name, *result.URL)
			if err != nil || len(foundLogo) == 0 || foundLogo[0] == nil {
				continue
			}

			source := models.Source{
				Category: result.Category,
				Url:      result.URL,
				Name:     result.Name,
				ID:       result.ID,
				Icon:     foundLogo[0].Icon,
			}

			insertModel := mongo.NewInsertOneModel().SetDocument(source)
			bulkOps = append(bulkOps, insertModel)
		}

		if len(bulkOps) > 0 {
			_, collErr := mgm.Coll(&models.Source{}).BulkWrite(ctx, bulkOps)
			if collErr != nil {
				log.Printf("Error inserting news articles: %v", collErr)
			}
		}

		for _, result := range response.Results {
			source := &model.Source{
				ID:       result.ID,
				URL:      result.URL,
				Name:     result.Name,
				Category: result.Category,
				Icon:     result.Icon,
			}
			sourceResult = append(sourceResult, source)
		}

		return sourceResult, nil
	} else {
		if er := json.Unmarshal(res.Body(), &errorResponse); er != nil {
			log.Printf("Error decoding sources: %v", er)
			return nil, er
		}

		resErr := errors.New(fmt.Sprintf("%v", errorResponse.Results.Message))
		return nil, resErr
	}
}

func (h Handler) FetchSources() ([]*model.Source, error) {
	var results []models.Source
	var sources []*model.Source

	dbErr := mgm.Coll(&models.Source{}).SimpleFind(&results, bson.M{})
	if dbErr != nil {
		log.Printf("Failed to fetch collection: %v", dbErr)
		return nil, dbErr
	}

	for _, result := range results {
		source := &model.Source{
			ID:       result.ID,
			Icon:     result.Icon,
			Name:     result.Name,
			Category: result.Category,
			URL:      result.Url,
		}

		sources = append(sources, source)
	}

	if len(sources) > 0 {
		return sources, nil
	}

	return nil, errors.New("no sources found")
}

func (h Handler) FetCategories() ([]*model.Category, error) {
	var results []models.Category
	var categories []*model.Category

	if err := mgm.Coll(&models.Category{}).SimpleFind(&results, bson.M{}); err != nil {
		return nil, err

	}

	for _, result := range results {
		category := &model.Category{
			ID:   result.ID,
			Name: result.Name,
		}

		categories = append(categories, category)

	}

	if len(categories) > 0 {
		return categories, nil

	} else {
		e := errors.New("no categories found")
		return nil, e
	}

}

func (h Handler) Logout(input model.Logout, ctx context.Context) (*model.GenericResponse, error) {

	result := redis.NewsCacheService{}.GetAppToken(ctx, "appToken")

	if result == nil {
		er := errors.New("operation failed")
		return nil, er
	}

	for i, token := range result {
		if token == input.Token {
			index := i
			result = append(result[:index], result[index+1:]...)
		}
	}
	redis.NewsCacheService{}.SetAppToken(ctx, "appToken", result)

	response := &model.GenericResponse{
		Message: "logout successful",
	}

	return response, nil

}

func (h Handler) FetchSavedNews(id interface{}, ctx context.Context) ([]*model.Article, error) {
	userIDString, err := convertToString(id)
	if err != nil {
		return nil, err
	}

	userID, err := convertToObjectId(userIDString)
	if err != nil {
		return nil, err
	}

	cursor, err := mgm.Coll(&models.SavedNews{}).Find(ctx, bson.M{"user_id": userID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var savedNewsList []models.SavedNews
	err = cursor.All(ctx, &savedNewsList)
	if err != nil {
		return nil, err
	}

	newsIDs := make([]primitive.ObjectID, len(savedNewsList))
	for i, savedNews := range savedNewsList {
		newsIDs[i] = savedNews.NewsID
	}

	newsCursor, err := mgm.Coll(&models.News{}).Find(ctx, bson.M{"_id": bson.M{"$in": newsIDs}})
	if err != nil {
		return nil, err
	}
	defer newsCursor.Close(ctx)

	var articles []*models.News
	err = newsCursor.All(ctx, &articles)
	if err != nil {
		return nil, err
	}

	var articleResponse []*model.Article
	for _, article := range articles {
		articleData := &model.Article{
			ID:          article.ID,
			Title:       article.Title,
			SourceID:    article.SourceID,
			Link:        article.Link,
			PubDate:     article.PubDate,
			Creator:     article.Creator,
			Category:    article.Category,
			ImageURL:    article.ImageURL,
			Description: article.Description,
			Content:     article.Content,
			Likes:       article.Likes,
		}
		articleResponse = append(articleResponse, articleData)
	}

	return articleResponse, nil
}

func (h Handler) SaveNews(newsID string, userID interface{}, ctx context.Context) (*bool, error) {
	var news models.News
	// check if news exist
	newsIDString, err := convertToObjectId(newsID)
	if err != nil {
		return nil, err
	}
	err = mgm.Coll(&models.News{}).First(bson.M{"_id": newsIDString}, &news)

	if err != nil && err == mongo.ErrNoDocuments {
		return nil, errors.New("news not found")
	}
	//check if saved, if saved remove is not add

	userIDString, err := convertToString(userID)

	if err != nil {
		return nil, err
	}

	userIDHex, err := convertToObjectId(userIDString)

	if err != nil {
		return nil, err
	}

	var savedNews models.SavedNews

	err = mgm.Coll(&models.SavedNews{}).First(bson.M{"news_id": newsIDString, "user_id": userIDHex}, &savedNews)

	if err != nil && err == mongo.ErrNoDocuments {
		// add
		addNewsToFavorite := models.SavedNews{
			NewsID: newsIDString,
			UserID: userIDHex,
		}
		err := mgm.Coll(&models.SavedNews{}).Create(&addNewsToFavorite)
		if err != nil {
			return nil, err
		}
		response := true
		return &response, nil
	}

	// delete
	res, err := mgm.Coll(&models.SavedNews{}).DeleteOne(ctx, bson.M{"news_id": newsIDString, "user_id": userIDHex})

	if err != nil {
		return nil, err
	}

	log.Printf("delete count : %v", res.DeletedCount)

	response := false
	return &response, nil

}

func (h Handler) GetNewsById(id string, userID interface{}) (*model.Article, error) {

	userIDString, err := convertToString(userID)

	if err != nil {
		return nil, err
	}

	userIDHex, err := convertToObjectId(userIDString)

	if err != nil {
		return nil, err
	}
	news := &models.News{}

	err = mgm.Coll(news).FindByID(id, news)

	if err != nil {
		return nil, err
	}

	userLike := containsValue(news.Likes, userIDHex.Hex())

	log.Printf("result :%v", userLike)

	article := &model.Article{
		ID:          news.ID,
		Title:       news.Title,
		SourceID:    news.SourceID,
		Link:        news.Link,
		PubDate:     news.PubDate,
		Creator:     news.Creator,
		Category:    news.Category,
		ImageURL:    news.ImageURL,
		Description: news.Description,
		Content:     news.Content,
		Likes:       news.Likes,
		IsLiked:     &userLike,
	}

	return article, nil
}

func (h Handler) LikeNews(userID interface{}, newsID string, ctx context.Context) (*bool, error) {
	newsIDString, err := convertToObjectId(newsID)

	if err != nil {
		return nil, err
	}

	userIDString, err := convertToString(userID)

	if err != nil {
		return nil, err
	}

	userIDHex, err := convertToObjectId(userIDString)

	if err != nil {
		return nil, err
	}
	// first check if user already liked news, if liked news remove if not like news
	result, err := checkValueExistenceInNews(userIDString, newsIDString, ctx)

	if err != nil {
		return nil, err
	}

	log.Printf("status :%v", result)

	if result == true {
		filter := bson.M{"_id": newsIDString}
		update := bson.M{"$pull": bson.M{"likes": userIDHex.Hex()}}

		_, err = mgm.Coll(&models.News{}).UpdateOne(ctx, filter, update)

		if err != nil {
			return nil, err
		}
		response := false
		return &response, nil
	}
	filter := bson.M{"_id": newsIDString}
	update := bson.M{"$push": bson.M{"likes": userIDHex.Hex()}}

	_, err = mgm.Coll(&models.News{}).UpdateOne(ctx, filter, update)

	if err != nil {
		return nil, err
	}

	response := true
	return &response, nil
}

func findLogo(brandName string, domain string) ([]*model.SourceLogo, error) {
	encodedBrandName := url.QueryEscape(brandName)

	url := domain
	trimmedDomain := strings.TrimPrefix(url, "https://")
	trimmedDomain = strings.TrimPrefix(trimmedDomain, "www.")

	var results []model.SourceLogo

	res, err := httpClient.R().
		SetHeader("Referer", fmt.Sprintf("%v", trimmedDomain)).
		SetHeader("Accept", "application/json").
		Get(fmt.Sprintf("https://api.brandfetch.io/v2/search/%v", encodedBrandName))

	if err != nil {
		log.Printf("Error fetching sources: %v", err)
		return nil, err
	}

	if res.StatusCode() != http.StatusOK {
		return nil, errors.New("failed to fetch logo")
	}

	if err := json.Unmarshal(res.Body(), &results); err != nil {
		log.Printf("Error decoding logo: %v", err)
		return nil, err
	}

	sourceLogo := make([]*model.SourceLogo, len(results))
	for _, result := range results {
		if *result.Domain == trimmedDomain && *result.Name == brandName {
			logo := &model.SourceLogo{
				Name:   result.Name,
				Domain: result.Domain,
				Icon:   result.Icon,
			}
			//log.Printf("icon %v", logo)
			sourceLogo = append(sourceLogo, logo)
		} else {
			//log.Printf("Skipping logo with nil or empty icon: %v", result)
			logo := &model.SourceLogo{
				Name:   result.Name,
				Domain: result.Domain,
				Icon:   nil,
			}
			//log.Printf("icon %v", logo)
			sourceLogo = append(sourceLogo, logo)
		}

	}

	if len(sourceLogo) > 0 {
		filteredSourceLogo := make([]*model.SourceLogo, 0, len(sourceLogo))
		for _, logo := range sourceLogo {
			if logo != nil {
				filteredSourceLogo = append(filteredSourceLogo, logo)
			}
		}

		return filteredSourceLogo, nil
	}

	return nil, errors.New("empty result")

}
func convertToObjectId(id string) (primitive.ObjectID, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return primitive.ObjectID{}, err
	}
	return objectID, nil
}

func convertToString(value interface{}) (string, error) {
	// Check if the underlying value is already a string
	if str, ok := value.(string); ok {
		return str, nil
	}

	// If the underlying value is not a string, convert it
	switch v := value.(type) {
	case []byte:
		return string(v), nil
	case fmt.Stringer:
		return v.String(), nil
	default:
		return "", fmt.Errorf("unable to convert %T to string", value)
	}
}

func checkValueExistenceInNews(value string, newsID primitive.ObjectID, ctx context.Context) (bool, error) {

	var result models.News

	err := mgm.Coll(&models.News{}).FindByID(newsID, &result)

	if err != nil {
		return false, err
	}

	if result.Likes == nil {
		return false, nil
	}

	filter := bson.M{"likes": bson.M{"$in": []string{value}}, "_id": newsID}

	count, err := mgm.Coll(&models.News{}).CountDocuments(ctx, filter)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func getCurrentTime() time.Time {
	//timeDiff := time.Hour * -1
	currentTime := time.Now().UTC()
	currentTimeZone, _ := time.Now().Zone()

	fmt.Println("Current Timezone:", currentTimeZone)

	return currentTime

}

func containsValue(slice []*string, value string) bool {
	for _, element := range slice {
		if *element == value {
			return true
		}
	}
	return false
}
