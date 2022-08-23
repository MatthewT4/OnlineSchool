package http

import (
	"OnlineSchool/internal/blogic"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/acme/autocert"
	"net/http"
)

const (
	Domain     string = "https://lk.lyc15.ru"
	ServDomain string = "serv.lyc15.ru"
	//Domain string = "http://localhost:3000"
	//ServDomain string = "localhost"
)

type Router struct {
	BLogic blogic.IBLogic
}

func NewRouter(db *mongo.Database) *Router {
	return &Router{BLogic: blogic.NewBLogic(db)}
}

func (r *Router) Start() {
	ro := mux.NewRouter()

	cloudPaymentRou := ro.PathPrefix("/pay_serv").Subrouter()
	cloudPaymentRou.HandleFunc("/check", r.CheckPayment)
	cloudPaymentRou.HandleFunc("/register_approved_pay", r.RegisterApprovedPayment)

	rAuth := ro.PathPrefix("/auth").Subrouter()
	rAuth.HandleFunc("/login", r.Login)
	rService := ro.PathPrefix("/service").Subrouter()
	rService.HandleFunc("/available_periods", r.AvailablePaymentPeriods)
	rService.HandleFunc("/create_payment", r.CreatePayment)
	rService.HandleFunc("/check_auth", r.CheckAuth)
	rService.Use(r.servProm)

	rou := ro.PathPrefix("/").Subrouter()
	rou.HandleFunc("/cal_amount_in_promo_code", r.CalculateTotalAmountInPromoCode)
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
	rou.HandleFunc("/linking_payment", r.LinkingPaymentToUser)
	rou.HandleFunc("/connecting_groups", r.ConnectingCourseGroups)
	rou.HandleFunc("/invitation_vk_link", r.InvitationLinkVkGroup)
	rou.Use(r.UserAuthentication)

	http.Serve(autocert.NewListener("serv.lyc15.ru"), ro)
	/*
		srv := &http.Server{
			Handler: ro,
			Addr:    ":80",
			// Good practice: enforce timeouts for servers you create!
			WriteTimeout: 15 * time.Second,
			ReadTimeout:  15 * time.Second,
		}
		log.Fatal(srv.ListenAndServe())*/
}
