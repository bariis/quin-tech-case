package task

import (
	"reflect"
	"testing"
)

func TestSaveParentEntry(t *testing.T) {
	tt := []struct {
		name  string
		entry *Entry
		want  *Entry
	}{
		{
			"single",
			&Entry{Name: "develop x feature", Assignee: "Baris", Tags: []string{"red"}, Subs: map[string]*Entry{}},
			&Entry{Id: 1, Code: "TEST-1", Name: "develop x feature", Assignee: "Baris", Tags: []string{"red"}, Subs: map[string]*Entry{}},
		},
		{
			"empty-name",
			&Entry{Name: "", Assignee: "Baris", Tags: []string{"red"}, Subs: map[string]*Entry{}},
			nil,
		},
	}

	ts := TaskService{
		Entries: make(map[string]*Entry),
	}

	for _, tc := range tt {
		got, _ := ts.Save(tc.entry)
		if !reflect.DeepEqual(got, tc.want) {
			t.Errorf("got %v, want %v", got, tc.want)
		}
	}
}

func TestSaveSubEntry(t *testing.T) {

	t.Run("existing-parent", func(t *testing.T) {
		ts := TaskService{
			Entries: make(map[string]*Entry),
		}
		tt := struct {
			entry      *Entry
			parentCode string
			want       *Entry
		}{
			&Entry{Name: "develop y feature", Assignee: "Baris", Tags: []string{"green"}, Subs: map[string]*Entry{}},
			"TEST-1",
			&Entry{Id: 1, Code: "TEST-1", Name: "develop x feature", Assignee: "Baris", Tags: []string{"red"}, Subs: map[string]*Entry{
				"TEST-2": {
					Id:       2,
					Code:     "TEST-2",
					Name:     "develop y feature",
					Assignee: "Baris",
					Tags:     []string{"green"},
					Subs:     map[string]*Entry{},
				},
			}},
		}

		ts.Save(&Entry{Name: "develop x feature", Assignee: "Baris", Tags: []string{"red"}, Subs: map[string]*Entry{}})
		got, _ := ts.SaveSubEntry(tt.entry, tt.parentCode)
		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("got %v, want %v", got, tt.want)
		}
	})

	t.Run("non-existing-parent", func(t *testing.T) {
		tt := struct {
			entry      *Entry
			parentCode string
			want       *Entry
		}{
			&Entry{Name: "develop y feature", Assignee: "Baris", Tags: []string{"green"}, Subs: map[string]*Entry{}},
			"TEST-5",
			nil,
		}
		ts := TaskService{
			Entries: make(map[string]*Entry),
		}
		got, _ := ts.SaveSubEntry(tt.entry, tt.parentCode)
		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("got %v, want %v", got, tt.want)
		}
	})

}
