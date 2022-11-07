package verifier

import (
	"fmt"
	"io/ioutil"
	"strings"

	http "github.com/bogdanfinn/fhttp"
)

func (in *Instance) AddEmail() error {
	var data = strings.NewReader(fmt.Sprintf(`{"email":"%s","email_token":null,"password":"%s"}`, in.NewEmail, in.OldPassword))

	req, err := http.NewRequest("PATCH", "https://discord.com/api/v9/users/@me", data)

	if err != nil {
		return err
	}

	req = in.SetHeaders(true, map[string]string{
		"authorization":      in.Token,
		"x-debug-options":    "bugReporterEnabled",
		"x-discord-locale":   "en-US",
		"content-type":       "application/json",
		"x-super-properties": in.XSuperProps,
	}, req)

	resp, err := in.Client.Do(req)

	if err != nil {
		return err
	}

	defer resp.Body.Close()
	bodyText, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return err
	}

	if !Ok(resp.StatusCode) {
		if strings.Contains(string(bodyText), "EMAIL_CHANGE_UPGRADE_CLIENT") {
			return fmt.Errorf("Token is already verified | " + in.Token)
		}

		return fmt.Errorf("Could not add new email, status code: %d", resp.StatusCode)
	}

	return nil
}
