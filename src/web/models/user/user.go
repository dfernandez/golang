package user

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"log"
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

func initDb(dbConfig interface{}) {

	config := dbConfig.(*struct {
		Connector string
		Dns       string
	})

	db, err = sql.Open(config.Connector, config.Dns)

	if err != nil {
		panic(err)
	}
}

// Update saves user profile to database.
func (p *Profile) Upsert(dbConfig interface{}) {
	initDb(dbConfig)
}

func (p *Profile) Unread() int {
	random := rand.Int()
	return random % 20
}
