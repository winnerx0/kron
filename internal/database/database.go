package database

import (
	"fmt"
	"time"



	
	"github.com/winnerx0/kron/internal/execution"
	"github.com/winnerx0/kron/internal/job"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Database struct {
	dbHost     string
	dbUser     string
	dbPassword string
	dbPort     string
	dbName     string
}

func NewDatabase(host, user, password, port, name string) *Database {
	return &Database{
		dbHost:     host,
		dbUser:     user,
		dbPassword: password,
		dbPort:     port,
		dbName:     name,
	}
}

func (d *Database) Start() *gorm.DB {
	dbUrl := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Africa/Lagos",
		d.dbHost, d.dbUser, d.dbPassword, d.dbName, d.dbPort)

	db, err := gorm.Open(postgres.Open(dbUrl), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	
	sqlDB, err := db.DB()
	
	if err != nil {
		panic(err)
	}
	
	sqlDB.SetMaxOpenConns(10)
	
	sqlDB.SetMaxIdleConns(5)
	
	sqlDB.SetConnMaxIdleTime(time.Minute * 5)
	
	sqlDB.SetConnMaxLifetime(time.Minute * 30)
	
	err = sqlDB.Ping()
	
	if err != nil {
		panic(err)
	}
	
	err = db.AutoMigrate(&job.Job{}, &execution.Execution{})
	
	if err != nil {
		panic(err)
	}
	
	return db
}
