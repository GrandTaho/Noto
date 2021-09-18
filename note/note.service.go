package note

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/GrandTaho/noto/cors"
)

const path = "notes"

func SetupRoutes(apiBasePath string) {
	notesHandler := http.HandlerFunc(handleNotes)

	http.Handle(fmt.Sprintf("%s/%s", apiBasePath, path), cors.Middleware(notesHandler))
}

func handleNotes(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		notes, err := getNotes()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		notesJson, err := json.Marshal(notes)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(notesJson)
	case http.MethodPost:
		var note Note
		err := json.NewDecoder(r.Body).Decode(&note)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		noteId, err := insertNote(note)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(fmt.Sprintf(`{"noteId":%d}`, noteId)))
	case http.MethodOptions:
		return
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
}

func handleNote(w http.ResponseWriter, r *http.Request) {
	urlPathSegments := strings.Split(r.URL.Path, fmt.Sprintf("%s/", path))
	if len(urlPathSegments[1:]) > 1 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	noteId, err := strconv.Atoi(urlPathSegments[len(urlPathSegments)-1])
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	switch r.Method {
	case http.MethodGet:
		note, err := getNote(noteId)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if note == nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		j, err := json.Marshal(note)
		if err != nil {
			log.Println("Error marshalling json", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		_, err = w.Write(j)
		if err != nil {
			log.Println("Unable to write json to response", err)
			return
		}
	case http.MethodPut:
		var note Note
		err := json.NewDecoder(r.Body).Decode(&note)
		if err != nil {
			log.Println("Unable to decode the body of the request", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		err = updateNote(note)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	case http.MethodDelete:
		err := removeNote(noteId)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	case http.MethodOptions:
		return
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
}
