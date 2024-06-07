package db

import (
	"context"
	"fmt"

	"github.com/SpectralJager/spender/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Dropper interface {
	Drop(ctx context.Context) error
}

type Updater interface {
	ToBsonDoc() (*bson.D, error)
}

type AllGetStorer[T any] interface {
	GetAll(context.Context) ([]T, error)
}

type GetStorer[T any] interface {
	GetByID(context.Context, string) (T, error)
}

type CreateStorer[T any] interface {
	Create(context.Context, T) (string, error)
}

type UpdateStorer interface {
	Update(context.Context, string, Updater) error
}

type DeleteStorer interface {
	Delete(context.Context, string) error
}

type BaseCRUDStore[T any] interface {
	AllGetStorer[T]
	GetStorer[T]
	CreateStorer[T]
	UpdateStorer
	DeleteStorer
}

type DefaultMongoStore[T any] struct {
	defaultMongoDropStore
	defaultMongoAllGetStore[T]
	defaultMongoGetStore[T]
	defaultMongoCreateStore[T]
	defaultMongoUpdateStore
	defaultMongoDeleteStore
}

func NewDefaultMongoStore[T any](coll *mongo.Collection) DefaultMongoStore[T] {
	return DefaultMongoStore[T]{
		defaultMongoDropStore:   defaultMongoDropStore{coll},
		defaultMongoAllGetStore: defaultMongoAllGetStore[T]{coll},
		defaultMongoGetStore:    defaultMongoGetStore[T]{coll},
		defaultMongoCreateStore: defaultMongoCreateStore[T]{coll},
		defaultMongoUpdateStore: defaultMongoUpdateStore{coll},
		defaultMongoDeleteStore: defaultMongoDeleteStore{coll},
	}
}

type defaultMongoDropStore struct {
	coll *mongo.Collection
}

func (st defaultMongoDropStore) Drop(ctx context.Context) error {
	return st.coll.Drop(ctx)
}

type defaultMongoAllGetStore[T any] struct {
	coll *mongo.Collection
}

func (st defaultMongoAllGetStore[T]) GetAll(ctx context.Context) ([]T, error) {
	filter := bson.M{}
	if ctxOwnerID, ok := ctx.Value("ownerID").(string); ok {
		filter["ownerid"] = ctxOwnerID
	}
	cur, err := st.coll.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	entities := []T{}
	err = cur.All(ctx, &entities)
	if err != nil {
		return nil, err
	}
	return entities, nil
}

type defaultMongoGetStore[T any] struct {
	coll *mongo.Collection
}

func (st defaultMongoGetStore[T]) GetByID(ctx context.Context, id string) (T, error) {
	var entity T
	filter := bson.M{}
	filter["_id"] = utils.ToObjectID(id)
	if ctxOwnerID, ok := ctx.Value("ownerID").(string); ok {
		filter["ownerid"] = ctxOwnerID
	}
	err := st.coll.FindOne(ctx, filter).Decode(&entity)
	return entity, err
}

type defaultMongoCreateStore[T any] struct {
	coll *mongo.Collection
}

func (st defaultMongoCreateStore[T]) Create(ctx context.Context, newEntity T) (string, error) {
	res, err := st.coll.InsertOne(ctx, newEntity)
	if err != nil {
		return "", err
	}
	return res.InsertedID.(primitive.ObjectID).Hex(), nil
}

type defaultMongoUpdateStore struct {
	coll *mongo.Collection
}

func (st defaultMongoUpdateStore) Update(ctx context.Context, id string, updater Updater) error {
	entityBson, err := updater.ToBsonDoc()
	if err != nil {
		return err
	}
	filter := bson.M{}
	filter["_id"] = utils.ToObjectID(id)
	if ctxOwnerID, ok := ctx.Value("ownerID").(string); ok {
		filter["ownerid"] = ctxOwnerID
	}
	res, err := st.coll.UpdateOne(ctx, filter, bson.D{{Key: "$set", Value: entityBson}})
	if err != nil {
		return err
	}
	if res.ModifiedCount <= 0 || res.MatchedCount <= 0 {
		return fmt.Errorf("no changes for entity with id = %s", id)
	}
	return nil
}

type defaultMongoDeleteStore struct {
	coll *mongo.Collection
}

func (st defaultMongoDeleteStore) Delete(ctx context.Context, id string) error {
	filter := bson.M{}
	filter["_id"] = utils.ToObjectID(id)
	if ctxOwnerID, ok := ctx.Value("ownerID").(string); ok {
		filter["ownerid"] = ctxOwnerID
	}
	res, err := st.coll.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}
	if res.DeletedCount <= 0 {
		return fmt.Errorf("can't delete entity with id = %s", id)
	}
	return nil
}
