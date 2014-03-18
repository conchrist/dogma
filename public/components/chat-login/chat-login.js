Polymer('chat-login', {
	ready: function () {
		console.log('Login form ready!');
		var form = this.$.loginForm;
		form.addEventListener('submit', function(event) {
			event.preventDefault();
			console.log('Do eet!');
			var ajax = document.createElement('polymer-ajax');

			ajax.url = form.action;
			ajax.method = form.method;
			var username = this.$.username.value;
			var password = this.$.password.value;
			var params = {
				username: username,
				password: password
			};
			console.log(params);
			ajax.params = params;


			ajax.addEventListener('polymer-response', this.handleLogin);
			ajax.addEventListener('polymer-error', this.failedToLogin);
			ajax.handleAs = 'json';

			console.log(ajax);
			ajax.go();



		}.bind(this));
	},
	handleLogin: function(event) {
		console.log(event.detail.response);
	},
	failedToLogin: function(event) {
		console.log('Failure');
		alert('Failed to login: \n'+event.detail.response);
	}
})