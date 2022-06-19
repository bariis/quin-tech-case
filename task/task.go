package task

import (
	"errors"
	"strconv"
)

var PROJECT_NAME string = "TEST"

// TODO stands for entries whose types are TODO
type TODO []*Entry

// InProgress stands for entries whose types are InProgress
type InProgress []*Entry

// TaskService represents the all related entities and fields related to tasks
type TaskService struct {
	Entries    map[string]*Entry `json:"entries"`
	LastId     int               `json:"-"`
	LastCode   int               `json:"-"`
	TODO       TODO              `json:"-"`
	InProgress InProgress        `json:"-"`
}

// Entry represents the "entry" entity
type Entry struct {
	Id       int      `json:"id"`
	Name     string   `json:"name"`
	Code     string   `json:"code"`
	Type     string   `json:"type"`
	Assignee string   `json:"assignee"`
	Tags     []string `json:"tags"`
	Subs     map[string]*Entry
	DueDate      string `json:"due_date"`
	CreationDate string `json:"creation_date"`
	UpdateDate   string `json:"update_date"`
}

// NewTaskService returns a new TaskService
func NewTaskService() *TaskService {
	return &TaskService{
		Entries:  make(map[string]*Entry),
		LastId:   0,
		LastCode: 0,
	}
}

// Save saves the given single-parent Entry to the map
func (t *TaskService) Save(entry *Entry) (*Entry, error) {
	if entry.Name == "" {
		return nil, errors.New("empty name for the entry")
	}
	entry.Code = t.getUniqueCode()
	entry.Id = t.getUniqueId()
	// entry.CreationDate = getCurrentTime()
	// entry.UpdateDate = getCurrentTime()
	entry.Subs = make(map[string]*Entry)
	t.Entries[entry.Code] = entry

	if entry.Type == "TODO" {
		t.TODO = append(t.TODO, entry)
	} else {
		t.InProgress = append(t.InProgress, entry)
	}

	return entry, nil
}

// SaveSubEntry saves the given entry within a parent entry based on given parentCode of a parent entry.
// If parentCode exists returns updated parent-entry, returns error otherwise.
func (t *TaskService) SaveSubEntry(entry *Entry, parentCode string) (*Entry, error) {
	if entry.Name == "" {
		return nil, errors.New("empty name for the entry")
	}
	if val, ok := t.Entries[parentCode]; ok {
		entry.Code = t.getUniqueCode()
		entry.Id = t.getUniqueId()
		// entry.CreationDate = getCurrentTime()
		// entry.UpdateDate = getCurrentTime()

		val.Subs[entry.Code] = entry
		return val, nil
	}
	return nil, errors.New("parent code does not exist")
}

// All returns all saved entries
func (t *TaskService) All() (*TaskService, error) {

	// entries := make([]time.Time, len(t.Entries))

	// for _, entry := range t.Entries {
	// 	entries = append(entries, entry.CreationDate)
	// }

	// sort.Sort(TimeSlice(entries))

	// kk := TaskService{}

	// for  _, value := range entries {
	// 	entry := t.
	// 	kk.Entries[""]
	// }

	// -----------------------
	// keys := make([]string, 0, len(t.Entries))

	// for key := range t.Entries {
	// 	keys = append(keys, key)
	// }

	// sort.SliceStable(keys, func(i, j int) bool {
	// 	return t.Entries[keys[i]].DueDate < t.Entries[keys[j]].DueDate
	// })

	return t, nil
}

// Update updates the specific entry whether be sub or parent entry based on entry-code
func (t *TaskService) Update(updatedEntry *Entry) (*Entry, error) {
	entryCode := updatedEntry.Code
	for key, entry := range t.Entries {
		if key == entryCode {
			// get the sub-entries if exist then save them again
			if len(t.Entries[entryCode].Subs) > 0 {
				var subEntries map[string]*Entry
				subEntries = t.Entries[entryCode].Subs
				t.Entries[entryCode] = updatedEntry
				// t.Entries[entryCode].UpdateDate = getCurrentTime()
				t.Entries[entryCode].Subs = subEntries
			} else {
				t.Entries[entryCode] = updatedEntry
			}
			return t.Entries[entryCode], nil
		}
		if _, ok := entry.Subs[entryCode]; ok {
			entry.Subs[entryCode] = updatedEntry
			// t.Entries[entryCode].UpdateDate = getCurrentTime()
			return entry.Subs[entryCode], nil
		}
	}
	return &Entry{}, errors.New("no such entry")
}

// Delete deletes the specific entry based on its entry code.
// It returns an error if entry-code does not exist.
func (t *TaskService) Delete(entryCode string) error {
	for key, entry := range t.Entries {
		if key == entryCode {
			delete(t.Entries, key)
			return nil
		}
		if len(entry.Subs) > 0 {
			if _, ok := entry.Subs[entryCode]; ok {
				delete(entry.Subs, entryCode)
				return nil
			}
		}
	}
	return errors.New("no such entry-code")
}

// GetUniqueId returns an unique id after incrementing based on last used value
func (t *TaskService) getUniqueId() int {
	t.LastId += 1
	return t.LastId
}

// GetUniqueCode returns an unique code after incrementing based on last used value
// TEST stands for default code. e.g. TEST-.
// e.g. last_one => TEST-1 -> TEST-2.
func (t *TaskService) getUniqueCode() string {
	t.LastCode++
	return PROJECT_NAME + "-" + strconv.Itoa(t.LastCode)
}

// addToList adds entry to the specified list based on its type.
func (t *TaskService) addToList(entry *Entry) {
	if entry.Type == "TODO" {
		t.TODO = append(t.TODO, entry)
	} else {
		t.InProgress = append(t.InProgress, entry)
	}
}
