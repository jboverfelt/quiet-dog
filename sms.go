package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

type smsSender interface {
	SendSMS(to []phoneNumber, from phoneNumber, body io.Reader) error
}

type twilioSmsSender struct {
	baseURL string
	sid     string
	token   string
}

func (t twilioSmsSender) constructRequest(to, from, msgBody string) (*http.Request, error) {
	twilioURL := t.baseURL + t.sid + "/Messages.json"
	form := url.Values{}
	form.Add("To", string(to))
	form.Add("From", string(from))
	form.Add("Body", string(msgBody))

	req, err := http.NewRequest(http.MethodPost, twilioURL, strings.NewReader(form.Encode()))

	if err != nil {
		return nil, err
	}

	req.SetBasicAuth(t.sid, t.token)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	return req, nil
}

func (t twilioSmsSender) SendSMS(to []phoneNumber, from phoneNumber, body io.Reader) error {
	msgBody, err := ioutil.ReadAll(body)

	if err != nil {
		return err
	}

	for _, number := range to {
		req, err := t.constructRequest(string(number), string(from), string(msgBody))

		if err != nil {
			return err
		}

		resp, err := http.DefaultClient.Do(req)

		if err != nil {
			return err
		}

		respBody, err := ioutil.ReadAll(resp.Body)
		resp.Body.Close()

		if err != nil {
			return err
		}

		if resp.StatusCode >= 300 || resp.StatusCode < 200 {
			return fmt.Errorf("sms: Error sending sms, status code %d, body %v", resp.StatusCode, string(respBody))
		}
	}

	return nil
}
