package atheme

import "testing"

const (
	ACCNAME = "fooicus"                       // Test account name
	PASS    = "foo"                           // Test password
	SERVER  = "http://proxy.xeserv.us/xmlrpc" // Test server
)

func Test_AthemeCreate(t *testing.T) {
	a, err := NewAtheme(SERVER)
	if err != nil {
		t.Fatal(err)
	}

	if a.Account != "*" {
		t.Fatal("Account is not \"*\" ", a.Account)
	}

	t.Logf("#%v\n", a)
}

func Test_AthemeLogin(t *testing.T) {
	a, err := NewAtheme(SERVER)
	if err != nil {
		t.Fatal(err)
	}

	success, err := a.Login(ACCNAME, PASS)

	t.Logf("%#v\n", a)

	if !success {
		t.Fatal(err)
	}

	if a.Account != ACCNAME {
		t.Fatalf("Account is %v not %v", a.Account, ACCNAME)
	}

	if a.authcookie == "" {
		var res string
		err := a.serverProxy.Call("atheme.login", []string{ACCNAME, PASS}, &res)
		t.Logf("%#v %#v", res, err)

		t.Fail()
	}
}
