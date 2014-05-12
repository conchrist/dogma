(function(window, document, undefined) {
    "use strict";
    Polymer('chat-room', {
        userid: null,
        username: null,
        messages: [],
        //START OMIT
        connect: function(host, port, path) {
            if (Array.isArray(path)) {
                path = pathArray.join('/');
            }
            if (path[0] === '/') {
                path = path.slice(1, path.length);
            }
            var ws = new WebSocket('wss://' + host + ':' + port + '/' + path);
            return ws;
        },
        //END OMIT
        formSubmit: function(event) {
            event.preventDefault();
            this.handleText();
        },
        init: function() {
            this.messages = [];
            var socket = this.connect(window.location.hostname, window.location.port || 4000, '/chat');
            socket.onopen = function(e) {
                socket.send(JSON.stringify({
                    type: 'contact_list',
                    body: '',
                    time: Date.now(),
                    from: this.username
                }));
            }.bind(this);
            socket.onmessage = function(e) {
                var data = e.data;
                var object = JSON.parse(data);
                if (object.type === 'message') {
                    this.messages.push(object);
                } else if (object.type === 'contacts') {
                    this.contacts = object.contacts;
                }
            }.bind(this);
            socket.onerror = function(e) {
                console.error(e.data);
            }.bind(this);
            this.socket = socket;
        },
        handleText: function() {
            if (!this.username) {
                return;
            }
            var text = this.$.textField.value;
            console.log('Sending', text);
            this.sendMessage(text);
            this.$.textField.value = '';
        },
        //START SEND OMIT
        sendMessage: function(message) {
            var object = {
                body: message,
                time: Date.now(),
                from: this.username,
                type: 'message'
            };
            this.socket.send(JSON.stringify(object));
        },
        //END SEND OMIT
        toggleDrawer: function() {
            this.$.drawerpanel.togglePanel();
        }
    });
}(this, document));