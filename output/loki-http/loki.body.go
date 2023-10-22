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

var StdFormat = "{msg} line={line} func={func} traceID={traceID} spanID={spanID} req={req} rsp={rsp} fields={fields}"

func MakeLokiBody(env, app, instance string, data []*log_data.LogData) *LokiBody {
	streams := make([]*LokiStreamBody, len(data))
	for i, v := range data {
		text := StdFormat
		text = strings.ReplaceAll(text, "{msg}", v.Msg)
		text = strings.ReplaceAll(text, "{line}", v.Line)
		text = strings.ReplaceAll(text, "{func}", v.Func)
		text = strings.ReplaceAll(text, "{traceID}", v.TraceID)
		text = strings.ReplaceAll(text, "{spanID}", v.SpanID)
		text = strings.ReplaceAll(text, "{req}", v.Req)
		text = strings.ReplaceAll(text, "{rsp}", v.Rsp)
		text = strings.ReplaceAll(text, "{fields}", v.Fields)

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
