package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
)

// Task ...
type Task struct {
	ID           string   `json:"id"`
	Description  string   `json:"description"`
	Note         string   `json:"note"`
	Applications []string `json:"applications"`
}

var tasks = map[string]Task{
	"1": {
		ID:          "1",
		Description: "Сделать финальное задание темы REST API",
		Note:        "Если сегодня сделаю, то завтра будет свободный день. Ура!",
		Applications: []string{
			"VS Code",
			"Terminal",
			"git",
		},
	},
	"2": {
		ID:          "2",
		Description: "Протестировать финальное задание с помощью Postmen",
		Note:        "Лучше это делать в процессе разработки, каждый раз, когда запускаешь сервер и проверяешь хендлер",
		Applications: []string{
			"VS Code",
			"Terminal",
			"git",
			"Postman",
		},
	},
}

// Ниже напишите обработчики для каждого эндпоинта

// Получить все задания
func getTasks(w http.ResponseWriter, r *http.Request) {
	resp, err := json.Marshal(tasks)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	//IDE подсказывает сделать обработку ошибок в строчке ниже, но не особо понял как ее осуществить
	//Если есть возможность описать как сделать, буду благодарен
	w.Write(resp)
}

// Получить задание по id
func getTask(w http.ResponseWriter, r *http.Request) {
	paramId := chi.URLParam(r, "id")
	task, ok := tasks[paramId]
	if !ok {
		//В ТЗ было написано сделать статус BadRequest, мне кажется лучше будет NotFound
		http.Error(w, "task not found", http.StatusNotFound)
	}
	resp, err := json.Marshal(task)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	//IDE подсказывает сделать обработку ошибок в строчке ниже, но не особо понял как ее осуществить
	//Если есть возможность описать как сделать, буду благодарен
	w.Write(resp)
}

// Создать задание
func postTask(w http.ResponseWriter, r *http.Request) {
	var task Task
	var buffer bytes.Buffer
	_, err := buffer.ReadFrom(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := json.Unmarshal(buffer.Bytes(), &task); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	tasks[task.ID] = task
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
}

// Удалить задание из map
func deleteTask(w http.ResponseWriter, r *http.Request) {
	paramId := chi.URLParam(r, "id")

	_, ok := tasks[paramId]
	if !ok {
		//В ТЗ было написано сделать статус BadRequest, мне кажется лучше будет NotFound
		http.Error(w, "task not found", http.StatusNotFound)
		return
	}
	delete(tasks, paramId)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

func main() {
	r := chi.NewRouter()
	r.Get("/tasks", getTasks)           //Запрос всех заданий
	r.Post("/tasks", postTask)          //Запрос на создание задания
	r.Get("/tasks/{id}", getTask)       //Запрос на задание по id
	r.Delete("/tasks/{id}", deleteTask) //Запрос на удаление задания

	if err := http.ListenAndServe(":8080", r); err != nil {
		fmt.Printf("Ошибка при запуске сервера: %s", err.Error())
		return
	}
}
