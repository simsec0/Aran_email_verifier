package email

import (
	utils "EmailVerifier/Modules"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

func (in *Instance) GetEmail(domain string) (string, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("https://api.hotmailbox.me/mail/buy?apikey=%s&mailcode=HOTMAIL.TRUSTED&quantity=1", in.ApiKey, domain), nil)

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
		Code    int    `json:"Code"`
		Message string `json:"Message"`
		Data    struct {
			TotalAmountUsd float64 `json:"TotalAmountUsd"`
			Emails         []struct {
				Email    string `json:"Email"`
				Password string `json:"Password"`
			} `json:"Emails"`
		} `json:"Data"`
	}

	err = json.Unmarshal(bodyText, &jsonBody)

	if err != nil {
		return "", err
	}

	if jsonBody.Code != 0 {
		return "", fmt.Errorf(fmt.Sprintf("could not get email, error code: %d", jsonBody.Code))
	}

	emails := jsonBody.Data.Emails[0]

	utils.Logger(true, fmt.Sprintf("Retrieved %s account %s:%s | $%v", strings.ToLower(domain), emails.Email, emails.Password, jsonBody.Data.TotalAmountUsd))
	in.Email = emails.Email
	in.Password = emails.Password
	return fmt.Sprintf("%s|%s", emails.Email, emails.Password), nil
}
