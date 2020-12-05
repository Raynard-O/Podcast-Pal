package database

import (
	"context"
	"errors"
	"github.com/raynard2/backend/models"
	"go.mongodb.org/mongo-driver/bson"
)

func (d *mongoStore) SavePodcast(podcast interface{}, Collection string) (models.Podcast, error) {

	var output models.Podcast
	_result, err := d.Database.Collection(Collection).InsertOne(context.Background(), podcast)
	if err != nil {
		return output, err
	}

	err = d.Database.Collection(Collection).FindOne(nil, bson.M{"_id": _result.InsertedID}).Decode(&output)

	if err != nil {
		return output, err
	}
	return output, nil
}



//Find One by

func (d *mongoStore) FindOnePodcast(collection string, key, pair string) (*models.Podcast, error) {
	result := new(models.Podcast)
	filter := bson.M{key: pair}
	//ops := options.FindOne()
	//ops.Projection = projection
	if err := d.Database.Collection(collection).FindOne(nil, filter).Decode(&result); err != nil {
		if err.Error() == "mongo: no documents in result" {
			return nil, errors.New("document not found")
		}
		return nil, err
	}
	return result, nil
}