Dogma
============

Chat application written in Go and Javascript.
In order to use this software you must have an mongoDB server up and running. You can specify name and address in the config file SocketServer/config.gcfg

```
Example config file
; Config file for dogma chat
; Don't forget to change the port in bower_components/chat-room/chat-room.js
[Server]
bind = 0.0.0.0 #Default is 127.0.0.1
port = 4000 #Port to use.
[DB]
address = localhost #Default is localhost
name = Dogma #Default is Dogma
```

How to install:
```
go get -u code.google.com/p/go.net/websocket
go get -u github.com/codegangsta/martini-contrib/binding
go get -u github.com/codegangsta/martini-contrib/render
go get -u github.com/codegangsta/martini-contrib/sessions
go get -u github.com/russross/blackfriday
go get -u labix.org/v2/mgo
```
or use the getDependencies.sh script which does the exact same thing.

#AUTHORS:

- Christopher Lillthors christopher.lillthors@gmail.com @github [christopherL91](http://github.com/christopherL91)
- Viktor Kronvall viktor.kronvall@gmail.com @github [considerate](http://github.com/considerate)


##LICENSE
[MIT License](./public/markdown/license.md)