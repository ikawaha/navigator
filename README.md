Navigator
===

Navigator is a WebDriver client library for Go. This webdriver is created by modifying the [sclevine/agouti](https://github.com/sclevine/agouti).

## Usage

### Prepare WebDriver

Navigator is a client to operate WebDriver, so you need to prepare your favorite WebDriver (chromedriver, geckodriver, or elenium-server-standalone, etc). 

e.g. For MacOS (using homebrew):
```
brew install chromedriver
```

### Getting Started

To understand how the library works, a simple sample is shown below.
In preparation for running the sample, [httpbin.org](https://httpbin.org) should be running locally.

Run localy: httpbin.org
```
docker run -p 80:80 kennethreitz/httpbin
```
#### Sample code

```go
package main

import (
	"context"
	"fmt"
	"os"

	"github.com/ikawaha/navigator"
)

func main() {
	if err := run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func run() error {
	driver := navigator.ChromeDriver(navigator.Browser("chrome"), navigator.Debug)
	defer driver.Stop()

	ctx := context.Background()
	if err := driver.Start(ctx); err != nil {
		return fmt.Errorf("driver.Start() failed: %w", err)
	}

	page, err := driver.NewPage()
	if err != nil {
		return fmt.Errorf("dirver.NewPage() failed: %w", err)
	}

	page.Navigate("http://localhost:80/forms/post")

	if err := page.FindByName("custname").Fill("John"); err != nil {
		return err
	}
	if err := page.FindByName("custtel").Fill("0000000000"); err != nil {
		return err
	}
	if err := page.FindByName("custemail").Fill("example@example.com"); err != nil {
		return err
	}
	if err := page.All("fieldset").FindByLabel("Medium").Click(); err != nil {
		return err
	}
	if err := page.All("fieldset").FindByLabel("Bacon").Check(); err != nil {
		return err
	}
	if err := page.All("fieldset").FindByLabel("Onion").Check(); err != nil {
		return err
	}
	if err := page.FindByName("delivery").Fill("12:00"); err != nil {
		return err
	}

	if err := page.FindByButton("Submit order").Submit(); err != nil {
		return err
	}

	fmt.Println(page.Find("body").Text())

	// time.Sleep(10 * time.Second)

	return nil
}
```

[![](https://user-images.githubusercontent.com/4232165/205639642-66c032ad-f52c-4b84-84c8-75d44225c5bd.png)](https://user-images.githubusercontent.com/4232165/205639642-66c032ad-f52c-4b84-84c8-75d44225c5bd.png)

---
MIT