package user

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"math/rand"
	"time"
)

const dateFormat = "02 Apr 2006"
const timeFormat = "02 Apr 2006 15:04:01"
const dbDateTime = "2006-01-02 15:04:01"

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
	ID                 int
	Name               string
	Email              string
	Profile            string
	Picture            string
	FirstLogin         time.Time
	LastLogin          time.Time
	FormatedFirstLogin string
	FormatedLastLogin  string
	Admin              bool
}

type dbParams struct {
	Connector string
	Dns       string
}

var db *sql.DB
var err error

func GetProfiles(db *sql.DB) map[int]Profile {
	profiles := make(map[int]Profile)

	rows, err := db.Query("select id, name, email, profile, picture, firstLogin, lastLogin, isAdmin from user order by isAdmin desc, name asc")

	if err != nil {
		panic(err)
	}

	defer rows.Close()

	i := 0

	for rows.Next() {

		var id int
		var name string
		var email string
		var profile string
		var picture string
		var firstLogin string
		var lastLogin string
		var admin bool

		rows.Scan(&id, &name, &email, &profile, &picture, &firstLogin, &lastLogin, &admin)

		p := Profile{ID: id, Name: name, Email: email, Profile: profile, Picture: picture, Admin: admin}

		firstLoginTime, err1 := time.Parse(dbDateTime, firstLogin)
		lastLoginTime, err2 := time.Parse(dbDateTime, lastLogin)

		if err1 == nil {
			p.FormatedFirstLogin = firstLoginTime.Format(dateFormat)
		}
		if err2 == nil {
			p.FormatedLastLogin = lastLoginTime.Format(dateFormat)
		}

		profiles[i] = p

		i++
	}

	return profiles
}

func (p *Profile) IsAdmin(db *sql.DB) bool {
	var isAdmin bool

	err = db.QueryRow("select isAdmin from user where email like ?", p.Email).Scan(&isAdmin)

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
	var firstLogin string

	err = db.QueryRow("select id, firstLogin, isAdmin from user where email like ?", p.Email).Scan(&id, &firstLogin, &admin)

	switch {
	case err == sql.ErrNoRows:

		firstLogin = time.Now().Local().Format(dbDateTime)

		result, err = db.Exec("insert into user set name = ?, email = ?, profile = ?, picture = ?, firstLogin = ?", p.Name, p.Email, p.Profile, p.Picture, firstLogin)
		if err != nil {
			panic(err)
		}

		firstLogin := time.Now().Local().Format(dbDateTime)
		firstLoginTime, _ := time.Parse(dbDateTime, firstLogin)
		p.FormatedFirstLogin = firstLoginTime.Format(dateFormat)
		p.FormatedLastLogin = "-"

		lastInsertId, _ := result.LastInsertId()
		p.ID = int(lastInsertId)
	case err != nil:
		panic(err)
	default:
		p.ID = id

		lastLogin := time.Now().Local().Format(dbDateTime)

		firstLoginTime, _ := time.Parse(dbDateTime, firstLogin)
		lastLoginTime, _ := time.Parse(dbDateTime, lastLogin)

		p.FormatedFirstLogin = firstLoginTime.Format(dateFormat)
		p.FormatedLastLogin = lastLoginTime.Format(dateFormat)

		db.Exec("update user set lastLogin = ? where id = ?", lastLoginTime.Format(dbDateTime), p.ID)

		p.Admin = admin
	}

}

func (p *Profile) Unread() int {
	random := rand.Int()
	return random % 20
}
