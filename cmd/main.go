package main

import (
	"chatappv2/config"
	"chatappv2/internal"
	"database/sql"
	"github.com/gin-gonic/gin"
	"log"
)

type App struct {
	Config *config.Config
	DB     *sql.DB
	Router *gin.Engine
}

func main() {
	app := InitApp()
	app.Start()
}

func InitApp() *App {
	cfg, err := config.LoadApplicationConfig()
	if err != nil {
		log.Fatal("cannot load config:", err)
	}

	db, err := config.InitPostgresConnection()
	if err != nil {
		log.Fatal("cannot connect to database:", err)
	}
	r := internal.InitRouter(db)
	return &App{
		Config: cfg,
		DB:     db,
		Router: r,
	}
}

func (app *App) Start() {
	defer app.DB.Close()
	if err := app.Router.Run(app.Config.Server.Port); err != nil {
		log.Fatal("server start error:", err)
	}
}
