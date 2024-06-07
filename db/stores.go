package db

import (
	"github.com/SpectralJager/spender/types"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserStore interface {
	Dropper
	BaseCRUDStore[types.User]
}

type SpendStor[T any] interface {
	BaseCRUDStore[T]
}

type MongoUserStore struct {
	DefaultMongoStore[types.User]
}

func NewMongoUserStore(cl *mongo.Client, dbname, collname string) *MongoUserStore {
	return &MongoUserStore{
		DefaultMongoStore: NewDefaultMongoStore[types.User](
			cl.Database(dbname).Collection(collname),
		),
	}
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
