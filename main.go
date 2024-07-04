package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5"
	_ "github.com/lib/pq"
)

type ActiveCall struct {
	Id          int        `json:"id"`
	Address     string     `json:"address"`
	PatientName string     `json:"patient_name"`
	OpenTs      time.Time  `json:"open_ts"`
	ClosedTs    *time.Time `json:"closed_ts,omitempty"`
	UpdTs       time.Time  `json:"upd_ts"`
	ResponderId *int       `json:"responder_id,omitempty"`
	Status      string     `json:"status"`
	Notes       string     `json:"notes"`
}

type OnCall struct {
	Id        int       `json:"id"`
	PriestId  int       `json:"priest_id"`
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
	CreatedTs time.Time `json:"created_ts"`
	UpdTs     time.Time `json:"upd_ts"`
	Status    string    `json:"status"`
}

type Priest struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

var db *pgx.Conn

func initDB(db *pgx.Conn) error {
	script, err := os.ReadFile("init.sql")
	if err != nil {
		return err
	}

	_, err = db.Exec(context.Background(), string(script))
	return err
}

// ActiveCalls CRUD
func getActiveCalls(w http.ResponseWriter, r *http.Request) {
	var activeCalls []ActiveCall
	rows, err := db.Query(context.Background(), "SELECT * FROM OCM.ACTIVE_CALLS")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var ac ActiveCall
		err := rows.Scan(
			&ac.Id,
			&ac.Address,
			&ac.PatientName,
			&ac.OpenTs,
			&ac.ClosedTs,
			&ac.UpdTs,
			&ac.ResponderId,
			&ac.Status,
			&ac.Notes,
		)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		activeCalls = append(activeCalls, ac)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(activeCalls)
}

func createActiveCall(w http.ResponseWriter, r *http.Request) {
	var ac ActiveCall
	err := json.NewDecoder(r.Body).Decode(&ac)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ac.OpenTs = time.Now()
	ac.UpdTs = time.Now()

	err = db.QueryRow(context.Background(),
		"INSERT INTO OCM.ACTIVE_CALLS (ADDRESS, PATIENT_NAME, OPEN_TS, STATUS, NOTES) VALUES ($1, $2, $3, $4, $5) RETURNING ID",
		ac.Address, ac.PatientName, ac.OpenTs, ac.Status, ac.Notes).Scan(&ac.Id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ac)
}

func updateActiveCall(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	var ac ActiveCall
	err = json.NewDecoder(r.Body).Decode(&ac)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ac.UpdTs = time.Now()

	_, err = db.Exec(context.Background(),
		"UPDATE OCM.ACTIVE_CALLS SET ADDRESS = $1, PATIENT_NAME = $2, STATUS = $3, NOTES = $4, UPD_TS = $5 WHERE ID = $6",
		ac.Address, ac.PatientName, ac.Status, ac.Notes, ac.UpdTs, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func deleteActiveCall(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	_, err = db.Exec(context.Background(), "DELETE FROM OCM.ACTIVE_CALLS WHERE ID = $1", id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// OnCall CRUD
func getOnCalls(w http.ResponseWriter, r *http.Request) {
	var onCalls []OnCall
	rows, err := db.Query(context.Background(), "SELECT * FROM OCM.ON_CALL")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var oc OnCall
		err := rows.Scan(&oc.Id, &oc.PriestId, &oc.StartTime, &oc.EndTime, &oc.CreatedTs, &oc.UpdTs, &oc.Status)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		onCalls = append(onCalls, oc)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(onCalls)
}

func createOnCall(w http.ResponseWriter, r *http.Request) {
	var oc OnCall
	err := json.NewDecoder(r.Body).Decode(&oc)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	oc.CreatedTs = time.Now()
	oc.UpdTs = time.Now()

	err = db.QueryRow(context.Background(),
		"INSERT INTO OCM.ON_CALL (PRIEST_ID, START_TIME, END_TIME, CREATED_TS, UPD_TS, STATUS) VALUES ($1, $2, $3, $4, $5, $6) RETURNING ID",
		oc.PriestId, oc.StartTime, oc.EndTime, oc.CreatedTs, oc.UpdTs, oc.Status).Scan(&oc.Id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(oc)
}

func updateOnCall(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	var oc OnCall
	err = json.NewDecoder(r.Body).Decode(&oc)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	oc.UpdTs = time.Now()

	_, err = db.Exec(context.Background(),
		"UPDATE OCM.ON_CALL SET PRIEST_ID = $1, START_TIME = $2, END_TIME = $3, UPD_TS = $4, STATUS = $5 WHERE ID = $6",
		oc.PriestId, oc.StartTime, oc.EndTime, oc.UpdTs, oc.Status, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func deleteOnCall(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	_, err = db.Exec(context.Background(), "DELETE FROM OCM.ON_CALL WHERE ID = $1", id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// Priests CRUD
func getPriests(w http.ResponseWriter, r *http.Request) {
	var priests []Priest
	rows, err := db.Query(context.Background(), "SELECT * FROM OCM.PRIESTS")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var p Priest
		err := rows.Scan(&p.Id, &p.Name)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		priests = append(priests, p)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(priests)
}

func createPriest(w http.ResponseWriter, r *http.Request) {
	var p Priest
	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = db.QueryRow(context.Background(),
		"INSERT INTO OCM.PRIESTS (NAME) VALUES ($1) RETURNING ID", p.Name).Scan(&p.Id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(p)
}

func updatePriest(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	var p Priest
	err = json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err = db.Exec(context.Background(), "UPDATE OCM.PRIESTS SET NAME = $1 WHERE ID = $2", p.Name, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func deletePriest(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	_, err = db.Exec(context.Background(), "DELETE FROM OCM.PRIESTS WHERE ID = $1", id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func main() {
	var err error
	db, err = pgx.Connect(context.Background(), os.Getenv("POSTGRES_URL"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer db.Close(context.Background())

	err = initDB(db)
	if err != nil {
		fmt.Println(err)
	}

	r := mux.NewRouter()

	// ActiveCalls routes
	r.HandleFunc("/active-calls", getActiveCalls).Methods("GET")
	r.HandleFunc("/active-calls", createActiveCall).Methods("POST")
	r.HandleFunc("/active-calls/{id}", updateActiveCall).Methods("PUT")
	r.HandleFunc("/active-calls/{id}", deleteActiveCall).Methods("DELETE")

	// OnCall routes
	r.HandleFunc("/on-calls", getOnCalls).Methods("GET")
	r.HandleFunc("/on-calls", createOnCall).Methods("POST")
	r.HandleFunc("/on-calls/{id}", updateOnCall).Methods("PUT")
	r.HandleFunc("/on-calls/{id}", deleteOnCall).Methods("DELETE")

	// Priests routes
	r.HandleFunc("/priests", getPriests).Methods("GET")
	r.HandleFunc("/priests", createPriest).Methods("POST")
	r.HandleFunc("/priests/{id}", updatePriest).Methods("PUT")
	r.HandleFunc("/priests/{id}", deletePriest).Methods("DELETE")

	fmt.Println("Server is running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
