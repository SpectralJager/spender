package utils

import (
	"encoding/json"
	"io"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
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

func EncryptPassword(password string, cost int) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), cost)
}

func DecodeBody[T any](r io.Reader) (T, error) {
	var content T
	err := json.NewDecoder(r).Decode(&content)
	return content, err
}
