package utils

import "fmt"

func MountURL(baseURL string, limit uint, offset uint) string {
	args := []interface{}{
		baseURL, limit, offset,
	}
	return fmt.Sprintf("%s?limit=%d&offset=%d", args...)
}
