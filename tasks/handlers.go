package tasks

import (
	"mzaky/simple-crud-with-login/util"

	"net/http"
	"io/ioutil"
	"encoding/json"
	"database/sql"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"github.com/Nerzal/gocloak/v11"
)

func NewHandler(database *sql.DB, keycloak gocloak.GoCloak) *handler {
	return &handler{database: database, keycloak: keycloak}
}

type handler struct {
	database *sql.DB
	keycloak gocloak.GoCloak
}

func (h *handler) GetTasks(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	ctx := r.Context()
	select {
	case <-ctx.Done():
		util.WriteErrorResponse(w, "Timeout", http.StatusRequestTimeout)
		return
	default:
	}

	userInfo, err := util.AuthorizeUser(ctx, h.keycloak, r.Header.Get("Authorization"))
	if err != nil {
		util.WriteErrorResponse(w, "User not authorized to use this function: " + err.Error(), http.StatusUnauthorized)
		return
	}

	rows, errQuery := h.database.Query("SELECT id, description, difficulty, done FROM tasks WHERE owner_identifier = $1", userInfo.PreferredUsername)
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

	userInfo, err := util.AuthorizeUser(ctx, h.keycloak, r.Header.Get("Authorization"))
	if err != nil {
		util.WriteErrorResponse(w, "User not authorized to use this function: " + err.Error(), http.StatusUnauthorized)
		return
	}

	idStr := params.ByName("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		util.WriteErrorResponse(w, "Invalid parameter: " + params.ByName("id"), http.StatusBadRequest)
		return
	}

	rows, errQuery := h.database.Query("SELECT id, description, difficulty, done FROM tasks WHERE id = $1 AND owner_identifier = $2", id, userInfo.PreferredUsername)
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

	userInfo, err := util.AuthorizeUser(ctx, h.keycloak, r.Header.Get("Authorization"))
	if err != nil {
		util.WriteErrorResponse(w, "User not authorized to use this function: " + err.Error(), http.StatusUnauthorized)
		return
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

	_, errQuery := h.database.Exec("INSERT INTO tasks (description, difficulty, done, owner_identifier) VALUES ($1, $2, $3, $4)", taskReq.Description, taskReq.Difficulty, taskReq.Done, userInfo.PreferredUsername)
	if errQuery != nil {
		util.WriteErrorResponse(w, "Error storing task to database: " + errQuery.Error(), http.StatusInternalServerError)
		return
	}

	util.WriteJSONResponse(w, util.MessageResponse{Message: "Task created"}, http.StatusOK)
}
