package http

import (
	"strconv"
	"strings"

	"github.com/zly-app/honey/log_data"
)

type LokiStreamBody struct {
	Stream map[string]string `json:"stream"`
	Values [][]string        `json:"values"`
}
type LokiBody struct {
	Streams []*LokiStreamBody `json:"streams"`
}

var StdFormat = "{msg} line={line} fields={fields} traceID={traceID} spanID={spanID}"

func MakeLokiBody(env, app, instance string, data []*log_data.LogData) *LokiBody {
	streams := make([]*LokiStreamBody, len(data))
	for i, v := range data {
		text := StdFormat
		text = strings.ReplaceAll(text, "{msg}", v.Msg)
		text = strings.ReplaceAll(text, "{fields}", v.Fields)
		text = strings.ReplaceAll(text, "{line}", v.Line)
		text = strings.ReplaceAll(text, "{traceID}", v.TraceID)
		text = strings.ReplaceAll(text, "{spanID}", v.SpanID)

		streamBody := &LokiStreamBody{
			Stream: map[string]string{
				"env":      env,
				"app":      app,
				"instance": instance,
				"level":    strings.ToLower(v.Level),
			},
			Values: [][]string{
				{strconv.FormatInt(v.T, 10), text},
			},
		}
		streams[i] = streamBody
	}

	body := &LokiBody{
		Streams: streams,
	}
	return body
}
