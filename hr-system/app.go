package main

import (
	"context"
	"fmt"
	"log"

	"hr-system/backend/bootstrap"

	"github.com/jmoiron/sqlx"
)

// App struct
type App struct {
	ctx        context.Context
	db         *sqlx.DB
	auth       *bootstrap.AuthFacade
	employees  *bootstrap.EmployeesFacade
	leave      *bootstrap.LeaveFacade
	startupErr error
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx

	runtime, err := bootstrap.Initialize(ctx)
	if err != nil {
		a.startupErr = err
		log.Printf("startup warning: %v", a.startupErr)
		return
	}
	a.db = runtime.DB
	a.auth = runtime.Auth
	a.employees = runtime.Employees
	a.leave = runtime.Leave

	log.Printf("backend foundation initialized")
}

func (a *App) shutdown(ctx context.Context) {
	_ = ctx
	if a.db != nil {
		if err := a.db.Close(); err != nil {
			log.Printf("shutdown warning: close database: %v", err)
		}
	}
}

// Greet returns a greeting for the given name
func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, It's show time!", name)
}
