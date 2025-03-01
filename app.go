package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	_ "github.com/go-sql-driver/mysql"
)

type App struct {
	DB     *sql.DB
	Router *mux.Router
}

// Initialize the DB and Router
func (app *App) Initialize(DbUser string, DbPassword string, DbName string) error {

	//Initialize DB
	dataSourceName := fmt.Sprintf("%v:%v@tcp(127.0.0.1:3306)/%v", DbUser, DbPassword, DbName)
	var err error
	app.DB, err = sql.Open("mysql", dataSourceName)

	if err != nil {
		return err
	}

	//Initialize Router

	app.Router = mux.NewRouter().StrictSlash(true)
	app.handleRoutes()
	return nil

}

func (app *App) Run(address string) {

	log.Fatal(http.ListenAndServe(address, app.Router))
}

func sendResponse(w http.ResponseWriter, status_code int, response_body interface{}) {

	response, _ := json.Marshal(response_body)
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(status_code)
	w.Write(response)
}

func sendError(w http.ResponseWriter, status_code int, err string) {

	error_message := map[string]string{"error": err}
	sendResponse(w, status_code, error_message)
}

func (app *App) getProducts(w http.ResponseWriter, r *http.Request) {

	products, err := getProducts(app.DB)

	if err != nil {
		sendError(w, http.StatusInternalServerError, "err")
		return
	}

	sendResponse(w, http.StatusOK, products)

}

// Get Product
func (app *App) getProduct(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	key, err := strconv.Atoi(vars["id"])
	if err != nil {
		sendError(w, http.StatusBadRequest, "Invalid Product ID")
		return
	}
	p := product{ID: key}
	err = p.getProduct(app.DB)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			sendError(w, http.StatusNotFound, "Product not Found")
		default:
			sendError(w, http.StatusInternalServerError, err.Error())

		}
		return

	}
	sendResponse(w, http.StatusOK, p)

}

func (app *App) createProduct(w http.ResponseWriter, r *http.Request) {
	var p product

	err := json.NewDecoder(r.Body).Decode(&p)
	log.Println("Product", p)
	if err != nil {
		sendError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	err = p.createProduct(app.DB)
	if err != nil {

		sendError(w, http.StatusInternalServerError, err.Error())
		return
	}

	sendResponse(w, http.StatusCreated, p)

}

//Update Product

func (app *App) updateProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key, err := strconv.Atoi(vars["id"])
	if err != nil {
		sendError(w, http.StatusBadRequest, "Invalid Product ID")
		return
	}

	var p product

	err = json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		sendError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	p.ID = key
	err = p.updateProduct(app.DB)
	if err != nil {
		sendError(w, http.StatusInternalServerError, err.Error())
		return
	}
	sendResponse(w, http.StatusOK, p)
}

// Delete Product
func (app *App) deleteProduct(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	key, err := strconv.Atoi(vars["id"])
	if err != nil {
		sendError(w, http.StatusBadRequest, "Invalid Product ID")
		return
	}

	p := product{ID: key}

	err = p.deleteProduct(app.DB)
	if err != nil {
		sendError(w, http.StatusInternalServerError, err.Error())
		return
	}
	sendResponse(w, http.StatusOK, map[string]string{"result": "Deleted successfully"})

}

// Handle Routes
func (app *App) handleRoutes() {

	app.Router.HandleFunc("/products", app.getProducts).Methods("GET")
	app.Router.HandleFunc("/product/{id}", app.getProduct).Methods("GET")
	app.Router.HandleFunc("/product", app.createProduct).Methods("POST")
	app.Router.HandleFunc("/product/{id}", app.updateProduct).Methods("PUT")
	app.Router.HandleFunc("/product/{id}", app.deleteProduct).Methods("DELETE")

}
