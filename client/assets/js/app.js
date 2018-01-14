const ENDPOINT = "ws://127.0.0.1:8192/wsc"
const STATE_LOGIN_PAGE = 0
const STATE_CHAT_PAGE = 1
const colors = [
    "#3366CC", "#DC3912", "#FF9900", "#109618",
    "#990099", "#3B3EAC", "#0099C6", "#DD4477",
    "#66AA00", "#B82E2E", "#316395", "#994499",
    "#22AA99", "#AAAA11", "#6633CC", "#E67300",
    "#8B0707", "#329262", "#5574A6", "#3B3EAC"
];

var ws = null
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
            // send message via websocket
        },
        appendIncomingMessage(jsonObj){
            // parse jsonObj
            // push the resulted page to messages
        },
        setPageState(state) {
            // change currentPageState & clearup some variables
            this.currentPageState = state
            switch(state) {
            case STATE_LOGIN_PAGE:
                this.handle = ""
                this.messages = {}
            case STATE_CHAT_PAGE:
                this.messageInput = ""
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
                if (handle === this.handle) {
                    color = 'green'
                } else {
                    color = colors[Math.floor(Math.random() * colors.length)]
                }
                this.handleColors[handle] = color
            }
            return this.handleColors[handle]
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
    ws = new WebSocket(ENDPOINT)
    ws.onopen = function(){
        // send {join} packet
        app.setPageState(STATE_CHAT_PAGE)
    }
}