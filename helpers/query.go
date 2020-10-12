package modulorgo

import (
	"net/http"
	"strconv"
)

//==========================================
//this function generate query to operation
//==========================================

// GenerateSortedNew should be exported"
func GenerateSortedNew(column string) string {
	return " order by " + column + " ASC "
}

// GeneratePagination should be exported"
func GeneratePagination(r *http.Request) (string, int, int) {
	if r.URL.Query().Get("limit") != "" && r.URL.Query().Get("page") != "" {
		page, err := strconv.ParseInt(r.URL.Query().Get("page"), 10, 64)
		if err != nil {
			return "", 0, 0
		}
		limit, err := strconv.ParseInt(r.URL.Query().Get("limit"), 10, 64)
		if err != nil {
			return "", 0, 0
		}
		start := 0
		if page == 0 || page == 1 {
			start = 0
		} else {
			start = int(page) - 1
		}
		start = start * int(limit)
		return " LIMIT " + strconv.Itoa(int(start)) + ", " + strconv.Itoa(int(limit)) + " ", int(page), int(limit)
	}
	return "", 0, 0
}

// GenerateSearch should be exported"
func GenerateSearch(r *http.Request, column []string) string {
	if r.URL.Query().Get("search") != "" {
		query := "("
		for i := 0; i < len(column); i++ {
			query += " " + column[i] + " like '%" + r.URL.Query().Get("search") + "%' "
			if i < len(column)-1 {
				query += " or "
			}
		}
		return query + ")"
	}
	return ""
}
