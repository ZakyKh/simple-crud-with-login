package tasks

import (
	"mzaky/simple-crud-with-login/util"

	"net/http"
	"io/ioutil"
	"encoding/json"
	"database/sql"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

func NewHandler(database *sql.DB) *handler {
	return &handler{database: database}
}

type handler struct {
	database *sql.DB
}

func (h *handler) GetTasks(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	ctx := r.Context()
	select {
	case <-ctx.Done():
		util.WriteErrorResponse(w, "Timeout", http.StatusRequestTimeout)
		return
	default:
	}

	rows, errQuery := h.database.Query("SELECT id, description, difficulty, done FROM tasks")
	if errQuery != nil {
		util.WriteErrorResponse(w, "Error retrieving tasks from database: " + errQuery.Error(), http.StatusInternalServerError)
		return
	}

	tasks := []Task{}
	for rows.Next() {
		var task Task
		errScan := rows.Scan(&task.Id, &task.Description, &task.Difficulty, &task.Done)
		if errScan != nil {
			util.WriteErrorResponse(w, "Error processing tasks retrieved from database", http.StatusInternalServerError)
			return
		}
		tasks = append(tasks, task)
	}

	util.WriteJSONResponse(w, tasks, http.StatusOK)
}

func (h *handler) GetTask(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	ctx := r.Context()
	select {
	case <-ctx.Done():
		util.WriteErrorResponse(w, "Timeout", http.StatusRequestTimeout)
		return
	default:
	}

	idStr := params.ByName("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		util.WriteErrorResponse(w, "Invalid parameter: " + params.ByName("id"), http.StatusBadRequest)
		return
	}

	rows, errQuery := h.database.Query("SELECT id, description, difficulty, done FROM tasks WHERE id = $1", id)
	if errQuery != nil {
		util.WriteErrorResponse(w, "Error retrieving tasks from database: " + errQuery.Error(), http.StatusInternalServerError)
		return
	}

	var task Task
	if rows.Next() {
		errScan := rows.Scan(&task.Id, &task.Description, &task.Difficulty, &task.Done)
		if errScan != nil {
			util.WriteErrorResponse(w, "Error processing tasks retrieved from database", http.StatusInternalServerError)
			return
		}
	}

	if task.Id != id {
		util.WriteErrorResponse(w, "Task not found", http.StatusNotFound)
		return
	}

	util.WriteJSONResponse(w, task, http.StatusOK)
}

func (h *handler) CreateTask(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	ctx := r.Context()
	select {
	case <-ctx.Done():
		util.WriteErrorResponse(w, "Timeout", http.StatusRequestTimeout)
		return
	default:
	}

	reqBodyBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		util.WriteErrorResponse(w, "Error parsing request body: " + err.Error(), http.StatusBadRequest)
		return
	}

	var taskReq TaskRequest
	err = json.Unmarshal(reqBodyBytes, &taskReq)
	if err != nil {
		util.WriteErrorResponse(w, "Error parsing request body: " + err.Error(), http.StatusBadRequest)
		return
	}

	_, errQuery := h.database.Exec("INSERT INTO tasks (description, difficulty, done) VALUES ($1, $2, $3)", taskReq.Description, taskReq.Difficulty, taskReq.Done)
	if errQuery != nil {
		util.WriteErrorResponse(w, "Error storing task to database: " + errQuery.Error(), http.StatusInternalServerError)
		return
	}

	util.WriteJSONResponse(w, util.MessageResponse{Message: "Task created"}, http.StatusOK)
}
