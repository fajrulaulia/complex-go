package modulorgo

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

// ResponseOnError should be exported"
func ResponseOnError(headerStatus int, err error, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(headerStatus)
	log.Println("[Error Message] : ", err)
	log.Println("[Error Time]    : ", time.Now().String())
	msg := "Error Unknwon, Please contact adminstrator"
	switch headerStatus {
	case http.StatusInternalServerError:
		msg = "Internal server error"
	case http.StatusNotFound:
		msg = "You request not found"
	case http.StatusUnauthorized:
		msg = "You're not allowed access this endpoint"
	}
	w.Write([]byte(fmt.Sprintf("{\"error_message\":\"%v\",\"error_datetime\":\"%v\"}", msg, time.Now().String())))
	return
}

// ResponseWithlastInsertID should be exported"
func ResponseWithlastInsertID(lastInsertID int64, w http.ResponseWriter) {
	w.WriteHeader(200)
	w.Write([]byte(fmt.Sprintf("{\"last_insert_id\":\"%v\",\"status_code\":\"%v\"}", lastInsertID, "success")))
	return
}

// ResponseAuth should be exported"
func ResponseAuth(email string, token string, w http.ResponseWriter) {
	w.Header().Set("Authorization", "Bearer "+token)
	w.WriteHeader(200)
	w.Write([]byte(fmt.Sprintf("{\"email\":\"%v\",\"status_code\":\"%v\"}", email, 200)))
	return
}
