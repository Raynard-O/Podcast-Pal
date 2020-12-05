package database

import (
	"context"
	"fmt"
	"github.com/raynard2/backend/config"
	"github.com/raynard2/backend/models"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// DataDB ...
type DataDB interface {
	// users
	SaveUser(user interface{}, Collection string) (models.User, error)
	FindOneUser(collection string, key, pair string, result interface{}) error
	FindUserByID(id interface{}, projection map[string]interface{}, result interface{}) error
	UpdateOneUser(fields map[string]interface{}, changes interface{}) error
	SavePodcast(podcast interface{}, Collection string) (models.Podcast, error)
	FindOnePodcast(collection string, key, pair string) (*models.Podcast, error)
	Save(user interface{}, Collection string) (models.User, error)
	//Save(payload interface{}) (map[string]interface{}, error)
	FindByID(id interface{}, projection map[string]interface{}, result interface{}) error
	//FindOne(fields, projection map[string]interface{}, result interface{}) error
	FindOne(collection string, key, pair string) (*models.User, error)
	FindMany(fields, projection, sort map[string]interface{}, limit, skip int64, results interface{}) error
	UpdateByID(id interface{}, payload interface{}) error
	UpdateOne(fields map[string]interface{}, payload interface{}) error
	UpdateMany(fields, payload map[string]interface{}) error
	DeleteOne(fields map[string]interface{}) error
	DeleteMany(fields map[string]interface{}) error
	FindAll(collection string) (interface{}, error)
}

type mongoStore struct {
	IsConnected    bool
	CollectionName string
	Collection     *mongo.Collection
	Database       *mongo.Database
}

/**
 * NewMongoConn
 * This initialises a new MongoDB mongoStore
 * param: string databaseURl
 * param: string databaseName
 * param: string collection
 * return: *mongoStore
 */

var mongoDb, err = config.LoadSecrets()

func NewMongoConn(databaseName string, collection string) (DataDB, error) {

	databaseURL := fmt.Sprintf("mongodb+srv://%v:%v@%v/<%v>?retryWrites=true&w=majority", mongoDb.MONGO_USER, mongoDb.MONGO_PASS, mongoDb.MONGO_HOST, mongoDb.MONGO_DB)
	//fmt.Print(databaseURL)
	clientOptions := options.Client().ApplyURI(databaseURL)

	// Connect to MongoDB
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		return nil, err
	}

	// Check the connection
	if err = client.Ping(context.TODO(), nil); err != nil {
		//log.Fatal(err)
		return nil, err
	}

	db := client.Database(databaseName)

	mongoStore := mongoStore{
		IsConnected:    true,
		CollectionName: collection,
		Collection:     db.Collection(collection),
		Database:       db,
	}
	return &mongoStore, nil
}
