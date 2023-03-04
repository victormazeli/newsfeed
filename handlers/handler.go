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
)

var httpClient = resty.New()

type Handler struct{}

func (h Handler) NewUser(input model.CreateUser) *models.User {
	var topics []models.Topic
	topicsArray := input.Topics
	for i := 0; i < len(topicsArray); i++ {
		topics = append(topics, models.Topic{
			Topic: *topicsArray[i],
		})
	}
	newUser := models.User{
		Email:    input.Email,
		FullName: input.FullName,
		UserId:   input.UserID,
		Picture:  input.Picture,
		Topics:   topics,
	}
	err := mgm.Coll(&models.User{}).Create(&newUser)

	if err != nil {
		return nil
	}

	return &newUser

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

func (h Handler) GetUserByAuth0Id(id string) *models.User {
	user := &models.User{}

	err := mgm.Coll(user).First(bson.M{"userid": id}, user)

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
func (h Handler) NewsFeed(env *config.Env, email string, ctx context.Context) ([]*model.Article, error) {
	user := &models.User{}

	err := mgm.Coll(user).First(bson.M{"email": email}, user)

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
