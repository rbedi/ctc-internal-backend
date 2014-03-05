// Copyright 2013 Google Inc. All Rights Reserved
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may not
// use this file except in compliance with the License. You may obtain a copy of
// the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
// WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
// License for the specific language governing permissions and limitations under
// the License.

// Package main provides a simple server to demonstrate how to use Google+
// Sign-In and make a request via your own server.
package auth

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"

	"code.google.com/p/goauth2/oauth"
	_ "code.google.com/p/google-api-go-client/plus/v1"
	"github.com/gorilla/sessions"
	"github.com/gorilla/mux"
)

// Update your Google API project information here.

const (
	clientID        = "296723349350-3r0ruujesbqtps4bi6nr63ok9jn0lfc8.apps.googleusercontent.com"
	clientSecret    = "f8TnFLms31yfzL2wlg-TZmiw"
	applicationName = "CTC Projects Site"
)

// config is the configuration specification supplied to the OAuth package.
var config = &oauth.Config{
	ClientId:     clientID,
	ClientSecret: clientSecret,
	// Scope determines which API calls you are authorized to make
	Scope:    "https://www.googleapis.com/auth/plus.login",
	AuthURL:  "https://accounts.google.com/o/oauth2/auth",
	TokenURL: "https://accounts.google.com/o/oauth2/token",
	// Use "postmessage" for the code-flow for server side apps
	RedirectURL: "postmessage",
}

// store initializes the Gorilla session store.
var store = sessions.NewCookieStore([]byte(randomString(32)))

// indexTemplate is the HTML template we use to present the index page.
var indexTemplate = template.Must(template.ParseFiles("auth/example.html"))

// Token represents an OAuth token response.
type Token struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
	IdToken     string `json:"id_token"`
}

// ClaimSet represents an IdToken response.
type ClaimSet struct {
	Sub string
}

// exchange takes an authentication code and exchanges it with the OAuth
// endpoint for a Google API bearer token and a Google+ ID
func exchange(code string) (accessToken string, idToken string, err error) {
	log.Println("exchange called")
	// Exchange the authorization code for a credentials object via a POST request
	addr := "https://accounts.google.com/o/oauth2/token"
	values := url.Values{
		"Content-Type":  {"application/x-www-form-urlencoded"},
		"code":          {code},
		"client_id":     {clientID},
		"client_secret": {clientSecret},
		"redirect_uri":  {config.RedirectURL},
		"grant_type":    {"authorization_code"},
	}
	resp, err := http.PostForm(addr, values)
	if err != nil {
		return "", "", fmt.Errorf("Exchanging code: %v", err)
	}
	defer resp.Body.Close()

	// Decode the response body into a token object
	var token Token
	err = json.NewDecoder(resp.Body).Decode(&token)
	if err != nil {
		return "", "", fmt.Errorf("Decoding access token: %v", err)
	}

	return token.AccessToken, token.IdToken, nil
}

// decodeIdToken takes an ID Token and decodes it to fetch the Google+ ID within
func decodeIdToken (idToken string) (gplusID string, err error) {
	// An ID token is a cryptographically-signed JSON object encoded in base 64.
	// Normally, it is critical that you validate an ID token before you use it,
	// but since you are communicating directly with Google over an
	// intermediary-free HTTPS channel and using your Client Secret to
	// authenticate yourself to Google, you can be confident that the token you
	// receive really comes from Google and is valid. If your server passes the ID
	// token to other components of your app, it is extremely important that the
	// other components validate the token before using it.
	var set ClaimSet
	if idToken != "" {
		// Check that the padding is correct for a base64decode
		parts := strings.Split(idToken, ".")
		if len(parts) < 2 {
			return "", fmt.Errorf("Malformed ID token")
		}
		// Decode the ID token
		b, err := base64Decode(parts[1])
		if err != nil {
			return "", fmt.Errorf("Malformed ID token: %v", err)
		}
		err = json.Unmarshal(b, &set)
		if err != nil {
			return "", fmt.Errorf("Malformed ID token: %v", err)
		}
	}
	return set.Sub, nil
}

// index sets up a session for the current user and serves the index page
func index(w http.ResponseWriter, r *http.Request) *appError {
	log.Println("index called")
	// This check prevents the "/" handler from handling all requests by default
	/*if r.URL.Path != "/" {
		log.Println("bar")
		http.NotFound(w, r)
		return nil
	}*/

	// Create a state token to prevent request forgery and store it in the session
	// for later validation
	session, err := store.Get(r, "sessionName")
	if err != nil {
		log.Println("error fetching session:", err)
		// Ignore the initial session fetch error, as Get() always returns a
		// session, even if empty.
		//return &appError{err, "Error fetching session", 500}
	}
	state := randomString(64)
	session.Values["state"] = state
	session.Save(r, w)

	stateURL := url.QueryEscape(session.Values["state"].(string))

	// Fill in the missing fields in index.html
	var data = struct {
		ApplicationName, ClientID, State string
	}{applicationName, clientID, stateURL}

	// Render and serve the HTML
	err = indexTemplate.Execute(w, data)
	if err != nil {
		log.Println("error rendering template:", err)
		return &appError{err, "Error rendering template", 500}
	}
	return nil
}

// connect exchanges the one-time authorization code for a token and stores the
// token in the session
func connect(w http.ResponseWriter, r *http.Request) *appError {
	// Ensure that the request is not a forgery and that the user sending this
	// connect request is the expected user
	log.Println("connect called")
	session, err := store.Get(r, "sessionName")
	if err != nil {
		log.Println("error fetching session:", err)
		return &appError{err, "Error fetching session", 500}
	}
	if r.FormValue("state") != session.Values["state"].(string) {
		m := "Invalid state parameter"
		return &appError{errors.New(m), m, 401}
	}
	// Normally, the state is a one-time token; however, in this example, we want
	// the user to be able to connect and disconnect without reloading the page.
	// Thus, for demonstration, we don't implement this best practice.
	// session.Values["state"] = nil

	// Setup for fetching the code from the request payload
	x, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return &appError{err, "Error reading code in request body", 500}
	}
	code := string(x)

	accessToken, idToken, err := exchange(code)
	if err != nil {
	        return &appError{err, "Error exchanging code for access token", 500}
	}
        gplusID, err := decodeIdToken(idToken)
	if err != nil {
	        return &appError{err, "Error decoding ID token", 500}
	}

	// Check if the user is already connected
	storedToken := session.Values["accessToken"]
	storedGPlusID := session.Values["gplusID"]
	if storedToken != nil && storedGPlusID == gplusID {
		log.Println(storedToken)
		m := "Current user already connected"
		fmt.Fprintf(w, "{logged_in: true}")
		return &appError{errors.New(m), m, 200}
	}

	// Store the access token in the session for later use
	session.Values["accessToken"] = accessToken
	session.Values["gplusID"] = gplusID
	session.Save(r, w)
	fmt.Fprintf(w, "{logged_in: true}")
	return nil
}

// disconnect revokes the current user's token and resets their session
func disconnect(w http.ResponseWriter, r *http.Request) *appError {
	log.Println("disconnect called")
	// Only disconnect a connected user
	session, err := store.Get(r, "sessionName")
	if err != nil {
		log.Println("error fetching session:", err)
		return &appError{err, "Error fetching session", 500}
	}
	token := session.Values["accessToken"]
	if token == nil {
		m := "Current user not connected"
		return &appError{errors.New(m), m, 401}
	}

	// Execute HTTP GET request to revoke current token
	url := "https://accounts.google.com/o/oauth2/revoke?token=" + token.(string)
	resp, err := http.Get(url)
	if err != nil {
		m := "Failed to revoke token for a given user"
		return &appError{errors.New(m), m, 400}
	}
	defer resp.Body.Close()

	// Reset the user's session
	session.Values["accessToken"] = nil
	session.Save(r, w)
	return nil
}

// appHandler is to be used in error handling
type AppHandler func(http.ResponseWriter, *http.Request) *appError

type appError struct {
	Err     error
	Message string
	Code    int
}

// serveHTTP formats and passes up an error
func (fn AppHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if e := fn(w, r); e != nil { // e is *appError, not os.Error.
		log.Println(e.Err)
		http.Error(w, e.Message, e.Code)
	}
}

// randomString returns a random string with the specified length
func randomString(length int) (str string) {
	b := make([]byte, length)
	rand.Read(b)
	return base64.StdEncoding.EncodeToString(b)
}

func base64Decode(s string) ([]byte, error) {
	// add back missing padding
	switch len(s) % 4 {
	case 2:
		s += "=="
	case 3:
		s += "="
	}
	return base64.URLEncoding.DecodeString(s)
}

func Register(r *mux.Router) {
	r = r.StrictSlash(true) // allow for optional trailing slashes
	r.Handle("/example", AppHandler(index))
	r.Handle("/login", AppHandler(connect))
	r.Handle("/logout", AppHandler(disconnect))
}