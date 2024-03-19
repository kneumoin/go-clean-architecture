package mongo

import (
	"context"
	"github.com/kneumoin/go-clean-architecture/link"
	"github.com/kneumoin/go-clean-architecture/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Link struct {
	ID     primitive.ObjectID `bson:"_id,omitempty"`
	UserID primitive.ObjectID `bson:"userId"`
	Secret string             `bson:"secret"`
}

type LinkRepository struct {
	db *mongo.Collection
}

func NewLinkRepository(db *mongo.Database, collection string) *LinkRepository {
	return &LinkRepository{
		db: db.Collection(collection),
	}
}

func (r LinkRepository) CreateLink(ctx context.Context, user *models.User, secret string) (*models.Link, error) {
	bm := &models.Link{
		Secret: secret,
		UserID: user.ID,
	}

	model := toModel(bm)
	res, err := r.db.InsertOne(ctx, model)
	if err != nil {
		return nil, err
	}

	bm.ID = res.InsertedID.(primitive.ObjectID).Hex()
	return bm, nil
}

func (r LinkRepository) GetLink(ctx context.Context, user *models.User, id string) (*models.Link, error) {
	objID, _ := primitive.ObjectIDFromHex(id)
	uID, _ := primitive.ObjectIDFromHex(user.ID)
	cur := r.db.FindOne(ctx, bson.M{"_id": objID, "userId": uID})

	// TODO Check if the result was found
	if err := cur.Err(); err != nil {
		return nil, link.ErrLinkNotFound
	}

	out := new(Link)
	if err := cur.Decode(out); err != nil {
		return nil, err
	}

	return toLink(out), nil
}

func (r LinkRepository) DeleteLink(ctx context.Context, user *models.User, id string) error {
	objID, _ := primitive.ObjectIDFromHex(id)
	uID, _ := primitive.ObjectIDFromHex(user.ID)
	_, err := r.db.DeleteOne(ctx, bson.M{"_id": objID, "userId": uID})
	return err
}

func toModel(b *models.Link) *Link {
	uid, _ := primitive.ObjectIDFromHex(b.UserID)

	return &Link{
		UserID: uid,
		Secret: b.Secret,
	}
}

func toLink(b *Link) *models.Link {
	return &models.Link{
		ID:     b.ID.Hex(),
		UserID: b.UserID.Hex(),
		Secret: b.Secret,
	}
}
