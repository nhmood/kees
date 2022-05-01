package device

import (
	"bytes"
	"net/http"
	"os"
	"time"

	"github.com/Masterminds/log-go"

	"kees-client/constants"
	"kees-client/helpers"
)

func (c *Client) Authenticate() *AuthResponse {
	log.Info("Authenticating " + c.Device.Name)
	jsonData, err := helpers.Format(c.Device)
	if err != nil {
		log.Warn("Failed to format Device info")
		os.Exit(1)
	}

	request, err := http.NewRequest(
		http.MethodPost,
		c.getAuthURL(),
		bytes.NewBuffer(jsonData),
	)
	request.Header.Set("User-Agent", "kees-client/"+constants.Version)
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("X-Kees-MC-Token", c.Token)

	httpClient := http.Client{Timeout: time.Duration(5 * time.Second)}
	log.Info("Making request for: " + c.getAuthURL())
	resp, err := httpClient.Do(request)

	if err != nil {
		log.Warn("Failed to authenticate " + c.Device.Name)
		log.Error(err)
		os.Exit(1)
	}

	// TODO: add handling of non200 response
	authResp := AuthResponse{}
	helpers.Parse(resp, &authResp)
	helpers.Debug(authResp)

	log.Info("Authentication successful - DeviceID:" + authResp.Device.ID)

	return &authResp
}