Polymer('chat-register', {
    username: '',
    password: '',
    createUser: function(event) {
        event.preventDefault();
        var form = this.$.registerForm;
        //Uppercase HTTP method
        var ajax = this.$.registerajax;
        var username = this.$.username.value;
        var password = this.$.password.value;
        var params = {
            username: username,
            password: password
        };
        ajax.params = params;
        ajax.go();
    },
    handleRegister: function(event) {
        var user = event.detail.response;
        this.asyncFire('register', {
            _id: user._id,
            name: user.name
        });
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