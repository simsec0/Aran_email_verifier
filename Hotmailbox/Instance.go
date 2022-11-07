package email

import "net/http"

type Instance struct {
	Client   *http.Client
	Email    string
	Password string
	ApiKey   string
	Timeout  int
}
