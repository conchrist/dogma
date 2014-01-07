(function (window) {
  "use strict";

  var username = 'USERNAME';
  var textField = document.getElementById('textField');
  var sendButton = document.getElementById('sendButton');
  var pingButton = document.getElementById('pingButton');
  var socket;

  function run () {
    socket = connect('localhost', 4000,'/echo');

    socket.onmessage = function (evt) {
      var data = evt.data;
      var object = JSON.parse(data);
      //Don't add your own messages to the list.
      if(object.from !== username) {
        messages.push(object);
      }
      console.log("Message received", object);
    };

    socket.onerror = function (e) {
      console.err(e);
    };
  }

  var messages = [];

  Object.observe(messages, function(changes) {
    renderMessages();
  });

  function renderMessages() {
    var messageElem = document.getElementById('messages');
    messageElem.innerHTML = '';
    messages.forEach(function (message) {
      var li = document.createElement('li');
      var textNode = document.createTextNode(message.from+': '+message.message);
      li.appendChild(textNode);
      messageElem.appendChild(li);
    });
  }


  textField.addEventListener('keyup', function (event) {
    var keyCode = event.keyCode;
    console.log(keyCode);
    if(keyCode === 13) {
      handleText();
    }
  }, false);

  sendButton.addEventListener('click', function () {
    handleText();
  });

  pingButton.addEventListener('click', function () {
    sendMessage('Pong', socket);
  });

  function connect(host,port,path) {
    if(Array.isArray(path)) {
      path = pathArray.join('/');
    }
    if(path[0] === '/') {
      path = path.slice(1,path.length);
    }
    var ws = new WebSocket('ws://'+host+':'+port+'/'+path);
    return ws; 
  }

  function handleText() {
    var text = textField.value;
    console.log('Sending', text);
    sendMessage(text,socket);
    textField.value = '';
  }

  function sendMessage(message, socket) {
    var object = {
      message: message,
      time: Date.now(),
      from: username
    };
    messages.push(object);
    socket.send(JSON.stringify(object));
  }

  window.main = {
    run: run
  };

})(this);