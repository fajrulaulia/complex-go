package modulorgo

import (
	"database/sql"
	"log"
	"os"

	//allow to blank
	_ "github.com/go-sql-driver/mysql"
)

// InitDB should be exported"
func InitDB() *sql.DB {
	if os.Getenv("DRIVER_TYPE") == "" || os.Getenv("DRIVER_TYPE") != "mysql" {
		log.Print("Envar DRIVER_TYPE not found !")
	}
	if os.Getenv("MYSQL_HOST") == "" {
		log.Print("Envar MYSQL_HOST not found !")
	}
	if os.Getenv("MYSQL_USER") == "" {
		log.Print("Envar MYSQL_USER not found !")
	}
	if os.Getenv("MYSQL_PORT") == "" {
		log.Print("Envar MYSQL_PORT not found !")
	}
	if os.Getenv("MYSQL_PWD") == "" {
		log.Print("Envar MYSQL_PWD not found !")
	}
	if os.Getenv("MYSQL_DB") == "" {
		log.Print("Envar MYSQL_DB not found !")
	}
	db, err := sql.Open(os.Getenv("DRIVER_TYPE"), os.Getenv("MYSQL_USER")+":"+os.Getenv("MYSQL_PWD")+"@tcp("+os.Getenv("MYSQL_HOST")+":"+os.Getenv("MYSQL_PORT")+")/"+os.Getenv("MYSQL_DB"))
	if err != nil {
		log.Println(err)
		panic(err.Error())
	}
	db.SetMaxIdleConns(0)
	return db
}
