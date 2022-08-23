package http

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

const (
	UserId string = "UserId"
)

func (rou *Router) UserAuthentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", Domain)
		//w.Header().Set("Access-Control-Allow-Origin", "https://lk.lyc15.ru")
		//w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "*")
		w.Header().Set("Access-Control-Allow-Credentials", "true") /*cookie, er := r.Cookie("authTokentest")
		if er != nil {
			http.Error(w, "not cookie", 401)
			return
		}
		group, pToken, err := rou.ser.UserSer.Authentication(cookie.Value)
		if err != nil {
			http.Error(w, "Authentication error"+err.Error(), 404)
		} else {
			ctx := context.WithValue(r.Context(), UserKey, pToken)
			r = r.WithContext(ctx)
			ctxt := context.WithValue(r.Context(), GroupKey, group)
			r = r.WithContext(ctxt)
			next.ServeHTTP(w, r)
		}*/
		cookie, er := r.Cookie("authToken")
		if er != nil {
			fmt.Println("not cookie")
			http.Error(w, "not cookie", 401)
			return
		}

		userId, _, err := rou.BLogic.Authentication(cookie.Value)
		if err != nil {
			fmt.Println("Authentication error: ", err.Error())
			http.Error(w, "not cookie", 401)
			return
		}
		ctxt := context.WithValue(r.Context(), UserId, userId)
		r = r.WithContext(ctxt)
		next.ServeHTTP(w, r)

	})
}

func (rou *Router) Login(w http.ResponseWriter, r *http.Request) {
	fmt.Println(1)
	var cVK string
	cVK = r.URL.Query().Get("code")

	var red string
	red = r.URL.Query().Get("redirect_uri")

	code, mes, token := rou.BLogic.Login(cVK, red)
	fmt.Println(2)
	w.Header().Set("Access-Control-Allow-Origin", Domain)
	//w.Header().Set("Access-Control-Allow-Origin", "https://lk.lyc15.ru")
	//w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "*")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	if code != 200 {
		fmt.Println("er 3")
		http.Error(w, string(mes), code)
		return
	}
	//cookie := http.Cookie{Name: "authToken", Value: token, Expires: time.Now().Add(time.Hour * 24 * 30), SameSite: 4, Secure: true, Path: "/", Domain: "serv.lyc15.ru"}
	cookie := http.Cookie{Name: "authToken", Value: token, Expires: time.Now().Add(time.Hour * 24 * 30), SameSite: 4, Secure: true, Path: "/", Domain: ServDomain}
	http.SetCookie(w, &cookie)
	var vr struct {
		Body string `json:"body"`
	}
	coc := http.Cookie{Name: "authToken", Value: token, Expires: time.Now().Add(time.Hour * 24 * 30), Path: "/", Domain: ServDomain} //https://serv.lyc15.ru
	vr.Body = coc.String()
	d, e := json.Marshal(&vr)
	if e != nil {
		http.Error(w, "server error", 500)
		return
	}
	w.Write(d)
}

func (rou *Router) CheckAuth(w http.ResponseWriter, r *http.Request) {

	cookie, er := r.Cookie("authToken")
	if er != nil {
		w.Write([]byte("{\"status\":false}"))
		return
	}

	_, _, err := rou.BLogic.Authentication(cookie.Value)
	if err != nil {
		w.Write([]byte("{\"status\":false}"))
		return
	}
	w.Write([]byte("{\"status\":true}"))
}
