package gee

import (
	"reflect"
	"testing"
)

const (
	checkMark = "\u2713"
	ballotX   = "\u2717"
)

func addRouteTest() *Router {
	router := NewRouter()
	router.addRoute("GET", "/hello", nil)
	router.addRoute("GET", "/user/:name", nil)
	router.addRoute("GET", "/admin/:name/page", nil)
	router.addRoute("GET", "/fils/*filepath", nil)

	return router
}

func TestParsePattern(t *testing.T) {
	ok := reflect.DeepEqual(parsePattern("/p/:name"), []string{"p", ":name"})
	ok = ok && reflect.DeepEqual(parsePattern("/p/*"), []string{"p", "*"})
	if ok {
		t.Logf("\tTesting function: parsePattern %v", checkMark)
	} else {
		t.Fatalf("\tTesting function: parsePattern %v", ballotX)
	}
}

func TestGetRoute(t *testing.T) {
	t.Log("\tTrying to GET /fils/*user/path")

	r := addRouteTest()
	t.Log("\t\tAdd route ", checkMark)

	n, ps := r.getRoute("GET", "/fils/*user/path")

	if n == nil {
		t.Fatal("\t\tnil shouldn't be returned")
	}
	t.Logf("\t\tShould get returned node %v", checkMark)

	if n.pattern != "/fils/*filepath" {
		t.Fatalf("\t\tshould match /fils/*filepath %v", ballotX)
	}
	t.Logf("\t\tRequest path should match /fils/*filepath %v", checkMark)

	if ps["filepath"] != "user/path" {
		t.Fatalf("\t\tFilepath should be equal to 'user/path' %v", ballotX)
	}
	t.Logf("\t\tFilepath should be equal to 'user/path' %v", checkMark)

}
