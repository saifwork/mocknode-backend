package mongo

const (
	usersCol       = "users"
	collectionsCol = "collections"
	projectsCol    = "projects"
	recordsCol     = "records"
)

// Collections exposes read-only grouped names.
var Collections = struct {
	Users      string
	Collection string
	Projects   string
	Records    string
}{
	Users:      usersCol,
	Collection: collectionsCol,
	Projects:   projectsCol,
	Records:    recordsCol,
}
