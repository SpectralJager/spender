package db

import (
	"context"

	"github.com/SpectralJager/spender/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	DBNAME   = "spender"
	USERCOLL = "users"
)

type UserStore interface {
	GetUserByID(context.Context, string) (*types.User, error)
	GetAllUsers(context.Context) ([]*types.User, error)
}

func ToObjectID(id string) primitive.ObjectID {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return primitive.ObjectID{}
	}
	return oid
}

type MongoUserStore struct {
	client *mongo.Client
	coll   *mongo.Collection
}

func NewMongoUserStore(cl *mongo.Client) *MongoUserStore {
	return &MongoUserStore{
		client: cl,
		coll:   cl.Database(DBNAME).Collection(USERCOLL),
	}
}

func (st MongoUserStore) GetUserByID(ctx context.Context, id string) (*types.User, error) {
	var user types.User
	if err := st.coll.FindOne(ctx, bson.M{"_id": ToObjectID(id)}).Decode(&user); err != nil {
		return nil, err
	}
	return &user, nil
}

func (st MongoUserStore) GetAllUsers(ctx context.Context) ([]*types.User, error) {
	cur, err := st.coll.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	users := []*types.User{}
	for cur.Next(ctx) {
		var user types.User
		err := cur.Decode(&user)
		if err != nil {
			return nil, err
		}
		users = append(users, &user)
	}
	return users, nil
}
