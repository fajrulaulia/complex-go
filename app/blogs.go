package modulorgo

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	helper "modulorgo/helpers"
)

// BlogsAPI should be exported
type BlogsAPI struct {
	Db             *sql.DB
	Router         *mux.Router
	TableName      string
	Endpoint       string
	EndpointPlural string
}

type blogResponse struct {
	ID        int    `json:"id"`
	Title     string `json:"title"`
	Body      string `json:"body"`
	Publisher string `json:"publisher"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type blogRequest struct {
	Title string `json:"title"`
	Body  string `json:"body"`
}

type blogsResponse struct {
	Total int            `json:"total"`
	Page  int            `json:"page"`
	Limit int            `json:"limit"`
	Data  []blogResponse `json:"data"`
}

// Register should be exported
func (api *BlogsAPI) Register() {
	api.Router.Handle("/api/"+api.EndpointPlural, helper.VerifyToken(http.HandlerFunc(api.createRequest))).Methods("POST")
	api.Router.Handle("/api/"+api.Endpoint+"/{id}", helper.VerifyToken(http.HandlerFunc(api.updateRequest))).Methods("PUT")
	api.Router.Handle("/api/"+api.EndpointPlural, helper.VerifyToken(http.HandlerFunc(api.listRequest))).Methods("GET")
	api.Router.Handle("/api/"+api.Endpoint+"/{id}", helper.VerifyToken(http.HandlerFunc(api.readRequest))).Methods("GET")
	api.Router.Handle("/api/"+api.Endpoint+"/{id}", helper.VerifyToken(http.HandlerFunc(api.deleteRequest))).Methods("DELETE")

}

func (api *BlogsAPI) listRequest(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var (
		blog  blogResponse
		blogs []blogResponse
		resp  blogsResponse
	)
	currentuser := helper.IndetifyToken(r, w)["public"].(map[string]interface{})
	query := "SELECT b.id, b.title, b.body, u.name , b.updated_at, b.created_at FROM blogs b "
	query += "INNER JOIN users u on b.id_user = u.id where u.id ='" + currentuser["id"].(string) + "'"
	rows, err := api.Db.Query(query)
	if err != nil {
		helper.ResponseOnError(500, err, w)
		return
	}
	count := 0
	for rows.Next() {
		count++
	}
	searchQuery := helper.GenerateSearch(r, []string{"b.title", "b.body"})
	if searchQuery != "" {
		query += " and " + searchQuery
	}
	query += helper.GenerateSortedNew("b.updated_at")
	var limitQuery string
	limitQuery, resp.Page, resp.Limit = helper.GeneratePagination(r)
	if limitQuery == "" {
		resp.Page = 1
		resp.Limit = count
	}
	query += limitQuery
	rows, err = api.Db.Query(query)
	if err != nil {
		helper.ResponseOnError(500, err, w)
		return
	}
	for rows.Next() {
		if err := rows.Scan(&blog.ID, &blog.Title, &blog.Body, &blog.Publisher, &blog.UpdatedAt, &blog.CreatedAt); err != nil {
			helper.ResponseOnError(500, err, w)
			return
		}
		blogs = append(blogs, blog)
	}
	resp.Total = count
	resp.Data = blogs
	output, err := json.Marshal(resp)
	if err != nil {
		helper.ResponseOnError(500, err, w)
		return
	}
	w.Write(output)
}

func (api *BlogsAPI) readRequest(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var blog blogResponse

	currentuser := helper.IndetifyToken(r, w)["public"].(map[string]interface{})

	query := "SELECT b.id, b.title, b.body, u.name , b.updated_at, b.created_at FROM blogs b "
	query += "INNER JOIN users u on b.id_user = u.id where  "
	query += "u.id = '" + currentuser["id"].(string) + "' "
	query += "AND b.id = '" + mux.Vars(r)["id"] + "' "
	rows := api.Db.QueryRow(query)
	if err := rows.Scan(&blog.ID, &blog.Title, &blog.Body, &blog.Publisher, &blog.UpdatedAt, &blog.CreatedAt); err != nil {
		helper.ResponseOnError(404, err, w)
		return
	}
	output, err := json.Marshal(blog)
	if err != nil {
		helper.ResponseOnError(500, err, w)
		return
	}
	w.Write(output)
}

func (api *BlogsAPI) createRequest(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var req blogRequest
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&req)
	if err != nil {
		helper.ResponseOnError(500, err, w)
		return
	}
	currentuser := helper.IndetifyToken(r, w)["public"].(map[string]interface{})
	query := "INSERT INTO blogs (id_user, title, body, created_at, updated_at) VALUES(?, ?,?, now(), now());"
	result, err := api.Db.Exec(query, currentuser["id"], req.Title, req.Body)
	if err != nil {
		helper.ResponseOnError(500, err, w)
		return
	}
	LastInsertId, err := result.LastInsertId()
	if err != nil {
		helper.ResponseOnError(500, err, w)
		return
	}

	helper.ResponseWithlastInsertID(LastInsertId, w)
}

func (api *BlogsAPI) updateRequest(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if !api.isAllowedRequest(w, r) {
		helper.ResponseOnError(404, errors.New("Data not found"), w)
		return
	}
	var req blogRequest
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&req)
	if err != nil {
		helper.ResponseOnError(500, err, w)
		return
	}
	query := "UPDATE blogs set title=?, body=?, updated_at= now() where id=?"
	_, err = api.Db.Exec(query, req.Title, req.Body, mux.Vars(r)["id"])
	if err != nil {
		helper.ResponseOnError(500, err, w)
		return
	}
	n, err := strconv.ParseInt(mux.Vars(r)["id"], 10, 64)
	if err != nil {
		helper.ResponseOnError(500, err, w)
		return
	}
	helper.ResponseWithlastInsertID(n, w)
}

func (api *BlogsAPI) deleteRequest(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if !api.isAllowedRequest(w, r) {
		helper.ResponseOnError(404, errors.New("Data not found"), w)
		return
	}

	query := "delete from blogs where id=?"
	_, err := api.Db.Exec(query, mux.Vars(r)["id"])
	if err != nil {
		helper.ResponseOnError(500, err, w)
		return
	}
	n, err := strconv.ParseInt(mux.Vars(r)["id"], 10, 64)
	if err != nil {
		helper.ResponseOnError(500, err, w)
		return
	}
	helper.ResponseWithlastInsertID(n, w)
}

func (api *BlogsAPI) isAllowedRequest(w http.ResponseWriter, r *http.Request) bool {
	var blog blogResponse
	currentuser := helper.IndetifyToken(r, w)["public"].(map[string]interface{})
	query := "SELECT b.id, b.title, b.body, u.name , b.updated_at, b.created_at FROM blogs b "
	query += "INNER JOIN users u on b.id_user = u.id where  "
	query += "u.id = '" + currentuser["id"].(string) + "' "
	query += "AND b.id = '" + mux.Vars(r)["id"] + "' "
	rows := api.Db.QueryRow(query)
	if err := rows.Scan(&blog.ID, &blog.Title, &blog.Body, &blog.Publisher, &blog.UpdatedAt, &blog.CreatedAt); err != nil {
		return false
	}
	return true
}
