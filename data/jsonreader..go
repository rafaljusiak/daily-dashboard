package data

import (
	"encoding/json"
    "log"
)

func ReadJson(bytes []byte) (map[string]interface{}, error) {
	var output map[string]interface{}

	if err := json.Unmarshal(bytes, &output); err != nil {
		log.Printf("Invalid JSON: %v", err)
		if e, ok := err.(*json.SyntaxError); ok {
			log.Printf("Syntax error at byte offset %d", e.Offset)
		}
		log.Printf("Response was: %s", string(bytes))
		return nil, err
	}
	
	return output, nil
}
