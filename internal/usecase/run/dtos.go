package run

type RunInputDTO struct {
	Url         string `json:"url"`
	Requests    int    `json:"requests"`
	Concurrency int    `json:"concurrency"`
	ShowData    bool   `json:"show_data"`
}

type RunOutputDTO struct {
	Id                    string            `json:"id"`
	Url                   string            `json:"url"`
	Requests              int               `json:"requests"`
	Concurrency           int               `json:"concurrency"`
	TimestampStart        string            `json:"timestamp_start"`
	TimestampEnd          string            `json:"timestamp_end"`
	TestDurationInSeconds int               `json:"test_duration_in_seconds"`
	Data                  []DataOutputDTO   `json:"data"`
	Report                []StatusReportDTO `json:"report"`
}

type DataOutputDTO struct {
	StatusCode            int    `json:"status_code"`
	DurationInMs          int    `json:"duration_in_ms"`
	RequestStartTimestamp string `json:"request_start_timestamp"`
	RequestEndTimestamp   string `json:"request_end_timestamp"`
}

type StatusReportDTO struct {
	Status      string  `json:"status"`
	Count       int     `json:"count"`
	MinTime     int     `json:"min_time_in_ms"`
	MaxTime     int     `json:"max_time_in_ms"`
	TotalTime   int     `json:"total_time_in_ms"`
	AverageTime float64 `json:"average_time_in_ms"`
}
