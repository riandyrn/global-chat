# Global Chat Room API

## Joining Chat

To join chat in global room, user must first send `{join}` packet to server. 

```json
{
    "join": {
        "handle": "riandyrn"
    }
}
```

The value of `handle` need to be unique on the server. Otherwise it would be fail.

### Success Response:
```json
{
    "ctrl": {
        "code": 200,
        "what": "join"
    }
}
```

### Error Responses:

```json
{
    "ctrl": {
        "code": 304,
        "err": "ERR_ALREADY_JOIN"
    }
}
```

```json
{
    "ctrl": {
        "code": 409,
        "err": "ERR_HANDLE_TAKEN"
    }
}
```

## Sending Message

To send a message, user use `{pub}` packet to convey the message.

```json
{
    "pub": {
        "content": "Hello World!"
    }
}
```

### Success Response:
```json
{
    "ctrl": {
        "code": 200,
        "what": "pub"
    }
}
```

### Error Response:

Packet `{pub}` currently doesn't have special error response.

## Generic Errors

For every command, there are additional error responses. The responses are following:

```json
{
    "ctrl": {
        "code": 400,
        "err": "ERR_BAD_REQUEST"
    }
}
```

```json
{
    "ctrl": {
        "code": 409,
        "err": "ERR_COMMAND_OUT_OF_SEQUENCE"
    }
}
```

```json
{
    "ctrl": {
        "code": 500,
        "err": "ERR_UNKNOWN"
    }
}
```