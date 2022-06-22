package http

import (
	"OnlineSchool/internal/blogic"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"net/http"
	"time"
)

type Router struct {
	BLUser blogic.IBUser
}

func NewRouter(db *mongo.Database) *Router {
	return &Router{BLUser: blogic.NewBUser(db)}
}

func (r *Router) Start() {
	rou := mux.NewRouter()
	rou.HandleFunc("/get_courses", r.GetCourses)

	rou.Use(r.UserAuthentication)

	srv := &http.Server{
		Handler: rou,
		Addr:    ":80",
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}
