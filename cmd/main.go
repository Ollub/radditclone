package main

import (
	"github.com/gorilla/mux"
	"golang-stepik-2022q1/reditclone/config"
	"golang-stepik-2022q1/reditclone/pkg/db"
	"golang-stepik-2022q1/reditclone/pkg/handlers"
	"golang-stepik-2022q1/reditclone/pkg/log"
	"golang-stepik-2022q1/reditclone/pkg/middleware"
	"golang-stepik-2022q1/reditclone/pkg/posts/delivery"
	post_repo "golang-stepik-2022q1/reditclone/pkg/posts/repo"
	post_uc "golang-stepik-2022q1/reditclone/pkg/posts/usecase"
	session_repo "golang-stepik-2022q1/reditclone/pkg/session/repo"
	session_uc "golang-stepik-2022q1/reditclone/pkg/session/usecase"
	user_delivery "golang-stepik-2022q1/reditclone/pkg/users/delivery"
	user_repo "golang-stepik-2022q1/reditclone/pkg/users/repo"
	user_uc "golang-stepik-2022q1/reditclone/pkg/users/usecase"
	"net/http"
	"time"
)

func NewServer(addr string) http.Server {

	postRepo := post_repo.NewMongoRepo(db.NewMongo())
	postManager := post_uc.NewManager(postRepo)
	postHandler := delivery.NewHandler(postManager)

	userRepo := user_repo.NewSql(db.GetPostgres())
	userManager := user_uc.NewManager(userRepo)

	sessionRepo := session_repo.NewRedis(db.NewRedis())
	sessionManager := session_uc.NewManager(sessionRepo)
	userHandler := user_delivery.NewHandler(userManager, sessionManager)

	apiHandler := mux.NewRouter()
	auth := middleware.Authentication(sessionManager)

	apiHandler.HandleFunc("/api/register", userHandler.Register).Methods("POST")
	apiHandler.HandleFunc("/api/login", userHandler.Login).Methods("POST")
	// POSTS
	apiHandler.HandleFunc("/api/post/{id}", postHandler.Get).Methods("GET")
	apiHandler.HandleFunc("/api/posts/", postHandler.List).Methods("GET")
	apiHandler.HandleFunc("/api/users/{username}", postHandler.GetByUser).Methods("GET")
	apiHandler.Handle("/api/posts", auth(http.HandlerFunc(postHandler.Create))).Methods("POST")
	apiHandler.Handle("/api/post/{postId}", auth(http.HandlerFunc(postHandler.Delete))).Methods("DELETE")
	// POST COMMENTS
	apiHandler.Handle("/api/post/{postId}", auth(http.HandlerFunc(postHandler.AddComment))).Methods("POST")
	apiHandler.Handle("/api/post/{postId}/{commentId}", auth(http.HandlerFunc(postHandler.DeleteComment))).Methods("DELETE")
	// POST VOTES
	apiHandler.Handle("/api/post/{postId}/upvote", auth(http.HandlerFunc(postHandler.Upvote))).Methods("GET")
	apiHandler.Handle("/api/post/{postId}/downvote", auth(http.HandlerFunc(postHandler.Downvote))).Methods("GET")
	apiHandler.Handle("/api/post/{postId}/unvote", auth(http.HandlerFunc(postHandler.Unvote))).Methods("GET")

	apiHandler.Use(
		middleware.SetupReqID,
		middleware.InjectLogger,
		middleware.SetupAccessLog,
	)

	siteMux := http.NewServeMux()
	siteMux.Handle("/", handlers.StaticHandler)
	siteMux.Handle("/api/", apiHandler)

	return http.Server{
		Addr:         addr,
		Handler:      siteMux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
}

func Init() {
	config.Load()

	logLevel := log.InfoLevel
	if config.Cfg.Debug {
		logLevel = log.DebugLevel
	}
	log.SetupLogger(logLevel)
}

func main() {
	Init()
	server := NewServer(":8008")
	log.Info("Start server")
	server.ListenAndServe()
}
