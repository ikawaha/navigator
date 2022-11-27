package navigator

import (
	"net/http"
	"time"
)

type config struct {
	// web driver config
	httpClient *http.Client
	debug      *bool
	timeout    *time.Duration

	// capabilities
	browserName         string
	rejectInvalidSSL    bool
	chromeOptions       map[string]any // chrome driver config
	desiredCapabilities Capabilities
}

func NewConfig(options []Option) config {
	var c config
	for _, option := range options {
		option(&c)
	}
	return c
}

func NewMergedConfig(c config, options []Option) config {
	for _, option := range options {
		option(&c)
	}
	return c
}

func (c *config) Capabilities() Capabilities {
	merged := Capabilities{"acceptSslCerts": true}
	for feature, value := range c.desiredCapabilities {
		merged[feature] = value
	}
	if c.browserName != "" {
		merged.Browser(c.browserName)
	}
	if c.chromeOptions != nil {
		merged["chromeOptions"] = c.chromeOptions
	}
	if c.rejectInvalidSSL {
		merged.Without("acceptSslCerts")
	}
	return merged
}
