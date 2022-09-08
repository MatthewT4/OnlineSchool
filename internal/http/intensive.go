package http

import (
	"net/http"
)

func (rou *Router) GetIntensive(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value(UserId).(int64)
	intensiveTag := r.URL.Query().Get("intensive_tag")
	code, mes := rou.BLogic.GetIntensive(intensiveTag, userId)
	if code/100 != 2 {
		http.Error(w, string(mes), code)
		return
	}
	w.Write(mes)
}

func (rou *Router) AddUserIntensive(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value(UserId).(int64)
	intensiveTag := r.URL.Query().Get("intensive_tag")
	code, mes := rou.BLogic.AddUserIntensive(intensiveTag, userId)
	http.Error(w, mes, code)
}
