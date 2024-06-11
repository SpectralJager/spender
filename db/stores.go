package db

import (
	"context"

	"github.com/SpectralJager/spender/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserStore interface {
	Dropper
	GetStorer[types.User]
	CreateStorer[types.User]
	DeleteStorer
	UpdateStorer

	GetByEmail(ctx context.Context, email string) (types.User, error)
}

type SpendStor[T any] interface {
	BaseCRUDStore[T]
}

type MongoUserStore struct {
	DefaultMongoDropStore
	DefaultMongoGetStore[types.User]
	DefaultMongoCreateStore[types.User]
	DefaultMongoUpdateStore
	DefaultMongoDeleteStore
	coll *mongo.Collection
}

func NewMongoUserStore(cl *mongo.Client, dbname, collname string) *MongoUserStore {
	coll := cl.Database(dbname).Collection(collname)
	return &MongoUserStore{
		DefaultMongoDropStore:   DefaultMongoDropStore{coll},
		DefaultMongoGetStore:    DefaultMongoGetStore[types.User]{coll},
		DefaultMongoCreateStore: DefaultMongoCreateStore[types.User]{coll},
		DefaultMongoUpdateStore: DefaultMongoUpdateStore{coll},
		DefaultMongoDeleteStore: DefaultMongoDeleteStore{coll},
		coll:                    coll,
	}
}

func (st MongoUserStore) GetByEmail(ctx context.Context, email string) (types.User, error) {
	res := st.coll.FindOne(ctx, bson.M{"email": email})
	var user types.User
	err := res.Decode(&user)
	if err != nil {
		return types.User{}, err
	}
	return user, nil
}

type MongoTimespendStore struct {
	DefaultMongoStore[types.Timespend]
}

func NewMongoTimespendStore(cl *mongo.Client, dbname string, collname string) MongoTimespendStore {
	return MongoTimespendStore{
		DefaultMongoStore: NewDefaultMongoStore[types.Timespend](
			cl.Database(dbname).Collection(collname),
		),
	}
}

type MongoMoneyspendStore struct {
	DefaultMongoStore[types.Moneyspend]
}

func NewMongoMoneyspendStore(cl *mongo.Client, dbname string, collname string) MongoMoneyspendStore {
	return MongoMoneyspendStore{
		DefaultMongoStore: NewDefaultMongoStore[types.Moneyspend](
			cl.Database(dbname).Collection(collname),
		),
	}
}
