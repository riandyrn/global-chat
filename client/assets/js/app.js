const ENDPOINT = `ws://${window.location.host}/wsc`
const STATE_LOGIN_PAGE = 0
const STATE_CHAT_PAGE = 1
const colors = [
    "#3366CC", "#DC3912", "#FF9900", "#109618",
    "#990099", "#3B3EAC", "#0099C6", "#DD4477",
    "#66AA00", "#B82E2E", "#316395", "#994499",
    "#22AA99", "#AAAA11", "#6633CC", "#E67300",
    "#8B0707", "#329262", "#5574A6", "#3B3EAC"
];

var client = {ws: null, closeReason: null}

var app = new Vue({
    el: '#app',
    data: {
        currentPageState: STATE_LOGIN_PAGE,
        handleColors: {},
        messages: [],
        handle: "",
        messageInput: "",
    },
    methods: {
        submitHandle(){
            // connect to websocket
            connectWebsocket()
        },
        submitMessageInput(){
            // pack messageInput to {pub}
            var pubPkt = buildPubPkt(this.messageInput)
            // send message via websocket
            client.ws.send(pubPkt)
            // clear messageInput
            this.messageInput = ""
        },
        appendIncomingMessage(jsonObj){
            // parse jsonObj of {pres} & {data}
            var msg = {}
            if (jsonObj.data) {
                var from = jsonObj.data.from
                if(from === this.handle) {
                    from = "You"
                }
                msg = {
                    "what": "data", 
                    "from": from,
                    "content": jsonObj.data.content,
                    "ts": jsonObj.data.ts,
                    "isMessage": true
                }
            } else if (jsonObj.pres) {
                msg = {
                    "what": jsonObj.pres.what,
                    "from": jsonObj.pres.from,
                    "ts": jsonObj.pres.ts,
                    "isEvent": true
                }
            }
            // push the resulted page to messages
            this.messages.push(msg)
        },
        setPageState(state) {
            // change currentPageState & clearup some variables
            this.currentPageState = state
            switch(state) {
            case STATE_LOGIN_PAGE:
                this.handle = ""
                this.messages = []
                break
            case STATE_CHAT_PAGE:
                this.messageInput = ""
                break
            }
        },
        applyCurrentPageStateToDOM() {
            // necessary for updating the DOM as consequence of page state changing
            // will be called on mounted() & update()
            if (this.currentPageState === STATE_LOGIN_PAGE) {
                // set focus to handle input
                this.$refs.handle.focus()
            } else if (this.currentPageState === STATE_CHAT_PAGE) {
                // scroll to bottom
                var container = this.$refs.messages
                container.scrollTop = container.scrollHeight
                // set focus to messageInput
                this.$refs.messageInput.focus()
            }
        },
        formatTimestamp(timestamp) {
            return moment(timestamp, moment.ISO_8601).format("MMMM DD, YYYY - HH:mm")
        },
        getHandleColor(handle) {
            if(!this.handleColors[handle]) {
                if (handle === "You") {
                    color = 'green'
                } else {
                    color = colors[Math.floor(Math.random() * colors.length)]
                }
                this.handleColors[handle] = color
            }
            return this.handleColors[handle]
        },
        getContainerClass(what) {
            if(what === "left" || what === "join") {
                return "container-event-box"
            } else if(what === "data") {
                return "container-message-box"
            }
        }
    },
    computed: {
        shouldShowLoginPage() {
            return this.currentPageState === STATE_LOGIN_PAGE
        },
        shouldShowChatPage() {
            return this.currentPageState === STATE_CHAT_PAGE
        }
    },
    mounted() {
        // will be called once page initialized
        this.applyCurrentPageStateToDOM()
    },
    updated() {
        // will be called everytime data changed
        this.applyCurrentPageStateToDOM()
    }
})

function connectWebsocket() {
    client.ws = new WebSocket(ENDPOINT)
    client.ws.onopen = function(){
        // send {join} packet
        client.ws.send(buildJoinPkt(app.handle))
    }
    client.ws.onmessage = function(serverMsg) {
        // parse message from server
        var jsonObj
        try {jsonObj = JSON.parse(serverMsg.data)} catch(e){ return }

        // react accordingly
        switch(app.currentPageState) {
        case STATE_LOGIN_PAGE:
            // check for response of {join}
            if (jsonObj.ctrl && (jsonObj.ctrl.id === "join chat")) {
                var statusCode = jsonObj.ctrl.code
                if(statusCode >= 200 && statusCode < 300) {
                    // change page state
                    app.setPageState(STATE_CHAT_PAGE)
                } else {
                    // close connection
                    client.closeReason = "unable to join due: " + jsonObj.ctrl.err
                    client.ws.close()
                }
            }
            break
        case STATE_CHAT_PAGE:
            // forward {pres} & {data} to messages
            if (jsonObj.pres || jsonObj.data) {
                app.appendIncomingMessage(jsonObj)
            }
            break
        }
    }

    client.ws.onclose = function(e) {
        // show error
        closeReason = client.closeReason
        if (!closeReason) {
            closeReason = getReasonWebsocketError(e.code)
        }
        alert(closeReason)

        // change page state
        app.setPageState(STATE_LOGIN_PAGE)

        // reset client state
        client.closeReason = null
        client.ws = null
    }
}

function buildJoinPkt(handle) {
    return JSON.stringify({"join": {"id": "join chat", "handle": handle}})
}

function buildPubPkt(content) {
    return JSON.stringify({"pub":{"id": "publish message", "content": content}})
}

function getReasonWebsocketError(errCode) {
    var reason = 'Unknown error';
    switch(errCode) {
    case 1000:
        reason = 'Normal closure';
        break;
    case 1001:
        reason = 'An endpoint is going away';
        break;
    case 1002:
        reason = 'An endpoint is terminating the connection due to a protocol error.';
        break;
    case 1003:
        reason = 'An endpoint is terminating the connection because it has received a type of data it cannot accept';
        break;
    case 1004:
        reason = 'Reserved. The specific meaning might be defined in the future.';
        break;
    case 1005:
        reason = 'No status code was actually present';
        break;
    case 1006:
        reason = 'Cannot connect to server';
        break;
    case 1007:
        reason = 'The endpoint is terminating the connection because a message was received that contained inconsistent data';
        break;
    case 1008:
        reason = 'The endpoint is terminating the connection because it received a message that violates its policy';
        break;
    case 1009:
        reason = 'The endpoint is terminating the connection because a data frame was received that is too large';
        break;
    case 1010:
        reason = 'The client is terminating the connection because it expected the server to negotiate one or more extension, but the server didn\'t.';
        break;
    case 1011:
        reason = 'The server is terminating the connection because it encountered an unexpected condition that prevented it from fulfilling the request.';
        break;
    case 1012:
        reason = 'The server is terminating the connection because it is restarting';
        break;
    case 1013:
        reason = 'The server is terminating the connection due to a temporary condition';
        break;
    case 1015:
        reason = 'The connection was closed due to a failure to perform a TLS handshake';
        break;
    }
    return reason
}