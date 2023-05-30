package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"log"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"
)

var filename = "dataset.xml"

var serverAccessToken = "7777"

type UserData struct {
	ID        int    `xml:"id"`
	FirstName string `xml:"first_name" json:"-"`
	LastName  string `xml:"last_name" json:"-"`
	Name      string
	Age       int    `xml:"age"`
	About     string `xml:"about"`
	Gender    string `xml:"gender"`
}

type Users struct {
	Users []UserData `xml:"row"`
}

var usersData Users

func ServerSearch(w http.ResponseWriter, r *http.Request) {
	jsonError := validateAccessToken(r)
	if jsonError != "" {
		http.Error(w, jsonError, http.StatusUnauthorized)
		return
	}

	if len(usersData.Users) == 0 {
		http.Error(w, `{"error": "Данные о пользователях не загружены"`, http.StatusInternalServerError)
		return
	}

	searchReq := &SearchRequest{}
	jsonError = validateParams(r, searchReq)
	if jsonError != "" {
		http.Error(w, jsonError, http.StatusBadRequest)
		return
	}

	data := make([]UserData, len(usersData.Users))
	copy(data, usersData.Users)

	data = filterByQuery(data, searchReq.Query)
	data = sortUsers(data, searchReq.OrderField, searchReq.OrderBy)
	data = selectUsersRange(data, searchReq.Limit, searchReq.Offset)

	jsonData, err := json.Marshal(data)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error": "%s"}`, err.Error()), http.StatusInternalServerError)
		return
	}
	_, err = w.Write(jsonData)
	if err != nil {
		log.Printf("error writing data to the connection: %v", err)
	}
}

func validateAccessToken(r *http.Request) (jsonError string) {
	accessToken := r.Header.Get("AccessToken")
	if accessToken != serverAccessToken {
		return `{"error": "bad AccessToken"}`
	}
	return ""
}

func validateParams(r *http.Request, sr *SearchRequest) (jsonError string) {
	sr.Query = r.FormValue("query")

	sr.OrderField = strings.ToLower(r.FormValue("order_field"))
	allowedOrderField := []string{"", "id", "age", "name"}
	if !Contains(allowedOrderField, sr.OrderField) {
		return `{"error": "OrderField invalid"}`
	}

	var err error
	sr.OrderBy, err = strconv.Atoi(r.FormValue("order_by"))
	allowedOrderBy := []int{-1, 0, 1}
	if err != nil || !Contains(allowedOrderBy, sr.OrderBy) {
		return `{"error": "OrderBy invalid"}`
	}

	sr.Limit, err = strconv.Atoi(r.FormValue("limit"))
	if err != nil || sr.Limit < 1 {
		return `{"error": "Limit invalid"}`
	}

	sr.Offset, err = strconv.Atoi(r.FormValue("offset"))
	if err != nil || sr.Offset < 0 {
		return `{"error": "Offset invalid"}`
	}

	return ""
}

func filterByQuery(data []UserData, query string) []UserData {
	if query == "" {
		return data
	}
	filteredData := []UserData{}
	for _, v := range data {
		if strings.Contains(v.Name, query) || strings.Contains(v.About, query) {
			filteredData = append(filteredData, v)
		}
	}
	return filteredData
}

func sortUsers(data []UserData, orderField string, orderBy int) []UserData {
	if orderBy == 0 {
		return data
	}
	less := func(i, j int) bool {
		switch orderField {
		case "age":
			return Compare(data[i].Age, data[j].Age, orderBy)
		case "id":
			return Compare(data[i].ID, data[j].ID, orderBy)
		case "name":
			fallthrough
		case "":
			fallthrough
		default:
			return Compare(data[i].Name, data[j].Name, orderBy)
		}
	}
	sort.SliceStable(data, less)
	return data
}

func selectUsersRange(data []UserData, limit int, offset int) []UserData {
	if offset > len(data) {
		return []UserData{}
	}
	data = data[offset:]

	if limit > len(data) {
		return data
	}
	return data[:limit]
}

type Ordered interface {
	int | int8 | int16 | int32 | int64 | float32 | float64 | string
}

func Compare[T Ordered](i T, j T, orderBy int) bool {
	if orderBy == OrderByAsc {
		return i < j
	}
	return i > j
}

func Contains[T comparable](s []T, e T) bool {
	for _, v := range s {
		if v == e {
			return true
		}
	}
	return false
}

func PrepareUsersData() error {
	xmlData, err := os.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("read file error: %w", err)
	}
	err = xml.Unmarshal(xmlData, &usersData)
	if err != nil {
		return fmt.Errorf("unmarshal error: %w", err)
	}
	for idx, userData := range usersData.Users {
		usersData.Users[idx].Name = fmt.Sprintf("%s %s", userData.FirstName, userData.LastName)
	}
	return nil
}
