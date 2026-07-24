package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/labstack/echo/v4"
	echomw "github.com/labstack/echo/v4/middleware"
	"go.uber.org/dig"

	"github.com/99katedegree/bunkasairpg2/api/config"
	genapi "github.com/99katedegree/bunkasairpg2/api/gen/api"
	mw "github.com/99katedegree/bunkasairpg2/api/internal/adapter/middleware"
	"github.com/99katedegree/bunkasairpg2/api/internal/adapter/handler"
	"github.com/99katedegree/bunkasairpg2/api/internal/domain/battletoken"
	domainrepo "github.com/99katedegree/bunkasairpg2/api/internal/domain/repository"
	infradb "github.com/99katedegree/bunkasairpg2/api/internal/infrastructure/db"
	dbrepo "github.com/99katedegree/bunkasairpg2/api/internal/infrastructure/db/repository"
	"github.com/99katedegree/bunkasairpg2/api/internal/infrastructure/storage"
	"github.com/99katedegree/bunkasairpg2/api/internal/usecase"
)

func main() {
	// ── Config ──────────────────────────────────────────────────────────────
	cfg, err := config.Load()
	if err != nil {
		slog.Error("failed to load config", "err", err)
		os.Exit(1)
	}

	// ── Migration ───────────────────────────────────────────────────────────
	dsn := fmt.Sprintf("mysql://%s:%s@tcp(%s:%d)/%s?multiStatements=true",
		cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName)

	m, err := migrate.New("file://migrations", dsn)
	if err != nil {
		slog.Error("failed to create migrator", "err", err)
		os.Exit(1)
	}
	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		slog.Error("migration failed", "err", err)
		os.Exit(1)
	}
	slog.Info("migration completed")

	// ── DB ping ─────────────────────────────────────────────────────────────
	db, err := infradb.New(cfg)
	if err != nil {
		slog.Error("failed to connect to database", "err", err)
		os.Exit(1)
	}
	if err := db.PingContext(context.Background()); err != nil {
		slog.Error("database ping failed", "err", err)
		os.Exit(1)
	}
	slog.Info("database connected")

	// ── R2 ping ──────────────────────────────────────────────────────────────
	r2, err := storage.NewR2Client(cfg)
	if err != nil {
		slog.Error("failed to create R2 client", "err", err)
		os.Exit(1)
	}
	if err := r2.Ping(context.Background()); err != nil {
		slog.Error("R2 ping failed", "err", err)
		os.Exit(1)
	}
	slog.Info("R2 connected")

	// ── DI Container ────────────────────────────────────────────────────────
	c := dig.New()

	// Config
	if err := c.Provide(func() (*config.Config, error) {
		return cfg, nil
	}); err != nil {
		slog.Error("dig provide config", "err", err)
		os.Exit(1)
	}

	// DB (*sql.DB) — 起動前 ping 済みのインスタンスを再利用
	if err := c.Provide(func() *sql.DB { return db }); err != nil {
		slog.Error("dig provide db", "err", err)
		os.Exit(1)
	}

	// R2 Client — 起動前 ping 済みのインスタンスを再利用
	if err := c.Provide(func() *storage.R2Client { return r2 }); err != nil {
		slog.Error("dig provide r2client", "err", err)
		os.Exit(1)
	}

	// BattleToken
	if err := c.Provide(func() *battletoken.BattleToken {
		return battletoken.New(cfg.BattleTokenSecret)
	}); err != nil {
		slog.Error("dig provide battletoken", "err", err)
		os.Exit(1)
	}

	// Repositories — bridge concrete → domain interface
	if err := c.Provide(func(db *sql.DB) domainrepo.UserRepository {
		return dbrepo.NewUserRepository(db)
	}); err != nil {
		slog.Error("dig provide user repo", "err", err)
		os.Exit(1)
	}
	if err := c.Provide(func(db *sql.DB) domainrepo.WeaponRepository {
		return dbrepo.NewWeaponRepository(db)
	}); err != nil {
		slog.Error("dig provide weapon repo", "err", err)
		os.Exit(1)
	}
	if err := c.Provide(func(db *sql.DB) domainrepo.ItemRepository {
		return dbrepo.NewItemRepository(db)
	}); err != nil {
		slog.Error("dig provide item repo", "err", err)
		os.Exit(1)
	}
	if err := c.Provide(func(db *sql.DB) domainrepo.MonsterRepository {
		return dbrepo.NewMonsterRepository(db)
	}); err != nil {
		slog.Error("dig provide monster repo", "err", err)
		os.Exit(1)
	}
	if err := c.Provide(func(db *sql.DB) domainrepo.BattleRepository {
		return dbrepo.NewBattleRepository(db)
	}); err != nil {
		slog.Error("dig provide battle repo", "err", err)
		os.Exit(1)
	}
	if err := c.Provide(func(db *sql.DB) domainrepo.ImageRepository {
		return dbrepo.NewImageRepository(db)
	}); err != nil {
		slog.Error("dig provide image repo", "err", err)
		os.Exit(1)
	}
	if err := c.Provide(func(db *sql.DB) domainrepo.AdminRepository {
		return dbrepo.NewAdminRepository(db)
	}); err != nil {
		slog.Error("dig provide admin repo", "err", err)
		os.Exit(1)
	}

	// Usecases
	// AuthUsecase needs jwtSecret injected from config closure
	if err := c.Provide(func(userRepo domainrepo.UserRepository, adminRepo domainrepo.AdminRepository) *usecase.AuthUsecase {
		return usecase.NewAuthUsecase(userRepo, adminRepo, cfg.JWTSecret)
	}); err != nil {
		slog.Error("dig provide auth usecase", "err", err)
		os.Exit(1)
	}
	if err := c.Provide(usecase.NewImageUsecase); err != nil {
		slog.Error("dig provide image usecase", "err", err)
		os.Exit(1)
	}
	if err := c.Provide(usecase.NewMeUsecase); err != nil {
		slog.Error("dig provide me usecase", "err", err)
		os.Exit(1)
	}
	if err := c.Provide(usecase.NewItemUsecase); err != nil {
		slog.Error("dig provide item usecase", "err", err)
		os.Exit(1)
	}
	if err := c.Provide(usecase.NewWeaponUsecase); err != nil {
		slog.Error("dig provide weapon usecase", "err", err)
		os.Exit(1)
	}
	if err := c.Provide(usecase.NewMonsterUsecase); err != nil {
		slog.Error("dig provide monster usecase", "err", err)
		os.Exit(1)
	}
	if err := c.Provide(usecase.NewBattleUsecase); err != nil {
		slog.Error("dig provide battle usecase", "err", err)
		os.Exit(1)
	}
	if err := c.Provide(usecase.NewBossBattleUsecase); err != nil {
		slog.Error("dig provide boss battle usecase", "err", err)
		os.Exit(1)
	}
	if err := c.Provide(usecase.NewGameUsecase); err != nil {
		slog.Error("dig provide game usecase", "err", err)
		os.Exit(1)
	}

	// Handler Server
	if err := c.Provide(handler.NewServer); err != nil {
		slog.Error("dig provide server", "err", err)
		os.Exit(1)
	}

	// ── Echo Server ─────────────────────────────────────────────────────────
	if err := c.Invoke(func(srv *handler.Server) {
		e := echo.New()
		e.HideBanner = true

		// Global middleware
		e.Use(echomw.Logger())
		e.Use(echomw.CORS())
		e.Use(mw.InjectEchoContext())

		// Auth middleware with skipper for public routes.
		// The Auth middleware does not natively support a Skipper, so we wrap it.
		authSkipper := func(c echo.Context) bool {
			path := c.Path()
			return path == "/health" || path == "/auth/user-login" || path == "/auth/admin-login" || path == "/admin/images"
		}
		e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
			authMW := mw.Auth(cfg.JWTSecret)
			return func(c echo.Context) error {
				if authSkipper(c) {
					return next(c)
				}
				return authMW(next)(c)
			}
		})

		// Register routes via oapi-codegen generated helpers
		genapi.RegisterHandlers(e, genapi.NewStrictHandler(srv, []genapi.StrictMiddlewareFunc{}))

		// Start server in a goroutine so the main goroutine can wait for signals
		go func() {
			addr := fmt.Sprintf(":%d", cfg.Port)
			slog.Info("starting server", "addr", addr)
			if err := e.Start(addr); err != nil && !errors.Is(err, http.ErrServerClosed) {
				slog.Error("server error", "err", err)
				os.Exit(1)
			}
		}()

		// Graceful shutdown on SIGINT / SIGTERM
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
		<-quit
		slog.Info("shutting down server")

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := e.Shutdown(ctx); err != nil {
			slog.Error("shutdown error", "err", err)
		}
		slog.Info("server stopped")
	}); err != nil {
		slog.Error("dig invoke failed", "err", err)
		os.Exit(1)
	}
}
