package main

import "time"

type Mail struct {
	From      string    `gorethink:"from"`
	To        string    `gorethink:"to"`
	Recipient string    `schema:"recipient" gorethink:"recipient"`
	Subject   string    `gorethink:"subject"`
	BodyPlain string    `schema:"body-plain" gorethink:"body-plain"`
	BodyHtml  string    `schema:"body-html" gorethink:"body-html"`
	Date      time.Time `gorethink:"date"`
}

type Mails []Mail
