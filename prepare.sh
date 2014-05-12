#!/bin/bash

echo "Downloading dependencies"

npm install -g bower
bower install
go get -u code.google.com/p/go.net/websocket
go get -u code.google.com/p/gcfg
go get -u github.com/codegangsta/martini-contrib/binding
go get -u github.com/codegangsta/martini-contrib/render
go get -u github.com/codegangsta/martini-contrib/sessions
go get -u github.com/russross/blackfriday
go get -u labix.org/v2/mgo

echo "Done"