package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/AndikaPrasetia/wash-shoe/internal/config"
	"github.com/AndikaPrasetia/wash-shoe/internal/sqlc/user"
	"github.com/AndikaPrasetia/wash-shoe/internal/delivery/controller"
	"github.com/AndikaPrasetia/wash-shoe/internal/middleware"
	"github.com/AndikaPrasetia/wash-shoe/internal/repository"
	"github.com/AndikaPrasetia/wash-shoe/internal/usecase"
	utils "github.com/AndikaPrasetia/wash-shoe/internal/utils/services"
	"github.com/gin-gonic/gin"
)

type Server struct {
	engine  *gin.Engine
	server  *http.Server
	jwtSvc  utils.JwtService
	dbPool  *pgxpool.Pool
	querier user.Querier
	authUC  usecase.AuthUserUsecase
	host    string
	port    string
}

func NewServer() *Server {
	cfg, err := config.NewConfig()
	if err != nil {
		panic(fmt.Errorf("failed to read config: %v", err.Error()))
	}

	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.Host,
		cfg.Port,
		cfg.Username,
		cfg.Password,
		cfg.Database,
	)

	dbPool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		panic(fmt.Errorf("failed to connect to database: %v", err))
	}

	queries := user.New(dbPool) // *user.Queries, implements user.Querier

	authRepo := repository.NewAuthUserRepo(queries)
	userRepo := repository.NewUserRepo(queries)
	authUC := usecase.NewAuthUserUsecase(authRepo, userRepo)

	// misalnya lanjutkan setup Server
	s := &Server{
		engine:  gin.Default(),
		jwtSvc:  utils.NewJwtService(cfg.TokenConfig),
		querier: queries,
		authUC:  authUC,
		host:    cfg.APIHost,
		port:    cfg.APIPort,
		dbPool:  dbPool,
	}
	return s
}

func (s *Server) initRoute() {
	publicGroup := s.engine.Group("/api/v1")

	controller.NewAuthController(s.authUC, publicGroup).Route()

	protectedGroup := s.engine.Group("/api/v1")
	authMiddleware := middleware.NewAuthMiddleware(s.jwtSvc)
	protectedGroup.Use(authMiddleware.Middleware())
}

func (s *Server) Run() {
	s.initRoute()

	addr := fmt.Sprintf("%s:%s", s.host, s.port)
	s.server = &http.Server{
		Addr:    addr,
		Handler: s.engine,
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	go func() {
		fmt.Printf("Server running on %s\n", s.host)
		if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			panic(fmt.Errorf("failed to start server: %v ", err))
		}
	}()

	<-quit
	s.dbPool.Close()
	fmt.Println("\nShutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := s.server.Shutdown(ctx); err != nil {
		fmt.Printf("Server forced to shutdown: %v", err)
	}

	fmt.Println("Server gracefully stopped 󱠡 ")
}
