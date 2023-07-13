package main

import (
	"context"
	"fmt"
	"github.com/POMBNK/restAPI/internal/user"
	"github.com/POMBNK/restAPI/internal/user/db"
	"github.com/POMBNK/restAPI/pkg/client/mongodb"
	"github.com/POMBNK/restAPI/pkg/client/postgresql"
	"github.com/POMBNK/restAPI/pkg/config"
	"github.com/POMBNK/restAPI/pkg/logger"
	"github.com/julienschmidt/httprouter"
	"net"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"
)

const (
	port   = "port"
	tcp    = "tcp"
	socket = "app.sock"
)

func main() {
	logs := logger.GetLogger()
	logs.Println("Logger initialized.")

	logs.Println("Config initialization...")
	cfg := config.GetCfg()
	logs.Println("Config initialized.")

	logs.Println("Router initialization...")
	router := httprouter.New()
	logs.Println("Router initialized.")
	storage := getStorage(cfg, logs)
	service := user.NewService(storage, logs)
	handler := user.NewHandler(logs, service)
	handler.Register(router)

	logs.Infof("Starting app...")
	start(logs, router, cfg)

}

func start(logs *logger.Logger, router *httprouter.Router, cfg *config.Config) {
	var listener net.Listener
	var listenErr error

	if cfg.Listen.Type == port {
		logs.Infof("Listen tcp")
		listener, listenErr = net.Listen(tcp, fmt.Sprintf("%s:%s", cfg.Listen.BindIP, cfg.Listen.Port))
		if listenErr != nil {
			logs.Fatal(listenErr)
		}
	} else {
		dirPath, err := filepath.Abs(filepath.Dir(os.Args[0]))
		if err != nil {
			logs.Fatal(err)
		}
		socketPath := filepath.Join(dirPath, socket)
		logs.Debugf("Create a socket at path: %s", socketPath)
		logs.Info("Listen socket")
		listener, listenErr = net.Listen("unix", socketPath)
		if listenErr != nil {
			logs.Fatal(listenErr)
		}

		// Remove socket after shutdown
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)
		go func() {
			<-c
			os.Remove(socketPath)
			os.Exit(1)
		}()

	}

	server := http.Server{
		Handler:      router,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	if err := server.Serve(listener); err != nil {
		logs.Fatalf("Server error:%s", err)
	}

}

func getStorage(cfg *config.Config, logs *logger.Logger) user.Storage {
	storage := cfg.Storage
	switch storage.Type {
	case "mongo":
		mongoDatabase, err := mongodb.NewClient(context.Background(), storage.MongoDB.Host, storage.MongoDB.Port, storage.MongoDB.User,
			storage.MongoDB.Password, storage.MongoDB.Database, storage.MongoDB.AuthDB)
		if err != nil {
			logs.Fatalln(err)
		}
		mongoStorage := db.NewMongoDB(mongoDatabase, cfg.Storage.MongoDB.Collection, logs)
		return mongoStorage
	case "sql":
		client, err := postgresql.NewClient(context.Background(), cfg)
		if err != nil {
			logs.Fatalln(err)
		}
		sqlStorage := db.NewPostgresDB(client, logs)
		return sqlStorage
	default:
		logs.Fatalln("incorrect database type")
		return nil
	}
}
