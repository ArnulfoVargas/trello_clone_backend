package main

import (
	"os"

	"github.com/ArnulfoVargas/trello_clone_backend.git/cmd/database"
	"github.com/ArnulfoVargas/trello_clone_backend.git/cmd/routes"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()

	if err != nil {
		panic(err)
	}

	var conn = &database.Database{}
	conn.ConnectDB(os.Getenv("DB_STRING"))

	app := routes.NewRouter()
	app.SetDb(conn)

	app.BindAuthRoutes()
	app.BindBoardsRoutes()
	app.BindUserRoutes()

	app.ServeHttp()
}
