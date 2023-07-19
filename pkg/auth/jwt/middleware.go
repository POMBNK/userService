package jwt

import (
	"context"
	"net/http"
)

func Middleware(h http.HandlerFunc, secretJWT string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		helper := NewTokenizer(secretJWT)

		cook, err := r.Cookie(AccessTokenName)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("no cookie"))
			return
		}
		jwtToken := cook.Value

		tokenMC, err := helper.ParseToken(jwtToken)
		if err != nil {
			cook, err = r.Cookie(RefreshTokenName)
			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte("bad access cookie. no refresh cookie"))
				return
			}
			refreshT := cook.Value

			mapClaims, err := helper.ParseToken(refreshT)
			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte("bad access and refresh cookies"))
				return
			}
			claims := helper.ParseMapClaims(mapClaims)

			pair, err := helper.GeneratePair(claims.UserID)
			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte("bad access and refresh cookies"))
				return
			}

			accessCook, refreshCook := helper.PrepareCookies(pair)
			http.SetCookie(w, accessCook)
			http.SetCookie(w, refreshCook)

			jwtToken = pair.AccessToken
			tokenMC, err = helper.ParseToken(jwtToken)
		}

		tokenClaims := helper.ParseMapClaims(tokenMC)

		ctx := context.WithValue(r.Context(), "user_id", tokenClaims.UserID)
		h(w, r.WithContext(ctx))
	}
}
