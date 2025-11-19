package logger

import (
	"encoding/json"
	"log"
)

// Выводит структуру в красивом JSON формате
func PrettyStructurePrint(prefix string, v interface{}) {
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		log.Printf("❌ Failed to marshal %s: %v", prefix, err)
		return
	}

	log.Printf("%s\n%s\n", prefix, string(b))
}
