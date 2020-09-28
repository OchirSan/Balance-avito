package main

import (
	"avito/Balance-avito/middleware"
	"avito/Balance-avito/models"
	"avito/Balance-avito/prometheus"
	"avito/Balance-avito/service"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

type application struct {
	servicePort   string
	statRepository service.Stat
	s             *mux.Router
	pr            *prometheus.Prometheus
}

var conf models.Config

func init() {
	models.LoadConfig(&conf)
}


func main() {
	models.LoadConfig(&conf)
	app := NewApplication(conf)
	app.initServer()

	log.Fatal(http.ListenAndServe(app.servicePort, app.s))

}

func NewApplication(conf models.Config) *application {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		conf.SQLDataBase.Server,conf.SQLDataBase.Port, conf.SQLDataBase.UserID, conf.SQLDataBase.Password, conf.SQLDataBase.Database)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	log.Println("Connect with database")
	err = db.Ping()
	if err != nil{
		panic(err)
	}
	ords := service.NewStatRepository(db)

	return &application{
		servicePort:   ":8080",
		statRepository: ords,
		pr:            prometheus.New("public-stat"),
	}
}

func (app *application) initServer() {
	app.s = mux.NewRouter().StrictSlash(true)

	app.s.Use(middleware.Metrics(app.pr))
	app.s.Handle("/metrics", promhttp.Handler())
	app.s.HandleFunc("/health", StatusHandler).Name("health").Methods(http.MethodGet)
	app.s.HandleFunc("/api/v1/balance", app.GetBalanceHandler).Name("get-balance").
		Methods("GET")
	app.s.HandleFunc("/api/v1/transactions", app.GetTransactionsHandler).Name("get-transactions").
		Methods("GET")
	app.s.HandleFunc("/api/v1/balance", app.AddBalanceHandler).Name("add-balance").
		Methods("POST")
	app.s.HandleFunc("/api/v1/accrual", app.AccrualHandler).Name("accrual").
		Methods("PUT")
	app.s.HandleFunc("/api/v1/debit", app.DebitHandler).Name("debit").
		Methods("PUT")
	app.s.HandleFunc("/api/v1/transfer/{id:[0-9]+}", app.TransferHandler).Name("transfer").
		Methods("PUT")
	app.s.HandleFunc("/api/v1/balance/{id:[0-9]+}", app.DeleteBalanceHandler).Name("delete-balance").
		Methods("DELETE")

}

// statusHandler return 200 for zabbix
func StatusHandler(w http.ResponseWriter, r *http.Request) {

	w.WriteHeader(http.StatusOK)
}

func (app *application) GetBalanceHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	id := r.URL.Query().Get("id")
	p, err := strconv.Atoi(id)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	currency := r.URL.Query().Get("currency")
	byteValue1, err := app.statRepository.GetBalance(p, currency)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	_, err = w.Write(byteValue1)
	if err != nil {
		log.Println(err)
		return
	}

}

func (app *application) AddBalanceHandler(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var ques models.Balance
	err = json.Unmarshal(body, &ques)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = app.statRepository.AddBalance(ques)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (app *application) AccrualHandler(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var rout models.Balance
	err = json.Unmarshal(body, &rout)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = app.statRepository.Accrual(rout.UserId, rout.Amount)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (app *application) DebitHandler(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var rout models.Balance
	err = json.Unmarshal(body, &rout)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = app.statRepository.Debit(rout.UserId, rout.Amount)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (app *application) TransferHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var rout models.Balance
	err = json.Unmarshal(body, &rout)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = app.statRepository.Transfer(rout.UserId,id, rout.Amount)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (app *application) DeleteBalanceHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = app.statRepository.DeleteBalance(id)
	if err != nil  {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (app *application) GetTransactionsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	id := r.URL.Query().Get("id")
	p, err := strconv.Atoi(id)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	count := r.URL.Query().Get("count")
	t, err := strconv.Atoi(count)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	sort := r.URL.Query().Get("sort")
	onSort := r.URL.Query().Get("onSort")

	byteValue1, err := app.statRepository.ListTransactions(p,t, sort, onSort)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	_, err = w.Write(byteValue1)
	if err != nil {
		log.Println(err)
		return
	}

}