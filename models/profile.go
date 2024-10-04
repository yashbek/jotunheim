package models

import (
	"context"
	"fmt"
	"time"

	"github.com/yashbek/jotunheim/db/client"
	maindb "github.com/yashbek/jotunheim/db/models/main"
)

type defaultProfileKey struct {}

type ProfileStore struct {}

type Profile struct {
	ID string
	Elo int
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
	return Profile{}
}

func (ProfileStore) GetProfile(ctx context.Context, id string) (Profile, error) {
	q := db.GetDBClient(ctx).Query()

	dbProfile, err := q.GetProfile(ctx, 1)
	if err != nil {
		return Profile{}, err
	}
	
	profile := Profile{
		ID: fmt.Sprint(dbProfile.ID),
		Elo: dbProfile.Elo,
	}

	return profile, nil

}

func (ProfileStore) InsertProfile(ctx context.Context, id string) error {
	q := db.GetDBClient(ctx).Query()

	err := q.InsertProfile(ctx, maindb.InsertProfileParams{
		Email: "yazbek@test.com",
		Elo: 1000000,
		PhoneNumber: "96122333444",
		Username: "Baloot",
		DateJoined: time.Now(),
	})
	
	return err
}
