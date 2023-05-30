package app

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"redditclone/internal/comment"
	"redditclone/internal/config"
	"redditclone/internal/handlers"
	"redditclone/internal/middleware"
	"redditclone/internal/post"
	"redditclone/internal/session"
	"redditclone/internal/user"
	"redditclone/internal/vote"
	"syscall"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

func Run(configPath string) error {
	cfg, err := config.NewConfig(configPath)
	if err != nil {
		return fmt.Errorf("NewConfig failed: %w", err)
	}

	zapCfg := zap.NewProductionConfig()
	setLoggerLevel(&zapCfg, cfg.Logger.Level)
	zapLogger, err := zapCfg.Build()
	if err != nil {
		return fmt.Errorf("zap logger build failed: %w", err)
	}

	defer zapLogger.Sync()

	logger := zapLogger.Sugar()

	sm := session.NewJWTSessionManager(cfg.Session.SecretKey, cfg.Session.AccessTokenTTL)

	userRepo := user.NewMemoryRepo()
	postRepo := post.NewMemoryRepo()
	voteRepo := vote.NewMemoryRepo()
	commentRepo := comment.NewMemoryRepo()

	userHandler := &handlers.UserHandler{
		Logger:   logger,
		UserRepo: userRepo,
		Sessions: sm,
	}
	postHandler := &handlers.PostHandler{
		Logger:      logger,
		PostRepo:    postRepo,
		UserRepo:    userRepo,
		CommentRepo: commentRepo,
		VoteRepo:    voteRepo,
	}

	auth := middleware.NewAuthMiddleware(logger, sm)

	r := mux.NewRouter()

	api := r.PathPrefix("/api").Subrouter()
	api.HandleFunc("/register", userHandler.Register).Methods(http.MethodPost)
	api.HandleFunc("/login", userHandler.Login).Methods(http.MethodPost)

	api.Handle("/posts", auth(postHandler.AddPost)).Methods(http.MethodPost)
	api.HandleFunc("/posts/", postHandler.GetAll).Methods(http.MethodGet)
	api.HandleFunc("/posts/{category}", postHandler.GetByCategory).Methods(http.MethodGet)

	api.HandleFunc("/user/{username}", postHandler.GetByUser).Methods(http.MethodGet)
	api.HandleFunc("/post/{post_id:[0-9]+}", postHandler.GetByID).Methods(http.MethodGet)
	api.HandleFunc("/post/{post_id:[0-9]+}", auth(postHandler.DeletePost)).Methods(http.MethodDelete)

	api.HandleFunc("/post/{post_id:[0-9]+}/upvote", auth(postHandler.Upvote)).Methods(http.MethodGet)
	api.HandleFunc("/post/{post_id:[0-9]+}/unvote", auth(postHandler.Unvote)).Methods(http.MethodGet)
	api.HandleFunc("/post/{post_id:[0-9]+}/downvote", auth(postHandler.Downvote)).Methods(http.MethodGet)

	api.HandleFunc("/post/{post_id:[0-9]+}", auth(postHandler.AddComment)).Methods(http.MethodPost)
	api.HandleFunc("/post/{post_id:[0-9]+}/{comment_id:[0-9]+}", auth(postHandler.DeleteComment)).Methods(http.MethodDelete)

	r.PathPrefix("/static/").Handler(handlers.Static)
	r.PathPrefix("/").HandlerFunc(handlers.Index)

	mux := middleware.AccessLog(logger, r)
	mux = middleware.Panic(logger, mux)

	srv := &http.Server{Addr: ":" + cfg.HTTP.Port, Handler: mux}
	go func() {
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			logger.Errorf("ListenAndServe failed: %v", err)
		}
	}()
	logger.Infoln("Server is running on port " + cfg.HTTP.Port)

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	<-c

	logger.Infoln("Server is shutting down")
	if err := srv.Shutdown(context.Background()); err != nil {
		return fmt.Errorf("server shutdown error: %w", err)
	}
	logger.Infoln("Server shut down")
	return nil
}

func setLoggerLevel(zapCfg *zap.Config, level string) {
	switch level {
	case "fatal":
		zapCfg.Level = zap.NewAtomicLevelAt(zap.FatalLevel)
	case "panic":
		zapCfg.Level = zap.NewAtomicLevelAt(zap.PanicLevel)
	case "dpanic":
		zapCfg.Level = zap.NewAtomicLevelAt(zap.DPanicLevel)
	case "error":
		zapCfg.Level = zap.NewAtomicLevelAt(zap.ErrorLevel)
	case "warn":
		zapCfg.Level = zap.NewAtomicLevelAt(zap.WarnLevel)
	case "info":
		zapCfg.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
	case "debug":
		zapCfg.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	default:
		zapCfg.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
	}
}
