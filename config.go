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

func newConfig(options []Option) config {
	var c config
	for _, option := range options {
		option(&c)
	}
	return c
}

func newMergedConfig(config config, options []Option) config {
	for _, option := range options {
		option(&config)
	}
	return config
}

func (c *config) capabilities() Capabilities {
	cb := Capabilities{"acceptSslCerts": true}
	for feature, value := range c.desiredCapabilities {
		cb[feature] = value
	}
	if c.browserName != "" {
		cb.Browser(c.browserName)
	}
	if c.chromeOptions != nil {
		cb["chromeOptions"] = c.chromeOptions
	}
	if c.rejectInvalidSSL {
		cb.Without("acceptSslCerts")
	}
	return cb
}
