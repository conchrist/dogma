Polymer('chat-register', {
    username: '',
    password: '',
    createUser: function(event) {
        event.preventDefault();
        var ajax = document.createElement('polymer-ajax');
        var form = this.$.registerForm;
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
        ajax.addEventListener('polymer-response', this.handleRegister.bind(this));
        ajax.addEventListener('polymer-error', this.failedToRegister.bind(this));
        ajax.go();
    },
    handleRegister: function(event) {
        var user = event.detail.response;
        this.asyncFire('register', {
            _id: user._id,
            name: user.name
        })
    },
    failedToRegister: function(event) {
        var message = event.detail.response;
        alert('Failed to create user: ' + message);
    },
    goBack: function(event) {
        event.preventDefault();
        this.fire('back');
    }
});