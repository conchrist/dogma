package SocketServer

import (
	"code.google.com/p/gcfg"
	"code.google.com/p/go.net/websocket"
	"crypto/tls"
	"fmt"
	. "github.com/smartystreets/goconvey/convey"
	"labix.org/v2/mgo"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"net/url"
	"testing"
)

func createHTTPClient() (*http.Client, *cookiejar.Jar) {
	j, err := cookiejar.New(nil)
	if err != nil {
		panic(err)
	}

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	client := &http.Client{
		Jar:       j,
		Transport: tr,
	}

	return client, j
}

func login(client *http.Client, serverURL *url.URL, jar *cookiejar.Jar) (*http.Cookie, *http.Response, error) {
	resp, err := client.PostForm(serverURL.String(), url.Values{"username": {"viktor"}, "password": {"hej"}})
	cookie := jar.Cookies(serverURL)[0]
	return cookie, resp, err
}

func createWSConnection(serverURL *url.URL, cookie *http.Cookie) *websocket.Conn {
	origin := "https://127.0.0.1/"
	socketURL := "wss://" + serverURL.Host + "/chat"
	socketConfig, _ := websocket.NewConfig(socketURL, origin)
	socketConfig.TlsConfig = &tls.Config{InsecureSkipVerify: true}
	socketConfig.Header.Add("Cookie", cookie.String())
	ws, err := websocket.DialConfig(socketConfig)
	if err != nil {
		fmt.Println(err.Error())
	}
	return ws
}

func TestWebsocket(t *testing.T) {
	var err error
	config := new(Config)
	gcfg.ReadFileInto(config, "../test.gcfg")
	conn, err := mgo.Dial(config.DB.Address)
	db := conn.DB(config.DB.Name)

	db.C("Users").DropCollection()
	db.C("Messages").DropCollection()

	app, _, _, _ := App("../test.gcfg")
	ts := httptest.NewTLSServer(app)

	address := ts.URL + "/users"
	serverURL, _ := url.Parse(address)
	client, jar := createHTTPClient()
	cookie, _, err := login(client, serverURL, jar)
	Convey("Websocket Test", t, func() {
		Convey("Should be able to login", func() {
			ShouldBeNil(err)

		})
	})

	defer ts.Close()
	ws := createWSConnection(serverURL, cookie)
	websocket.JSON.Send(ws, &messageStruct{
		From:    "viktor",
		Message: "Hej",
		Time:    0,
		Type:    "message",
	})
}
