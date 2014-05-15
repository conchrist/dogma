package SocketServer

/*
Test code for dogma chat.
*/

import (
	"code.google.com/p/gcfg"
	"code.google.com/p/go.net/websocket"
	"crypto/tls"
	. "github.com/smartystreets/goconvey/convey"
	"labix.org/v2/mgo"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"net/url"
	"testing"
)

func createHTTPClient() (*http.Client, *cookiejar.Jar, error) {
	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, nil, err
	}

	transport := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true},
	}

	client := &http.Client{
		Jar:       jar,
		Transport: transport,
	}
	return client, jar, nil
}

func createUser(client *http.Client, serverURL *url.URL, jar *cookiejar.Jar) (*http.Cookie, *http.Response, error) {
	values := url.Values{
		"username": {"viktor"},
		"password": {"123"},
	} //the actual values we want to send to the server.
	resp, err := client.PostForm(serverURL.String(), values)
	cookie := jar.Cookies(serverURL)[0] //We suppose we only get one cookie from the server.
	return cookie, resp, err
}

func createWSConnection(serverURL *url.URL, cookie *http.Cookie) (*websocket.Conn, error) {
	origin := "https://127.0.0.1/" //server address
	socketURL := "wss://" + serverURL.Host + "/chat"
	socketConfig, _ := websocket.NewConfig(socketURL, origin)
	socketConfig.TlsConfig = &tls.Config{InsecureSkipVerify: true}
	socketConfig.Header.Add("Cookie", cookie.String())
	ws, err := websocket.DialConfig(socketConfig)
	return ws, err
}

func TestWebsocket(t *testing.T) {
	var db *mgo.Database
	var server *httptest.Server
	var serverURL *url.URL
	var client *http.Client
	var jar *cookiejar.Jar
	var cookie *http.Cookie
	Convey("Websocket Test", t, func() {
		Convey("Should connect to DB", func() {
			config := new(Config)
			gcfg.ReadFileInto(config, "../test.gcfg")

			conn, err := mgo.Dial(config.DB.Address)
			So(err, ShouldBeNil)
			db = conn.DB(config.DB.Name)

			db.C("Users").DropCollection()
			db.C("Messages").DropCollection()
		})

		Convey("Should parse mockup server url", func() {
			var err error
			app, _, _, _ := App("../test.gcfg")
			server = httptest.NewTLSServer(app)

			address := server.URL + "/users"
			serverURL, err = url.Parse(address)
			So(err, ShouldBeNil)
			So(serverURL, ShouldNotBeNil)
		})

		Convey("Should create client", func() {
			var err error
			client, jar, err = createHTTPClient()
			So(err, ShouldBeNil)
		})

		Convey("Should be able to create new user", func() {
			var err error
			So(serverURL, ShouldNotBeNil)
			cookie, _, err = createUser(client, serverURL, jar)
			So(err, ShouldBeNil)
		})

		Convey("Should be able to create ws", func() {
			//defer conn.Close()
			ws, err := createWSConnection(serverURL, cookie)
			So(err, ShouldBeNil)
			websocket.JSON.Send(ws, &messageStruct{
				From:    "viktor",
				Message: "Hej",
				Time:    0,
				Type:    "message",
			})
		})
		Convey("Tear down", func() {
			//ws.Close()
			//server.CloseClientConnections()
			//server.Close()
		})
	})
}
