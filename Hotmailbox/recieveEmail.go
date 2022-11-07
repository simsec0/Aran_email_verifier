package email

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

func (in *Instance) GetVerificationEmail() (string, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("https://getcode.hotmailbox.me/discord?email=%s&password=%s&timeout=%d", in.Email, in.Password, in.Timeout), nil)

	if err != nil {
		return "", err
	}

	resp, err := in.Client.Do(req)

	if err != nil {
		return "", err
	}

	defer resp.Body.Close()
	bodyText, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return "", err
	}

	var jsonBody struct {
		Success          bool   `json:"Success,omitempty"`
		Message          string `json:"Message,omitempty"`
		VerificationCode string `json:"VerificationCode,omitempty"`
	}

	err = json.Unmarshal(bodyText, &jsonBody)

	if err != nil {
		return "", err
	}

	if !jsonBody.Success {
		return "", fmt.Errorf(fmt.Sprintf("could not get email verification link, error message: %s", jsonBody.Message))
	}

	return strings.TrimSuffix(jsonBody.VerificationCode, "\r"), nil
}
