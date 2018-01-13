package main

import (
	"fmt"
	"net/http"
	"strconv"
	"time"
)

// MsgClient represents client input
type MsgClient struct {
	Join *JoinPayload `json:"join,omitempty"`
	Pub  *PubPayload  `json:"pub,omitempty"`
}

// JoinPayload represents join command payload
type JoinPayload struct {
	ID     string `json:"id"`
	Handle string `json:"handle"`
}

// PubPayload represents publish command payload
type PubPayload struct {
	ID      string `json:"id"`
	Content string `json:"content"`
}

// MsgServer represents server response
type MsgServer struct {
	Ctrl *CtrlPayload `json:"ctrl,omitempty"`
	Pres *PresPayload `json:"pres,omitempty"`
	Data *DataPayload `json:"data,omitempty"`

	skipHandle string
}

// CtrlPayload represents status of user request
type CtrlPayload struct {
	ID         string    `json:"id,omitempty"`
	StatusCode int       `json:"code"`
	What       string    `json:"what,omitempty"`
	ErrCode    string    `json:"err,omitempty"`
	Timestamp  time.Time `json:"ts"`
}

// PresPayload represents important event
type PresPayload struct {
	What      string    `json:"what"`
	From      string    `json:"from"`
	Timestamp time.Time `json:"ts"`
}

// DataPayload is result of publish command
type DataPayload struct {
	From      string    `json:"from"`
	Content   string    `json:"content"`
	Timestamp time.Time `json:"ts"`
}

func timeNow() time.Time {
	// the idea of this function is always return non-zero
	// last digit of millisecond to keep resulted json-marshalled
	// RFC3339 string always retain its milliseconds value
	var oneBillion int64 = 1000 * 1000 * 1000
	tNano := time.Now().Round(time.Millisecond).UnixNano()
	tSeconds := tNano / oneBillion
	strNano, strSeconds := fmt.Sprintf("%v", tNano), fmt.Sprintf("%v", tSeconds)
	idxLastDigitMillis := len(strSeconds) - 1 + 3
	if strNano[idxLastDigitMillis] == '0' {
		strNano = strNano[:idxLastDigitMillis] + "1" + strNano[idxLastDigitMillis+1:]
		tNano, _ = strconv.ParseInt(strNano, 10, 64)
	}
	tRemainder := tNano % oneBillion
	now := time.Unix(tSeconds, tRemainder).UTC()
	return now
}

const (
	errMalformed = iota
	errCommandOutOfSequence
	errUnknown
	errAlreadyJoin
	errHandleTaken
)

func newCtrl(id string, code int, what string, err string, ts time.Time) *MsgServer {
	return &MsgServer{
		Ctrl: &CtrlPayload{
			ID:         id,
			StatusCode: code,
			What:       what,
			ErrCode:    err,
			Timestamp:  ts,
		},
	}
}

func resolveErrCode(errCode int) (statusCode int, errStr string) {
	switch errCode {
	case errAlreadyJoin:
		statusCode = http.StatusNotModified // 304
		errStr = "ERR_ALREADY_JOIN"
	case errMalformed:
		statusCode = http.StatusBadRequest // 400
		errStr = "ERR_MALFORMED"
	case errHandleTaken:
		statusCode = http.StatusConflict // 409
		errStr = "ERR_HANDLE_TAKEN"
	case errCommandOutOfSequence:
		statusCode = http.StatusConflict // 409
		errStr = "ERR_COMMAND_OUT_OF_SEQUENCE"
	case errUnknown:
		statusCode = http.StatusInternalServerError // 500
		errStr = "ERR_UNKNOWN"
	}
	return statusCode, errStr
}

func newErr(id string, errCode int, ts time.Time) *MsgServer {
	statusCode, errStr := resolveErrCode(errCode)
	return newCtrl(id, statusCode, "", errStr, ts)
}

// NoErr returns ok response
func NoErr(id string, what string, ts time.Time) *MsgServer {
	return newCtrl(id, http.StatusOK, what, "", ts)
}

// NoErrAccepted returns accepted response (request received, but still processed)
func NoErrAccepted(id string, what string, ts time.Time) *MsgServer {
	return newCtrl(id, http.StatusAccepted, what, "", ts)
}

// ErrMalformed returns bad request error response
func ErrMalformed(id string, ts time.Time) *MsgServer {
	return newErr(id, errMalformed, ts)
}

// ErrCommandOutOfSequence returns out of sequence error response
func ErrCommandOutOfSequence(id string, ts time.Time) *MsgServer {
	return newErr(id, errCommandOutOfSequence, ts)
}

// ErrUnknown returns internal server error response
func ErrUnknown(id string, ts time.Time) *MsgServer {
	return newErr(id, errUnknown, ts)
}

// ErrAlreadyJoin returns already join error response
func ErrAlreadyJoin(id string, ts time.Time) *MsgServer {
	return newErr(id, errAlreadyJoin, ts)
}

// ErrHandleTaken returns handle taken error response
func ErrHandleTaken(id string, ts time.Time) *MsgServer {
	return newErr(id, errHandleTaken, ts)
}
