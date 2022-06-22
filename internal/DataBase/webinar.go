package DataBase

import (
	"go.mongodb.org/mongo-driver/mongo"
)

type WebinarDB struct {
	collection *mongo.Collection
}
type IWebinarDB interface {
}

func NewWebinarDB(db *mongo.Database) *WebinarDB {
	return &WebinarDB{collection: db.Collection(nameWebinarDB)}
}
func (w *WebinarDB) GetWebinar(ff string) {

}
