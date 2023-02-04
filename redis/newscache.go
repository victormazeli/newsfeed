package redis

import (
	"context"
	"encoding/json"
	"github.com/go-redis/redis/v8"
	"log"
	"newsfeedbackend/config"
	"newsfeedbackend/graph/model"
)

var redisClient *redis.Client = nil

type NewsCacheService struct{}

func (n NewsCacheService) Setup(env *config.Env) {
	redisClient = redis.NewClient(&redis.Options{
		Addr:     env.RedisAddr,
		Password: "", // no password set
		DB:       0,  // use default DB
	})

}

func (n NewsCacheService) SetNews(ctx context.Context, topic string, articles []*model.Article) {
	jsonStr, err := json.Marshal(articles)
	if err != nil {
		log.Printf("Error converting article to json, %s", err)
		panic(err)
	}
	err = redisClient.Set(ctx, topic, jsonStr, 10800).Err()
	if err != nil {
		panic(err)
	}

}

func (n NewsCacheService) GetNews(ctx context.Context, topic string) []*model.Article {
	value, err := redisClient.Get(ctx, topic).Result()
	if err != nil {
		return nil
	}
	var article []*model.Article
	err = json.Unmarshal([]byte(value), &article)
	if err != nil {
		log.Printf("Error converting json string to article object, %s", err)
		return nil
	}
	return article

}
