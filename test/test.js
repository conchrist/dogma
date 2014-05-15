var spawn = require('child_process').spawn;
var WebSocket = require('ws');
var ws = new WebSocket('wss://localhost:8001/chat');

var go = spawn('/usr/bin/env', ['go','run','--config="test.gcfg"','dogma.go']);

go.stdout.on('data', function (data) {
  console.log('stdout: ' + data);
});

var ls    = spawn('ls', ['-lh', '/usr']);

ls.stdout.on('data', function (data) {
  console.log('stdout: ' + data);
});


ws.on('open', function() {
    console.log('Connected');
});
ws.on('message', function(message) {
    console.log('received: %s', message);
});