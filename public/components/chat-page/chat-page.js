Polymer('chat-page', {
    page: 'login',
    handleLogin: function(event, data) {
        var room = this.$.room;
        room.username = data.name;
        room.userid = data._id;
        room.init();
        this.page = 'room';
    },
    register: function(event, data) {
        this.$.register.username = data.username;
        this.$.register.password = data.password;
        this.page = 'register';
    },
    registered: function(event, data) {
        this.handleLogin(event, data);
    },
    registrationAborted: function(event) {
        this.page = 'login';
    }
});