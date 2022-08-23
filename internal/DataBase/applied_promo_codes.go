package DataBase

import (
	"OnlineSchool/internal/structs"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type AppliedPromoСodeDB struct {
	collection *mongo.Collection
}
type IAppliedPromoСodeDB interface {
	GetUserAppliedThePC(ctx context.Context, userId int64, promoCode string) ([]structs.ApplyPromoCode, error)
	AddHistoryElem(ctx context.Context, code structs.ApplyPromoCode) error
}

func NewAppliedPromoСodeDB(db *mongo.Database) *AppliedPromoСodeDB {
	return &AppliedPromoСodeDB{collection: db.Collection(nameAppliedPromoCodes)}
}

func (a *AppliedPromoСodeDB) GetUserAppliedThePC(ctx context.Context, userId int64, promoCode string) ([]structs.ApplyPromoCode, error) {
	filter := bson.M{
		"owner":      userId,
		"promo_code": promoCode,
	}
	var ret []structs.ApplyPromoCode
	cur, err := a.collection.Find(ctx, filter)
	if err != nil {
		return ret, err
	}
	for cur.Next(context.TODO()) {
		var elem structs.ApplyPromoCode
		err = cur.Decode(&elem)
		if err != nil {
			return ret, err
		}
		ret = append(ret, elem)
	}
	err = cur.Err()
	if err != nil {
		return ret, err
	}
	cur.Close(context.TODO())
	return ret, err
}

func (a *AppliedPromoСodeDB) AddHistoryElem(ctx context.Context, code structs.ApplyPromoCode) error {
	_, err := a.collection.InsertOne(ctx, code)
	return err
}
