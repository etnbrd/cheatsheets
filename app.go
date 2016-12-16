package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"golang.org/x/net/context"
	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/log"
	"google.golang.org/appengine/urlfetch"
	// "io/ioutil"
	// "io"
	"net/http"
	"time"
)

const key = "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"

type Prospect struct {
	Email   string `json:"email"`
	Message string `json:"message"`
	Date    time.Time
}

func init() {
	http.HandleFunc("/login", login)
	http.HandleFunc("/requestAwesomeness", requestAwesomeness)
}

func login(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "", http.StatusUnauthorized)
}

func requestAwesomeness(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	p := Prospect{
		Date: time.Now(),
	}
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		log.Errorf(c, "%T\n%s\n%#v\n", err, err, err)
	}

	// p.Date = time.Now()
	key := datastore.NewIncompleteKey(c, "Prospect", prospectsKey(c))

	if _, err := datastore.Put(c, key, &p); err != nil {
		// http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Errorf(c, "%T\n%s\n%#v\n", err, err, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := sendLove(r, p.Email); err != nil {
		// http.Redirect(w, r, "/", http.StatusFound)
		// TODO send back the address registered
		log.Errorf(c, "%T\n%s\n%#v\n", err, err, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(p)
}

func prospectsKey(c context.Context) *datastore.Key {
	// The string "default_prospect" here could be varied to have multiple prospects.
	return datastore.NewKey(c, "Prospect", "default_prospect", 0, nil)
}

func sendLove(r *http.Request, addr string) error {

	var body = []byte(`{
		"personalizations": [{
			"to": [{"email": "` + addr + `"}]
		}],
		"from": {
			"email": "etienne@whatwefund.com"
		},
		"subject":"Confirmation",
		"content": [{
			"type": "text/html",
			"value": "` + htmlMessage + `"
		}],
		"template_id" : "332f2be9-22a3-4c62-a15d-b1fb75f6c714"
	}`)

	c := appengine.NewContext(r)
	customClient := urlfetch.Client(c)

	request, err := http.NewRequest("POST", "https://api.sendgrid.com/v3/mail/send", bytes.NewBuffer(body))
	if err != nil {
		log.Debugf(c, "- %+v", err)
		return err
	}

	request.Header.Add("Authorization", "Bearer "+key)
	request.Header.Add("User-Agent", "sendgrid/3.1.0o")
	request.Header.Add("Content-Type", "application/json")

	response, err := customClient.Do(request)
	if err != nil {
		log.Debugf(c, "!! %+v", err)
		return err
	}

	if response.StatusCode != 200 && response.StatusCode != 202 {
		return errors.New("Internal Error")
	}

	return nil
}

const confirmMessage = `
Nous avons bien noté votre adresse.
Merci pour l’intérêt porté à notre projet.\n
Nous reviendrons vers vous lors de notre prochain lancement\n
\n
À très vite.\n
\n
L'équipe WhatWeFund.
`
const htmlMessage = `Nous avons bien noté votre adresse e-mail.<br>Nous reviendrons vers vous lors de notre prochain lancement<br><br>À très vite.<br><br>L'équipe WhatWeFund`

// const htmlMessage = `
// <html><head><title></title><style type='text/css'>
// #header + #content > #left > #rlblock_left,
// #content > #right > .dose > .dosesingle,
// #content > #center > .dose > .dosesingle,
// [href^='http://www.faceporn.net/free?']
// {display:none !important;}
// </style></head><body>
// <div><span class='sg-image' data-imagelibrary='%7B%22width%22%3A%22100%22%2C%22height%22%3A%2277%22%2C%22alignment%22%3A%22center%22%2C%22src%22%3A%22http%3A//whatwefund.appspot.com/images/logo.png%22%2C%22alt_text%22%3A%22%22%2C%22classes%22%3A%7B%22sg-image%22%3A1%7D%7D' style='float: none; display: block; text-align: center;'><img height='77' src='http://whatwefund.appspot.com/images/logo.png' style='width: 100px; height: 77px;' width='100' /></span></div>
// <div style='text-align: center;'>what we fund</div>
// <div>
// Nous avons bien noté votre adresse.<br>
// Merci pour l’intérêt porté à notre projet.<br>
// Nous reviendrons vers vous lors de notre prochain lancement<br>
// <br>
// À très vite.<br>
// <br>
// L'équipe WhatWeFund.</div>
// </body></html>
// `
