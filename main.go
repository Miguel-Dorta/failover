package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"
)

var varPath, url, settingsPath string
var internetStatus, serverStatus bool

func main() {
	readArgs(os.Args)
	readVars()

	if checkServer() {
		setInternetStatus(true)
		setServerStatus(true)
	} else {
		if checkInternet() {
			setInternetStatus(true)
			setServerStatus(false)
		} else {
			setInternetStatus(false)
		}
	}
}

func checkServer() bool {
	_, err := http.Get(url)
	return err == nil
}

func checkInternet() bool {
	err := exec.Command("ping", "-c5", "-W2", "1.1.1.1").Run()
	return err == nil
}

func setInternetStatus(b bool) {
	if internetStatus != b {
		internetStatus = b
		msg := " the Internet"
		if b {
			msg = "Connected to" + msg
		} else {
			msg = "Disconnected from" + msg
		}
		logEvent(msg)
		writeVars()
	}
}

func setServerStatus(b bool) {
	if serverStatus == b {
		return
	}

	doubleCheck := checkServer()
	if doubleCheck != b {
		return
	}

	serverStatus = b
	msg := "Server "
	if b {
		msg += "up"
	} else {
		msg += "down"
	}
	logEvent(msg)
	writeVars()

	if serverStatus {
		return
	}

	thisIP, err := getPublicIP()
	if err != nil {
		logEvent("Error getting public IP: " + err.Error())
		return
	}

	sett, err := readSettings(settingsPath)
	if err != nil {
		logEvent("Error reading settings: " + err.Error())
		return
	}

	err = generateRequests(sett, thisIP)
	if err != nil {
		logEvent(err.Error())
	}
}

func getPublicIP() (string, error) {
	txts, err := net.LookupTXT("o-o.myaddr.l.google.com")
	if err != nil {
		return "", err
	}
	if len(txts) != 1 {
		return "", errors.New("too many responses")
	}
	return txts[0], nil
}

func generateRequests(sett setting, thisIP string) error {
	client := &http.Client{Timeout: time.Second * 2}
	var errs []error

	for _, z := range sett.Zone {
		for _, r := range z.Record {
			err := requestChangeDNS(
				sett.Email,
				sett.Key,
				fmt.Sprintf("https://api.cloudflare.com/client/v4/zones/%s/dns_records/%s", z.ID, r.ID),
				fmt.Sprintf("{\"type\":\"%s\",\"name\":\"%s\",\"content\":\"%s\",\"ttl\":%d,\"proxied\":%t}", r.RecordType, r.Name, thisIP, r.TTL, r.Proxied),
				client,
			)

			if err != nil {
				errs = append(errs, errors.New(fmt.Sprintf("%s: %v", r.Name, err)))
			}
		}
	}

	if len(errs) != 0 {
		msg := "Error updating DNS records:\n"
		for _, e := range errs {
			msg = fmt.Sprintf("%s    %v\n", msg, e)
		}
		return errors.New(msg)
	}
	return nil
}

func requestChangeDNS(email, key, url, data string, client *http.Client) error {
	req, err := http.NewRequest("PUT", url, strings.NewReader(data))
	if err != nil {
		return errors.New("Error creating PUT request: " + err.Error())
	}

	req.ContentLength = int64(len(data))
	req.Header.Set("X-Auth-Email", email)
	req.Header.Set("X-Auth-Key", key)
	req.Header.Set("Content-Type", "application/json")
	res, err := client.Do(req)
	if err != nil {
		return errors.New("Error doing request to Cloudflare: " + err.Error())
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return errors.New("Error reading response body: " + err.Error())
	}

	var resMap map[string]interface{}
	json.Unmarshal(body, &resMap)
	if !resMap["success"].(bool) {
		return errors.New(fmt.Sprintf("Cloudflare Errors: %+v - Cloudflare Messages: %+v", resMap["errors"], resMap["messages"]))
	}
	return nil
}
