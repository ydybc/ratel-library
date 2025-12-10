package datasource

import (
	"github.com/jmoiron/sqlx"
	//"fmt"

	_ "github.com/go-sql-driver/mysql"
)

var DB *sqlx.DB

func init() {
	var err error
	//	DBW = DbWorker{
	//		Dsn: "root:@tcp(localhost:3306)/jieqi?charset=utf8",
	//	}
	Dsn := "root:@tcp(localhost:3306)/jieqi?charset=utf8"
	DB, err = sqlx.Open("mysql", Dsn)
	if err != nil {
		panic(err)
		return
	}
	//defer DBW.Db.Close()
}
