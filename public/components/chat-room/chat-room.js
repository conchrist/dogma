(function(window, document, undefined) {
      "use strict";
      Polymer('chat-room', {
            userid: null,
            username: null,
            messages: [],
            reversed: [],
            showAddOverlay: false,
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
                this.contacts = [];
                var socket = this.connect(window.location.hostname, window.location.port || 4000, '/chat');
                socket.onopen = function(e) {
                    socket.send(JSON.stringify({
                        type: 'contact_list',
                        body: '',
                        time: Date.now(),
                        from: this.username
                    }));
                }.bind(this);
                var contacts = {};
                socket.onmessage = function(e) {
                    var data = e.data;
                    var object = JSON.parse(data);
                    if (['message', 'image'].indexOf(object.type) !== -1) {
                        this.messages.push(object);
                    } else if (object.type === 'contacts') {
                        object.contacts.forEach(function(contact) {
                            contacts[contact] = true;
                        }, this);
                    } else if(object.type === 'client joined') {
                        contacts[object.body] = true;
                    } else if (object.type === 'client left') {
                        delete contacts[object.body];
                    }
                    this.contacts = Object.keys(contacts);
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
                this.socket.send(JSON.stringify({
                    type: 'contact_list',
                    body: '',
                    time: Date.now(),
                    from: this.username
                }));
            },
            logout: function() {
                this.$.logoutajax.go();
                this.socket.close();
                this.fire('logout');
            },
            contactsChanged: function() {
                console.log('changed');
            },
            plusbutton: function() {
                console.log('Button tapped');
                this.showAddOverlay = true;
            },
            sendImage: function(filename, data) {
                  this.socket.send(JSON.stringify({
                      type: 'image',
                      body: data,
                      time: Date.now(),
                      from: this.username
                  }));
            },
            readImages: function(files, handleFile) {
                  var self = this;
                  // Loop through the FileList and render image files as thumbnails.
                  for (var i = 0, f; f = files[i]; i++) {

                        // Only process image files.
                        if (!f.type.match('image.*')) {
                            continue;
                        }

                        var reader = new FileReader();
                        // Closure to capture the file information.
                        reader.onload = (function(theFile) {
                           return function(e) {
                              // Render thumbnail.
                              var result =e.target.result; 
                              var filename = theFile.name;
                              handleFile.call(self, filename, result);
                           };
                        })(f);

                        // Read in the image file as a data URL.
                        reader.readAsDataURL(f);
                  }
            },
            addFile: function() {
                var files = this.$.mediafile.files;
                this.readImages(files, this.sendImage);
            },
            messagesChanged: function() {
                var array = this.messages;
                var copy = array.slice();
                copy.reverse();
                this.reversed = copy;
            }
      });
}(this, document));