new Vue({
    el: '#app',

    data: {
        ws: null, // Our websocket
        newMsg: '', // Holds new messages to be sent to the server
        chatContent: '', // A running list of chat messages displayed on the screen
        // email: null, // Email address used for grabbing an avatar
        timestamp: 'placeholder',
        action: 'msg',
        username: null, // Our username
        joined: false // True if email and username have been filled in
    },

    created: function() {
        var self = this;
        this.ws = new WebSocket('ws://' + window.location.host + '/ws');
        // Add code to retrieve all previous messages here and display them
        this.ws.addEventListener('message', function(e) {
            var msg = JSON.parse(e.data);
            self.chatContent += '<div class="chip">'
                    + '<img src="http://api.adorable.io/avatar/32/' + msg.username + '">' // Avatar
                    + msg.username
                + '</div>'
                + emojione.toImage(msg.contents) + '<br/>'; // Parse emojis
                // Uses https://github.com/adorableio/avatars-api for random avatar generation
            var element = document.getElementById('chat-messages');
            element.scrollTop = element.scrollHeight; // Auto scroll to the bottom
        });
    },

    methods: {
    send: function () {
        if (this.newMsg != '') {
            this.ws.send(
                JSON.stringify({
                    timestamp: this.timestamp,
                    username: this.username,
                    action: this.action,
                    contents: $('<p>').html(this.newMsg).text() // Strip out html
                }
            ));
            this.newMsg = ''; // Reset newMsg
        }
    },

    join: function () {
    if (!this.username) {
        Materialize.toast('You must choose a username', 2000);
        return
    }
    this.username = $('<p>').html(this.username).text();
    this.joined = true;
    },
  }
});
