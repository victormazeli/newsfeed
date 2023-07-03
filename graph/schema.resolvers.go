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

func (r *mutationResolver) CreateNewUser(ctx context.Context, input model.CreateUser) (*model.User, error) {
	checkUser := handlers.Handler{}.GetUserByEmail(input.Email)
	if checkUser != nil {
		err := errors.New("user already exist")
		return nil, err

	}
	newUser, err := handlers.Handler{}.NewUser(input, r.Env)
	if newUser == nil {
		//err := errors.New("unable to create user")
		return nil, err
	}

	//var topicsFromDB []*model.Topic
	//for _, t := range newUser.Topics {
	//	item := model.Topic{
	//		Topic: t.Topic,
	//	}
	//	topicsFromDB = append(topicsFromDB, &item)
	//}

	user := &model.User{
		Email:           &newUser.Email,
		IsVerified:      &newUser.IsVerified,
		IsOtpVerified:   &newUser.IsOtpVerified,
		IsPasswordReset: &newUser.IsPasswordReset,
		Picture:         &newUser.Picture,
		Topics:          newUser.Topics,
		FullName:        &newUser.FullName,
		ID:              newUser.ID.Hex(),
		CreatedAt:       newUser.CreatedAt.String(),
		UpdatedAt:       newUser.UpdatedAt.String(),
	}
	return user, nil
}

func (r *mutationResolver) Login(ctx context.Context, input model.Login) (*model.LoginResponse, error) {
	response, err := handlers.Handler{}.Login(input, r.Env, ctx)

	if err != nil {
		return nil, err
	}

	return response, nil
}

func (r *mutationResolver) CompleteRegistration(ctx context.Context, input model.CompleteRegistration) (*model.GenericResponse, error) {
	gc, err := middlewares.GinContextFromContext(ctx)
	if err != nil {
		return nil, err
	}
	sub, er := middlewares.Auth(gc, r.Env)
	if er != nil {
		return nil, er
	}

	response, e := handlers.Handler{}.CompleteRegistration(input, sub, ctx)

	if e != nil {
		return nil, e
	}

	return response, nil
}

func (r *mutationResolver) ForgotPassword(ctx context.Context, input model.ForgotPassword) (*model.GenericResponse, error) {
	response, err := handlers.Handler{}.ForgotPassword(input, r.Env, ctx)

	if err != nil {
		return nil, err
	}

	return response, nil
}

func (r *mutationResolver) GoogleLogin(ctx context.Context, input model.GoogleAuth) (*model.LoginResponse, error) {
	response, err := handlers.Handler{}.GoogleLogin(input, r.Env, ctx)

	if err != nil {
		return nil, err
	}

	return response, nil
}

func (r *mutationResolver) ResetPassword(ctx context.Context, input model.ResetPassword) (*model.GenericResponse, error) {
	response, err := handlers.Handler{}.ResetPassword(input, ctx)

	if err != nil {
		return nil, err
	}

	return response, nil
}

func (r *mutationResolver) VerifyEmail(ctx context.Context, input model.VerifyOtp) (*model.GenericResponse, error) {
	response, err := handlers.Handler{}.VerifyEmail(input, ctx)

	if err != nil {
		return nil, err
	}

	return response, nil
}

func (r *mutationResolver) VerifyResetOtp(ctx context.Context, input model.VerifyOtp) (*model.GenericResponse, error) {
	response, err := handlers.Handler{}.VerifyResetOtp(input, ctx)

	if err != nil {
		return nil, err
	}

	return response, nil
}

func (r *mutationResolver) Logout(ctx context.Context, input model.Logout) (*model.GenericResponse, error) {
	response, err := handlers.Handler{}.Logout(input, ctx)

	if err != nil {
		return nil, err
	}

	return response, nil
}

func (r *mutationResolver) AskKora(ctx context.Context, input model.PromptContent) (*model.PromptResponse, error) {
	gc, err := middlewares.GinContextFromContext(ctx)
	if err != nil {
		return nil, err
	}
	_, er := middlewares.Auth(gc, r.Env)
	if er != nil {
		return nil, er
	}

	response, err := handlers.Handler{}.AskChatGPT(input, r.Env, ctx)

	if err != nil {
		return nil, err
	}

	return response, nil
}

func (r *mutationResolver) SaveNews(ctx context.Context, newsID string) (*bool, error) {
	gc, err := middlewares.GinContextFromContext(ctx)
	if err != nil {
		return nil, err
	}
	sub, er := middlewares.Auth(gc, r.Env)
	if er != nil {
		return nil, er
	}

	response, err := handlers.Handler{}.SaveNews(newsID, sub, ctx)

	if err != nil {
		return nil, err
	}

	return response, nil
}

func (r *mutationResolver) LikeNews(ctx context.Context, newsID string) (*bool, error) {
	gc, err := middlewares.GinContextFromContext(ctx)
	if err != nil {
		return nil, err
	}
	sub, er := middlewares.Auth(gc, r.Env)
	if er != nil {
		return nil, er
	}

	response, err := handlers.Handler{}.LikeNews(sub, newsID, ctx)

	if err != nil {
		return nil, err
	}

	return response, nil
}

func (r *mutationResolver) ResendOtp(ctx context.Context, email string) (*model.GenericResponse, error) {
	response, err := handlers.Handler{}.ResendOtp(email, ctx, r.Env)

	if err != nil {
		return nil, err
	}

	return response, nil
}

func (r *mutationResolver) EditUserProfile(ctx context.Context, input model.UpdateProfile) (*model.GenericResponse, error) {
	gc, err := middlewares.GinContextFromContext(ctx)
	if err != nil {
		return nil, err
	}
	sub, er := middlewares.Auth(gc, r.Env)
	if er != nil {
		return nil, er
	}

	response, err := handlers.Handler{}.EditUserProfile(sub, input, ctx)

	if err != nil {
		return nil, err
	}

	return response, nil
}

func (r *mutationResolver) DeleteProfile(ctx context.Context) (*model.GenericResponse, error) {
	gc, err := middlewares.GinContextFromContext(ctx)
	if err != nil {
		return nil, err
	}
	sub, er := middlewares.Auth(gc, r.Env)
	if er != nil {
		return nil, er
	}

	response, err := handlers.Handler{}.DeleteUserProfile(sub, ctx)

	if err != nil {
		return nil, err
	}

	return response, nil
}

func (r *mutationResolver) EditUserInterest(ctx context.Context, topics []string) (*model.GenericResponse, error) {
	gc, err := middlewares.GinContextFromContext(ctx)
	if err != nil {
		return nil, err
	}
	sub, er := middlewares.Auth(gc, r.Env)
	if er != nil {
		return nil, er
	}

	response, err := handlers.Handler{}.EditUserInterest(sub, topics, ctx)

	if err != nil {
		return nil, err
	}

	return response, nil
}

func (r *mutationResolver) ChangePassword(ctx context.Context, input model.ChangePassword) (*model.GenericResponse, error) {
	gc, err := middlewares.GinContextFromContext(ctx)
	if err != nil {
		return nil, err
	}
	sub, er := middlewares.Auth(gc, r.Env)
	if er != nil {
		return nil, er
	}
	id := sub
	response, err := handlers.Handler{}.ChangeUserPassword(id, input, ctx)

	if err != nil {
		return nil, err
	}

	return response, nil
}

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
	//log.Print(id)
	getUser := handlers.Handler{}.GetUserById(id)

	if getUser == nil {
		err := errors.New("unable to fetch user")
		return nil, err
	}

	//var topicsFromDB []*model.Topic
	//for _, t := range getUser.Topics {
	//	item := model.Topic{
	//		Topic: t.Topic,
	//	}
	//	topicsFromDB = append(topicsFromDB, &item)
	//}

	user := &model.User{
		Email:           &getUser.Email,
		IsVerified:      &getUser.IsVerified,
		IsOtpVerified:   &getUser.IsOtpVerified,
		IsPasswordReset: &getUser.IsPasswordReset,
		Picture:         &getUser.Picture,
		PhoneNumber:     &getUser.PhoneNumber,
		FullName:        &getUser.FullName,
		ID:              getUser.ID.Hex(),
		Topics:          getUser.Topics,
		UpdatedAt:       getUser.UpdatedAt.String(),
		CreatedAt:       getUser.CreatedAt.String(),
	}

	return user, nil
}

func (r *queryResolver) GetLatestAndTrendingNews(ctx context.Context, query model.NewsQuery) ([]*model.Article, error) {
	gc, err := middlewares.GinContextFromContext(ctx)
	if err != nil {
		return nil, err
	}
	sub, er := middlewares.Auth(gc, r.Env)
	if er != nil {
		return nil, er
	}
	id := sub
	response, err := handlers.Handler{}.NewsFeed(id, query, ctx)

	if err != nil {
		return nil, err
	}

	return response, nil
}

func (r *queryResolver) GetNewsSources(ctx context.Context) ([]*model.Source, error) {
	response, err := handlers.Handler{}.FetchSources()

	if err != nil {
		return nil, err
	}

	return response, nil
}

func (r *queryResolver) SeedNewsSources(ctx context.Context) ([]*model.Source, error) {
	response, err := handlers.Handler{}.SeedSources(ctx, r.Env)

	if err != nil {
		return nil, err
	}

	return response, nil
}

func (r *queryResolver) GetNewsCategories(ctx context.Context) ([]*model.Category, error) {
	response, err := handlers.Handler{}.FetCategories()

	if err != nil {
		return nil, err
	}

	return response, nil
}

func (r *queryResolver) GetSingleNews(ctx context.Context, newsID string) (*model.Article, error) {
	gc, err := middlewares.GinContextFromContext(ctx)
	if err != nil {
		return nil, err
	}
	sub, er := middlewares.Auth(gc, r.Env)
	if er != nil {
		return nil, er
	}
	id := sub
	response, err := handlers.Handler{}.GetNewsById(newsID, id)

	if err != nil {
		return nil, err
	}

	return response, nil
}

func (r *queryResolver) GetSavedNews(ctx context.Context) ([]*model.Article, error) {
	gc, err := middlewares.GinContextFromContext(ctx)
	if err != nil {
		return nil, err
	}
	sub, er := middlewares.Auth(gc, r.Env)
	if er != nil {
		return nil, er
	}
	id := sub
	response, err := handlers.Handler{}.FetchSavedNews(id, ctx)

	if err != nil {
		return nil, err
	}

	return response, nil
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
