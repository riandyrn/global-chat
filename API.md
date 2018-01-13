# Global Chat Room API

## Joining Chat

To join chat in global room, user must first send `{join}` packet to server. 

```json
{
    "join": {
        "id": "join chat",
        "handle": "riandyrn"
    }
}
```

The value of `handle` need to be unique on the server. Otherwise it would be fail.

### Success Response:
```json
{
    "ctrl": {
        "id": "join chat",
        "code": 200,
        "what": "join",
        "ts": "2018-01-12T23:40:12.426Z"
    }
}
```

### Error Responses:

```json
{
    "ctrl": {
        "id": "join chat",
        "code": 304,
        "err": "ERR_ALREADY_JOIN",
        "ts": "2018-01-12T23:40:12.426Z"
    }
}
```

```json
{
    "ctrl": {
        "id": "join chat",
        "code": 409,
        "err": "ERR_HANDLE_TAKEN",
        "ts": "2018-01-12T23:40:12.426Z"
    }
}
```

### Other Users Reponse:

When user is successfully joining the chat, other users will receive following packet:

```json
{
    "pres": {
        "what": "join",
        "from": "riandyrn",
        "ts": "2018-01-12T23:40:12.426Z"
    }
}
```

## Sending Message

To send a message, user use `{pub}` packet to convey the message.

```json
{
    "pub": {
        "id": "publish message",
        "content": "Hello World!"
    }
}
```

### Success Response:
```json
{
    "ctrl": {
        "id": "publish message",
        "code": 202,
        "what": "pub",
        "ts": "2018-01-12T23:40:12.426Z"
    }
}
```

### Error Response:

Packet `{pub}` currently doesn't have special error response.

### Other Users Response:

When sender is successfully sending a message, then all users (including sender) will receive following packet:

```json
{
    "data": {
        "from": "riandyrn",
        "content": "Hello World!",
        "ts": "2018-01-12T23:40:12.426Z"
    }
}
```

## Leaving Chat

To leave chat user's client could just destroy the session.

### Other Users Response:

When user's session is destroyed, other users will receive following packet:

```json
{
    "pres": {
        "what": "left",
        "from": "riandyrn",
        "ts": "2018-01-12T23:40:12.426Z"
    }
}
```

## Generic Errors

For every command, there are additional error responses. The responses are following:

```json
{
    "ctrl": {
        "code": 400,
        "err": "ERR_BAD_REQUEST",
        "ts": "2018-01-12T23:40:12.426Z"
    }
}
```

```json
{
    "ctrl": {
        "code": 409,
        "err": "ERR_COMMAND_OUT_OF_SEQUENCE",
        "ts": "2018-01-12T23:40:12.426Z"
    }
}
```

```json
{
    "ctrl": {
        "code": 500,
        "err": "ERR_UNKNOWN",
        "ts": "2018-01-12T23:40:12.426Z"
    }
}
```