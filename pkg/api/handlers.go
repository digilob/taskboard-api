package api

import (
	"encoding/json"
	"net/http"

	r "github.com/digilob/taskboard-api/pkg/repository"
	"github.com/go-chi/chi"
)

// r.Tasks is a collection of r.Task.
var Tasks []r.Task

func GetTasksHandler(repo r.Repository, w http.ResponseWriter, r *http.Request) {
	tasks, err := repo.GetAll.Tasks()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(tasks)
}

func CreateTaskHandler(repo r.Repository, w http.ResponseWriter, r *http.Request) {
	var task r.Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	createdTask, err := repo.Create.Task(task)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createdTask)
}

func UpdateTaskHandler(repo r.Repository, w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var task r.Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	updatedTask, err := repo.UpdateTask(id, task)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(updatedTask)
}

func DeleteTaskHandler(repo r.Repository, w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	if err := repo.DeleteTask(id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
