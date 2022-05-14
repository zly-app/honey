package http

import (
	"strconv"
	"strings"

	jsoniter "github.com/json-iterator/go"

	"github.com/zly-app/honey/log_data"
)

type LokiStreamBody struct {
	Stream map[string]string `json:"stream"`
	Values [][]string        `json:"values"`
}
type LokiBody struct {
	Streams []*LokiStreamBody `json:"streams"`
}

func MakeLokiBody(env, app, instance string, data []*log_data.LogData) *LokiBody {
	streams := make([]*LokiStreamBody, len(data))
	for i, v := range data {
		msg := map[string]string{
			"msg":      v.Msg,
			"fields":   v.Fields,
			"line":     v.Line,
			"trace_id": v.TraceID,
		}
		msgData, _ := jsoniter.ConfigCompatibleWithStandardLibrary.MarshalToString(msg)
		streamBody := &LokiStreamBody{
			Stream: map[string]string{
				"env":      env,
				"app":      app,
				"instance": instance,
				"level":    strings.ToLower(v.Level),
			},
			Values: [][]string{
				{strconv.FormatInt(v.T, 10), msgData},
			},
		}
		streams[i] = streamBody
	}

	body := &LokiBody{
		Streams: streams,
	}
	return body
}
