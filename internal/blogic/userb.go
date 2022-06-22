package blogic

import (
	"OnlineSchool/internal/DataBase"
	"context"
	"encoding/json"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
)

type BUser struct {
	DBUser DataBase.IUserDB
}

func NewBUser(db *mongo.Database) *BUser {
	return &BUser{DBUser: DataBase.NewUserDB(db)}
}

type IBUser interface {
	GetCouses(user_id int) (int, string)
}

func (b *BUser) GetCouses(user_id int) (int, string) {
	res, err := b.DBUser.GetCourses(context.TODO(), user_id)
	if err != nil {
		fmt.Println(err.Error())
		return 404, "not found"
	}
	jr, err := json.Marshal(res)
	if err != nil {
		log.Fatal(err)
	}
	return 200, string(jr)
}
