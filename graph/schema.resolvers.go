package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"errors"
	"log"
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
	newUser, err := handlers.Handler{}.NewUser(input)
	if newUser == nil {
		//err := errors.New("unable to create user")
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
		Email:           newUser.Email,
		IsVerified:      newUser.IsVerified,
		IsOtpVerified:   newUser.IsOtpVerified,
		IsPasswordReset: newUser.IsPasswordReset,
		Picture:         newUser.Picture,
		Topics:          topicsFromDB,
		ID:              newUser.ID.String(),
	}
	return user, nil
}

// Login is the resolver for the Login field.
func (r *mutationResolver) Login(ctx context.Context, input model.Login) (*model.LoginResponse, error) {
	response, err := handlers.Handler{}.Login(input, r.Env)

	if err != nil {
		return nil, err
	}

	return response, nil
}

// ForgotPassword is the resolver for the ForgotPassword field.
func (r *mutationResolver) ForgotPassword(ctx context.Context, input model.ForgotPassword) (*model.GenericResponse, error) {
	response, err := handlers.Handler{}.ForgotPassword(input)

	if err != nil {
		return nil, err
	}

	return response, nil
}

// ResetPassword is the resolver for the ResetPassword field.
func (r *mutationResolver) ResetPassword(ctx context.Context, input model.ResetPassword) (*model.GenericResponse, error) {
	response, err := handlers.Handler{}.ResetPassword(input)

	if err != nil {
		return nil, err
	}

	return response, nil
}

// VerifyEmail is the resolver for the VerifyEmail field.
func (r *mutationResolver) VerifyEmail(ctx context.Context, input model.VerifyOtp) (*model.GenericResponse, error) {
	response, err := handlers.Handler{}.VerifyEmail(input)

	if err != nil {
		return nil, err
	}

	return response, nil
}

// VerifyResetOtp is the resolver for the VerifyResetOtp field.
func (r *mutationResolver) VerifyResetOtp(ctx context.Context, input model.VerifyOtp) (*model.GenericResponse, error) {
	response, err := handlers.Handler{}.VerifyResetOtp(input)

	if err != nil {
		return nil, err
	}

	return response, nil
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
	id := sub
	log.Print(id)
	getUser := handlers.Handler{}.GetUserById(id)

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
		Email:           getUser.Email,
		IsVerified:      getUser.IsVerified,
		IsOtpVerified:   getUser.IsOtpVerified,
		IsPasswordReset: getUser.IsPasswordReset,
		Picture:         getUser.Picture,
		FullName:        getUser.FullName,
		ID:              getUser.ID.String(),
		Topics:          topicsFromDB,
		UpdatedAt:       getUser.UpdatedAt,
		CreatedAt:       getUser.CreatedAt,
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
	id := sub

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
