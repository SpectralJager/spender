package utils

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func ToObjectID(id string) primitive.ObjectID {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return primitive.ObjectID{}
	}
	return oid
}

func ToBsonDoc(v any) (*bson.D, error) {
	data, err := bson.Marshal(v)
	if err != nil {
		return nil, err
	}
	var doc *bson.D
	err = bson.Unmarshal(data, &doc)
	if err != nil {
		return nil, err
	}
	return doc, nil
}
