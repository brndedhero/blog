package helpers

import (
	"bytes"
	"encoding/json"
	"errors"
)

func PrepareString(statusCode int, message []byte) (string, error) {
	var buffer bytes.Buffer
	switch statusCode {
	case 200:
		buffer.WriteString(`{"Status": "OK", "StatusCode": 200, "Message":`)
		buffer.Write(message)
		buffer.WriteString(`}`)

		return buffer.String(), nil
	case 201:
		buffer.WriteString(`{"Status": "Created", "StatusCode": 201, "Message":`)
		buffer.Write(message)
		buffer.WriteString(`}`)

		return buffer.String(), nil
	case 404:
		buffer.WriteString(`{"Status": "Not Found", "StatusCode": 404}`)

		return buffer.String(), nil
	case 405:
		buffer.WriteString(`{"Status": "Method Not Allowed", "StatusCode": 405}`)

		return buffer.String(), nil
	default:
		return "", errors.New("status code is not implemented")
	}
}

func PrepareErrorString(statusCode int, err error) string {
	var buffer bytes.Buffer
	switch statusCode {
	case 500:
		buffer.WriteString(`{"Status": "Internal Server Error", "StatusCode": 500, "Message": "`)
		buffer.WriteString(err.Error())
		buffer.WriteString(`"}`)

		return buffer.String()
	default:
		return ""
	}
}

func PrepareCreateJson(id int, rowsAffected int) ([]byte, error) {
	m := make(map[string]int)
	m["id"] = id
	m["RowsAffected"] = rowsAffected
	data, err := json.Marshal(m)
	return data, err
}
