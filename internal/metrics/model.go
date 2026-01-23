package metrics

type Snapshot struct {
	Node      string  `json:"node"`
	CPU       float64 `json:"cpu"`
	Memory    float64 `json:"memory"`
	Disk      float64 `json:"disk"`
	Timestamp int64   `json:"timestamp"`
}
