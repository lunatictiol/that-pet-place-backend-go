package config

import (
	"context"

	firebase "firebase.google.com/go"
	"google.golang.org/api/option"
)

func InitFirebaseApp() (*firebase.App, error) {
	opt := option.WithCredentialsFile("path/to/serviceAccountKey.json")
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		return nil, err
	}
	return app, nil
}
