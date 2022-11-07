package email

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func (in *Instance) GetBalance() (float64, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("https://api.hotmailbox.me/user/balance?apikey=%s", in.ApiKey), nil)

	if err != nil {
		return -1, err
	}

	resp, err := in.Client.Do(req)

	if err != nil {
		return -1, err
	}

	defer resp.Body.Close()
	bodyText, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return -1, err
	}

	var jsonBody struct {
		Code       int     `json:"Code"`
		Message    string  `json:"Message"`
		BalanceUsd float64 `json:"BalanceUsd"`
	}

	err = json.Unmarshal(bodyText, &jsonBody)

	if err != nil {
		return -1, err
	}

	if jsonBody.Code != 0 {
		return -1, fmt.Errorf(fmt.Sprintf("could not get balance, error code: %d", jsonBody.Code))
	}

	return jsonBody.BalanceUsd, nil
}
