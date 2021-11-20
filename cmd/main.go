package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"os/signal"
	"syscall"

	"github.com/mhdiiilham/POS/api"
	"github.com/mhdiiilham/POS/config"
	"github.com/mhdiiilham/POS/database"
	"github.com/mhdiiilham/POS/pkg/hasher"
	"github.com/mhdiiilham/POS/pkg/logger"
	"github.com/mhdiiilham/POS/pkg/server"
	"github.com/mhdiiilham/POS/pkg/token"
	userrepository "github.com/mhdiiilham/POS/repository/user"
	"github.com/mhdiiilham/POS/service"
	"github.com/sirupsen/logrus"
)

func init() {
	logrus.SetFormatter(&logrus.TextFormatter{
		DisableColors: true,
		FullTimestamp: true,
		ForceQuote:    true,
	})
}

func main() {
	const ops = "main"
	ctx, done := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer func() {
		done()
		if r := recover(); r != nil {
			panic(fmt.Sprintf("application panic: %v", r))
		}
	}()
	logger.Info(ctx, ops, "starting api service")

	env := flag.String("env", "local", "To Set Service Environment Mode")
	logger.Info(ctx, ops, "starting service in %s mode", *env)

	dbConn, err := realMain(ctx, *env)
	if err != nil {
		panic(err)
	}

	logger.Info(ctx, ops, "closing connection to db, error: %v", dbConn.Close())
	logger.Info(ctx, ops, "successfully shutdown")
}

func realMain(ctx context.Context, env string) (*sql.DB, error) {
	const ops = "main.realMain"
	cfg, cfgErr := config.ReadConfig(env)
	if cfgErr != nil {
		return nil, cfgErr
	}

	dbDNS := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.DBName,
	)

	logger.Info(ctx, ops, "connecting to postgresql")
	db, dbErr := database.NewPostgreSQLConnection(ctx, dbDNS)
	if dbErr != nil {
		return nil, dbErr
	}

	pwdHasher := hasher.NewHasher()
	tokenService := token.NewJWTService(cfg.JwtSecret, cfg.JwtIssuer)
	userRepository := userrepository.NewRepository(db)
	userService := service.NewAPIService(userRepository, pwdHasher, tokenService)

	restAPI := api.NewServer(userService)
	srv, err := server.New(cfg.Port)
	if err != nil {
		return nil, err
	}

	logger.Info(ctx, ops, "server is listening on port: %s", cfg.Port)
	return db, srv.ServeHTTPHandler(ctx, restAPI.CORS(restAPI.HandlerLogging(restAPI.Routes(ctx))))
}
