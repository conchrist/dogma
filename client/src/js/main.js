/* Copyright (C) 2013 Christopher Lillthors and Viktor Kronvall
 * This file is part of GoWebSocket.
 *
 * GoWebSocket is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * GoWebSocket is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with GoWebSocket.  If not, see <http://www.gnu.org/licenses/>.
 */

(function (window, document) {
  "use strict";

  var username;
  var form = document.getElementById('sendForm');
  var textField = document.getElementById('textField');
  var sendButton = document.getElementById('sendButton');
  var socket;


  function run () {
    socket = connect(window.location.hostname, 4000,'/echo');

    socket.onopen = function () {
      requestUsername();
    };

    socket.onmessage = function (evt) {
      var data = evt.data;
      var object = JSON.parse(data);
      if(object.type === 'message') {
        messages.push(object);
        renderMessages();
        console.log("Message received", object);
      } else if (object.type === 'user') {
        username = object.body;
      }
    };

    socket.onerror = function (e) {
      console.err(e);
    };
  }

  var messages = [];

  /*Object.observe(messages, function(changes) {
    renderMessages();
  });*/

  function renderMessages() {
    var listMessages = messages.slice(0);
    if(document.documentElement.clientWidth <= 480) {
      window.scrollTo(0,document.body.scrollHeight)
    }
    else {
      listMessages = listMessages.reverse(); 
    }
    var messageElem = document.getElementById('messages');
    messageElem.innerHTML = '';
    listMessages.forEach(function (message) {
      var li = document.createElement('li');
      var textNode = document.createTextNode(message.from+': '+message.body);
      li.appendChild(textNode);
      messageElem.appendChild(li);
    });
    console.log(document.documentElement.clientWidth, document.body.scrollHeight)

  }

  function requestUsername() {
    var object = {
      type: 'user',
      time: Date.now()
    };
    var string = JSON.stringify(object);
    socket.send(string);
  }

  sendButton.addEventListener('click', function () {
    handleText();
  });

  form.addEventListener('submit', function (event) {
    event.preventDefault();
    handleText();
  });

  function connect(host,port,path) {
    if(Array.isArray(path)) {
      path = pathArray.join('/');
    }
    if(path[0] === '/') {
      path = path.slice(1,path.length);
    }
    var ws = new WebSocket('wss://'+host+':'+port+'/'+path);
    return ws; 
  }

  function handleText() {
    if(!username) {
      return;
    }
    var text = textField.value;
    console.log('Sending', text);
    sendMessage(text,socket);
    textField.value = '';
  }

  function sendMessage(message, socket) {
    var object = {
      body: message,
      time: Date.now(),
      from: username,
      type: 'message'
    };
    socket.send(JSON.stringify(object));
  }

  window.main = {
    run: run,
    requestUsername: requestUsername
  };

})(this, document);