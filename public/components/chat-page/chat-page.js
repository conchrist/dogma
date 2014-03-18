Polymer('chat-page', {
    handleLogin: function(event, data) {
        this.$.login.hidden = true;
        var room = this.$.room;
        room.username = data.name;
        room.userid = data._id;
        room.hidden = false;
        room.init();
    }
})