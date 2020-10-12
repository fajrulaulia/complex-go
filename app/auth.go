package modulorgo

import (
	"database/sql"
	"encoding/json"
	"errors"
	helper "modulorgo/helpers"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// AuthAPI should be exported
type AuthAPI struct {
	Db             *sql.DB
	Router         *mux.Router
	TableName      string
	Endpoint       string
	EndpointPlural string
}

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type registerRequest struct {
	Name        string `json:"name"`
	Username    string `json:"username"`
	Phonenumber string `json:"phonenumber"`
	Email       string `json:"email"`
	Password    string `json:"password"`
}

type userResponse struct {
	Name        string `json:"name"`
	Username    string `json:"username"`
	Phonenumber string `json:"phonenumber"`
	Email       string `json:"email"`
	CreatedAt   string `json:"created_at"`
}

// Register should be exported
func (api *AuthAPI) Register() {
	api.Router.Handle("/api/login", http.HandlerFunc(api.login)).Methods("POST")
	api.Router.Handle("/api/register", http.HandlerFunc(api.register)).Methods("POST")
	api.Router.Handle("/api/me", helper.VerifyToken(http.HandlerFunc(api.userCurrent))).Methods("GET")
}

func (api *AuthAPI) login(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	req := loginRequest{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&req)
	if err != nil {
		helper.ResponseOnError(500, err, w)
		return
	}

	query := "select id, username, password  from users u where u.email=?"
	rows, err := api.Db.Query(query, req.Email)
	if err != nil {
		helper.ResponseOnError(500, err, w)
		return
	}
	count := 0
	var (
		id       string
		username string
		password string
	)
	for rows.Next() {
		count++
		rows.Scan(&id, &username, &password)
	}
	if count == 1 && helper.ComparedPassword(password, req.Password) {
		token, err := helper.GenerateToken(map[string]string{"id": id, "username": username, "email": req.Email})
		if err != nil {
			helper.ResponseOnError(500, err, w)
			return
		}
		helper.ResponseAuth(req.Email, token, w)
		return
	}
	helper.ResponseOnError(401, errors.New("Account not found in our database"), w)
}

func (api *AuthAPI) register(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	req := registerRequest{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&req)
	if err != nil {
		helper.ResponseOnError(500, err, w)
		return
	}
	query := "INSERT INTO users (name, username, phonenumber, email, password, status, created_at, updated_at) VALUES(?, ?, ?, ?, ?, 'DEACTIVE', now(), now())"
	result, err := api.Db.Exec(query, req.Name, req.Username, req.Phonenumber, req.Email, req.Password)
	if err != nil {
		helper.ResponseOnError(500, err, w)
		return
	}
	LastInsertId, err := result.LastInsertId()
	if err != nil {
		helper.ResponseOnError(500, err, w)
		return
	}
	token, err := helper.GenerateToken(map[string]string{"id": strconv.Itoa(int(LastInsertId)), "username": req.Username, "email": req.Email})
	if err != nil {
		helper.ResponseOnError(500, err, w)
		return
	}
	helper.ResponseAuth(req.Email, token, w)
}

func (api *AuthAPI) userCurrent(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	res := userResponse{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&res)
	currentuser := helper.IndetifyToken(r, w)["public"].(map[string]interface{})
	query := "SELECT name,username ,phonenumber ,email ,created_at from users where id=?;"
	err = api.Db.QueryRow(query, currentuser["id"]).Scan(&res.Name, &res.Username, &res.Email, &res.Phonenumber, &res.CreatedAt)
	if err != nil {
		helper.ResponseOnError(500, err, w)
		return
	}
	output, err := json.Marshal(res)
	if err != nil {
		helper.ResponseOnError(500, err, w)
		return
	}
	w.Write(output)
}
