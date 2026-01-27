package hub

type podMetric struct {
	Name        string `json:"name"`
	Namespace   string `json:"namespace"`
	CPUMilli    int64  `json:"cpuMilli"`
	MemoryBytes int64  `json:"memoryBytes"`
	Timestamp   int64  `json:"timestamp"`
}

type deploymentMetric struct {
	Name           string `json:"name"`
	Replicas       int32  `json:"replicas"`
	Available      int32  `json:"available"`
	CPUMilli       int64  `json:"cpuMilli"`
	MemoryBytes    int64  `json:"memoryBytes"`
	PodCount       int    `json:"podCount"`
	MissingMetrics int    `json:"missingMetrics"`
}
