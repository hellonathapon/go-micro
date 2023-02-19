package main

import "net/http"

func (app *Config) SendMail(w http.ResponseWriter, r *http.Request) {
	type mailMessage struct {
		From    string `json:"from"`
		To      string `json:"to"`
		Subject string `json:"subject"`
		Message string `json:"message"`
	}

	var requestPayload mailMessage

	err := app.readJSON(w, r, &requestPayload)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	msg := Message{
		From:    requestPayload.From,
		To:      requestPayload.To,
		Subject: requestPayload.Subject,
		Data:    requestPayload.Message,
	}

	// send email, is another protocal further than http api
	// this mail service first takes json rest api http protocal forwarding from broker service
	// and verify neccesary data than leaves it to mailer functionality to process mailing protocal
	// so in this service there are two explicitly server running concurrently http service (before implement rpc)
	// and mailer server.
	err = app.Mailer.SendSMTPMessage(msg)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	payload := jsonResponse{
		Error:   false,
		Message: "send to " + requestPayload.To,
	}

	app.writeJSON(w, http.StatusAccepted, payload)
}
