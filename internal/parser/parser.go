package parser

import (
	"encoding/json"
	"strconv"
	"strings"
	"time"

	"tailscope/internal/model"
)

type Parser struct {
}

func New() *Parser {
	return &Parser{}
}

func (p *Parser) Parse(raw model.RawLog) model.LogEntry {

	entry := model.LogEntry{

		Timestamp: time.Now(),

		Namespace: raw.Namespace,

		Pod: raw.Pod,

		Raw: strings.TrimSpace(raw.Line),

		Fields: make(map[string]string),
	}

	var jsonData map[string]any

	err := json.Unmarshal(
		[]byte(entry.Raw),
		&jsonData,
	)

	// Plain text log
	if err != nil {

		entry.Message = entry.Raw

		return entry
	}

	for key, value := range jsonData {

		switch v := value.(type) {

		case string:

			entry.Fields[key] = v

		case float64:

			entry.Fields[key] =
				strconv.FormatFloat(
					v,
					'f',
					-1,
					64,
				)

		case bool:

			entry.Fields[key] =
				strconv.FormatBool(v)
		}
	}

	entry.Level = findString(
		jsonData,
		"Level",
		"level",
		"Severity",
		"severity",
		"logLevel",
	)

	entry.Message = findString(
		jsonData,
		"Message",
		"message",
		"msg",
	)

	entry.TraceID = findString(
		jsonData,
		"TraceId",
		"TraceID",
		"traceId",
		"traceID",
	)

	entry.CorrelationID = findString(
		jsonData,
		"CorrelationId",
		"CorrelationID",
		"correlationId",
	)

	entry.RequestID = findString(
		jsonData,
		"RequestId",
		"RequestID",
		"requestId",
	)

	if entry.Message == "" {
		entry.Message = entry.Raw
	}

	return entry
}
