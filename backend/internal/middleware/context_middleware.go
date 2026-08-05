package middleware

import (
	"fmt"
	"net/http"

	"alwis.dev/selectify/internal/model"
	"alwis.dev/selectify/internal/repo"
)

type requestContextKey uint

const (
	sessionKey requestContextKey = iota + 4783478488
)

type sessionContext struct {
	*model.UserSession
	Valid bool
}

//func Session(context context.Context) http.Handler {
//	if sCtx, ok := context.Value(sessionKey).(*sessionContext); ok {
//		fmt.Println()
//		fmt.Println()
//		fmt.Println()
//		fmt.Println("here is the error")
//		return sCtx.UserSession
//	} else {
//		return nil
//	}
//}

func ContextMiddleware(sessionRepo repo.UserSessionRepo) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			c, err := r.Cookie("slf")
			if err == nil {
				session, _ := sessionRepo.GetBySessionId(ctx, c.Value)
				if session != nil && !session.IsExpired() {

					fmt.Println()
					fmt.Println()
					fmt.Println(session.UserId)
					fmt.Println(session.IpAddress)
					//ctx = appctx.WithSession(ctx, session)

					//user, _ := userRepo.FindByID(ctx, session.UserID)
					//if user != nil {
					//	ctx = appctx.WithUser(ctx, user)
					//
					//	admin, _ := adminRepo.FindByUser(ctx, user)
					//	if admin != nil {
					//		ctx = appctx.WithAdmin(ctx, admin)
					//	}
					//}
				}
			}

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}

}
