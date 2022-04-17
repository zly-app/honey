package log_data

// 收集数据
type CollectData struct {
	Service  string `json:"service"`
	Instance string `json:"instance"`
	*LogData
}

func MakeCollectData(service, instance string, log *LogData) *CollectData {
	return &CollectData{
		Service:  service,
		Instance: instance,
		LogData:  log,
	}
}
