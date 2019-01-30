package driver

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

//DbConnectorConfig sql credentials
type DbConnectorConfig struct {
	User string `json:"user"`
	Pass string `json:"pass"`
	Host string `json:"host"`
	Port string `json:"port"`
	Name string `json:"name"`
}

//NewDbConnector get 1 new db handler
func NewDbConnector(cfg *DbConnectorConfig) (*sql.DB, error) {
	//get handle
	var err error
	var dbh *sql.DB
	var connstr = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s",
		cfg.User,
		cfg.Pass,
		cfg.Host,
		cfg.Port,
		cfg.Name,
	)
	//try to limit the try
	for i := 0; i < 5; i++ {
		dbh, err = sql.Open("mysql", connstr)
		if err != nil {
			log.Println("SQL:", err)
			time.Sleep(1000 * time.Millisecond)
		} else {
			log.Println("Db driver connected.")
			break
		}
	}
	//important tweak is here :-)
	if dbh != nil {
		dbh.SetConnMaxLifetime(time.Minute * 5)
		dbh.SetMaxIdleConns(0)
		dbh.SetMaxOpenConns(5)
	}

	if err != nil {
		return nil, err
	}
	return dbh, nil

}
