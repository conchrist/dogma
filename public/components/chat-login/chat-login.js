Polymer('chat-login', {
	ready: function () {
		var form = this.$.loginForm;
		form.addEventListener('submit', function(event) {
			event.preventDefault();
			var ajax = document.createElement('polymer-ajax');
			ajax.url = form.action;
			//Uppercase HTTP method
			ajax.method = String(form.method).toUpperCase();
			var username = this.$.username.value;
			var password = this.$.password.value;
			var params = {
				username: username,
				password: password
			};
			ajax.params = params;
			ajax.handleAs = 'json';
			ajax.addEventListener('polymer-response', this.handleLogin);
			ajax.addEventListener('polymer-error', this.failedToLogin);
			ajax.go();
		}.bind(this));
	},
	handleLogin: function(event) {
		console.log('Logged in', event.detail.response);
	},
	failedToLogin: function(event) {
		alert('Failed to login: \n'+event.detail.response);
	}
})