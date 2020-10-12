package modulorgo

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	app "modulorgo/app"
	helper "modulorgo/helpers"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func RegisterTestBlogsAPI(t *testing.T) *app.BlogsAPI {
	router := mux.NewRouter()
	db := helper.InitDB()
	api := &app.BlogsAPI{Db: db, Router: router, TableName: "blogs", Endpoint: "blog", EndpointPlural: "blogs"}
	api.Register()
	return api
}

func RegisterTestAuthAPI(t *testing.T) *app.AuthAPI {
	router := mux.NewRouter()
	db := helper.InitDB()
	api := &app.AuthAPI{Db: db, Router: router, EndpointPlural: "users", Endpoint: "user"}
	api.Register()
	return api
}

var token string

//access letter without token
func TestRequestNotesWithoutToken(t *testing.T) {
	api := RegisterTestBlogsAPI(t)
	ts := httptest.NewServer(api.Router)
	client := &http.Client{}
	defer ts.Close()
	req, err := http.NewRequest("GET", ts.URL+"/api/"+api.Endpoint+"/1", nil)
	if err != nil {
		t.Fatal(err)
	}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode, "Should be error")
	bodyBytes, _ := ioutil.ReadAll(resp.Body)
	bodyString := string(bodyBytes)
	var obj map[string]interface{}
	if err := json.Unmarshal([]byte(bodyString), &obj); err != nil {
		assert.Equal(t, true, false, "Should not be error")
	}
	assert.Equal(t, "You're not allowed access this endpoint", obj["error_message"], "Should be equals.")
}

//Register new Account (arianda@gmail.com)

//login from existing account
func TestLoginNewUser(t *testing.T) {
	api := RegisterTestAuthAPI(t)
	ts := httptest.NewServer(api.Router)
	client := &http.Client{}
	jpayload := []byte(`{"email":"auliafajrul7@gmail.com", "password":"123456"}`)
	req, err := http.NewRequest("POST", ts.URL+"/api/login", bytes.NewBuffer(jpayload))
	if err != nil {
		assert.Equal(t, true, false, "Should not be error")
	}
	resp, err := client.Do(req)
	if err != nil {
		assert.Equal(t, true, false, "Should not be error")
	}
	assert.Equal(t, http.StatusOK, resp.StatusCode, "Should not be error")
	token = resp.Header.Get("Authorization")
}

func TestRequestNotesWithToken(t *testing.T) {
	api := RegisterTestBlogsAPI(t)
	ts := httptest.NewServer(api.Router)
	client := &http.Client{}
	defer ts.Close()
	req, err := http.NewRequest("GET", ts.URL+"/api/"+api.Endpoint+"/1", nil)
	req.Header.Set("Authorization", token)
	if err != nil {
		t.Fatal(err)
	}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, http.StatusOK, resp.StatusCode, "Should not be error")
	bodyBytes, _ := ioutil.ReadAll(resp.Body)
	bodyString := string(bodyBytes)
	var obj map[string]interface{}
	if err := json.Unmarshal([]byte(bodyString), &obj); err != nil {
		assert.Equal(t, true, false, "Should not be error")
	}
	assert.Equal(t, "Fajrul Aulia", obj["publisher"], "Should be equals.")
}

func TestRegisterNewAccount(t *testing.T) {
	api := RegisterTestAuthAPI(t)
	ts := httptest.NewServer(api.Router)
	client := &http.Client{}

	defer ts.Close()
	jpayload := []byte(`
		{
			"name":"Arianda",
			"username":"Arianda",
			"phonenumber":"082166773",
			"email":"arianda@gmail.com",
			"password":"$2y$12$tNKsfjSrVU/SHfdc/yNBP.YDogpFbCye9qpbESUPtYlQV1LgVCk9S"
	}`)
	req, err := http.NewRequest("POST", ts.URL+"/api/register", bytes.NewBuffer(jpayload))
	if err != nil {
		t.Fatal(err)
	}
	resp, err := client.Do(req)
	if err != nil {
		assert.Equal(t, true, false, "Should not be error")
	}
	assert.Equal(t, http.StatusOK, resp.StatusCode, "Should not be error")
	bodyIO, _ := ioutil.ReadAll(resp.Body)
	var bodyItem map[string]interface{}
	if err := json.Unmarshal([]byte(string(bodyIO)), &bodyItem); err != nil {
		assert.Equal(t, true, false, "Should not be error")
	}
	email := bodyItem["email"]
	assert.Equal(t, "arianda@gmail.com", email, "Should be equals.")
	token = resp.Header.Get("Authorization")

}

func TestRequestNotesAriandaCantAcccessIt(t *testing.T) {
	api := RegisterTestBlogsAPI(t)
	ts := httptest.NewServer(api.Router)
	client := &http.Client{}
	defer ts.Close()
	req, err := http.NewRequest("GET", ts.URL+"/api/"+api.Endpoint+"/1", nil)
	req.Header.Set("Authorization", token)
	if err != nil {
		t.Fatal(err)
	}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, 404, resp.StatusCode, "Should not be error")
	bodyBytes, _ := ioutil.ReadAll(resp.Body)
	bodyString := string(bodyBytes)
	var obj map[string]interface{}
	if err := json.Unmarshal([]byte(bodyString), &obj); err != nil {
		assert.Equal(t, true, false, "Should not be error")
	}
	assert.Equal(t, "You request not found", obj["error_message"], "Should be error.")
}
