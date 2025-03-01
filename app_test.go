package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

var a App

func TestMain(m *testing.M) {

	err := a.Initialize(DbUser, DbPassword, "test")
	if err != nil {
		log.Fatal(err)
	}
	createTestTable()
	m.Run()

}

func createTestTable() {

	createTableQuery := `CREATE TABLE IF NOT EXISTS products (
		      id int NOT NULL AUTO_INCREMENT,
			  name varchar(255) NOT NULL,
			  quantity int,
		      price float(10,7),
			  PRIMARY KEY (id)
	);`

	_, err := a.DB.Exec(createTableQuery)
	if err != nil {
		log.Fatal(err)
	}
}

func clearTable() {
	a.DB.Exec("DELETE from products")
	a.DB.Exec("ALTER table products AUTO_INCREMENT=1")
}

func addProduct(name string, quantity int, price float64) {
	query := fmt.Sprintf("INSERT into products(name,quantity,price) VALUES('%v','%v','%v')", name, quantity, price)
	_, err := a.DB.Exec(query)
	if err != nil {
		log.Println(err)
	}

}

func TestGetProduct(t *testing.T) {

	clearTable()
	addProduct("pen", 100, 50)
	request, _ := http.NewRequest("GET", "/product/1", nil)
	response := sendRequest(request)
	compareResponse(t, http.StatusOK, response.Code)
}

func compareResponse(t *testing.T, expectedResponseCode int, actualResponseCode int) {

	if expectedResponseCode != actualResponseCode {

		t.Errorf("Expected Response %v but got response %v", expectedResponseCode, actualResponseCode)
	}

}

func sendRequest(request *http.Request) *httptest.ResponseRecorder {

	recorder := httptest.NewRecorder()

	a.Router.ServeHTTP(recorder, request)
	return recorder
}

func TestCreateProduct(t *testing.T) {
	clearTable()
	product := []byte(`{ "name":"boxes", "quantity":5, "price":20 }`)
	request, _ := http.NewRequest("POST", "/product", bytes.NewBuffer(product))
	response := sendRequest(request)
	compareResponse(t, http.StatusCreated, response.Code)

	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)

	if m["name"] != "boxes" {
		t.Errorf("Expected product %v but got %v", "boxes", m["name"])
	}

	if m["quantity"] != 5.0 {
		t.Errorf("Expected quantity %v but got %v", 5, m["quantity"])
	}
}

func TestDeleteProduct(t *testing.T) {

	clearTable()
	addProduct("laptop", 12, 100)

	request, _ := http.NewRequest("GET", "/product/1", nil)
	response := sendRequest(request)
	compareResponse(t, http.StatusOK, response.Code)

	request, _ = http.NewRequest("DELETE", "/product/1", nil)
	response = sendRequest(request)
	compareResponse(t, http.StatusOK, response.Code)

	request, _ = http.NewRequest("GET", "/product/1", nil)
	response = sendRequest(request)
	compareResponse(t, http.StatusNotFound, response.Code)

}

func TestUpdateProduct(t *testing.T) {
	clearTable()
	addProduct("laptop", 10, 100)
	request, _ := http.NewRequest("GET", "/product/1", nil)
	response := sendRequest(request)

	var Value map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &Value)

	payload := []byte(`{"name":"laptop", "quantity": 20, "price": 100}`)
	request, _ = http.NewRequest("PUT", "/product/1", bytes.NewBuffer(payload))
	response = sendRequest(request)

	var newValue map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &newValue)

	if Value["id"] != newValue["id"] {
		t.Errorf("Expeced id %v, but got id %v", Value["id"], newValue["id"])
	}

	if Value["quantity"] == newValue["quantity"] {
		t.Errorf("Expeced quantity %v, but got quantity %v", Value["quantity"], newValue["quantity"])
	}

	if Value["price"] != newValue["price"] {
		t.Errorf("Expeced price %v, but got price %v", Value["price"], newValue["price"])
	}
}
