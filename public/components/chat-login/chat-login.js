Polymer('chat-login', {
    login: function(event) {
        event.preventDefault();
        var ajax = this.$.loginajax;
        var username = this.$.username.value;
        var password = this.$.password.value;
        var params = {
            username: username,
            password: password
        };
        ajax.params = params;
        ajax.go();
    },
    register: function(event) {
        event.preventDefault();
        this.asyncFire('new-user', {
            username: this.$.username.value,
            password: this.$.password.value
        });
    },
    handleLogin: function(event) {
        console.log('Logged in', event.detail.response);
        var user = event.detail.response;
        this.asyncFire('login', {
            _id: user._id,
            name: user.name
        });
    },
    failedToLogin: function(event) {
        alert('Failed to login: \n' + event.detail.response);
    }
});