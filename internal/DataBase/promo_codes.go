package DataBase

import (
	"OnlineSchool/internal/structs"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type PromoСodeDB struct {
	collection *mongo.Collection
}
type IPromoСodeDB interface {
	GetPromoCode(ctx context.Context, promoCodes string) (structs.PromoCode, error)
}

func NewPromoСodeDB(db *mongo.Database) *PromoСodeDB {
	return &PromoСodeDB{collection: db.Collection(namePromoСodesDB)}
}

func (p *PromoСodeDB) GetPromoCode(ctx context.Context, promoCodes string) (structs.PromoCode, error) {
	filter := bson.M{
		"promo_code": promoCodes,
	}
	var pCode structs.PromoCode
	err := p.collection.FindOne(ctx, filter).Decode(&pCode)
	return pCode, err
}
