package main

import (
	"net/http"

	"crypto/hmac"
	"crypto/sha256"

	"encoding/hex"
	"io"

	"log"

	"os"

	"reflect"
	"time"

	"errors"

	"fmt"

	"github.com/gorilla/schema"
	r "gopkg.in/gorethink/gorethink.v3"
)

var (
	mgApiKey       string
	session        *r.Session
	decoder        *schema.Decoder
	gitCommit      string
	buildTimestamp string
)

func init() {
	mgApiKey = os.Getenv("MG_API_KEY")
	if mgApiKey == "" {
		log.Fatalln("Mailgun API key is missing")
	}

	var err error
	session, err = r.Connect(r.ConnectOpts{
		Address:  os.Getenv("R_ADDR"),
		Database: "mailgun",
	})
	if err != nil {
		log.Fatalln(err)
	}

	decoder = schema.NewDecoder()
	decoder.IgnoreUnknownKeys(true)
	decoder.RegisterConverter(time.Time{}, func(value string) reflect.Value {
		timeFormat := "Mon, 02 Jan 2006 15:04:05 -0700 (MST)"
		date, err := time.Parse(timeFormat, value)
		if err != nil {
			log.Fatalln(err)
		}

		return reflect.ValueOf(date)
	})
}

type Handler func(http.ResponseWriter, *http.Request) (error, int)

func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	start := time.Now()

	err, status := h(w, r)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), status)
	}

	defer log.Printf("%s\t%s\t%s", r.Method, r.RequestURI, time.Since(start))
}

func ReceiveEmail(_ http.ResponseWriter, request *http.Request) (error, int) {
	err := verifyRequest(request)
	if err != nil {
		return err, http.StatusNotAcceptable
	}

	err = request.ParseForm()
	if err != nil {
		return err, http.StatusBadRequest
	}

	var mail Mail
	err = decoder.Decode(&mail, request.PostForm)
	if err != nil {
		return err, http.StatusInternalServerError
	}

	_, err = r.Table("mails").Insert(&mail).RunWrite(session)
	if err != nil {
		return err, http.StatusInternalServerError
	}

	return nil, 0
}

func verifyRequest(r *http.Request) error {
	sig, err := hex.DecodeString(r.PostFormValue("signature"))
	if err != nil {
		return errors.New("Could not decode signature")
	}

	mac := hmac.New(sha256.New, []byte(mgApiKey))
	io.WriteString(mac, r.PostFormValue("timestamp"))
	io.WriteString(mac, r.PostFormValue("token"))
	expectedSig := mac.Sum(nil)

	if len(sig) == len(expectedSig) && hmac.Equal(sig, expectedSig) {
		return nil
	} else {
		return errors.New("Fake request")
	}
}

func Version(w http.ResponseWriter, r *http.Request) (error, int) {
	_, err := fmt.Fprintf(w, "git-commit: %s\nbuild-timestamp: %s\n", gitCommit, buildTimestamp)
	if err != nil {
		return err, http.StatusInternalServerError
	}

	return nil, 0
}
