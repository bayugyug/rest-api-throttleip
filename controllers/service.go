package controllers

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/bayugyug/rest-api-throttleip/config"
	"github.com/bayugyug/rest-api-throttleip/driver"
	"github.com/bayugyug/rest-api-throttleip/models"
	redis "gopkg.in/redis.v3"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/jwtauth"
	"github.com/go-chi/render"
)

const (
	svcOptionWithHandler   = "svc-opts-handler"
	svcOptionWithAddress   = "svc-opts-address"
	svcOptionWithRedisHost = "svc-opts-redis-host"
)

var ApiInstance *ApiService

type ApiService struct {
	Api        *ApiHandler
	Router     *chi.Mux
	Address    string
	RedisHost  string
	RedisCache *redis.Client
	Context    context.Context
	IPHistory  *models.TrackerIPHistory
}

//WithSvcOptHandler opts for handler
func WithSvcOptHandler(r *ApiHandler) *config.Option {
	return config.NewOption(svcOptionWithHandler, r)
}

//WithSvcOptAddress opts for port#
func WithSvcOptAddress(r string) *config.Option {
	return config.NewOption(svcOptionWithAddress, r)
}

//WithSvcOptRedisHost opts for db connector
func WithSvcOptRedisHost(r string) *config.Option {
	return config.NewOption(svcOptionWithRedisHost, r)
}

//NewApiService service new instance
func NewApiService(opts ...*config.Option) (*ApiService, error) {

	//default
	svc := &ApiService{
		Address: ":8989",
		Api:     &ApiHandler{},
		Context: context.Background(),
	}

	//add options if any
	for _, o := range opts {
		//chk opt-name
		switch o.Name() {
		case svcOptionWithHandler:
			if s, oks := o.Value().(*ApiHandler); oks && s != nil {
				svc.Api = s
			}
		case svcOptionWithAddress:
			if s, oks := o.Value().(string); oks && s != "" {
				svc.Address = s
			}
		case svcOptionWithRedisHost:
			if s, oks := o.Value().(string); oks && s != "" {
				svc.RedisHost = s
			}
		}
	} //iterate all opts

	//set the actual router
	svc.Router = svc.MapRoute()

	//get db
	client, err := driver.NewRedisConnector(svc.RedisHost)
	if err != nil {
		return svc, err
	}

	//save
	svc.RedisCache = client

	//q manager
	isready := make(chan bool, 1)
	svc.IPHistory = models.NewTrackerIPHistory()
	go svc.IPHistory.ManageQ(isready)
	<-isready

	isreadySave := make(chan bool, 1)
	go svc.IPHistory.ManageHistory(isreadySave, svc.RedisCache)
	<-isreadySave

	//good :-)
	return svc, nil
}

//Run run the http server based on settings
func (svc *ApiService) Run() {

	//gracious timing
	srv := &http.Server{
		Addr:         svc.Address,
		Handler:      svc.Router,
		ReadTimeout:  60 * time.Second,
		WriteTimeout: 60 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	//async run
	go func() {
		log.Println("Listening on port", svc.Address)
		if err := srv.ListenAndServe(); err != nil {
			log.Printf("listen: %s\n", err)
			os.Exit(0)
		}

	}()

	//watcher
	stopChan := make(chan os.Signal)
	signal.Notify(stopChan, os.Interrupt)

	<-stopChan
	log.Println("Shutting down service...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	srv.Shutdown(ctx)
	defer cancel()
	log.Println("Server gracefully stopped!")
}

//MapRoute route map all endpoints
func (svc *ApiService) MapRoute() *chi.Mux {

	// Multiplexer
	router := chi.NewRouter()

	// Basic settings
	router.Use(
		render.SetContentType(render.ContentTypeJSON),
		middleware.Logger,
		middleware.DefaultCompress,
		middleware.StripSlashes,
		middleware.Recoverer,
		middleware.RequestID,
		middleware.RealIP,
	)

	// Basic gracious timing
	router.Use(middleware.Timeout(60 * time.Second))

	// Basic CORS
	cors := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		AllowCredentials: true,
	})

	router.Use(cors.Handler)

	router.Get("/", svc.Api.IndexPage)

	/*
		@end-points

		GET     /v1/api/request/{dummy}
		POST    /v1/api/request/{dummy}
		PUT 	/v1/api/request/{dummy}
		DELETE  /v1/api/request/{dummy}




	*/

	//end-points-mapping
	router.Route("/v1", func(r chi.Router) {
		r.Mount("/api/request",
			func(api *ApiHandler) *chi.Mux {
				sr := chi.NewRouter()
				sr.Post("/{dummy}", api.DummyReqPost)
				sr.Put("/{dummy}", api.DummyReqPut)
				sr.Get("/{dummy}", api.DummyReqGet)
				sr.Delete("/{dummy}", api.DummyReqDelete)
				return sr
			}(svc.Api))
	})

	return router
}

//SetContextKeyVal version context
func (svc *ApiService) SetContextKeyVal(k, v string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			r = r.WithContext(context.WithValue(r.Context(), k, v))
			next.ServeHTTP(w, r)
		})
	}
}

//BearerChecker check token
func (svc *ApiService) BearerChecker(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, _, err := jwtauth.FromContext(r.Context())

		if err != nil {
			switch err {
			default:
				log.Println("ERROR:", err)
				svc.Api.ReplyErrContent(w, r, http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
				return
			case jwtauth.ErrExpired:
				log.Println("ERROR: Expired")
				http.Error(w, "Expired", http.StatusUnauthorized)
				svc.Api.ReplyErrContent(w, r, http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
				return
			case jwtauth.ErrUnauthorized:
				log.Println("ERROR: ErrUnauthorized")
				svc.Api.ReplyErrContent(w, r, http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
				return
			}
		}

		if token == nil || !token.Valid {
			svc.Api.ReplyErrContent(w, r, http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
			return
		}

		// Token is authenticated, pass it through
		next.ServeHTTP(w, r)
	})

}
