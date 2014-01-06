(function () {
  var textField = document.getElementById('textField');
  var sendButton = document.getElementById('sendButton');
  var pingButton = document.getElementById('pingButton');
  var socket = connect('localhost', 4000,'/echo');

  var messages = [];

  Object.observe(messages, function(changes) {
    //changes.forEach(whatHappened);
    renderMessages();
  });

  function whatHappened(change) {
    console.log(change.name + " was " + change.type + " and is now " + change.object[change.name]);
  }

  function renderMessages() {
    var messageElem = document.getElementById('messages');
    messageElem.innerHTML = '';
    messages.forEach(function (message) {
      var li = document.createElement('li');
      var textNode = document.createTextNode(message.user+': '+message.text);
      li.appendChild(textNode);
      messageElem.appendChild(li);
    })
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

  socket.onmessage = function (evt) {
    var data = evt.data;
    console.log("Message received", data)
    messages.push({
      user: 'you',
      text: data
    });
  };

  socket.onerror = function (e) {
    console.err(e);
  }

  function handleText() {
    var text = textField.value;
    console.log('Sending', text);
    sendMessage(text,socket);
    messages.push({
      user: 'me',
      text: text
    });
    textField.value = '';
  }

  function sendMessage(message, socket) {
    socket.send(message);
  }

  
})();