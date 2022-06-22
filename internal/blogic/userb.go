package blogic

import (
	"OnlineSchool/internal/DataBase"
	"go.mongodb.org/mongo-driver/mongo"
)

type BUser struct {
	DBUser DataBase.IUserDB
}

func NewBUser(db *mongo.Database) *BUser {
	return &BUser{DBUser: DataBase.NewUserDB(db)}
}
