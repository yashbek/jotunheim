package models

import (
	"context"
)

type defaultProfileKey struct{}

type ProfileStore struct{}

type User struct {
	ID         string
	Username   string
	Hash       string
	Email      string
	Elo        int
	DateJoined string
}

func FromCtx(ctx context.Context) User {
	profile, ok := ctx.Value(defaultProfileKey{}).(User)

	if !ok {
		return GetAnonymousProfile()
	}

	return profile
}

func (p *User) ToCtx(ctx context.Context) context.Context {
	return context.WithValue(ctx, defaultProfileKey{}, *p)
}

func GetAnonymousProfile() User {
	return User{}
}
