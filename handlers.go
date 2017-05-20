package main

import (
	"net/http"

	"crypto/hmac"
	"crypto/sha256"

	"encoding/hex"
	"io"

	"log"

	"os"

	"github.com/gorilla/schema"
	r "gopkg.in/gorethink/gorethink.v3"
)

var schemaDecoder *schema.Decoder

func init() {
	schemaDecoder = schema.NewDecoder()
}

func ReceiveEmail(responseWriter http.ResponseWriter, request *http.Request) {
	verified := verifyRequest(request)
	if !verified {
		responseWriter.WriteHeader(http.StatusNotAcceptable)
		log.Panicln("Fake request!")
	}

	err := request.ParseForm()
	if err != nil {
		http.Error(responseWriter, err.Error(), http.StatusBadRequest)
		log.Panicln(err)
	}

	session, err := r.Connect(r.ConnectOpts{
		Address:  os.Getenv("R_ADDR"),
		Database: "mailgun",
	})
	if err != nil {
		http.Error(responseWriter, err.Error(), http.StatusInternalServerError)
		log.Panicln(err)
	}

	var mail Mail
	schemaDecoder.Decode(&mail, request.PostForm)

	_, err = r.Table("mails").Insert(mail).RunWrite(session)
	if err != nil {
		http.Error(responseWriter, err.Error(), http.StatusInternalServerError)
		log.Panicln(err)
	}
}

func verifyRequest(request *http.Request) bool {
	sig, err := hex.DecodeString(request.FormValue("signature"))
	if err != nil {
		return false
	}

	mac := hmac.New(sha256.New, []byte(os.Getenv("MG_API_KEY")))
	io.WriteString(mac, request.FormValue("timestamp"))
	io.WriteString(mac, request.FormValue("token"))
	expectedSig := mac.Sum(nil)

	return len(sig) == len(expectedSig) && hmac.Equal(sig, expectedSig)
}
