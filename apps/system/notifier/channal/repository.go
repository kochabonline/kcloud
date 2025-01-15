package channal

import (
	"context"

	"github.com/kochabonline/kcloud/internal/util"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/mgo.v2/bson"
)

type Repository struct {
	coll *mongo.Collection
}

func NewRepository(client *mongo.Client) *Repository {
	return &Repository{
		coll: client.Database("kcloud").Collection("channal"),
	}
}

func (repo *Repository) Create(ctx context.Context, channal *Channal) error {
	_, err := repo.coll.InsertOne(ctx, channal)
	return err
}

func (repo *Repository) FindByName(ctx context.Context, name string) (*Channal, error) {
	var channal Channal
	if err := repo.coll.FindOne(ctx, bson.M{"name": name}).Decode(&channal); err != nil {
		return nil, err
	}

	return &channal, nil
}

func (repo *Repository) FindByApiKey(ctx context.Context, apiKey string) (*Channal, error) {
	var channal Channal
	if err := repo.coll.FindOne(ctx, bson.M{"api_key": apiKey}).Decode(&channal); err != nil {
		return nil, err
	}

	return &channal, nil
}

func (repo *Repository) FindAll(ctx context.Context, req *FindAllRequest) (*Channels, error) {
	var channals Channels

	// 根据请求参数构建查询条件
	filter := bson.M{}
	if req.ProviderType != "" {
		filter["provider_type"] = req.ProviderType
	}
	if req.Status != "" {
		filter["status"] = req.Status
	}
	if req.CreatedAt != 0 {
		filter["created_at"] = bson.M{"$gte": req.CreatedAt}
	}
	if req.Keyword != "" {
		filter["name"] = bson.M{"$regex": req.Keyword}
	}

	total, err := repo.coll.CountDocuments(ctx, filter)
	if err != nil {
		return nil, err
	}
	channals.Total = total

	offset, limit := util.Paginate(req.Page, req.Size)

	findOptions := options.Find()
	findOptions.SetSort(bson.M{"created_at": -1})
	findOptions.SetSkip(int64(offset))
	findOptions.SetLimit(int64(limit))

	cursor, err := repo.coll.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	if err := cursor.All(ctx, &channals.Items); err != nil {
		return nil, err
	}

	return &channals, nil
}

func (repo *Repository) Delete(ctx context.Context, id string) error {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	_, err = repo.coll.DeleteOne(ctx, bson.M{"_id": oid})

	return err
}
