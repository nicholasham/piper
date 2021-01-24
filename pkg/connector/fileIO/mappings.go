package fileIO

func ByteString(value interface{}) (interface{}, error) {
	return []byte(value.(string)), nil
}
