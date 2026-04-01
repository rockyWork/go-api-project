package repository

import (
	"go-api-project/config"
	"go-api-project/internal/model"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Repository struct {
	DB   *gorm.DB
	User *UserRepository
}

func New(cfg *config.DatabaseConfig) (*Repository, error) {
	db, err := gorm.Open(mysql.Open(cfg.DSN()), &gorm.Config{
		Logger:                 logger.Default.LogMode(logger.Info),
		PrepareStmt:            true,
		SkipDefaultTransaction: true,
	})
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	sqlDB.SetConnMaxLifetime(cfg.ConnMaxLifetime)

	return &Repository{
		DB:   db,
		User: NewUserRepository(db),
	}, nil
}

func (r *Repository) AutoMigrate() error {
	return r.DB.AutoMigrate(
		&model.User{},
		&model.RefreshToken{},
	)
}

func (r *Repository) Close() error {
	sqlDB, err := r.DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}
