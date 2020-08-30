// authenticator_test tests the token interaction bit of the api
// as well as export/import functionality
package authenticator_test

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/jcelliott/lumber"

	"github.com/thordin9/logvac/api"
	"github.com/thordin9/logvac/authenticator"
	"github.com/thordin9/logvac/collector"
	"github.com/thordin9/logvac/config"
	"github.com/thordin9/logvac/core"
	"github.com/thordin9/logvac/drain"
)

func TestMain(m *testing.M) {
	// clean test dir
	os.RemoveAll("/tmp/authTest")

	// manually configure
	err := initialize()
	if err != nil {
		os.Exit(1)
	}

	// start api
	go api.Start(collector.CollectHandler)
	<-time.After(1 * time.Second)
	rtn := m.Run()

	// clean test dir
	os.RemoveAll("/tmp/authTest")

	os.Exit(rtn)
}

// test removing an auth token
func TestRemoveToken(t *testing.T) {
	err := authenticator.Remove("nobody")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
}

// test adding an auth token
func TestAddToken(t *testing.T) {
	err := authenticator.Add("user")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
}

// test get logs // doesn't provide any "coverage" but ensures auth works
func TestGetLogs(t *testing.T) {
	body, err := rest("GET", "/logs?type=app", "")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	if string(body) != "[]\n" {
		t.Errorf("%q doesn't match expected out", body)
		t.FailNow()
	}
}

// test importing auth tokens
func TestImport(t *testing.T) {
	token := &bytes.Buffer{}
	token.Write([]byte("[\"user2\"]"))
	err := authenticator.ImportLogvac(token)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	if string(token.Bytes()) != "" {
		t.Errorf("%q doesn't match expected out", token)
		t.FailNow()
	}
}

// test exporting auth tokens
func TestExport(t *testing.T) {
	token := &bytes.Buffer{}
	err := authenticator.ExportLogvac(token)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	if string(token.Bytes()) != "[\"user\",\"user2\"]" {
		t.Errorf("%q doesn't match expected out", token)
		t.FailNow()
	}
}

// hit api and return response body
func rest(method, route, data string) ([]byte, error) {
	body := bytes.NewBuffer([]byte(data))

	req, _ := http.NewRequest(method, fmt.Sprintf("https://%s%s", config.ListenHttp, route), body)
	req.Header.Add("X-AUTH-TOKEN", "secret")
	req.Header.Add("X-USER-TOKEN", "user")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Unable to %s %s - %s", method, route, err)
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return nil, fmt.Errorf("Status '200' expected, got '%d'", res.StatusCode)
	}

	b, _ := ioutil.ReadAll(res.Body)

	return b, nil
}

// manually configure and start internals
func initialize() error {
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	config.ListenHttp = "127.0.0.1:3234"
	config.ListenTcp = "127.0.0.1:3235"
	config.ListenUdp = "127.0.0.1:3234"
	config.DbAddress = "boltdb:///tmp/authTest/logvac.bolt"
	config.Log = lumber.NewConsoleLogger(lumber.LvlInt("ERROR"))

	// initialize logvac
	logvac.Init()

	// setup authenticator
	config.AuthAddress = ""
	err := authenticator.Init()
	if err != nil {
		return fmt.Errorf("Authenticator failed to initialize - %s", err)
	}
	config.AuthAddress = "file:///tmp/authTest/logvac-auth.bolt"
	err = authenticator.Init()
	if err != nil {
		return fmt.Errorf("Authenticator failed to initialize - %s", err)
	}
	config.AuthAddress = "~!@#$%^&*()_"
	err = authenticator.Init()
	if err == nil {
		return fmt.Errorf("Authenticator failed to initialize - %s", err)
	}
	config.AuthAddress = "boltdb:///tmp/authTest/logvac-auth.bolt"
	err = authenticator.Init()
	if err != nil {
		return fmt.Errorf("Authenticator failed to initialize - %s", err)
	}

	// initialize drains
	err = drain.Init()
	if err != nil {
		return fmt.Errorf("Drain failed to initialize - %s", err)
	}

	// initializes collectors
	err = collector.Init()
	if err != nil {
		return fmt.Errorf("Collector failed to initialize - %s", err)
	}

	return nil
}
