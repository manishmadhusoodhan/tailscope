package parser

import (
	"fmt"
	"strings"
)

func findString(data map[string]any, keys ...string) string {

	for _, wanted := range keys {

		for key, value := range data {

			if !strings.EqualFold(key, wanted) {
				continue
			}

			switch v := value.(type) {

			case string:
				return v

			default:
				return fmt.Sprint(v)
			}
		}
	}

	return ""
}
