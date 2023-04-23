package database

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/spf13/viper"
)

func NewMySQLConn() (*sql.DB, error) {
	connString := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?timeout=%ds",
		viper.GetString("mysql.username"),
		os.Getenv("MYSQL_PASSWORD"),
		viper.GetString("mysql.host"),
		viper.GetString("mysql.port"),
		viper.GetString("mysql.database"),
		viper.GetInt("mysql.timeout"),
	)

	// fmt.Println(connString)
	db, err := sql.Open("mysql", connString)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
