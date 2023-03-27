package mongo

import (
	"context"
	"m3game/plugins/db"
	"m3game/plugins/log"
	"m3game/runtime/plugin"

	"github.com/go-playground/validator/v10"
	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"google.golang.org/protobuf/proto"
)

var (
	_mongodb *MongoDB
	_factory = &Factory{}
	_codec   = &db.CacheDBCodec{}
)

const (
	_name = "db_mongo"
)

func init() {
	plugin.RegisterFactory(_factory)
}

type MongoCfg struct {
	URI      string `mapstructure:"URI" validate:"required"`
	DataBase string `mapstructure:"DataBase" validate:"required"`
}

type Factory struct {
}

func (f *Factory) Type() plugin.Type {
	return plugin.DB
}
func (f *Factory) Name() string {
	return _name
}

func (f *Factory) Setup(ctx context.Context, c map[string]interface{}) (plugin.PluginIns, error) {
	if _mongodb != nil {
		return _mongodb, nil
	}
	var cfg MongoCfg
	if err := mapstructure.Decode(c, &cfg); err != nil {
		return nil, errors.Wrap(err, "RedisDB Decode Cfg")
	}
	validate := validator.New()
	if err := validate.Struct(&cfg); err != nil {
		return nil, err
	}
	_mongodb = &MongoDB{}
	if client, err := mongo.Connect(ctx, options.Client().ApplyURI(cfg.URI)); err != nil {
		return nil, err
	} else {
		_mongodb.client = client
	}
	if err := _mongodb.client.Ping(ctx, readpref.Primary()); err != nil {
		return nil, err
	}
	_mongodb.database = _mongodb.client.Database(cfg.DataBase)
	if _, err := db.New(_mongodb); err != nil {
		return nil, err
	}
	return _mongodb, nil
}

func (f *Factory) Destroy(plugin.PluginIns) error {
	return nil
}

func (f *Factory) Reload(plugin.PluginIns, map[string]interface{}) error {
	return nil
}

func (f *Factory) CanUnload(plugin.PluginIns) bool {
	return false
}

type MongoDB struct {
	client   *mongo.Client
	database *mongo.Database
}

func (c *MongoDB) Factory() plugin.Factory {
	return _factory
}

func (c *MongoDB) Read(ctx context.Context, meta db.DBMetaInter, key interface{}, flags ...int32) (proto.Message, error) {
	coll := c.database.Collection(meta.Table())
	if len(flags) == 0 {
		flags = meta.AllFlags()
	}
	projection := bson.M{}
	for _, flag := range flags {
		projection[meta.FlagName(flag)] = 1
	}
	obj := meta.New()
	if err := coll.FindOne(ctx, bson.M{meta.FlagName(meta.KeyFlag()): key}, options.FindOne().SetProjection(projection)).Decode(obj); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, db.Err_KeyNotFound
		}
		return nil, err
	}
	return obj, nil
}
func (c *MongoDB) Update(ctx context.Context, meta db.DBMetaInter, key interface{}, obj proto.Message, flags ...int32) error {
	coll := c.database.Collection(meta.Table())
	if len(flags) == 0 {
		flags = meta.AllFlags()
	}
	update := bson.M{}
	for _, flag := range flags {
		update[meta.FlagName(flag)] = meta.Getter(obj, flag)
	}
	_, err := coll.UpdateOne(ctx, bson.M{meta.FlagName(meta.KeyFlag()): key}, bson.M{"$set": update})
	return err
}

func (c *MongoDB) Create(ctx context.Context, meta db.DBMetaInter, key interface{}, obj proto.Message) error {
	coll := c.database.Collection(meta.Table())
	insert := bson.M{}
	for _, flag := range meta.AllFlags() {
		insert[meta.FlagName(flag)] = meta.Getter(obj, flag)
	}
	if _, err := coll.InsertOne(ctx, insert); err != nil {
		log.Error("%v ", err)
		return err
	}
	return nil
}

func (c *MongoDB) Delete(ctx context.Context, meta db.DBMetaInter, key interface{}) error {
	coll := c.database.Collection(meta.Table())
	_, err := coll.DeleteOne(ctx, bson.M{meta.FlagName(meta.KeyFlag()): key})
	return err
}

func (c *MongoDB) ReadMany(ctx context.Context, meta db.DBMetaInter, filters interface{}, flags ...int32) ([]proto.Message, error) {
	coll := c.database.Collection(meta.Table())
	if len(flags) == 0 {
		flags = meta.AllFlags()
	}
	projection := bson.M{}
	for _, flag := range flags {
		projection[meta.FlagName(flag)] = 1
	}
	cur, err := coll.Find(ctx, filters, options.Find().SetProjection(projection))
	if err != nil {
		return nil, err
	}
	var objs []proto.Message
	for cur.Next(ctx) {
		obj := meta.New()
		if err := cur.Decode(obj); err != nil {
			log.Error("%s", err.Error())
			continue
		}
		objs = append(objs, obj)
	}
	return objs, nil
}
