package db

import (
	"context"
	"fmt"

	"github.com/SpectralJager/spender/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type StoreDropper interface {
	Drop(context.Context) error
}

type UserStore interface {
	StoreDropper

	GetUserByID(context.Context, string) (types.User, error)
	GetAllUsers(context.Context) ([]types.User, error)
	CreateUser(context.Context, types.User) (string, error)
	DeleteUser(context.Context, string) error
	UpdateUser(context.Context, string, types.UpdateUserParams) error
}

func ToObjectID(id string) primitive.ObjectID {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return primitive.ObjectID{}
	}
	return oid
}

func toBsonDoc(v interface{}) (*bson.D, error) {
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

type MongoUserStore struct {
	coll *mongo.Collection
}

func NewMongoUserStore(cl *mongo.Client, dbname, userColl string) *MongoUserStore {
	return &MongoUserStore{
		coll: cl.Database(dbname).Collection(userColl),
	}
}

func (st MongoUserStore) Drop(ctx context.Context) error {
	return st.coll.Drop(ctx)
}

func (st MongoUserStore) GetUserByID(ctx context.Context, id string) (types.User, error) {
	var user types.User
	if err := st.coll.FindOne(ctx, bson.M{"_id": ToObjectID(id)}).Decode(&user); err != nil {
		return types.User{}, err
	}
	return user, nil
}

func (st MongoUserStore) GetAllUsers(ctx context.Context) ([]types.User, error) {
	cur, err := st.coll.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	users := []types.User{}
	err = cur.All(ctx, &users)
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (st MongoUserStore) CreateUser(ctx context.Context, user types.User) (string, error) {
	res, err := st.coll.InsertOne(ctx, user)
	if err != nil {
		return "", err
	}
	user.ID = res.InsertedID.(primitive.ObjectID).Hex()
	return user.ID, nil
}

func (st MongoUserStore) DeleteUser(ctx context.Context, id string) error {
	res, err := st.coll.DeleteOne(ctx, bson.M{"_id": ToObjectID(id)})
	if err != nil {
		return err
	}
	if res.DeletedCount <= 0 {
		return fmt.Errorf("can't delete user with id = %s", id)
	}
	return nil
}

func (st MongoUserStore) UpdateUser(ctx context.Context, id string, user types.UpdateUserParams) error {
	userBson, err := toBsonDoc(user)
	if err != nil {
		return err
	}
	res, err := st.coll.UpdateOne(ctx, bson.M{"_id": ToObjectID(id)}, bson.D{{Key: "$set", Value: userBson}})
	if err != nil {
		return err
	}
	if res.ModifiedCount <= 0 || res.MatchedCount <= 0 {
		return fmt.Errorf("no changes for user with id = %s", id)
	}
	return nil
}
