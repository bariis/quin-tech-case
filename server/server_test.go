package server

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/bariis/quin-tech-case/task"
	"github.com/gorilla/mux"
)

func TestSaveEntryHandler(t *testing.T) {
	srv := setRequiredServices()

	err := map[string]string{
		"error": "empty name for the entry",
	}
	case2Expected, _ := json.Marshal(err)
	caseExpected, _ := json.Marshal(&task.Entry{Id: 1, Code: "TEST-1", Name: "develop x feature", Assignee: "Baris", Tags: []string{"red"}, Subs: map[string]*task.Entry{}})
	tt := []struct {
		name  string
		entry *task.Entry
		want  []byte
		code  int
	}{
		{
			"single",
			&task.Entry{Name: "develop x feature", Assignee: "Baris", Tags: []string{"red"}, Subs: map[string]*task.Entry{}},
			caseExpected,
			http.StatusOK,
		},
		{
			"empty-name",
			&task.Entry{Name: "", Assignee: "Baris", Tags: []string{"red"}, Subs: map[string]*task.Entry{}},
			case2Expected,
			http.StatusBadRequest,
		},
	}

	var buf bytes.Buffer
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			json.NewEncoder(&buf).Encode(tc.entry)
			req, err := http.NewRequest(http.MethodPost, "/saveTask", &buf)
			if err != nil {
				t.Fatalf("could not create request: %v", err)
			}
			rec := httptest.NewRecorder()

			srv.handleSaveTask(rec, req)
			if rec.Code != tc.code {
				t.Fatalf("got status %d, want %v", rec.Code, tc.code)
			}

			if !reflect.DeepEqual(rec.Body.Bytes(), tc.want) {
				t.Fatalf("NAME:%v,  got %v, want %v", tc.name, rec.Body.Bytes(), tc.want)
			}
		})
	}
}

func TestSaveSubEntryHandler(t *testing.T) {
	srv := setRequiredServices()

	expected, _ := json.Marshal(&task.Entry{Id: 1, Code: "TEST-1", Name: "develop x feature", Assignee: "Baris", Tags: []string{"red"}, Subs: map[string]*task.Entry{
		"TEST-2": {
			Id:       2,
			Code:     "TEST-2",
			Name:     "develop y feature",
			Assignee: "Baris",
			Tags:     []string{"red"},
			Subs:     map[string]*task.Entry{},
		},
	}})

	err := map[string]string{
		"error": "parent code does not exist",
	}
	case2Expected, _ := json.Marshal(err)

	tt := []struct {
		name       string
		entry      *task.Entry
		want       []byte
		code       int
		parentCode map[string]string
	}{
		{
			"single",
			&task.Entry{Name: "develop y feature", Assignee: "Baris", Tags: []string{"red"}, Subs: map[string]*task.Entry{}},
			expected,
			http.StatusOK,
			map[string]string{"code": "TEST-1"},
		},
		{
			"no parent code",
			&task.Entry{Name: "develop y feature", Assignee: "Baris", Tags: []string{"red"}, Subs: map[string]*task.Entry{}},
			case2Expected,
			http.StatusBadRequest,
			map[string]string{"code": "TEST-10"},
		},
	}

	srv.TaskService.Save(&task.Entry{Name: "develop x feature", Assignee: "Baris", Tags: []string{"red"}, Subs: map[string]*task.Entry{}})
	var buf bytes.Buffer
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			json.NewEncoder(&buf).Encode(tc.entry)
			req, err := http.NewRequest(http.MethodPost, "/save/{code}", &buf)
			if err != nil {
				t.Fatalf("could not create request: %v", err)
			}
			rec := httptest.NewRecorder()
			req = mux.SetURLVars(req, tc.parentCode)

			srv.handleSaveSubTask(rec, req)
			if rec.Code != tc.code {
				t.Fatalf("got status %d, want %v", rec.Code, tc.code)
			}

			if !reflect.DeepEqual(rec.Body.Bytes(), tc.want) {
				t.Fatalf("NAME:%v,  got %v, want %v", tc.name, rec.Body.Bytes(), tc.want)
			}
		})
	}
}

func TestGetAllEntries(t *testing.T) {
	srv := setRequiredServices()

	exp, _ := json.Marshal(&task.TaskService{
		Entries: map[string]*task.Entry{
			"TEST-1": &task.Entry{
				Id:       1,
				Code:     "TEST-1",
				Name:     "develop x feature",
				Assignee: "Baris",
				Tags:     []string{"red"},
				Subs: map[string]*task.Entry{
					"TEST-2": {
						Id:       2,
						Code:     "TEST-2",
						Name:     "develop y feature",
						Assignee: "Ertas",
						Tags:     []string{"green"},
						Subs:     map[string]*task.Entry{},
					},
				},
			},
		},
	})

	tt := []struct {
		name string
		want []byte
		code int
	}{
		{
			"multiple entries",
			exp,
			http.StatusOK,
		},
	}

	srv.TaskService.Save(&task.Entry{Name: "develop x feature", Assignee: "Baris", Tags: []string{"red"}, Subs: map[string]*task.Entry{}})

	srv.TaskService.SaveSubEntry(&task.Entry{Name: "develop y feature", Assignee: "Ertas", Tags: []string{"green"}, Subs: map[string]*task.Entry{}}, "TEST-1")
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			req, err := http.NewRequest("GET", "/retrieve", nil)
			if err != nil {
				t.Fatalf("could not create request: %v", err)
			}
			rec := httptest.NewRecorder()
			srv.handleGetTasks(rec, req)

			if rec.Code != tc.code {
				t.Fatalf("got status %d, want %v", rec.Code, tc.code)
			}

			if !reflect.DeepEqual(rec.Body.Bytes(), tc.want) {
				t.Fatalf("NAME:%v,  got %v, want %v", tc.name, rec.Body.Bytes(), tc.want)
			}
		})
	}
}

func setRequiredServices() *Server {
	srv := NewServer()
	taskService := task.NewTaskService()
	srv.TaskService = *taskService

	return srv
}
