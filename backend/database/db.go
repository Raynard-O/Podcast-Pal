package database

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/raynard2/backend/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

/**
 * Save
 * Save is used to save a record in the mongoStore
 */

func (d *mongoStore) Save(user interface{}, Collection string) (models.User, error) {

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
 * FindByID
 * find a single record by id in the mongoStore
 * returns nil if record isn't found.
 *
 * param: interface{}            id
 * param: map[string]interface{} projection
 * return: map[string]interface{}
 */
func (d *mongoStore) FindByID(id interface{}, projection map[string]interface{}, result interface{}) error {
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
 * Find One by
 */
func (d *mongoStore) FindOne(collection string, key, pair string) (*models.User, error) {
	result := new(models.User)
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

func (d *mongoStore) FindMany(fields, projection, sort map[string]interface{}, limit, skip int64, results interface{}) error {
	ops := options.Find()
	if limit > 0 {
		ops.Limit = &limit
	}
	if skip > 0 {
		ops.Skip = &skip
	}
	if projection != nil {
		ops.Projection = projection
	}
	if sort != nil {
		ops.Sort = sort
	}
	cursor, err := d.Collection.Find(nil, fields, ops)
	if err != nil {
		return err
	}

	var output []map[string]interface{}
	for cursor.Next(nil) {
		var item map[string]interface{}
		_ = cursor.Decode(&item)
		output = append(output, item)
	}

	if b, e := json.Marshal(output); e == nil {
		_ = json.Unmarshal(b, &results)
	} else {
		return e
	}
	return nil
}

/**
 * UpdateByID
 * Updates a single record by id in the mongoStore
 *
 * param: interface{} id
 * param: interface{} payload
 * return: error
 */
func (d *mongoStore) UpdateByID(id interface{}, payload interface{}) error {
	var u map[string]interface{}
	opts := options.FindOneAndUpdate()
	up := true
	opts.Upsert = &up
	if err := d.Collection.FindOneAndUpdate(nil, bson.M{"_id": id}, bson.M{
		"$set": payload,
	}).Decode(&u); err != nil {
		return err
	}
	return nil
}
func (d *mongoStore) FindAll(collection string) (interface{}, error) {
	//var users []models.User
	ops := options.Find()
	cursor, err := d.Collection.Find(nil, bson.D{}, ops)
	var output []map[string]interface{}
	for cursor.Next(nil) {
		var item map[string]interface{}
		_ = cursor.Decode(&item)
		output = append(output, item)
	}
	var results interface{}
	if b, e := json.Marshal(output); e == nil {
		_ = json.Unmarshal(b, &results)
	} else {
		return e, err
	}
	return results, err
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
func (d *mongoStore) UpdateOne(fields map[string]interface{}, payload interface{}) error {
	var u map[string]interface{}
	if err := d.Collection.FindOneAndUpdate(nil, fields, bson.M{
		"$set": payload,
	}).Decode(&u); err != nil {
		return err
	}
	return nil
}

/**
 * UpdateMany
 * Updates many items in the collection
 * `fields` this is the search criteria
 * `payload` this is the update payload.
 *
 * param: map[string]interface{} fields
 * param: interface{}            payload
 * return: error
 */
func (d *mongoStore) UpdateMany(fields, payload map[string]interface{}) error {
	if _, err := d.Collection.UpdateMany(nil, fields, bson.M{
		"$set": payload,
	}); err != nil {
		return err
	}
	return nil
}

/**
 * DeleteOne
 * Deletes one item from the mongoStore using fields a hash map to properly filter what is to be deleted.
 *
 * param: map[string]interface{} fields
 * return: error
 */
func (d *mongoStore) DeleteOne(fields map[string]interface{}) error {
	_, err := d.Collection.DeleteOne(nil, fields)
	if err != nil {
		return err
	}

	return nil
}

/**
 * Delete Many items from the mongoStore
 *
 * param: map[string]interface{} fields
 * return: error
 */
func (d *mongoStore) DeleteMany(fields map[string]interface{}) error {
	_, err := d.Collection.DeleteMany(nil, fields)
	if err != nil {
		return err
	}

	return nil
}
