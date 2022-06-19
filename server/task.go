package server

import (
	"encoding/json"
	"net/http"

	"github.com/bariis/quin-tech-case/task"
	"github.com/gorilla/mux"
)

// sendResponse wraps the http response functionalities based on parameters.
func (s *Server) sendResponse(w http.ResponseWriter, data []byte, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(data)
}

// registerTaskRoutes is a helper function to register task routes.
func (s *Server) registerTaskRoutes(r *mux.Router) {
	r.HandleFunc("/save", s.handleSaveTask).Methods("POST")
	r.HandleFunc("/retrieve", s.handleGetTasks).Methods("GET")
	r.HandleFunc("/save/{code}", s.handleSaveSubTask).Methods("POST")
	r.HandleFunc("/edit", s.handleUpdateTask).Methods("PUT")
	r.HandleFunc("/delete/{code}", s.handleDeleteEntry).Methods("DELETE")
}

// handleSaveTasks handles the "POST /save" route. This route saves the entry.
func (s *Server) handleSaveTask(w http.ResponseWriter, r *http.Request) {
	var entry task.Entry
	err := json.NewDecoder(r.Body).Decode(&entry)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	
	newEntry, err := s.TaskService.Save(&entry)
	if err != nil {
		ser := map[string]string{
			"error": "empty name for the entry",
		}
		r, _ := json.Marshal(ser)
		s.sendResponse(w, r, 400)
		return
	}
	res, err := json.Marshal(newEntry)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	s.sendResponse(w, res, http.StatusOK)
}

// handleSaveSubTask handles the "PUT /save/{code}" route. This route saves the sub-entry within a parent.
func (s *Server) handleSaveSubTask(w http.ResponseWriter, r *http.Request) {
	var entry task.Entry
	err := json.NewDecoder(r.Body).Decode(&entry)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	parentCode := mux.Vars(r)["code"]
	newEntry, err := s.TaskService.SaveSubEntry(&entry, parentCode)
	if err != nil {
		ser := map[string]string{
			"error": err.Error(),
		}
		r, _ := json.Marshal(ser)
		s.sendResponse(w, r, 400)
		return
	}

	res, err := json.Marshal(newEntry)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	s.sendResponse(w, res, http.StatusOK)
}

// handleGetTasks handles the "GET /retrieve" route. This route lists all the entries.
func (s *Server) handleGetTasks(w http.ResponseWriter, r *http.Request) {
	entries, err := s.TaskService.All()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	res, err := json.Marshal(entries)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	s.sendResponse(w, res, http.StatusOK)
}


// handleUpdateTasks handles the "PUT /edit" route. This route updates the specified entry.
func (s *Server) handleUpdateTask(w http.ResponseWriter, r *http.Request) {
	var entry task.Entry
	err := json.NewDecoder(r.Body).Decode(&entry)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	updatedEntry, err := s.TaskService.Update(&entry)
	if err != nil {
		ser := map[string]string{
			"error": "no such entry",
		}
		r, _ := json.Marshal(ser)
		s.sendResponse(w, r, 400)
		return
	}

	res, err := json.Marshal(updatedEntry)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	s.sendResponse(w, res, http.StatusOK)
}

// handleDeleteEntry handles the "DELETE /delete" route. This route deletes the specified entry.
func (s *Server) handleDeleteEntry(w http.ResponseWriter, r *http.Request) {
	entryCode := mux.Vars(r)["code"]

	err := s.TaskService.Delete(entryCode)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return 
	}
	s.sendResponse(w, []byte("entry successfully deleted"), http.StatusOK)
}
