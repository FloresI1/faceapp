package main

import (
	"fmt"
	"image"
	"log"
	"net/http"
	"path/filepath"
)

type ImageStore struct {
	Images map[string]image.Image
}

func main() {
	store := &ImageStore{
		Images: make(map[string]image.Image),
	}

	http.HandleFunc("/upload", store.uploadHandler)
	log.Println("Сервер запущен на порту 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func (store *ImageStore) uploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Только POST запросы поддерживаются", http.StatusMethodNotAllowed)
		return
	}

	logRequestHeaders(r)

	err := r.ParseMultipartForm(10 << 20) // Размер формы до 10 MB
	if err != nil {
		http.Error(w, "Ошибка обработки формы", http.StatusInternalServerError)
		return
	}

	file, fileHeader, err := r.FormFile("image")
	if err != nil {
		http.Error(w, "Ошибка получения файла", http.StatusBadRequest)
		return
	}
	defer file.Close()

	log.Printf("Получен файл с именем: %s\n", fileHeader.Filename)

	fileName := filepath.Base(fileHeader.Filename)
	if fileName == "" {
		http.Error(w, "Имя файла не найдено", http.StatusBadRequest)
		return
	}

	img, _, err := image.Decode(file)
	if err != nil {
		http.Error(w, "Ошибка декодирования изображения", http.StatusBadRequest)
		return
	}

	store.Images[fileName] = img

	log.Printf("Изображение успешно загружено и сохранено под именем: %s\n", fileName)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("Изображение успешно загружено и сохранено под именем %s", fileName)))
}

// Логгирование
func logRequestHeaders(r *http.Request) {
	log.Println("Получены заголовки запроса:")
	for name, values := range r.Header {
		for _, value := range values {
			log.Printf("%s: %s", name, value)
		}
	}
}
