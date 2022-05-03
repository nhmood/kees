package api

import (
	"net/http"

	"kees/server/helpers"
	"kees/server/web/responses"
)

type ClientCredentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Client   string `json:"client"`
}

type User struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Client   string `json:"client"`
}

type ClientAuthResponse struct {
	Message string      `json:"message"`
	User    User        `json:"user"`
	JWT     JWTResponse `json:"jwt"`
}

// TODO: move to common messages
type JWTResponse struct {
	ExpiresIn int64  `json:"expires_in"`
	Token     string `json:"token"`
}

func ClientAuthV1(w http.ResponseWriter, r *http.Request) {
	credentials := ClientCredentials{}
	err := helpers.Parse(r, &credentials)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	helpers.Debug(credentials)

	// TODO: make this a database lookup for username/password
	if credentials.Username != "kees" && credentials.Password != "cdplayer" {
		data, err := helpers.Format(responses.Generic{
			Message: "Unauthorized Client",
			Data:    map[string]interface{}{},
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusBadRequest)
		w.Write(data)
		return
	}

	user := User{
		ID:       "kees",
		Username: credentials.Username,
		Client:   credentials.Client,
	}

	// TODO: replace id with proper db user id
	jwt, expiresIn, err := helpers.GenerateJWT(map[string]string{
		"id":       user.ID,
		"username": user.Username,
		"client":   user.Client,
	})

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	helpers.Debug(jwt)

	jwtResponse := ClientAuthResponse{
		Message: "successfully authd " + user.Username,
		User:    user,
		JWT: JWTResponse{
			ExpiresIn: expiresIn,
			Token:     jwt,
		},
	}

	data, err := helpers.Format(jwtResponse)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(data)
	return
}
