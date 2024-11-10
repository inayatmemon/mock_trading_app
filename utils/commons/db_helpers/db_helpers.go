package db_helpers

import (
	"9DotTechnology/trading_platform/connections/mongodb"
	"9DotTechnology/trading_platform/constants/db_constants"
	"9DotTechnology/trading_platform/utils/logwrapper"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetOneDocument[T any](collectionName string, qry bson.M) (T, error) {
	collection := mongodb.GetCollection(collectionName)
	dbctx, cancel := mongodb.DbContext(db_constants.MongoTimeout)
	defer cancel()

	if qry == nil {
		qry = bson.M{}
	}

	var data T
	err := collection.FindOne(dbctx, qry).Decode(&data)
	if err != nil {
		logwrapper.Logger.Errorln("error in getting single document from db: ", err)
		return data, err

	}

	return data, nil
}

func GetManyDocumentsWithOptions[T any](collectionName string, qry bson.M, skip, limit int64, sort bson.D, projection bson.M) ([]T, error) {
	collection := mongodb.GetCollection(collectionName)
	dbctx, cancel := mongodb.DbContext(db_constants.MongoTimeout)
	defer cancel()

	if qry == nil {
		qry = bson.M{}
	}

	options := options.Find()

	if skip != 0 {
		options.SetSkip(skip)
	}

	if limit != 0 {
		options.SetLimit(limit)
	}

	if sort != nil {
		options.SetSort(sort)
	}

	if projection != nil {
		options.SetProjection(projection)
	}

	var data []T
	cursor, err := collection.Find(dbctx, qry, options)
	if err != nil {
		logwrapper.Logger.Errorln("error in getting multiple document from db: ", err)
		return data, err
	}

	for cursor.Next(dbctx) {
		var doument T
		err = cursor.Decode(&doument)
		if err != nil {
			logwrapper.Logger.Errorln("error in decoding single document in Find Many: ", err)
			return data, err
		}
		data = append(data, doument)
	}

	return data, nil
}

func GetManyDocuments[T any](collectionName string, qry bson.M) ([]T, error) {
	collection := mongodb.GetCollection(collectionName)
	dbctx, cancel := mongodb.DbContext(db_constants.MongoTimeout)
	defer cancel()

	if qry == nil {
		qry = bson.M{}
	}

	var data []T
	cursor, err := collection.Find(dbctx, qry)
	if err != nil {
		logwrapper.Logger.Errorln("error in getting multiple document from db: ", err)
		return data, err
	}

	for cursor.Next(dbctx) {
		var doument T
		err = cursor.Decode(&doument)
		if err != nil {
			logwrapper.Logger.Errorln("error in decoding single document in Find Many: ", err)
			return data, err
		}
		data = append(data, doument)
	}

	return data, nil
}

func GetCount(collectionName string, qry bson.M) (int64, error) {
	collection := mongodb.GetCollection(collectionName)
	dbctx, cancel := mongodb.DbContext(db_constants.MongoTimeout)
	defer cancel()

	if qry == nil {
		qry = bson.M{}
	}

	count, err := collection.CountDocuments(dbctx, qry)
	if err != nil {
		logwrapper.Logger.Errorln("error in getting counts of document from db: ", err)
		return 0, err
	}

	return count, nil
}

func InsertOneDocument(collectionName string, data interface{}) error {
	collection := mongodb.GetCollection(collectionName)
	dbctx, cancel := mongodb.DbContext(db_constants.MongoTimeout)
	defer cancel()

	_, err := collection.InsertOne(dbctx, data)
	if err != nil {
		logwrapper.Logger.Errorln("error in inserting document: ", err)
		return err
	}

	return nil
}

func InsertManyDocuments(collectionName string, data []interface{}) error {
	collection := mongodb.GetCollection(collectionName)
	dbctx, cancel := mongodb.DbContext(db_constants.MongoTimeout)
	defer cancel()

	_, err := collection.InsertMany(dbctx, data)
	if err != nil {
		logwrapper.Logger.Errorln("error in inserting document: ", err)
		return err
	}

	return nil
}

func UpdateOneDocument(collectionName string, cond bson.M, data interface{}) error {
	collection := mongodb.GetCollection(collectionName)
	dbctx, cancel := mongodb.DbContext(db_constants.MongoTimeout)
	defer cancel()

	updData := bson.M{
		"$set": data,
	}

	_, err := collection.UpdateOne(dbctx, cond, updData)
	if err != nil {
		logwrapper.Logger.Errorln("error in updating document: ", err)
		return err
	}

	return nil
}

func UpdateOneDocumentWithData(collectionName string, cond bson.M, data interface{}) error {
	collection := mongodb.GetCollection(collectionName)
	dbctx, cancel := mongodb.DbContext(db_constants.MongoTimeout)
	defer cancel()

	_, err := collection.UpdateOne(dbctx, cond, data)
	if err != nil {
		logwrapper.Logger.Errorln("error in updating document: ", err)
		return err
	}

	return nil
}

func UpdateOneDocumentWithID(collectionName string, id primitive.ObjectID, data interface{}) error {
	collection := mongodb.GetCollection(collectionName)
	dbctx, cancel := mongodb.DbContext(db_constants.MongoTimeout)
	defer cancel()

	updData := bson.M{
		"$set": data,
	}

	_, err := collection.UpdateByID(dbctx, id, updData)
	if err != nil {
		logwrapper.Logger.Errorln("error in updating document: ", err)
		return err
	}

	return nil
}

func UpdateOneDocumentWithIDWithData(collectionName string, id primitive.ObjectID, data interface{}) error {
	collection := mongodb.GetCollection(collectionName)
	dbctx, cancel := mongodb.DbContext(db_constants.MongoTimeout)
	defer cancel()

	_, err := collection.UpdateByID(dbctx, id, data)
	if err != nil {
		logwrapper.Logger.Errorln("error in updating document: ", err)
		return err
	}

	return nil
}

func UpdateManyDocuments(collectionName string, cond bson.M, data interface{}) error {
	collection := mongodb.GetCollection(collectionName)
	dbctx, cancel := mongodb.DbContext(db_constants.MongoTimeout)
	defer cancel()

	updData := bson.M{
		"$set": data,
	}

	_, err := collection.UpdateMany(dbctx, cond, updData)
	if err != nil {
		logwrapper.Logger.Errorln("error in updating many documents: ", err)
		return err
	}

	return nil
}

func UpdateManyDocumentsWithData(collectionName string, cond bson.M, data interface{}) error {
	collection := mongodb.GetCollection(collectionName)
	dbctx, cancel := mongodb.DbContext(db_constants.MongoTimeout)
	defer cancel()

	_, err := collection.UpdateMany(dbctx, cond, data)
	if err != nil {
		logwrapper.Logger.Errorln("error in updating many documents: ", err)
		return err
	}

	return nil
}

func DeleteOneDocument(collectionName string, cond bson.M) error {
	collection := mongodb.GetCollection(collectionName)
	dbctx, cancel := mongodb.DbContext(db_constants.MongoTimeout)
	defer cancel()

	_, err := collection.DeleteOne(dbctx, cond)
	if err != nil {
		logwrapper.Logger.Errorln("error in deleting one document: ", err)
		return err
	}

	return nil
}

func DeleteManyDocuments(collectionName string, cond bson.M) error {
	collection := mongodb.GetCollection(collectionName)
	dbctx, cancel := mongodb.DbContext(db_constants.MongoTimeout)
	defer cancel()

	_, err := collection.DeleteMany(dbctx, cond)
	if err != nil {
		logwrapper.Logger.Errorln("error in deleting many documents: ", err)
		return err
	}

	return nil
}
