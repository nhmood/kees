package api

import (
	"net/http"

	"kees/server/helpers"
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
	User User        `json:"user"`
	JWT  JWTResponse `json:"jwt"`
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
		helpers.Halt(w, http.StatusBadRequest, "Unauthorized Client", nil)
		return
	}

	user := User{
		ID:       "kees",
		Username: credentials.Username,
		Client:   credentials.Client,
	}

	// TODO: replace id with proper db user id
	jwt, expiresIn, err := helpers.GenerateJWT(map[string]interface{}{
		"id":       user.ID,
		"username": user.Username,
		"client":   user.Client,
	})

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	helpers.Debug(jwt)

	jwtPayload := ClientAuthResponse{
		User: user,
		JWT: JWTResponse{
			ExpiresIn: expiresIn,
			Token:     jwt,
		},
	}

	data := helpers.ToInterface(jwtPayload)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	helpers.Halt(w, http.StatusOK, "successfully authd "+user.Username, data)
	return
}
