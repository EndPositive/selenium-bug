package main

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/go-rod/rod"
	"github.com/sirupsen/logrus"
	"github.com/tebeka/selenium"
)

const hubURL = "http://localhost:4444"

func main() {
	seleniumCapabilities := selenium.Capabilities{}
	seleniumCapabilities["browserName"] = "chrome"

	wd, err := selenium.NewRemote(seleniumCapabilities, fmt.Sprintf("%s/wd/hub", hubURL))
	if err != nil {
		logrus.Fatalf("Couldn't create Selenium Session: %w", err)
	}
	defer wd.Quit()

	if wd.SessionID() == "" {
		logrus.Fatalf("Couldn't create Selenium Session: %w", fmt.Errorf("no session ID returned"))
	}

	cdpURL, err := url.Parse(fmt.Sprintf("%s/session/%s/se/cdp", hubURL, wd.SessionID()))
	if err != nil {
		logrus.Fatalf("could not parse websocket url: %w", err)
	}
	cdpURL.Scheme = "ws"

	browser := rod.New().ControlURL(cdpURL.String())

	err = browser.Connect()
	if err != nil {
		logrus.Fatalf("failed to connect to browser: %w", err)
	}
	defer browser.MustClose()

	page := browser.MustPage("https://github.com")

	logrus.Infof("from now will hang... search selenium-bug-chrome-1 logs for `console.log`")
	eval, err := page.Eval(fmt.Sprintf("console.log('%s')", strings.Repeat("Test", 100000)))
	if err != nil {
		logrus.Fatalf("Error while evaluating large string over WS: %w", err)
	}

	logrus.Info(eval)

}
