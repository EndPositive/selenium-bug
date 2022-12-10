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

	logrus.Info("Requesting Selenium session (if this hangs, restart Selenium)")
	wd, err := selenium.NewRemote(seleniumCapabilities, fmt.Sprintf("%s/wd/hub", hubURL))
	if err != nil {
		logrus.Fatalf("Couldn't create Selenium Session: %s", err)
	}

	if wd.SessionID() == "" {
		logrus.Fatalf("Couldn't create Selenium Session: %s", fmt.Errorf("no session ID returned"))
	}

	cdpURL, err := url.Parse(fmt.Sprintf("%s/session/%s/se/cdp", hubURL, wd.SessionID()))
	if err != nil {
		logrus.Fatalf("could not parse websocket url: %s", err)
	}
	cdpURL.Scheme = "ws"

	browser := rod.New().ControlURL(cdpURL.String())
	logrus.Info("Connection to Selenium browser")
	err = browser.Connect()
	if err != nil {
		logrus.Fatalf("failed to connect to browser: %s", err.Error())
	}
	defer browser.MustClose()

	logrus.Info("Navigating to https://github.com")
	page := browser.MustPage("https://github.com")
	defer page.MustClose()

	logrus.Infof("from now will hang... search selenium-bug-chrome-1 logs for `console.log`")
	eval, err := page.Eval(fmt.Sprintf("() => { console.log('%s'); return 'worked!'; }", strings.Repeat("Test", 100000)))
	if err != nil {
		logrus.Fatalf("Error while evaluating large string over WS: %s", err)
	}

	logrus.Infof("Output should equal `worked!`: `%s`", eval.Value)

}
