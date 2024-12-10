package database

import (
	"errors"
	"fmt"
	"github.com/glebarez/sqlite"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"my-blog/internal/config"
	"my-blog/internal/database/repo"
	"my-blog/internal/utils/logger"
)

// InitGormDB initializes the database connection and repositories.
func InitGormDB(conf *config.Config) (*gorm.DB, error) {
	db, err := NewGormDB(conf)
	if err != nil {
		logger.Fatalf("Database connection failed: %v", err)
		return nil, err
	}

	repo.InitRepo(db)
	logger.Infof("Database is ready.")

	if conf.Server.DbAutoMigrate {
		if err := TryAutoMigrate(db); err != nil {
			logger.Fatalf("Database migration failed: %v", err)
			return nil, err
		}
		logger.Infof("Database migration completed successfully.")
	}

	return db, nil
}

// NewGormDB creates a new Gorm DB instance based on the configuration.
func NewGormDB(conf *config.Config) (*gorm.DB, error) {
	dbType := conf.DbType()
	dsn := conf.DbDSN()

	var level gormLogger.LogLevel
	switch conf.Server.DbLogMode {
	case "silent":
		level = gormLogger.Silent
	case "info":
		level = gormLogger.Info
	case "warn":
		level = gormLogger.Warn
	case "error":
		fallthrough
	default:
		level = gormLogger.Error
	}

	gormConfig := &gorm.Config{
		Logger:                                   gormLogger.Default.LogMode(level),
		DisableForeignKeyConstraintWhenMigrating: true, // Disable foreign key constraints
		SkipDefaultTransaction:                   true, // Improve performance by skipping default transactions
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true, // Use singular table names
		},
	}

	var db *gorm.DB
	var err error

	switch dbType {
	case "mysql":
		db, err = gorm.Open(mysql.Open(dsn), gormConfig)
	case "sqlite":
		db, err = gorm.Open(sqlite.Open(dsn), gormConfig)
	default:
		return nil, errors.New(fmt.Sprintf("Unsupported database type: %s", dbType))
	}

	if err != nil {
		return nil, errors.New(fmt.Sprintf("Failed to connect to database: %v", err))
	}

	logger.Infof("Database connected successfully: %s", dbType)
	return db, nil
}

// TryAutoMigrate performs database migrations.
func TryAutoMigrate(db *gorm.DB) error {
	// Example: Add your models for migration here
	// if err := db.AutoMigrate(&model.UserAuth{}, &model.UserInfo{}); err != nil {
	//     return err
	// }
	return nil
}
