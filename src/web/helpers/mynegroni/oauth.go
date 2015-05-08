package mynegroni

import (
	"encoding/gob"
	"encoding/json"
	"fmt"
	"github.com/codegangsta/negroni"
	sessions "github.com/goincremental/negroni-sessions"
	"github.com/gorilla/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/facebook"
	"golang.org/x/oauth2/google"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

const authToken = "oauth_token"
const oauthProfile = "oauth_profile"
const userProfile = "profile"

var config *Config

// Token struct.
type token struct {
	oauth2.Token
}

type OauthProfile struct {
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

var BasicOAuth = func() negroni.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		s := sessions.GetSession(r)

		if r.Method == "GET" {
			if r.URL.Path == "/logout" {
				s.Delete(authToken)
				s.Delete(userProfile)
				http.Redirect(rw, r, "/", http.StatusFound)
				return
			}
		}

		// Check token validity.
		tk := GetToken(r)
		if tk != nil {
			if !tk.Valid() && tk.RefreshToken == "" {

				s.Delete(authToken)
				s.Delete(userProfile)
				tk = nil

				http.Redirect(rw, r, "/login", http.StatusFound)
				return
			}
		}

		next(rw, r)
	}
}()

var GoogleOAuth = func() negroni.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {

		oauth := make(map[string]*oauth2.Config)
		oauth["google"] = &oauth2.Config{
			ClientID:     config.OAuth["google"].ClientId,
			ClientSecret: config.OAuth["google"].SecretId,
			Scopes:       []string{"openid email", "https://www.googleapis.com/auth/userinfo.profile"},
			RedirectURL:  config.Domain[os.Getenv("ENV")].Url + "/login/google/callback",
			Endpoint:     google.Endpoint,
		}

		if r.Method == "GET" {
			switch r.URL.Path {

			case "/login/google":
				http.Redirect(rw, r, oauth["google"].AuthCodeURL("/"), http.StatusFound)

			case "/login/google/callback":
				err := r.URL.Query().Get("error")
				s := sessions.GetSession(r)

				if err != "" {
					log.Printf("OAUTH ERROR: %s", err)
					s.AddFlash("Authentication error", "oauth")
					http.Redirect(rw, r, "/login", http.StatusFound)
					return
				}

				code := r.URL.Query().Get("code")
				t, oerr := oauth["google"].Exchange(oauth2.NoContext, code)

				if oerr != nil {
					log.Printf("OAUTH ERROR: %s", oerr)
					s.AddFlash("Authentication error", "oauth")
					http.Redirect(rw, r, "/login", http.StatusFound)
					return
				}

				client := oauth["google"].Client(oauth2.NoContext, t)

				// Get profile
				response, _ := client.Get("https://www.googleapis.com/oauth2/v1/userinfo")
				defer response.Body.Close()
				body, _ := ioutil.ReadAll(response.Body)

				var profile OauthProfile
				json.Unmarshal(body, &profile)

				// Store the credentials in the session.
				val, _ := json.Marshal(t)

				// Save session token.
				s.Set(authToken, val)

				// Save oauth profile to context.
				context.Set(r, oauthProfile, profile)

				if os.Getenv("ENV") == "production" {
					stathat.PostEZCountOne("users - oauth google", config.Stathat.Account)
				}

				next(rw, r)
			default:
				next(rw, r)
			}
		} else {
			next(rw, r)
		}
	}
}()

var FacebookOAuth = func() negroni.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		oauth := make(map[string]*oauth2.Config)
		oauth["facebook"] = &oauth2.Config{
			ClientID:     config.OAuth["facebook"].ClientId,
			ClientSecret: config.OAuth["facebook"].SecretId,
			Scopes:       []string{"email", "public_profile"},
			RedirectURL:  config.Domain[os.Getenv("ENV")].Url + "/login/facebook/callback",
			Endpoint:     facebook.Endpoint,
		}

		if r.Method == "GET" {
			switch r.URL.Path {

			case "/login/facebook":
				http.Redirect(rw, r, oauth["facebook"].AuthCodeURL("/"), http.StatusFound)

			case "/login/facebook/callback":
				err := r.URL.Query().Get("error")
				s := sessions.GetSession(r)

				if err != "" {
					log.Printf("OAUTH ERROR: %s", err)
					s.AddFlash("Authentication error", "oauth")
					http.Redirect(rw, r, "/login", http.StatusFound)
					return
				}

				code := r.URL.Query().Get("code")
				t, oerr := oauth["facebook"].Exchange(oauth2.NoContext, code)

				if oerr != nil {
					log.Printf("OAUTH ERROR: %s", oerr)
					s.AddFlash("Authentication error", "oauth")
					http.Redirect(rw, r, "/login", http.StatusFound)
					return
				}

				client := oauth["facebook"].Client(oauth2.NoContext, t)

				// Get profile
				accessToken := fmt.Sprintf("%s", t.Extra("access_token"))
				response, _ := client.Get("https://graph.facebook.com/me?access_token=" + accessToken)
				defer response.Body.Close()
				body, _ := ioutil.ReadAll(response.Body)

				var profile OauthProfile
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

				// Save session token.
				s.Set(authToken, val)

				// Save oauth profile to context.
				context.Set(r, oauthProfile, profile)

				if os.Getenv("ENV") == "production" {
					stathat.PostEZCountOne("users - oauth facebook", config.Stathat.Account)
				}

				next(rw, r)
			default:
				next(rw, r)
			}
		} else {
			next(rw, r)
		}
	}
}()

func init() {
	var oauth_profile OauthProfile
	gob.Register(oauth_profile)

	config = LoadConfig()
}

func GetToken(r *http.Request) (t *token) {
	s := sessions.GetSession(r)

	if s.Get(authToken) == nil {
		return
	}
	data := s.Get(authToken).([]byte)
	var tk oauth2.Token
	json.Unmarshal(data, &tk)

	return &token{tk}
}
