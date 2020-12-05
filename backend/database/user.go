package database

import (
	"context"
	"errors"
	"github.com/raynard2/backend/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

/**
 * Save User
 * Save is used to save a record in the mongoStore
 */

func (d *mongoStore) SaveUser(user interface{}, Collection string) (models.User, error) {

	var output models.User
	_result, err := d.Database.Collection(Collection).InsertOne(context.Background(), user)
	if err != nil {
		return output, err
	}

	err = d.Database.Collection(Collection).FindOne(nil, bson.M{"_id": _result.InsertedID}).Decode(&output)

	if err != nil {
		return output, err
	}
	return output, nil
}

/**
* Find One by
 */
func (d *mongoStore) FindOneUser(collection string, key, pair string, result interface{}) error {
	//result := new(interface{})
	filter := bson.M{key: pair}
	//ops := options.FindOne()
	//ops.Projection = projection
	if err := d.Database.Collection(collection).FindOne(nil, filter).Decode(result); err != nil {
		if err.Error() == "mongo: no documents in result" {
			return errors.New("document not found")
		}
		return err
	}
	return nil
}

/**
* FindByID
* find a single record by id in the mongoStore
* returns nil if record isn't found.
*
* param: interface{}            id
* param: map[string]interface{} projection
* return: map[string]interface{}
 */
func (d *mongoStore) FindUserByID(id interface{}, projection map[string]interface{}, result interface{}) error {
	ops := options.FindOne()
	if projection != nil {
		ops.Projection = projection
	}
	if err := d.Collection.FindOne(nil, bson.M{"_id": id}, ops).Decode(result); err != nil {
		if err.Error() == "mongo: no documents in result" {
			return errors.New("document not found")
		}
		return err
	}
	return nil
}

/**
 * UpdateOne
 *
 * Updates one item in the mongoStore using fields as the criteria.
 *
 * param: map[string]interface{} fields
 * param: interface{}            payload
 * return: error
 */
func (d *mongoStore) UpdateOneUser(fields map[string]interface{}, changes interface{}) error {
	var u map[string]interface{}
	if err := d.Collection.FindOneAndUpdate(nil, fields, bson.M{
		"$set": changes,
	}).Decode(&u); err != nil {
		return err
	}
	return nil
}
