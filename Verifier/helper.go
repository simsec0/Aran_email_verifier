package verifier

import (
	"encoding/base64"
	"fmt"
	"io"
	"regexp"
	"strconv"
	"strings"
	"time"

	shttp "net/http"

	http "github.com/bogdanfinn/fhttp"
)

func (in *Instance) SetHeaders(useCommonHeaders bool, headers map[string]string, req *http.Request) *http.Request {
	if useCommonHeaders {
		for k, v := range map[string]string{
			"accept":                    "*/*",
			"accept-language":           "en-US,en;q=0.9",
			"dnt":                       "1",
			"origin":                    "https://discord.com",
			"referer":                   "https://discord.com/channels/@me",
			"sec-ch-ua":                 `\"Google Chrome\";v=\"105\", \"Not)A;Brand\";v=\"8\", \"Chromium\";v=\"105\"`,
			"sec-ch-ua-mobile":          "?0",
			"sec-ch-ua-platform":        `"Windows"`,
			"sec-fetch-dest":            "document",
			"sec-fetch-mode":            "navigate",
			"sec-fetch-site":            "none",
			"sec-fetch-user":            "?1",
			"upgrade-insecure-requests": "1",
			"user-agent":                in.UserAgent,
		} {
			req.Header.Set(k, v)
		}
	}

	for k, v := range headers {
		req.Header.Set(k, v)
	}

	return req
}

func (in *Instance) EmailHeaders(useCommonHeaders bool, headers map[string]string, req *http.Request) *http.Request {
	if useCommonHeaders {
		for k, v := range map[string]string{
			"authority":                 "emailtemp.org",
			"accept":                    "*/*",
			"accept-language":           "en-CA,en;q=0.9",
			"dnt":                       "1",
			"sec-ch-ua":                 `"Chromium";v="106", "Google Chrome";v="106", "Not;A=Brand";v="99"`,
			"sec-ch-ua-mobile":          "?0",
			"sec-ch-ua-platform":        `"Windows"`,
			"sec-fetch-dest":            "document",
			"sec-fetch-mode":            "navigate",
			"sec-fetch-site":            "none",
			"sec-fetch-user":            "?1",
			"upgrade-insecure-requests": "1",
			"user-agent":                in.UserAgent,
		} {
			req.Header.Set(k, v)
		}
	}

	for k, v := range headers {
		req.Header.Set(k, v)
	}

	return req
}

var (
	buildNumber string
)

func UpdateDiscordBuildInfo() error {
	defer func() {
		if r := recover(); r != nil {
			buildNumber = "149345"
		}
	}()

	jsFileRegex := regexp.MustCompile(`([a-zA-z0-9]+)\.js`)
	req, err := shttp.NewRequest("GET", "https://discord.com/app", nil)
	if err != nil {
		return err
	}
	req.Header = shttp.Header{
		"accept":             {`application/json, text/plain, */*`},
		"accept-language":    {`en-US,en;q=0.9`},
		"cache-control":      {`no-cache`},
		"origin":             {`https://discord.com`},
		"pragma":             {`no-cache`},
		"referer":            {`https://discord.com/`},
		"sec-ch-ua":          {`".Not/A)Brand";v="99", "Google Chrome";v="103", "Chromium";v="103"`},
		"sec-ch-ua-mobile":   {`?0`},
		"sec-ch-ua-platform": {`"Windows"`},
		"sec-fetch-dest":     {`empty`},
		"sec-fetch-mode":     {`cors`},
		"sec-fetch-site":     {`same-site`},
		"user-agent":         {`Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/103.0.0.0 Safari/537.36`},
	}
	client := &shttp.Client{
		Timeout:   10 * time.Second,
		Transport: shttp.DefaultTransport,
	}
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}
	index := 0
	if strings.Contains(string(body), "/alpha/invisible.js") {
		index = 2
	} else {
		index = 1
	}

	r := jsFileRegex.FindAllString(string(body), -1)
	asset := r[len(r)-index]
	resp, err := client.Get("https://discord.com/assets/" + asset)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	buildNumber = strings.Split(strings.Split(string(b), `build_number:"`)[1], `"`)[0]

	return nil
}

func GetDiscordBuildNumber() string {
	return buildNumber
}

func BuildSuperProps(userAgent string, browserVersion string) string {
	err := UpdateDiscordBuildInfo()

	if err != nil {
		return "eyJvcyI6IldpbmRvd3MiLCJicm93c2VyIjoiQ2hyb21lIiwiZGV2aWNlIjoiIiwic3lzdGVtX2xvY2FsZSI6ImVuLVVTIiwiYnJvd3Nlcl91c2VyX2FnZW50IjoiTW96aWxsYS81LjAgKFdpbmRvd3MgTlQgMTAuMDsgV09XNjQpIEFwcGxlV2ViS2l0LzUzNy4zNiAoS0hUTUwsIGxpa2UgR2Vja28pIENocm9tZS8xMDUuMC4wLjAgU2FmYXJpLzUzNy4zNiIsImJyb3dzZXJfdmVyc2lvbiI6IjEwNS4wLjAuMCIsIm9zX3ZlcnNpb24iOiIxMCIsInJlZmVycmVyIjoiIiwicmVmZXJyaW5nX2RvbWFpbiI6IiIsInJlZmVycmVyX2N1cnJlbnQiOiIiLCJyZWZlcnJpbmdfZG9tYWluX2N1cnJlbnQiOiIiLCJyZWxlYXNlX2NoYW5uZWwiOiJzdGFibGUiLCJjbGllbnRfYnVpbGRfbnVtYmVyIjowLCJjbGllbnRfZXZlbnRfc291cmNlIg=="
	}

	buildNum, _ := strconv.Atoi(GetDiscordBuildNumber())

	toEncode := fmt.Sprintf(`{"os":"Windows","browser":"Chrome","device":"","system_locale":"en-US","browser_user_agent":"%v","browser_version":"%v","os_version":"10","referrer":"","referring_domain":"","referrer_current":"","referring_domain_current":"","release_channel":"stable","client_build_number":%d,"client_event_source":null}`, userAgent, browserVersion, buildNum)

	return base64.StdEncoding.EncodeToString([]byte(toEncode))
}

func Ok(statusCode int) bool {
	for _, v := range []int{200, 201, 204} {
		if v == statusCode {
			return true
		}
	}

	return false
}
