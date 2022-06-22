package main

import (
	_ "github.com/jackc/pgx/stdlib"
	"github.com/labstack/echo/v4"
	"github.com/mailcourses/technopark-dbms-forum/api/init/db"
	"github.com/mailcourses/technopark-dbms-forum/api/init/logger"
	"github.com/mailcourses/technopark-dbms-forum/api/init/system"
	"github.com/mailcourses/technopark-dbms-forum/api/internal"
	"os"
)

const (
	portKey   = "port"
	dbTypeKey = "dbType"
	dsnKey    = "dsn"
)

func main() {
	port := os.Getenv(portKey)
	dbType := os.Getenv(dbTypeKey)
	e := echo.New()

	logs, err := logger.InitLogrus(port, dbType)
	if err != nil {
		e.Logger.Fatalf("error to init logrus:", err)
	}

	//e.Use(logs.ColoredLogMiddleware)
	//e.Use(logs.JsonLogMiddleware)
	//e.Logger.SetOutput(logs.Logrus.Writer())

	pools, err := InitDb.InitPostgres(dsnKey)
	if err != nil {
		logs.Logrus.Fatalln("error to init db, dsn:", os.Getenv(dsnKey), "err:", err)
	}

	pgxpools := internal.PgxPoolContainer{
		ForumPool:   pools,
		UserPool:    pools,
		ThreadPool:  pools,
		PostPool:    pools,
		ServicePool: pools,
	}

	if err := system.InitApi(e, pgxpools); err != nil {
		logs.Logrus.Fatalln("error to set routes, err:", err)
	}

	logs.Logrus.Warn("start listening on port: ", port)

	if err := e.Start("0.0.0.0:" + port); err != nil {
		logs.Logrus.Fatal("server error:", err)
	}

	logs.Logrus.Warn("shutdown")

}
