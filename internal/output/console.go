package output

import (
	"fmt"
	"time"

	"tailscope/internal/model"
)

func Print(
	entry model.LogEntry,
) {

	timestamp := entry.Timestamp

	if timestamp.IsZero() {
		timestamp = time.Now()
	}

	fmt.Printf(
		"%s %-8s %-45s | %s",
		timestamp.Format("15:04:05"),
		entry.Level,
		entry.Pod,
		entry.Message,
	)

	if entry.TraceID != "" {

		fmt.Printf(
			" | traceId=%s",
			entry.TraceID,
		)
	}

	fmt.Println()
}

func PrintError(
	err error,
) {

	if err == nil {
		return
	}

	fmt.Printf(
		"%s ERROR: %s\n",
		time.Now().Format("15:04:05"),
		err.Error(),
	)
}
