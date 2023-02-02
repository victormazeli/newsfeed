package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
	"gatewayservice/graph/generated"
	"gatewayservice/graph/model"
	"gatewayservice/handlers"
)

// ManualSignUp is the resolver for the ManualSignUp field.
func (r *mutationResolver) ManualSignUp(ctx context.Context, input model.UserSignUpDetail) (*model.User, error) {
	addedUser, err := handlers.AuthHandler{}.UserRegistration(input)
	if err == nil {
		return addedUser, nil
	} else {
		return nil, err
	}
}

// ManualLogin is the resolver for the ManualLogin field.
func (r *mutationResolver) ManualLogin(ctx context.Context, input model.LoginCredential) (*model.AuthResponse, error) {
	tokenPair, err := handlers.AuthHandler{}.UserLogin(input)
	if err == nil {
		return tokenPair, nil
	} else {
		return nil, err
	}
}

// User is the resolver for the user field.
func (r *queryResolver) User(ctx context.Context, id string) (*model.User, error) {
	panic(fmt.Errorf("not implemented: User - user"))
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
