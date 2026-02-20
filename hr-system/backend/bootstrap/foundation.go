package bootstrap

import (
	"context"
	"fmt"

	"hr-system/backend/internal/config"
	"hr-system/backend/internal/db"

	"github.com/jmoiron/sqlx"
)

type Runtime struct {
	DB        *sqlx.DB
	Auth      *AuthFacade
	Employees *EmployeesFacade
	Leave     *LeaveFacade
	Payroll   *PayrollFacade
	Users     *UsersFacade
}

func Initialize(ctx context.Context) (*Runtime, error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, fmt.Errorf("load config: %w", err)
	}

	if cfg.RunMigrations {
		if err := db.RunMigrations(cfg.DatabaseURL, cfg.MigrationsPath); err != nil {
			return nil, fmt.Errorf("run migrations: %w", err)
		}
	}

	conn, err := db.Connect(ctx, cfg)
	if err != nil {
		return nil, fmt.Errorf("connect database: %w", err)
	}

	authFacade, err := NewAuthFacade(
		conn,
		cfg.JWTSecret,
		int64(cfg.AccessTokenTTL.Seconds()),
		int64(cfg.RefreshTokenTTL.Seconds()),
	)
	if err != nil {
		_ = conn.Close()
		return nil, fmt.Errorf("initialize auth: %w", err)
	}

	employeesFacade, err := NewEmployeesFacade(conn)
	if err != nil {
		_ = conn.Close()
		return nil, fmt.Errorf("initialize employees: %w", err)
	}

	leaveFacade, err := NewLeaveFacade(conn)
	if err != nil {
		_ = conn.Close()
		return nil, fmt.Errorf("initialize leave: %w", err)
	}

	payrollFacade, err := NewPayrollFacade(conn)
	if err != nil {
		_ = conn.Close()
		return nil, fmt.Errorf("initialize payroll: %w", err)
	}

	usersFacade, err := NewUsersFacade(conn)
	if err != nil {
		_ = conn.Close()
		return nil, fmt.Errorf("initialize users: %w", err)
	}
	if err := usersFacade.SeedInitialAdmin(ctx, cfg.InitialAdminUsername, cfg.InitialAdminPassword, cfg.InitialAdminRole); err != nil {
		_ = conn.Close()
		return nil, fmt.Errorf("seed initial admin: %w", err)
	}

	return &Runtime{
		DB:        conn,
		Auth:      authFacade,
		Employees: employeesFacade,
		Leave:     leaveFacade,
		Payroll:   payrollFacade,
		Users:     usersFacade,
	}, nil
}
