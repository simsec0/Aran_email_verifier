package verifier

import (
	"fmt"
	"strings"

	http "github.com/bogdanfinn/fhttp"
)

func (in *Instance) VerifyEmail(link string) error {
	req, err := http.NewRequest("GET", link, nil)

	if err != nil {
		return err
	}

	resp, err := in.Client.Do(req)

	if err != nil {
		return err
	}

	if resp.StatusCode != 302 && !Ok(resp.StatusCode) {
		return fmt.Errorf("could not verify email, status code: %d", resp.StatusCode)
	}

	location := resp.Header.Get("location")

	emailToken := strings.Split(location, "token=")[1]

	var data = strings.NewReader(fmt.Sprintf(`{"token":"%s","captcha_key":null}`, emailToken))
	req, err = http.NewRequest("POST", "https://discord.com/api/v9/auth/verify", data)

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

	resp, err = in.Client.Do(req)

	if err != nil {
		return err
	}

	if !Ok(resp.StatusCode) {
		return fmt.Errorf("could not verify email, status code: %d", resp.StatusCode)
	}

	return nil
}
