package firebasedb

import (
	"context"
	"crypto/tls"
	"time"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
	"firebase.google.com/go/v4/db"
	"github.com/yashbek/jotunheim/models"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/api/option"
)

type Config struct {
	DatabaseURL string
}

type App struct {
	Auth   *auth.Client
	Client *db.Client
	ctx    context.Context
}

var FirebaseClient *App

const databaseURL = "https://lakeofnine-19bf4-default-rtdb.europe-west1.firebasedatabase.app"

func NewFirebaseApp(ctx context.Context, tlsConfig *tls.Config) (*App, error) {

	credOpt := option.WithCredentialsFile("/Users/myazbek/Documents/landing/jotunheim/db/firebasedb/lakeofnine.json")

	app, err := firebase.NewApp(ctx, nil,
		credOpt,
	)
	if err != nil {
		return nil, err
	}

	authClient, err := app.Auth(ctx)
	if err != nil {
		return nil, err
	}

	client, err := app.DatabaseWithURL(ctx, databaseURL)
	if err != nil {
		return nil, err
	}

	return &App{
		Auth:   authClient,
		Client: client,
		ctx:    ctx,
	}, nil
}

func (a *App) Create(path string, data interface{}) error {
	ref := a.Client.NewRef(path)
	return ref.Set(a.ctx, data)
}

func (a *App) CreateGame(path, gameID string, data interface{}) error {
	ref := a.Client.NewRef(path).Child(gameID)
	return ref.Set(a.ctx, data)
}

func (a *App) Read(path string, dest interface{}) error {
	ref := a.Client.NewRef(path)
	return ref.Get(a.ctx, dest)
}

func (a *App) ReadGame(path string, id string, dest interface{}) error {
	ref := a.Client.NewRef(path).Child(id)
	return ref.Get(a.ctx, dest)
}

func (a *App) ReadGameBoard(path string, id string, dest interface{}) error {
	ref := a.Client.NewRef(path).Child(id).Child("board")
	return ref.Get(a.ctx, dest)
}

func (a *App) ReadGameMoves(path string, id string, dest interface{}) error {
	ref := a.Client.NewRef(path).Child(id).Child("moves")
	return ref.Get(a.ctx, dest)
}

func (a *App) ReadUser(path string, id string, dest interface{}) error {
	ref := a.Client.NewRef(path).Child(id)
	return ref.Get(a.ctx, dest)
}

func (a *App) Update(path string, updates map[string]interface{}) error {
	ref := a.Client.NewRef(path)
	return ref.Update(a.ctx, updates)
}

func (a *App) UpdateGame(gameID string, data map[string]interface{}) error {
	ref := a.Client.NewRef("games").Child(gameID)
	return ref.Update(a.ctx, data)
}

func (a *App) Delete(path string) error {
	ref := a.Client.NewRef(path)
	return ref.Delete(a.ctx)
}

func (a *App) CreateUser(username, email, password string) (*auth.UserRecord, error) {
	params := (&auth.UserToCreate{}).
		Email(email).
		Password(password).
		DisplayName(username)

	userRecord, err := a.Auth.CreateUser(context.Background(), params)
	if err != nil {
		return &auth.UserRecord{}, err
	}

	newUser := a.Client.NewRef("users").Child(userRecord.UID)

	passKey, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	if err != nil {
		return &auth.UserRecord{}, err
	}

	values := models.User{
		ID:         userRecord.UID,
		Username:   username,
		Email:      email,
		Hash:       string(passKey),
		Elo:        0,
		DateJoined: time.Now().Format("2006-01-01"),
	}

	err = newUser.Set(a.ctx, values)

	return userRecord, err
}

func (a *App) SignInUser(email, password string) (*auth.UserRecord, error) {
	user, err := a.Auth.GetUserByEmail(context.Background(), email)
	if err != nil {
		return nil, err

	}

	userdb := a.Client.NewRef("users").Child(user.UID)

	userInfo := &models.User{}
	userdb.Get(a.ctx, userInfo)

	err = bcrypt.CompareHashAndPassword([]byte(userInfo.Hash), []byte(password))
	if err != nil {
		return &auth.UserRecord{}, err
	}

	return user, nil
}
