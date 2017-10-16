package db

import (
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	"github.com/vlad-s/gophircbot/apiconfig"
)

// AdminUser is the Gorm database model for an admin user
type AdminUser struct {
	gorm.Model
	Nick string
}

// IgnoredUser is the Gorm database model for an ignored user
type IgnoredUser struct {
	gorm.Model
	Nick string
}

// AutoReply is the Gorm database model for an automated reply
type AutoReply struct {
	gorm.Model
	Trigger string
	Reply   string
}

var database *gorm.DB

// Connect returns a pointer to a gorm.DB instance using the apiconfig.DBConfig parameters,
// or an error in case it fails to connect
func Connect(config apiconfig.DBConfig) (*gorm.DB, error) {
	mysqlConfig := mysql.Config{
		User:   config.Username,
		Passwd: config.Password,
		DBName: config.Database,
		Net:    "tcp",
		Addr:   config.Host,
		Params: map[string]string{
			"charset": "utf8",
		},
		ParseTime: true,
		Loc:       time.Local,
	}

	var err error
	database, err = gorm.Open("mysql", mysqlConfig.FormatDSN())
	if err != nil {
		return nil, errors.Wrap(err, "Error opening the database")
	}

	return database, nil
}

// Get returns the open database connection
func Get() *gorm.DB {
	return database
}

// AutoMigrate migrates the models specified
func AutoMigrate(db *gorm.DB) {
	models := []interface{}{
		&AdminUser{},
		&IgnoredUser{},
		&AutoReply{},
	}
	db.AutoMigrate(models...)
}
