package navigator

import (
	"net/http"
	"time"
)

// An Option specifies configuration for a new WebDriver or Page.
type Option func(*config)

// Debug is an Option that connects the running WebDriver to stdout and stdin.
var Debug Option = func(c *config) {
	b := true
	c.debug = &b
}

// HTTPClient provides an Option for specifying a *http.Client
func HTTPClient(client *http.Client) Option {
	return func(c *config) {
		c.httpClient = client
	}
}

// Timeout provides an Option for specifying a timeout in seconds.
func Timeout(seconds int) Option {
	return func(c *config) {
		t := time.Duration(seconds) * time.Second
		c.timeout = &t
	}
}

// Browser provides an Option for specifying a browser.
func Browser(name string) Option {
	return func(c *config) {
		c.browserName = name
	}
}

// ChromeOptions is used to pass additional options to Chrome via ChromeDriver.
// e.g.
// ChromeOptions("args", []strings{"--headless"}
// ChromeOptions("prefs", map[string]any{"download.default_directory": "/tmp"})
func ChromeOptions(opt string, value any) Option {
	return func(c *config) {
		if c.chromeOptions == nil {
			c.chromeOptions = map[string]any{}
		}
		c.chromeOptions[opt] = value
	}
}

// Desired provides an Option for specifying desired WebDriver Capabilities.
func Desired(capabilities Capabilities) Option {
	return func(c *config) {
		c.desiredCapabilities = capabilities
	}
}

// RejectInvalidSSL is an Option specifying that the WebDriver should reject
// invalid SSL certificates. All WebDrivers should accept invalid SSL certificates
// by default. See: http://www.w3.org/TR/webdriver/#invalid-ssl-certificates
var RejectInvalidSSL Option = func(c *config) {
	c.rejectInvalidSSL = true
}
