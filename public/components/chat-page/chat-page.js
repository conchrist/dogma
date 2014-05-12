Polymer('chat-page', {
    page: 'login',
    ready: function() {

    },
    checkLoggedIn: function(event, res) {
        var data = res.response;
        if(data.loggedIn) {
            this.goToChatRoom({
                "_id": data._id,
                "name": data.name
            });
        }
    },
    goToChatRoom: function(user) {
        var room = this.$.room;
        room.username = user.name;
        room.userid = user._id;
        room.init();
        this.page = 'room';
    },
    handleLogin: function(event, data) {
        this.goToChatRoom(data);
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
    },
    loggedOut: function() {
        this.page = 'login';
    }
});