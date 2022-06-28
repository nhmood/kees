package device

import (
	"bytes"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/Masterminds/log-go"

	"kees/controller/constants"
	"kees/controller/helpers"
)

type DeviceUpdate struct {
	Version      string   `json:"version"`
	Capabilities []string `json:"capabilities"`
}

func (c *Controller) getAuthURL() string {
	return c.baseURL("http") + "/ws/v1/auth"
}

func (c *Controller) getCapabilities() []string {
	capabilities := make([]string, 0)
	for capability, _ := range c.Device.Capabilities {
		capabilities = append(capabilities, capability)
	}

	return capabilities
}

func (c *Controller) Authenticate() *AuthResponse {
	deviceUpdate := DeviceUpdate{
		Version:      constants.Version,
		Capabilities: c.getCapabilities(),
	}

	jsonData, err := helpers.Format(deviceUpdate)
	if err != nil {
		log.Warn("Failed to format Device info")
		os.Exit(1)
	}

	request, err := http.NewRequest(
		http.MethodPost,
		c.getAuthURL(),
		bytes.NewBuffer(jsonData),
	)
	request.Header.Set("User-Agent", "kees/controller/"+constants.Version)
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

	if resp.StatusCode != http.StatusOK {
		log.Warn("Failed to authenticate " + c.Device.Name)

		defer resp.Body.Close()
		bs, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
			os.Exit(1)
		}
		log.Error(string(bs))
		os.Exit(1)
	}

	authResp := AuthResponse{}
	helpers.Parse(resp, &authResp)
	helpers.Debug(authResp)

	log.Info("Authentication successful - DeviceID:" + authResp.Device.ID)

	return &authResp
}
