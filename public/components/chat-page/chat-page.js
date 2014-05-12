Polymer('chat-page', {
    handleLogin: function(event, data) {
        this.$.login.hidden = true;
        this.$.register.hidden = true;
        var room = this.$.room;
        room.username = data.name;
        room.userid = data._id;
        room.hidden = false;
        room.init();
    },
    register: function(event, data) {
        this.$.register.username = data.username;
        this.$.register.password = data.password;
        this.$.login.hidden = true;
        this.$.register.hidden = false;
    },
    registered: function(event, data) {
        this.handleLogin(event, data);
    },
    registrationAborted: function(event) {
        this.$.login.hidden = false;
        this.$.register.hidden = true;
    }
})