package authentication

import "encoding/json"
import "encoding/gob"
import "net/http"
import "io/ioutil"
import "fmt"
import "github.com/go-martini/martini"
import "github.com/martini-contrib/sessions"
import "golang.org/x/oauth2"
import "golang.org/x/oauth2/google"
import "golang.org/x/oauth2/facebook"
import "web/models/user"
import "web/helpers/configuration"

const authToken = "oauth2_token"
const authProfile = "oauth2_profile"
const redirectHttpCode = 302

// PathLogin is the path to handle OAuth 2.0 logins.
var urlLogin = "/login"

// PathLogout is the path to handle OAuth 2.0 logouts.
var urlLogout = "/logout"

// PathLoginOK is the path to redirect when login success.
var urlProfile = "/profile"

// PathLogoutOK is the path to redirect when logout success.
var urlExit = "/"

// Token struct.
type token struct {
	oauth2.Token
}

type oauthProfile struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Email   string `json:"email"`
	Gender  string `json:"gender"`
	Locale  string `json:"locale"`
	Profile string `json:"link"`
	Picture string `json:"picture"`
}

type facebookProfilePicture struct {
	ID      string `json:"id"`
	Picture struct {
		Data struct {
			IsSilhouette bool   `json:"is_silhouette"`
			URL          string `json:"url"`
		} `json:"data"`
	} `json:"picture"`
}

var LoginRequired = func() martini.Handler {
	return func(s sessions.Session, c martini.Context, w http.ResponseWriter, r *http.Request) {
		token := getToken(s)

		if token == nil || !token.Valid() {
			http.Redirect(w, r, urlLogin, redirectHttpCode)
		}
	}
}()

var Basic = func() martini.Handler {
	return func(s sessions.Session, c martini.Context, w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/logout" {
			s.Delete(authToken)
			s.Delete(authProfile)
			http.Redirect(w, r, urlExit, redirectHttpCode)
		}

		// Check token validity.
		tk := getToken(s)
		if tk != nil {

			if !tk.Valid() && tk.RefreshToken == "" {
				s.Delete(authToken)
				s.Delete(authProfile)
				tk = nil

				http.Redirect(w, r, urlLogin, redirectHttpCode)
			}
		}

		// Inject user profile.
		up := getProfile(s)
		c.MapTo(up, (*user.Profiler)(nil))
	}
}()

var Google = func() martini.Handler {
	return func(s sessions.Session, c martini.Context, w http.ResponseWriter, r *http.Request, app configuration.Application) {
		config := make(map[string]*oauth2.Config)
		config["google"] = &oauth2.Config{
			ClientID:     app.OAuth["google"].ClientId,
			ClientSecret: app.OAuth["google"].SecretId,
			Scopes:       []string{"openid email", "https://www.googleapis.com/auth/userinfo.profile"},
			RedirectURL:  app.Domain[martini.Env].Url + "/login/google/callback",
			Endpoint:     google.Endpoint,
		}

		if r.Method == "GET" {
			switch r.URL.Path {
			case "/login/google":
				http.Redirect(w, r, config["google"].AuthCodeURL("/"), redirectHttpCode)
			case "/login/google/callback":
				err := r.URL.Query().Get("error")

				if err != "" {
					s.AddFlash(err)
					http.Redirect(w, r, "/login", redirectHttpCode)
					return
				}

				code := r.URL.Query().Get("code")
				t, _ := config["google"].Exchange(oauth2.NoContext, code)
				client := config["google"].Client(oauth2.NoContext, t)

				// Get profile
				response, _ := client.Get("https://www.googleapis.com/oauth2/v1/userinfo")
				defer response.Body.Close()
				body, _ := ioutil.ReadAll(response.Body)

				var profile oauthProfile
				json.Unmarshal(body, &profile)

				// Store the credentials in the session.
				val, _ := json.Marshal(t)

				// Save token and profile to session.
				s.Set(authToken, val)
				s.Set(authProfile, profile)

				c.Map(profile)

				http.Redirect(w, r, "/profile", redirectHttpCode)
			}
		}
	}
}()

var Facebook = func() martini.Handler {
	return func(s sessions.Session, c martini.Context, w http.ResponseWriter, r *http.Request, app configuration.Application) {
		config := make(map[string]*oauth2.Config)
		config["facebook"] = &oauth2.Config{
			ClientID:     app.OAuth["facebook"].ClientId,
			ClientSecret: app.OAuth["facebook"].SecretId,
			Scopes:       []string{"email", "public_profile"},
			RedirectURL:  app.Domain[martini.Env].Url + "/login/facebook/callback",
			Endpoint:     facebook.Endpoint,
		}

		if r.Method == "GET" {
			switch r.URL.Path {
			case "/login/facebook":
				http.Redirect(w, r, config["facebook"].AuthCodeURL("/"), redirectHttpCode)
			case "/login/facebook/callback":
				err := r.URL.Query().Get("error")

				if err != "" {
					s.AddFlash(err)
					http.Redirect(w, r, "/login", redirectHttpCode)
					return
				}

				code := r.URL.Query().Get("code")
				t, _ := config["facebook"].Exchange(oauth2.NoContext, code)
				client := config["facebook"].Client(oauth2.NoContext, t)

				// Get profile
				accessToken := fmt.Sprintf("%s", t.Extra("access_token"))
				response, _ := client.Get("https://graph.facebook.com/me?access_token=" + accessToken)
				defer response.Body.Close()
				body, _ := ioutil.ReadAll(response.Body)

				var profile oauthProfile
				json.Unmarshal(body, &profile)

				// Get profile picture
				response, _ = client.Get("https://graph.facebook.com/" + profile.ID + "?fields=picture.type(large)")
				defer response.Body.Close()
				body, _ = ioutil.ReadAll(response.Body)

				var profilePicture facebookProfilePicture
				json.Unmarshal(body, &profilePicture)

				profile.Picture = profilePicture.Picture.Data.URL

				// Store the credentials in the session.
				val, _ := json.Marshal(t)

				// Save token and profile to session.
				s.Set(authToken, val)
				s.Set(authProfile, profile)

				c.Map(profile)

				http.Redirect(w, r, "/profile", redirectHttpCode)
			}
		}
	}
}()

func init() {
	// Register oauthProfile struct in session handler.
	var profile oauthProfile
	gob.Register(profile)
}

func getProfile(s sessions.Session) (t *user.Profile) {
	if s.Get(authProfile) == nil {
		return
	}
	// data contains oauth profile information.
	data := s.Get(authProfile).(oauthProfile)

	// todo. store/retrieve user obj. from database.

	return &user.Profile{1, data.Name, data.Email, data.Profile, data.Picture}
}

func getToken(s sessions.Session) (t *token) {
	if s.Get(authToken) == nil {
		return
	}
	data := s.Get(authToken).([]byte)
	var tk oauth2.Token
	json.Unmarshal(data, &tk)
	return &token{tk}
}
