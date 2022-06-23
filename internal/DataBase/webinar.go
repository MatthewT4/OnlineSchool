package DataBase

import (
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

type WebinarDB struct {
	collection *mongo.Collection
}
type IWebinarDB interface {
}

func NewWebinarDB(db *mongo.Database) *WebinarDB {
	return &WebinarDB{collection: db.Collection(nameWebinarDB)}
}
func (w *WebinarDB) GetNextWebinars(date time.Time) {

}
