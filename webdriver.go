package navigator

import (
	"fmt"
	"runtime"

	"github.com/ikawaha/navi/webdriver"
)

// A WebDriver controls a WebDriver process. This struct embeds webdriver.WebDriver,
// which provides Start and Stop methods for starting and stopping the process.
type WebDriver struct {
	*webdriver.WebDriver
	defaultConfig config
}

// NewWebDriver returns an instance of a WebDriver specified by
// a templated URL and command. The URL should be the location of the
// WebDriver Wire Protocol web service brought up by the command. The
// command should be provided as a list of arguments (each of which are
// templated).
//
// The Timeout Option specifies how many seconds to wait for the web service
// to become available. The default timeout is 5 seconds.
//
// The HTTPClient Option specifies a *http.Client to use for all WebDriver
// communications. The default client is http.DefaultClient.
//
// Any other provided Options are treated as default Options for new pages.
//
// Valid template parameters are:
//
//	{{.Host}} - local address to bind to (usually 127.0.0.1)
//	{{.Port}} - arbitrary free port on the local address
//	{{.Address}} - {{.Host}}:{{.Port}}
//
// Selenium JAR example:
//
//	command := []string{"java", "-jar", "selenium-server.jar", "-port", "{{.Port}}"}
//	navigator.New("http://{{.Address}}/wd/hub", command)
func NewWebDriver(url string, command []string, options ...Option) *WebDriver {
	driver := webdriver.New(url, command)
	c := NewConfig(options)
	if c.timeout != nil {
		driver.Timeout = *c.timeout
	}
	if c.debug != nil {
		driver.Debug = *c.debug
	}
	if c.httpClient != nil {
		driver.HTTPClient = c.httpClient
	}
	return &WebDriver{
		WebDriver:     driver,
		defaultConfig: c,
	}
}

// ChromeDriver returns an instance of a ChromeDriver WebDriver.
//
// Provided Options will apply as default arguments for new pages.
// New pages will accept invalid SSL certificates by default. This
// may be disabled using the RejectInvalidSSL Option.
func ChromeDriver(options ...Option) *WebDriver {
	var binaryName string
	if runtime.GOOS == "windows" {
		binaryName = "chromedriver.exe"
	} else {
		binaryName = "chromedriver"
	}
	command := []string{binaryName, "--port={{.Port}}"}
	return NewWebDriver("http://{{.Address}}", command, options...)
}

// EdgeDriver returns an instance of a EdgeDriver WebDriver.
//
// Provided Options will apply as default arguments for new pages.
// New pages will accept invalid SSL certificates by default. This
// may be disabled using the RejectInvalidSSL Option.
func EdgeDriver(options ...Option) (*WebDriver, error) {
	var binaryName string
	if runtime.GOOS == "windows" {
		binaryName = "MicrosoftWebDriver.exe"
	} else {
		return nil, fmt.Errorf("not supported, windows only")
	}
	command := []string{binaryName, "--port={{.Port}}"}
	// Using {{.Address}} means using 127.0.0.1
	// But MicrosoftWebDriver only supports localhost, not 127.0.0.1
	return NewWebDriver("http://localhost:{{.Port}}", command, options...), nil
}

// GeckoDriver returns an instance of a geckodriver WebDriver which supports
// gecko based browser like Firefox.
//
// Provided Options will apply as default arguments for new pages.
//
// See https://github.com/mozilla/geckodriver for geckodriver details.
func GeckoDriver(options ...Option) *WebDriver {
	var binaryName string
	if runtime.GOOS == "windows" {
		binaryName = "geckodriver.exe"
	} else {
		binaryName = "geckodriver"
	}
	command := []string{binaryName, "--port={{.Port}}"}
	return NewWebDriver("http://{{.Address}}", command, options...)
}

// NewPage returns a *Page that corresponds to a new WebDriver session.
// Provided Options configure the page. For instance, to disable JavaScript:
//
//	capabilities := navigator.NewCapabilities().Without("javascriptEnabled")
//	driver.NewPage(navigator.Desired(capabilities))
//
// For Selenium, a Browser Option (or a Desired Option with Capabilities that
// specify a Browser) must be provided. For instance:
//
//	seleniumDriver.NewPage(navigator.Browser("safari"))
//
// Specific Options (such as Browser) have precedence over Capabilities
// specified by the Desired Option.
//
// The HTTPClient Option will be ignored if passed to this function. New pages
// will always use the *http.Client provided to their WebDriver, or
// http.DefaultClient if none was provided.
func (w *WebDriver) NewPage(options ...Option) (*Page, error) {
	c := NewMergedConfig(w.defaultConfig, options)
	session, err := w.Open(c.Capabilities())
	if err != nil {
		return nil, fmt.Errorf("failed to connect to WebDriver: %w", err)
	}

	return newPage(session), nil
}
