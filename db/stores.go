package db

import (
	"context"
	"fmt"

	"github.com/SpectralJager/spender/types"
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

type TimespendStore interface {
	Dropper

	GetAllTimes(ctx context.Context, ownerID string) ([]types.Timespend, error)

	CreateTimespend(ctx context.Context, timespend types.Timespend) (string, error)
	GetTimespendByID(ctx context.Context, ownerid, id string) (types.Timespend, error)
	DeleteTimespend(ctx context.Context, ownerid, id string) error
	UpdateTimespend(ctx context.Context, ownerid, id string, params types.UpdateTimespendParams) error
}

type MoneyspendStore interface {
	Dropper

	GetAllMonies(ctx context.Context, ownerID string) ([]types.Moneyspend, error)

	CreateMoneyspend(ctx context.Context, moneyspend types.Moneyspend) (string, error)
	GetMoneyspendByID(ctx context.Context, ownerid, id string) (types.Moneyspend, error)
	DeleteMoneyspend(ctx context.Context, ownerid, id string) error
	UpdateMoneyspend(ctx context.Context, ownerid, id string, params types.UpdateMoneyspendParams) error
}

func ToObjectID(id string) primitive.ObjectID {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return primitive.ObjectID{}
	}
	return oid
}

func toBsonDoc(v any) (*bson.D, error) {
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

type MongoTimespendStore struct {
	coll *mongo.Collection
}

func NewMongoTimespendStore(cl *mongo.Client, dbname, timespendColl string) *MongoTimespendStore {
	return &MongoTimespendStore{
		coll: cl.Database(dbname).Collection(timespendColl),
	}
}

func (st MongoTimespendStore) Drop(ctx context.Context) error {
	return st.coll.Drop(ctx)
}

func (st MongoTimespendStore) GetAllTimes(ctx context.Context, ownerID string) ([]types.Timespend, error) {
	cur, err := st.coll.Find(ctx, bson.M{"ownerid": ownerID})
	if err != nil {
		return nil, err
	}
	times := []types.Timespend{}
	err = cur.All(ctx, &times)
	if err != nil {
		return nil, err
	}
	return times, nil
}

func (st MongoTimespendStore) CreateTimespend(ctx context.Context, timepend types.Timespend) (string, error) {
	res, err := st.coll.InsertOne(ctx, timepend)
	if err != nil {
		return "", err
	}
	return res.InsertedID.(primitive.ObjectID).Hex(), nil
}

func (st MongoTimespendStore) GetTimespendByID(ctx context.Context, ownerid, id string) (types.Timespend, error) {
	var timespend types.Timespend
	if err := st.coll.FindOne(ctx, bson.M{"_id": ToObjectID(id), "ownerid": ownerid}).Decode(&timespend); err != nil {
		return types.Timespend{}, err
	}
	return timespend, nil
}

func (st MongoTimespendStore) DeleteTimespend(ctx context.Context, ownerid, id string) error {
	res, err := st.coll.DeleteOne(ctx, bson.M{"_id": ToObjectID(id), "ownerid": ownerid})
	if err != nil {
		return err
	}
	if res.DeletedCount <= 0 {
		return fmt.Errorf("can't delete timepend with id: %s", id)
	}
	return nil
}

func (st MongoTimespendStore) UpdateTimespend(ctx context.Context, ownerid, id string, params types.UpdateTimespendParams) error {
	paramsBson, err := toBsonDoc(params)
	if err != nil {
		return err
	}
	res, err := st.coll.UpdateOne(ctx, bson.M{"_id": ToObjectID(id), "ownerid": ownerid}, bson.D{{Key: "$set", Value: paramsBson}})
	if err != nil {
		return err
	}
	if res.MatchedCount <= 0 {
		return fmt.Errorf("can't update timespend with id: %s", id)
	}
	if res.ModifiedCount <= 0 {
		return fmt.Errorf("no changes timespend with id: %s", id)
	}
	return nil
}

type MongoMoneyspendStore struct {
	coll *mongo.Collection
}

func NewMongoMoneyspendStore(cl *mongo.Client, dbname, timespendColl string) *MongoMoneyspendStore {
	return &MongoMoneyspendStore{
		coll: cl.Database(dbname).Collection(timespendColl),
	}
}

func (st MongoMoneyspendStore) Drop(ctx context.Context) error {
	return st.coll.Drop(ctx)
}

func (st MongoMoneyspendStore) GetAllMonies(ctx context.Context, ownerID string) ([]types.Moneyspend, error) {
	cur, err := st.coll.Find(ctx, bson.M{"ownerid": ownerID})
	if err != nil {
		return nil, err
	}
	monies := []types.Moneyspend{}
	err = cur.All(ctx, &monies)
	if err != nil {
		return nil, err
	}
	return monies, nil
}

func (st MongoMoneyspendStore) CreateMoneyspend(ctx context.Context, moneyspend types.Moneyspend) (string, error) {
	res, err := st.coll.InsertOne(ctx, moneyspend)
	if err != nil {
		return "", err
	}
	return res.InsertedID.(primitive.ObjectID).Hex(), nil
}

func (st MongoMoneyspendStore) GetMoneyspendByID(ctx context.Context, ownerid, id string) (types.Moneyspend, error) {
	var moneyspend types.Moneyspend
	if err := st.coll.FindOne(ctx, bson.M{"_id": ToObjectID(id), "ownerid": ownerid}).Decode(&moneyspend); err != nil {
		return types.Moneyspend{}, err
	}
	return moneyspend, nil
}

func (st MongoMoneyspendStore) DeleteMoneyspend(ctx context.Context, ownerid, id string) error {
	res, err := st.coll.DeleteOne(ctx, bson.M{"_id": ToObjectID(id), "ownerid": ownerid})
	if err != nil {
		return err
	}
	if res.DeletedCount <= 0 {
		return fmt.Errorf("can't delete timepend with id: %s", id)
	}
	return nil
}

func (st MongoMoneyspendStore) UpdateMoneyspend(ctx context.Context, ownerid, id string, params types.UpdateMoneyspendParams) error {
	paramsBson, err := toBsonDoc(params)
	if err != nil {
		return err
	}
	res, err := st.coll.UpdateOne(ctx, bson.M{"_id": ToObjectID(id), "ownerid": ownerid}, bson.D{{Key: "$set", Value: paramsBson}})
	if err != nil {
		return err
	}
	if res.MatchedCount <= 0 {
		return fmt.Errorf("can't update timespend with id: %s", id)
	}
	if res.ModifiedCount <= 0 {
		return fmt.Errorf("no changes timespend with id: %s", id)
	}
	return nil
}
