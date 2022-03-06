package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"sort"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/unrolled/render"
	"github.com/urfave/negroni"
)

var renderer *render.Render

type Result struct {
	isSuccess bool `json:"isSuccess"`
}

type Task struct {
	Id          int    `json:"id,omitempty"`
	Title       string `json:"title"`
	IsCompleted bool   `json:"isCompleted,omitempty"`
}

type Tasks []Task

func (t Tasks) Len() int {
	return len(t)
}

func (t Tasks) Swap(i, j int) {
	t[i], t[j] = t[j], t[i]
}

func (t Tasks) Less(i, j int) bool {
	return t[i].Id > t[j].Id
}

var tasks map[int]Task
var lastTaskId int = 0

func main() {
	renderer = render.New()
	router := NewRouter()
	advancedRouter := negroni.Classic()
	advancedRouter.UseHandler(router)

	log.Println("Task Manager is started...")

	port := os.Getenv("PORT")
	err := http.ListenAndServe(":"+port, advancedRouter)
	if err != nil {
		panic(err)
	}
}

func NewRouter() http.Handler {
	tasks = make(map[int]Task)

	router := mux.NewRouter()
	router.Handle("/", http.FileServer(http.Dir("public")))
	router.HandleFunc("/tasks", ReadTasks).Methods("GET")
	router.HandleFunc("/tasks", CreateTask).Methods("POST")
	router.HandleFunc("/tasks/{id:[0-9]+}", DeleteTask).Methods("DELETE")
	router.HandleFunc("/tasks/{id:[0-9]+}", UpdateTask).Methods("PUT")

	return router
}

func ReadTasks(w http.ResponseWriter, r *http.Request) {
	targetTasks := make(Tasks, 0)

	for _, task := range tasks {
		targetTasks = append(targetTasks, task)
	}
	sort.Sort(targetTasks)

	renderer.JSON(w, http.StatusOK, targetTasks)
	return
}

func CreateTask(w http.ResponseWriter, r *http.Request) {
	var newTask Task

	err := json.NewDecoder(r.Body).Decode(&newTask)
	if err != nil {
		log.Fatal(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	lastTaskId++
	newTask.Id = lastTaskId
	tasks[lastTaskId] = newTask

	renderer.JSON(w, http.StatusCreated, newTask)
	return
}

func DeleteTask(w http.ResponseWriter, r *http.Request) {
	variables := mux.Vars(r)

	id, _ := strconv.Atoi(variables["id"])
	if _, ok := tasks[id]; ok {
		delete(tasks, id)
		renderer.JSON(w, http.StatusOK, Result{true})
		return
	}

	renderer.JSON(w, http.StatusNotFound, Result{false})
	return
}

func UpdateTask(w http.ResponseWriter, r *http.Request) {
	variables := mux.Vars(r)

	var newTask Task

	err := json.NewDecoder(r.Body).Decode(&newTask)
	if err != nil {
		log.Fatal(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	id, _ := strconv.Atoi(variables["id"])
	if targetTask, ok := tasks[id]; ok {
		targetTask.Title = newTask.Title
		targetTask.IsCompleted = newTask.IsCompleted

		renderer.JSON(w, http.StatusOK, Result{true})
		return
	}

	renderer.JSON(w, http.StatusBadRequest, Result{false})
	return
}
