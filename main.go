package main

import (
	Email "EmailVerifier/Hotmailbox"
	Utils "EmailVerifier/Modules"
	Verifier "EmailVerifier/Verifier"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/bogdanfinn/fhttp/cookiejar"
	"github.com/its-vichy/GoCycle"

	tls_client "github.com/bogdanfinn/tls-client"
)

var (
	Config  *Verifier.Config
	Tokens  []string
	Proxies *GoCycle.Cycle
	Success int
)

func loadConfig() error {
	configFile, err := os.Open("config.json")
	if err != nil {
		return err
	}
	defer configFile.Close()
	jsonParser := json.NewDecoder(configFile)
	jsonParser.Decode(&Config)
	return nil
}

func createIterator(file string) (*GoCycle.Cycle, error) {
	iterator, err := GoCycle.NewFromFile(file)

	if err != nil {
		return nil, err
	}

	return iterator, nil
}

func init() {
	Utils.Clear()

	Utils.SetTitle("Opti Boosts ^| Developed by Aran")

	err := loadConfig()

	if err != nil {
		Utils.Logger(false, fmt.Sprintf("An error occured while loading config, error: %v", err.Error()))
		time.Sleep(time.Second * 5)
		os.Exit(0)
	}

	if !Config.Proxyless {
		proxyIterator, err := createIterator("Data/proxies.txt")

		if err != nil {
			panic(fmt.Sprintf("Could not create a proxy cycle, error: %v", err.Error()))
		}

		Proxies = proxyIterator
	}

	Tokens, err = Utils.ReadLines("Data/tokens.txt")

	if err != nil {
		Utils.Logger(false, fmt.Sprintf("An error occured while loading tokens, error: %v", err.Error()))
		time.Sleep(time.Second * 5)
		os.Exit(0)
	}
}

func createThread(fullToken, proxy string) {
	defer func() {
		if r := recover(); r != nil {
			Utils.Logger(false, "Goroutine Panicked Error", r)
		}
	}()

	if !strings.Contains(Config.Fingerprint.UserAgent, "Chrome") {
		panic(fmt.Sprintf("Currently useragent %v is not supported, please use a chrome user agent version 103-105", Config.Fingerprint.UserAgent))
	} else if !strings.Contains(fullToken, ":") {
		Utils.Logger(false, "Token doesn't have password |", fullToken)
	}

	var (
		emailClient Email.Instance
		verifyLink  interface{}
		splitToken  []string = strings.Split(fullToken, ":")
		verifyErr   error
		password    string
		token       string
	)

	token, password = splitToken[2], splitToken[1]

	jar, _ := cookiejar.New(nil)

	options := []tls_client.HttpClientOption{
		tls_client.WithTimeout(Config.ClientTimeout),

		func() tls_client.HttpClientOption {
			switch strings.Split(Config.Fingerprint.BrowserVersion, ".")[0] {
			case "103":
				return tls_client.WithClientProfile(tls_client.Chrome_103)
			case "104":
				return tls_client.WithClientProfile(tls_client.Chrome_104)
			case "105":
				return tls_client.WithClientProfile(tls_client.Chrome_105)
			default:
				panic(fmt.Sprintf("Currently version %v is not supported, please use a chrome useragent/version", Config.Fingerprint.BrowserVersion))
			}
		}(),
		tls_client.WithNotFollowRedirects(),
		tls_client.WithInsecureSkipVerify(),
		tls_client.WithCookieJar(jar),
	}

	if proxy != "" {
		options = append(options, tls_client.WithProxyUrl(fmt.Sprintf("http://%v", proxy)))
	}

	httpClient, err := tls_client.NewHttpClient(tls_client.NewNoopLogger(), options...)

	if err != nil {
		Utils.Logger(false, "Error:", err.Error())
	}

	VClient := Verifier.Instance{
		OldPassword:    password,
		Token:          token,
		FullToken:      fullToken,
		UserAgent:      "Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/106.0.0.0 Safari/537.36",
		BrowserVersion: "106.0.0.0",
		Client:         httpClient,
	}

	if len(Config.ApiKey) > 0 {
		emailClient = Email.Instance{
			ApiKey:  Config.ApiKey,
			Timeout: Config.ClientTimeout,
			Client: &http.Client{
				Timeout: time.Second * time.Duration(Config.ClientTimeout+10),
			},
		}

		_, err := emailClient.GetEmail(Config.Domain)

		if err != nil {
			Utils.Logger(false, "Error:", err.Error())
			return
		}

		Utils.Logger(true, fmt.Sprintf("Retrieved %s.com account %s", Config.Domain, VClient.NewEmail))

		VClient.NewEmail = emailClient.Email
	} else {
		err := VClient.GetPage()

		if err != nil {
			Utils.Logger(false, "Error:", err.Error())
			return
		}

		_, err = VClient.GetNewEmail()

		if err != nil {
			Utils.Logger(false, "Error:", err.Error())
			return
		}

		Utils.Logger(true, fmt.Sprintf("Retrieved tormails.com account %s", VClient.NewEmail))
	}

	err = VClient.Session()

	if err != nil {
		Utils.Logger(false, "Error:", err.Error())
		return
	}

	err = VClient.AddEmail()

	if err != nil {
		if strings.Contains(err.Error(), "Token is already verified") {
			Utils.AppendLine("Data/success.txt", VClient.FullToken)
			Utils.RemoveFromFile("Data/tokens.txt", VClient.FullToken)
		}

		Utils.Logger(false, "Error:", err.Error())
		return
	}

	if len(Config.ApiKey) > 0 {
		verifyLink, verifyErr = emailClient.GetVerificationEmail()
	} else {
		verifyLink, verifyErr = VClient.WaitForMail(60, true)
	}

	if verifyErr != nil {
		Utils.Logger(false, "Error:", err.Error())
		return
	}

	err = VClient.VerifyEmail(fmt.Sprintf("%v", verifyLink))

	if err != nil {
		Utils.Logger(false, "Error:", err.Error())
		return
	} else {
		Utils.Logger(true, "Verified email |", VClient.Token)
		Utils.AppendLine("Data/success.txt", VClient.FullToken)
		Utils.RemoveFromFile("Data/tokens.txt", VClient.FullToken)
		Success++
	}

}

func main() {
	go func() {
		for {
			Utils.SetTitle(fmt.Sprintf("Developed by Aran ^| Successes: %d", Success))
			time.Sleep(10 * time.Millisecond)
		}
	}()

	limiter := make(chan struct{}, Config.Workers)

	for _, v := range Tokens {
		var Proxy string

		if !Config.Proxyless {
			Proxy = Proxies.Next()
		} else {
			Proxy = ""
		}

		go func(v string) {
			limiter <- struct{}{}
			createThread(v, Proxy)
			<-limiter
		}(v)
	}

	select {}
}
