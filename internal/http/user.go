package http

import "net/http"

func (rou *Router) GetCourses(w http.ResponseWriter, r *http.Request) {
	user_id := r.Context().Value(UserId).(int)

	code, body := rou.BLogic.GetUserCourses(user_id)
	if code != 200 {
		http.Error(w, body, code)
		return
	}
	w.Write([]byte(body))
}

func (rou *Router) GetNextWebinars(w http.ResponseWriter, r *http.Request) {
	user_id := r.Context().Value(UserId).(int)
	course_id := 1
	code, body := rou.BLogic.GetNextWebinars(user_id, course_id)
	if code != 200 {
		http.Error(w, body, code)
		return
	}
	w.Write([]byte(body))
}
