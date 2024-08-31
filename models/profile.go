package models

import "context"

type defaultProfileKey struct {}

type Profile struct {
	ID string
	Elo string
	Anonymous bool
}

func FromCtx(ctx context.Context) Profile {
	profile, ok := ctx.Value(defaultProfileKey{}).(Profile)

	if !ok {
		return GetAnonymousProfile()
	}

	return profile
}


func (p *Profile) ToCtx(ctx context.Context) context.Context {
	return context.WithValue(ctx, defaultProfileKey{}, *p)
}

func GetAnonymousProfile() Profile {
	return Profile{
		Anonymous: true,
	}
}