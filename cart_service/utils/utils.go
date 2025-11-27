package utils

// To get timestamps to fill created_at and updated_at fields
import (
	"time"
)

func GetCurrentTimestamp() string {
	return time.Now().Format(time.RFC3339)
}
