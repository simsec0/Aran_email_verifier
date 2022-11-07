package verifier

import (
	utils "EmailVerifier/Modules"
	"fmt"
	"io/ioutil"
	"strings"

	http "github.com/bogdanfinn/fhttp"
)

func (in *Instance) Session() error {
	req, err := http.NewRequest("GET", "https://discord.com/login", nil)

	if err != nil {
		return err
	}

	req = in.SetHeaders(true, map[string]string{}, req)

	_, err = in.Client.Do(req)

	if err != nil {
		return err
	}

	in.XSuperProps = BuildSuperProps(in.UserAgent, in.BrowserVersion)

	req, err = http.NewRequest("GET", "https://discord.com/api/v9/users/@me/library", nil)

	if err != nil {
		return err
	}

	req = in.SetHeaders(true, map[string]string{
		"authorization":      in.Token,
		"x-debug-options":    "bugReporterEnabled",
		"x-discord-locale":   "en-US",
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

	if resp.StatusCode == 403 && strings.Contains(string(bodyText), "You need to verify your account in order to perform this action.") {
		utils.RemoveFromFile("Data/tokens.txt", in.Token)
		return fmt.Errorf(fmt.Sprintf("Token is locked | %v", in.FullToken))

	} else if resp.StatusCode == 401 && strings.Contains(string(bodyText), "unauthorized") {
		utils.RemoveFromFile("Data/tokens.txt", in.Token)
		return fmt.Errorf(fmt.Sprintf("Token is invalid | %v", in.FullToken))

	} else if resp.StatusCode != 200 {
		return fmt.Errorf(fmt.Sprintf("could not get @me, status code: %v", resp.StatusCode))
	}

	return nil
}
