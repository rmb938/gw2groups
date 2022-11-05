package api

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/rmb938/gw2groups/pkg/api_clients/gw2"
	"github.com/rmb938/gw2groups/pkg/api_clients/playfab"
	"k8s.io/utils/pointer"
)

type contextKey string

func (c contextKey) String() string {
	return "context key " + string(c)
}

var (
	contextKeyGW2APIClient  = contextKey("gw2-api-client")
	contextKeyGW2APIAccount = contextKey("gw2-api-account")

	contextPlayfabClient        = contextKey("playfab-api-client")
	contextKeyPlayfabAPIAccount = contextKey("playfab-api-account")
)

func HTTPRouter() *chi.Mux {
	chiRouter := chi.NewRouter()

	playFabClient := playfab.NewPlayFabClient()

	chiRouter.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			authorization, ok := r.Header["Authorization"]
			if !ok {
				http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
				return
			}

			authorizationParts := strings.Split(authorization[0], " ")
			if authorizationParts[0] != "Bearer" {
				http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
				return
			}

			// TODO: we probably want to cache this so we aren't checking every request
			gw2Client := gw2.NewGW2APIClient(authorizationParts[1])
			ctx = context.WithValue(ctx, contextKeyGW2APIClient, gw2Client)
			account, err := gw2Client.GetAccount(ctx)
			if err != nil {
				log.Printf("error getting gw2 account %s\n", err)
				// TODO: check if instance of error and return correct response
				http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
				return
			}
			ctx = context.WithValue(ctx, contextKeyGW2APIAccount, account)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	})

	chiRouter.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			ctx = context.WithValue(ctx, contextPlayfabClient, playFabClient)

			// TODO: we probably want to cache this so we aren't getting the playfab account every request
			playfabAccount, err := playFabClient.LoginWithCustomID(ctx, &playfab.ServerLoginWithCustomIDRequest{
				CreateAccount:  pointer.Bool(true),
				ServerCustomId: fmt.Sprintf("gw2-api-%s", ctx.Value(contextKeyGW2APIAccount).(*gw2.AccountResponse).ID),
				InfoRequestParameters: &playfab.PlayerCombinedInfoRequestParams{
					GetUserData: pointer.Bool(true),
				},
			})
			if err != nil {
				log.Printf("error getting playfab account %s\n", err)
				// TODO: check if instance of error and return correct response
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return
			}

			ctx = context.WithValue(ctx, contextKeyPlayfabAPIAccount, playfabAccount)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	})

	return chiRouter
}
