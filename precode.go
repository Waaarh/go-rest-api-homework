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
func SerAll(r http.ResponseWriter, w *http.Request) {
	resp, err := json.Marshal(tasks)
	if err != nil {
		http.Error(r, err.Error(), http.StatusInternalServerError)
		return
	}
	r.Header().Set("Content-Type", "application/json")
	r.WriteHeader(http.StatusOK)
	r.Write(resp)
}
func SerPost(r http.ResponseWriter, w *http.Request) {
	var task Task
	var buf bytes.Buffer

	_, err := buf.ReadFrom(w.Body)
	if err != nil {
		http.Error(r, err.Error(), http.StatusBadRequest)
		return
	}
	if err = json.Unmarshal(buf.Bytes(), &task); err != nil {
		http.Error(r, err.Error(), http.StatusBadRequest)
		return
	}
	r.Header().Set("Content-Type", "application/json")
	r.WriteHeader(http.StatusCreated)
}
func searchID(r http.ResponseWriter, w *http.Request) {
	id := chi.URLParam(w, "id")

	artist, ok := tasks[id]
	if !ok {
		http.Error(r, "err", http.StatusBadRequest)
		return
	}
	resp, err := json.Marshal(artist)
	if err != nil {
		http.Error(r, err.Error(), http.StatusBadRequest)
		return
	}
	r.Header().Set("Content-Type", "application/json")
	r.WriteHeader(http.StatusOK)
	r.Write(resp)
}
func DelID(r http.ResponseWriter, w *http.Request) {
	// Извлекаем id из URL: /tasks/{id}
	id := chi.URLParam(w, "id")

	// Проверяем, существует ли задача с таким id
	_, exists := tasks[id]
	if !exists {
		http.Error(r, "Задача с таким ID не найдена", http.StatusBadRequest)
		return
	}

	// Удаляем задачу из мапы
	delete(tasks, id)

	// Возвращаем статус 200 OK
	r.WriteHeader(http.StatusOK)
	fmt.Fprint(r, `{"status": "deleted", "id": "`+id+`"}`)
}
func main() {
	r := chi.NewRouter()

	// здесь регистрируйте ваши обработчики
	r.Get("/tasks", SerAll)
	r.Post("/tasks", SerPost)
	r.Get("/search/{id}", searchID)
	r.Delete("/tasks/{id}", DelID)

	if err := http.ListenAndServe(":8080", r); err != nil {
		fmt.Printf("Ошибка при запуске сервера: %s", err.Error())
		return
	}
}
