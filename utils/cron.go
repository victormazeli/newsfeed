package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/kamva/mgm/v3"
	"github.com/robfig/cron/v3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"net/http"
	"newsfeedbackend/config"
	"newsfeedbackend/database/models"
	"newsfeedbackend/graph/model"
	"sync"
)

var httpClient = resty.New()
var categories = []string{
	"business",
	"entertainment",
	"environment",
	"food",
	"health",
	"politics",
	"science",
	"sports",
	"technology",
	"top",
	"tourism",
	"world",
}

func InitCron(ctx context.Context, env *config.Env) {
	c := cron.New(cron.WithSeconds())
	jobID, err := c.AddFunc("@daily", GetHeadlineFunc(ctx, env))
	if err != nil {
		log.Printf("error :%v", err)
		return
	}
	log.Printf("jobID : %v", jobID)
	c.Start()

}

func GetHeadlineFunc(ctx context.Context, env *config.Env) func() {
	return func() {
		fetchHeadlineNews(ctx, env)
	}
}
func fetchHeadlineNews(ctx context.Context, env *config.Env) {
	//var page = "1"
	//var pageSize = "10"
	var wg sync.WaitGroup

	for _, category := range categories {
		wg.Add(1)
		go func(category string) {
			defer wg.Done()
			res, err := httpClient.R().
				SetQueryString(fmt.Sprintf("apiKey=%v&category=%v&language=%s", env.NewsApiKey, category, "en")).
				SetHeader("Accept", "application/json").
				Get(env.NewsApiBaseUrl + "/news")

			if err != nil {
				log.Printf("Error fetching news for category %s: %v", category, err)
				return
			}

			if res.StatusCode() == http.StatusOK {
				var response model.Response
				if err := json.Unmarshal(res.Body(), &response); err != nil {
					log.Printf("Error decoding news response for category %s: %v", category, err)
					return
				}

				var bulkOps []mongo.WriteModel
				var result models.News
				for _, v := range response.Results {
					article := models.News{
						ImageURL:    v.ImageURL,
						Description: v.Description,
						Title:       v.Title,
						Creator:     v.Creator,
						Content:     v.Content,
						PubDate:     v.PubDate,
						Link:        v.Link,
						SourceID:    v.SourceID,
						Category:    v.Category,
						Likes:       []*string{},
					}
					collErr := mgm.Coll(&models.News{}).FindOne(ctx, bson.M{"title": article.Title}).Decode(&result)
					if collErr != nil && collErr != mongo.ErrNoDocuments {
						log.Printf("Error finding news article: %v", collErr)
						continue
					}
					if collErr == mongo.ErrNoDocuments {
						insertModel := mongo.NewInsertOneModel().SetDocument(article)
						bulkOps = append(bulkOps, insertModel)
					}
				}

				if len(bulkOps) > 0 {
					_, collErr := mgm.Coll(&models.News{}).BulkWrite(ctx, bulkOps)
					if collErr != nil {
						log.Printf("Error inserting news articles: %v", collErr)
					}
				}
			} else {
				var errorResponse model.ErrorResponse
				if jsonErr := json.Unmarshal(res.Body(), &errorResponse); jsonErr != nil {
					log.Printf("Error decoding error response for category %s: %v", category, jsonErr)
					return
				}
				log.Printf("API error for category %s: %v", category, *errorResponse.Results.Message)
			}
		}(category)
	}

	wg.Wait()
}
