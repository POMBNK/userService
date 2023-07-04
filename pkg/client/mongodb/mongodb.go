package mongodb

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func NewClient(ctx context.Context, host, port, username, password, database, authDB string) (*mongo.Database, error) {

	var isAuth bool
	var mongodbURL string

	if username == "" && password == "" {
		isAuth = false
		mongodbURL = fmt.Sprintf("mongodb://%s:%s", host, port)
	} else {
		isAuth = true
		mongodbURL = fmt.Sprintf("mongodb://%s:%s@%s:%s", username, password, host, port)
	}

	clientOpts := options.Client().ApplyURI(mongodbURL)
	if isAuth {
		if authDB == "" {
			authDB = database
		}
		credential := options.Credential{
			AuthSource: authDB,
			Username:   username,
			Password:   password,
		}
		clientOpts.SetAuth(credential)
	}

	client, err := mongo.Connect(ctx, clientOpts)
	if err != nil {
		return nil, fmt.Errorf("Can't connect to mongoDB due: %w ", err)
	}

	if err = client.Ping(ctx, nil); err != nil {
		return nil, fmt.Errorf("Can't ping mongoDB due: %w ", err)
	}

	return client.Database(database), nil
}
