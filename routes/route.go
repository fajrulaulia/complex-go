package modulorgo

import (
	"database/sql"
	app "modulorgo/app"

	"github.com/gorilla/mux"
)

//Register should be exported
func Register(db *sql.DB, router *mux.Router) {
	BlogsAPI := app.BlogsAPI{Db: db, Router: router, TableName: "blogs", Endpoint: "blog", EndpointPlural: "blogs"}
	BlogsAPI.Register()

	AuthAPI := app.AuthAPI{Db: db, Router: router, TableName: "users", Endpoint: "user", EndpointPlural: "users"}
	AuthAPI.Register()
}
