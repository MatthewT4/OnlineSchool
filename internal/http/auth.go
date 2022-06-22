package http

import (
	"context"
	"net/http"
)

const (
	UserId string = "UserId"
)

func (rou *Router) UserAuthentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		/*cookie, er := r.Cookie("authTokentest")
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
		user_id := 1
		ctxt := context.WithValue(r.Context(), UserId, user_id)
		r = r.WithContext(ctxt)
		next.ServeHTTP(w, r)

	})
}
