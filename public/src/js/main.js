(function (window, document) {
  "use strict";

  var username;
  var form = document.getElementById('sendForm');
  var textField = document.getElementById('textField');
  var sendButton = document.getElementById('sendButton');
  var socket;


  function run () {
    socket = connect('considerate.com',4000,'/chatroom');

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
      console.error(e.data);
    };
  }

  var messages = [];


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
    console.log('wss://'+host+':'+port+'/'+path);
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