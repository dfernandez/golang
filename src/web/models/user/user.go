package user

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"math/rand"
)

// UserProfile interface that will hold user struct.
type Profiler interface {
	Upsert() string
	Unread() int
}

// User profile struct.
type Profile struct {
	ID        int
	Name      string
	Email     string
	Profile   string
	Picture   string
	LastLogin string
}

type dbParams struct {
	Connector string
	Dns       string
}

var db *sql.DB
var err error

// Update saves user profile to database.
func (p *Profile) Upsert(db *sql.DB) {

}

func (p *Profile) Unread() int {
	random := rand.Int()
	return random % 20
}
