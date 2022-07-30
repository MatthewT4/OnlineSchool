package DataBase

import (
	"OnlineSchool/internal/structs"
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type PaymentDB struct {
	collection *mongo.Collection
}
type IPaymentDB interface {
	AddPayment(ctx context.Context, payment structs.Payment) (string, error)
	FindPayment(ctx context.Context, paymentID string) (structs.Payment, error)
}

func NewPaymentDBDB(db *mongo.Database) *PaymentDB {
	return &PaymentDB{collection: db.Collection(namePaymentDB)}
}

func (p *PaymentDB) AddPayment(ctx context.Context, payment structs.Payment) (string, error) {
	res, err := p.collection.InsertOne(ctx, payment)
	if err != nil {
		return "", err
	}
	ID := res.InsertedID.(primitive.ObjectID).Hex()
	filter := bson.M{
		"_id": res.InsertedID,
	}
	update := bson.M{
		"$set": bson.M{
			"payment_id": ID,
		},
	}
	upd, er := p.collection.UpdateOne(context.TODO(), filter, update)
	if er != nil {
		return "", er
	}
	if upd.ModifiedCount == 0 {
		return "", fmt.Errorf("No doc")
	}
	return ID, nil
}

func (p *PaymentDB) FindPayment(ctx context.Context, paymentID string) (structs.Payment, error) {
	filter := bson.M{
		"payment_id": paymentID,
	}
	var payment structs.Payment
	err := p.collection.FindOne(ctx, filter).Decode(&payment)
	return payment, err
}

func (p *PaymentDB) EditOwnerPayment(ctx context.Context, paymentId string, userId int64, historyElement structs.History) (int64, error) {
	filter := bson.M{
		"payment_id": paymentId,
	}
	update := bson.M{
		"$set": bson.M{
			"user_id": userId,
		},
		"$push": bson.M{
			"change_history": historyElement,
		},
	}

	res, err := p.collection.UpdateOne(ctx, filter, update)
	return res.ModifiedCount, err
}
