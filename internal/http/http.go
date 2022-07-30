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
	BLogic blogic.IBLogic
}

func NewRouter(db *mongo.Database) *Router {
	return &Router{BLogic: blogic.NewBLogic(db)}
}

func (r *Router) Start() {
	ro := mux.NewRouter()
	rAuth := ro.PathPrefix("/auth").Subrouter()
	rAuth.HandleFunc("/login", r.Login)
	rService := ro.PathPrefix("/service").Subrouter()
	rService.HandleFunc("/available_periods", r.AvailablePaymentPeriods)
	rService.HandleFunc("/create_payment", r.CreatePayment)

	rou := ro.PathPrefix("/").Subrouter()
	rou.HandleFunc("/get_courses", r.GetCourses)
	rou.HandleFunc("/get_next_webinars", r.GetNextWebinars)
	rou.HandleFunc("/get_today_webinars", r.GetTodayWebinars)
	rou.HandleFunc("/get_past_webinars", r.GetPastWebinars)
	rou.HandleFunc("/get_homework", r.GetHomework)
	rou.HandleFunc("/get_next_course_homeworks", r.GetNextCourseHomeworks)
	rou.HandleFunc("/get_past_course_homeworks", r.GetPastCourseHomeworks)
	rou.HandleFunc("/get_next_homeworks", r.GetNextHomeworks)
	rou.HandleFunc("/info_course", r.GetInfoCourse)
	rou.HandleFunc("/submit_homework", r.SubmitHomework)
	rou.Use(r.UserAuthentication)

	srv := &http.Server{
		Handler: ro,
		Addr:    ":80",
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}
