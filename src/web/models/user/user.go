package user

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"math/rand"
	"time"
)

const firstLoginFormat = "02 Apr 2006"
const lastLoginFormat = "02 Apr 2006 15:04:01"
const createdAtFormat = "2006-01-02 15:04:01"

// UserProfile interface that will hold user struct.
type Profiler interface {
	GetFirstLogin() string
	GetLastLogin() string
	IsAdmin() bool
	Upsert() string
	Unread() int
}

// User profile struct.
type Profile struct {
	ID         int
	Name       string
	Email      string
	Profile    string
	Picture    string
	FirstLogin string
	LastLogin  time.Time
	CreatedAt  time.Time
	Admin      bool
}

type dbParams struct {
	Connector string
	Dns       string
}

var db *sql.DB
var err error

func GetProfiles(db *sql.DB) map[int]Profile {
	profiles := make(map[int]Profile)

	rows, err := db.Query("select id, name, email, profile, picture, created_at, is_admin from user order by is_admin desc, name asc")

	if err != nil {
		panic(err)
	}

	defer rows.Close()

	i := 0

	var id int
	var name string
	var email string
	var profile string
	var picture string
	var createdAt string
	var admin bool

	for rows.Next() {
		rows.Scan(&id, &name, &email, &profile, &picture, &createdAt, &admin)
		time.Parse(createdAtFormat, createdAt)

		p := Profile{ID: id, Name: name, Email: email, Profile: profile, Picture: picture, Admin: admin}
		firstLogin, _ := time.Parse(createdAtFormat, createdAt)
		p.FirstLogin = firstLogin.Format(firstLoginFormat)
		profiles[i] = p

		i++
	}

	return profiles
}

func (p *Profile) GetFirstLogin() string {
	return p.CreatedAt.Format(firstLoginFormat)
}

func (p *Profile) GetLastLogin() string {
	return p.LastLogin.Format(lastLoginFormat)
}

func (p *Profile) IsAdmin(db *sql.DB) bool {
	var isAdmin bool

	err = db.QueryRow("select is_admin from user where email like ?", p.Email).Scan(&isAdmin)

	return isAdmin
}

// Update saves user profile to database.
func (p *Profile) Upsert(db *sql.DB) {

	var result sql.Result
	var err error

	err = db.Ping()

	if err != nil {
		panic(err)
	}

	var id int
	var admin bool
	var createdAt string

	err = db.QueryRow("select id, created_at, is_admin from user where email like ?", p.Email).Scan(&id, &createdAt, &admin)

	switch {
	case err == sql.ErrNoRows:
		result, err = db.Exec("insert into user set name = ?, email = ?, profile = ?, picture = ?, created_at = ?", p.Name, p.Email, p.Profile, p.Picture, time.Now().Local().Format(createdAtFormat))
		if err != nil {
			panic(err)
		}
		lastInsertId, _ := result.LastInsertId()
		p.ID = int(lastInsertId)
	case err != nil:
		panic(err)
	default:
		p.ID = id
		p.CreatedAt, _ = time.Parse(createdAtFormat, createdAt)
		p.Admin = admin
	}

}

func (p *Profile) Unread() int {
	random := rand.Int()
	return random % 20
}
