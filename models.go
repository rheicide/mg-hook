package main

import "time"

type Mail struct {
	From      string    `gorethink:"from"`
	To        string    `gorethink:"to"`
	Recipient string    `gorethink:"recipient" schema:"recipient"`
	Subject   string    `gorethink:"subject"`
	BodyPlain string    `gorethink:"body-plain" schema:"body-plain"`
	BodyHtml  string    `gorethink:"body-html" schema:"body-html"`
	Date      time.Time `gorethink:"date"`
}

type Mails []Mail
