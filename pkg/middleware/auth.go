package middleware

import (
	"context"
	"golang-stepik-2022q1/reditclone/pkg/log"
	"golang-stepik-2022q1/reditclone/pkg/session"
	sessionUC "golang-stepik-2022q1/reditclone/pkg/session/usecase"
	"golang-stepik-2022q1/reditclone/pkg/utils/http_utils"
	"net/http"
	"strings"
)

const AuthHeader = "Authorization"

func Authentication(sm *sessionUC.Manager) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			token := r.Header.Get(AuthHeader)
			if token == "" {
				log.Clog(ctx).Info("Authorization failed. No token provided")
				http_utils.HttpError(w, "Authorization failed", http.StatusUnauthorized)
				return
			}
			parts := strings.Fields(token)
			if len(parts) != 2 {
				log.Clog(ctx).Info("Authorization failed. Wrong token", log.Fields{"token": token})
				http_utils.HttpError(w, "Authorization failed", http.StatusUnauthorized)
				return
			}
			token = parts[1]
			sess, err := sm.Check(ctx, token)
			if err != nil {
				log.Clog(ctx).Info("Authorization failed", log.Fields{"error": err.Error(), "token": token})
				http_utils.HttpError(w, err.Error(), http.StatusUnauthorized)
				return
			}

			ctx = context.WithValue(ctx, session.SessionKey, sess)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
