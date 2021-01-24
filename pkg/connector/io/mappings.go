package io

func ByteString(value interface{}) (interface{}, error) {
	return []byte(value.(string)), nil
}
