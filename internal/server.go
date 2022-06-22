package internal

import (
	db2 "OnlineSchool/internal/DataBase"
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

func StartServer() {
	//client, err := mongodb.NewClient("mongodb+srv://cluster0.lbets.mongodb.net/myFirstDatabase", "Mathew", "829079")
	/*client, err := mongodb.NewClient("mongodb://adm:adm@127.0.0.1:27017")
	//client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://127.0.0.1:27017"))
	if err != nil {
		log.Fatal(err.Error())

	}*/
	const uri = "mongodb://adm:adm@127.0.0.1:27017/?maxPoolSize=20&w=majority"
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatal(err.Error())
	}
	fmt.Println("Connected to MongoDB!")
	name := "test"
	db := client.Database(name)
	d := db2.NewUserDB(db)
	d.GetCourses(context.TODO(), 1)

}
