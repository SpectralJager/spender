package db

import (
	"context"
	"fmt"

	"github.com/SpectralJager/spender/types"
	"github.com/SpectralJager/spender/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Dropper interface {
	Drop(ctx context.Context) error
}

type UserStore interface {
	Dropper

	GetAllUsers(ctx context.Context) ([]types.User, error)

	GetUserByID(ctx context.Context, id string) (types.User, error)
	CreateUser(ctx context.Context, user types.User) (string, error)
	DeleteUser(ctx context.Context, id string) error
	UpdateUser(ctx context.Context, id string, params types.UpdateUserParams) error
}

type UpdateParams interface {
	ToBsonDoc() (*bson.D, error)
}

type SpendStor[T any] interface {
	Dropper

	GetAllSpends(ctx context.Context, ownerID string) ([]T, error)
	CreateSpend(ctx context.Context, spend T) (string, error)
	GetSpendByID(ctx context.Context, ownerID, id string) (T, error)
	DeleteSpend(ctx context.Context, ownerID, id string) error
	UpdateSpend(ctx context.Context, ownerID, id string, params UpdateParams) error
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
	if err := st.coll.FindOne(ctx, bson.M{"_id": utils.ToObjectID(id)}).Decode(&user); err != nil {
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
	res, err := st.coll.DeleteOne(ctx, bson.M{"_id": utils.ToObjectID(id)})
	if err != nil {
		return err
	}
	if res.DeletedCount <= 0 {
		return fmt.Errorf("can't delete user with id = %s", id)
	}
	return nil
}

func (st MongoUserStore) UpdateUser(ctx context.Context, id string, user types.UpdateUserParams) error {
	userBson, err := utils.ToBsonDoc(user)
	if err != nil {
		return err
	}
	res, err := st.coll.UpdateOne(ctx, bson.M{"_id": utils.ToObjectID(id)}, bson.D{{Key: "$set", Value: userBson}})
	if err != nil {
		return err
	}
	if res.ModifiedCount <= 0 || res.MatchedCount <= 0 {
		return fmt.Errorf("no changes for user with id = %s", id)
	}
	return nil
}

type MongoSpendStore[T any] struct {
	coll *mongo.Collection
}

func NewMongoTimespendStore(client *mongo.Client, dbname string, collname string) MongoSpendStore[types.Timespend] {
	return MongoSpendStore[types.Timespend]{
		coll: client.Database(dbname).Collection(collname),
	}
}

func NewMongoMoneyspendStore(client *mongo.Client, dbname string, collname string) MongoSpendStore[types.Moneyspend] {
	return MongoSpendStore[types.Moneyspend]{
		coll: client.Database(dbname).Collection(collname),
	}
}

func (st MongoSpendStore[T]) Drop(ctx context.Context) error {
	return st.coll.Drop(ctx)
}

func (st MongoSpendStore[T]) GetAllSpends(ctx context.Context, ownerID string) ([]T, error) {
	cur, err := st.coll.Find(ctx, bson.M{"ownerid": ownerID})
	if err != nil {
		return nil, err
	}
	spends := []T{}
	err = cur.All(ctx, &spends)
	if err != nil {
		return nil, err
	}
	return spends, nil
}

func (st MongoSpendStore[T]) CreateSpend(ctx context.Context, spend T) (string, error) {
	res, err := st.coll.InsertOne(ctx, spend)
	if err != nil {
		return "", err
	}
	return res.InsertedID.(primitive.ObjectID).Hex(), nil
}

func (st MongoSpendStore[T]) GetSpendByID(ctx context.Context, ownerid, id string) (T, error) {
	var spend T
	err := st.coll.FindOne(ctx, bson.M{"_id": utils.ToObjectID(id), "ownerid": ownerid}).Decode(&spend)
	return spend, err
}

func (st MongoSpendStore[T]) DeleteSpend(ctx context.Context, ownerid, id string) error {
	res, err := st.coll.DeleteOne(ctx, bson.M{"_id": utils.ToObjectID(id), "ownerid": ownerid})
	if err != nil {
		return err
	}
	if res.DeletedCount <= 0 {
		return fmt.Errorf("can't delete spend with id: %s", id)
	}
	return nil
}

func (st MongoSpendStore[T]) UpdateSpend(ctx context.Context, ownerid, id string, params UpdateParams) error {
	paramsBson, err := params.ToBsonDoc()
	if err != nil {
		return err
	}
	res, err := st.coll.UpdateOne(ctx, bson.M{"_id": utils.ToObjectID(id), "ownerid": ownerid}, bson.D{{Key: "$set", Value: paramsBson}})
	if err != nil {
		return err
	}
	if res.MatchedCount <= 0 {
		return fmt.Errorf("can't update spend with id: %s", id)
	}
	if res.ModifiedCount <= 0 {
		return fmt.Errorf("no changes spend with id: %s", id)
	}
	return nil
}
