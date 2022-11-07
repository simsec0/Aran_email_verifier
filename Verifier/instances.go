package verifier

import (
	tls_client "github.com/bogdanfinn/tls-client"
)

type Instance struct {
	Client           tls_client.HttpClient
	FullToken        string
	NewEmail         string
	OldPassword      string
	NewEmailPassword string
	Token            string
	UserAgent        string
	BrowserVersion   string
	XSuperProps      string
	Csrf             string
}

type Config struct {
	ApiKey        string `json:"apikey,omitempty"`
	Domain        string `json:"domain,omitempty"`
	Proxyless     bool   `json:"proxyless,omitempty"`
	Workers       int    `json:"amt_of_workers,omitempty"`
	ClientTimeout int    `json:"client_timeout,omitempty"`
	Fingerprint   struct {
		UserAgent      string `json:"user_agent,omitempty"`
		BrowserVersion string `json:"browser_version,omitempty"`
	} `json:"fingerprint,omitempty"`
}
