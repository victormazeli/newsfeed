package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"errors"
	"newsfeedbackend/graph/generated"
	"newsfeedbackend/graph/model"
	"newsfeedbackend/handlers"
	"newsfeedbackend/middlewares"
)

// CreateNewUser is the resolver for the CreateNewUser field.
func (r *mutationResolver) CreateNewUser(ctx context.Context, input model.CreateUser) (*model.User, error) {
	checkUser := handlers.Handler{}.GetUserByEmail(input.Email)
	if checkUser != nil {
		err := errors.New("user already exist")
		return nil, err

	}
	newUser := handlers.Handler{}.NewUser(input)
	if newUser == nil {
		err := errors.New("unable to create user")
		return nil, err
	}

	var topicsFromDB []*model.Topic
	for _, t := range newUser.Topics {
		item := model.Topic{
			Topic: t.Topic,
		}
		topicsFromDB = append(topicsFromDB, &item)
	}

	user := &model.User{
		Email:  newUser.Email,
		UserID: newUser.UserId,
		Topics: topicsFromDB,
		ID:     newUser.ID.String(),
	}
	return user, nil
}

// GetUser is the resolver for the GetUser field.
func (r *queryResolver) GetUser(ctx context.Context) (*model.User, error) {
	gc, err := middlewares.GinContextFromContext(ctx)
	if err != nil {
		return nil, err
	}
	sub, er := middlewares.Auth(gc, r.Env)
	if er != nil {
		return nil, er
	}
	id := sub.UserID
	getUser := handlers.Handler{}.GetUserByAuth0Id(id)

	if getUser == nil {
		err := errors.New("unable to fetch user")
		return nil, err
	}

	var topicsFromDB []*model.Topic
	for _, t := range getUser.Topics {
		item := model.Topic{
			Topic: t.Topic,
		}
		topicsFromDB = append(topicsFromDB, &item)
	}

	user := &model.User{
		Email:     getUser.Email,
		UserID:    getUser.UserId,
		Picture:   getUser.Picture,
		FullName:  getUser.FullName,
		ID:        getUser.ID.String(),
		Topics:    topicsFromDB,
		UpdatedAt: getUser.UpdatedAt,
		CreatedAt: getUser.CreatedAt,
	}

	return user, nil
}

// GetNews is the resolver for the GetNews field.
func (r *queryResolver) GetNews(ctx context.Context, query string) ([]*model.Article, error) {
	getNews, err := handlers.Handler{}.FetchNews(query, r.Env, ctx)
	if err != nil {
		return nil, err
	}
	return getNews, nil
}

// NewsFeed is the resolver for the NewsFeed field.
func (r *queryResolver) NewsFeed(ctx context.Context) ([]*model.Article, error) {
	gc, err := middlewares.GinContextFromContext(ctx)
	if err != nil {
		return nil, err
	}
	sub, er := middlewares.Auth(gc, r.Env)
	if er != nil {
		return nil, er
	}
	id := sub.UserID
	newsFeed, er := handlers.Handler{}.NewsFeed(r.Env, id, ctx)
	if er != nil {
		return nil, er
	}
	return newsFeed, nil
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
