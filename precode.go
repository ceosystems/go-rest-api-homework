package main

import (
    "encoding/json"
    "fmt"
    "net/http"
	 "bytes"

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

func getTasks(w http.ResponseWriter, r *http.Request) {
    resp, err := json.Marshal(tasks)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    w.Write(resp)
}

func postTasks(w http.ResponseWriter, r *http.Request) {
    var task Task = Task{}
    var buf bytes.Buffer
	
    _, err := buf.ReadFrom(r.Body)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    if err = json.Unmarshal(buf.Bytes(), &task); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

	tasks[task.ID] = task

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
}

func getTask(w http.ResponseWriter, r *http.Request) {
    id := chi.URLParam(r, "id")

    tasks, ok := tasks[id]
    if !ok {
        http.Error(w, "Задача не найдена", http.StatusNoContent)
        return
    }

    resp, err := json.Marshal(tasks)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    w.Write(resp)
}


func deleteTask(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
    if _, ok := tasks[id]; !ok {
        http.Error(w, "Задача не найдена", http.StatusNotFound)
        return
    }
    delete(tasks, id)
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
}

func main() {
	r := chi.NewRouter()

	// здесь регистрируйте ваши обработчики
	r.Get("/tasks", getTasks)
	r.Post("/tasks", postTasks)
	r.Get("/task/{id}", getTask)
	r.Delete("/task/{id}", deleteTask)

	if err := http.ListenAndServe(":8080", r); err != nil {
		fmt.Printf("Ошибка при запуске сервера: %s", err.Error())
		return
	}
}
