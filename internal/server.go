package internal

import (
	"OnlineSchool/internal/http"
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

// user:iE9H4v7i626mZX.6
func StartServer() {
	//client, err := mongodb.NewClient("mongodb+srv://cluster0.lbets.mongodb.net/myFirstDatabase", "Mathew", "829079")
	/*client, err := mongodb.NewClient("mongodb://adm:adm@127.0.0.1:27017")
	//client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://127.0.0.1:27017"))
	if err != nil {
		log.Fatal(err.Error())

	}*/
	//const uri = "mongodb://adm:adm@127.0.0.1:27017/?maxPoolSize=20&w=majority"
	const uri = "mongodb://system:1hY8phB4f4q921a<@185.130.114.130:27017/production"
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatal(err.Error())
	}
	fmt.Println("Connected to MongoDB!")
	name := "production"
	db := client.Database(name)
	/*d := db2.NewUserDB(db)
	d.GetCourses(context.TODO(), 1)
	logic := blogic.NewBUser(db)
	code, res :=logic.GetCouses(1)
	fmt.Println(code, res)*/
	router := http.NewRouter(db)
	router.Start()
}
