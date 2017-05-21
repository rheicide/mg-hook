package main

type Route struct {
	Methods []string
	Path    string
	Handler Handler
}

type Routes []Route

var routes = Routes{
	Route{
		[]string{"POST"},
		"/",
		ReceiveEmail,
	},
}
