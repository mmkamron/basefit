package handler

import (
	"encoding/json"
	"net/http"
)

type trainer struct {
	Name           string
	Specialization string
	Description    string
	Availability   string
	Contact        string
}

func (h *handler) Read(w http.ResponseWriter, r *http.Request) {
	rows, err := h.db.Query("SELECT name, specialization, description, availability, contact FROM trainers")
	if err != nil {
		http.Error(w, "500 Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var trainers []trainer

	for rows.Next() {
		var t trainer
		err := rows.Scan(&t.Name, &t.Specialization, &t.Description, &t.Availability, &t.Contact)
		if err != nil {
			http.Error(w, "500 Internal Server Error", http.StatusInternalServerError)
			return
		}
		trainers = append(trainers, t)
	}

	jsonTrainer, err := json.Marshal(trainers)
	if err != nil {
		http.Error(w, "500 Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonTrainer)
}

func (h *handler) Create(w http.ResponseWriter, r *http.Request) {
	var t trainer
	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
		return
	}
	var id int
	err := h.db.QueryRow(`INSERT INTO trainers (name, specialization, description, availability, contact)
		VALUES ($1, $2, $3, $4, $5) RETURNING trainer_id`, t.Name, t.Specialization, t.Description, t.Availability, t.Contact).Scan(&id)
	if err != nil {
		if id == 0 {
			http.Error(w, "Trainer already exists", http.StatusConflict)
			return
		}
		http.Error(w, "500 Internal Server Error", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/trainers", http.StatusMovedPermanently)
}

func (h *handler) Update(w http.ResponseWriter, r *http.Request) {
	var t trainer
	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
		return
	}
	trainerId := r.PathValue("id")
	_, err := h.db.Exec(`UPDATE trainers SET name = $1, specialization = $2, description = $3,
	availability = $4, contact = $5 WHERE trainer_id = $6`, t.Name, t.Specialization, t.Description, t.Availability, t.Contact, trainerId)
	if err != nil {
		http.Error(w, "500 Internal Server Error", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/trainers", http.StatusMovedPermanently)
}

func (h *handler) Delete(w http.ResponseWriter, r *http.Request) {
	trainerId := r.PathValue("id")
	_, err := h.db.Exec(`DELETE FROM trainers WHERE trainer_id = $1`, trainerId)
	if err != nil {
		http.Error(w, "500 Internal Server Error", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/trainers", http.StatusMovedPermanently)
}
