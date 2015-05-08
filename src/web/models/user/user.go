package user

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/stathat/go"
	"math/rand"
	"os"
	"time"
	"web/helpers/mynegroni"
)

const dateFormat = "02 Jan 2006"
const timeFormat = "02 Jan 2006 15:04:01"
const dbDateTime = "2006-01-02 15:04:01"

// User profile struct.
type Profile struct {
	ID                 int
	Name               string
	Email              string
	Gender             string
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
var config *mynegroni.Config

func init() {

	config = mynegroni.LoadConfig()

	dbConfig := config.Database[os.Getenv("ENV")]

	db, err = sql.Open(dbConfig.Connector, dbConfig.Dns)

	if err != nil {
		panic(err)
	}
}

func GetProfiles() map[int]Profile {
	profiles := make(map[int]Profile)

	rows, err := db.Query("select id, name, email, gender, profile, picture, firstLogin, lastLogin, isAdmin from user order by isAdmin desc, name asc")

	if err != nil {
		panic(err)
	}

	defer rows.Close()

	i := 0

	for rows.Next() {

		var id int
		var name string
		var email string
		var gender string
		var profile string
		var picture string
		var firstLogin string
		var lastLogin string
		var admin bool

		rows.Scan(&id, &name, &email, &gender, &profile, &picture, &firstLogin, &lastLogin, &admin)

		p := Profile{ID: id, Name: name, Email: email, Gender: gender, Profile: profile, Picture: picture, Admin: admin}

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

func (p *Profile) IsAdmin() bool {
	var isAdmin bool

	err = db.QueryRow("select isAdmin from user where email like ?", p.Email).Scan(&isAdmin)

	return isAdmin
}

// Update saves user profile to database.
func (p *Profile) Upsert() {

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

		result, err = db.Exec("insert into user set name = ?, email = ?, gender = ?, profile = ?, picture = ?, firstLogin = ?", p.Name, p.Email, p.Gender, p.Profile, p.Picture, firstLogin)
		if err != nil {
			panic(err)
		}

		firstLogin := time.Now().Local().Format(dbDateTime)
		firstLoginTime, _ := time.Parse(dbDateTime, firstLogin)
		p.FormatedFirstLogin = firstLoginTime.Format(dateFormat)
		p.FormatedLastLogin = "-"

		lastInsertId, _ := result.LastInsertId()
		p.ID = int(lastInsertId)

		if os.Getenv("ENV") == "production" {
			stathat.PostEZCountOne("users - new user", config.Stathat.Account)
		}

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
