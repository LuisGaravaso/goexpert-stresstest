package run

import (
	"context"
	"net/http"
	"strconv"
	"stresstest/internal/entity"
	"stresstest/internal/repository"
	"sync"
	"time"
)

type RunUseCase struct {
	repo repository.RepositoryInterface
}

func NewRunUseCase(repo repository.RepositoryInterface) RunUseCase {
	return RunUseCase{repo: repo}
}

func (u *RunUseCase) Run(ctx context.Context, input RunInputDTO) (RunOutputDTO, error) {

	// Validate input
	testOpts := &entity.TestRunOptions{Requests: input.Requests, Concurrency: input.Concurrency}
	testRun, err := entity.NewTestRun(input.Url, testOpts)
	if err != nil {
		return RunOutputDTO{}, err
	}

	// Won't actually save anything, just a placeholder for future implementations
	err = u.repo.Save(ctx, testRun)
	if err != nil {
		return RunOutputDTO{}, err
	}

	// Run the Stress Test
	var wg sync.WaitGroup
	var mu sync.Mutex
	data := make([]DataOutputDTO, 0)
	reportMap := make(map[string]*StatusReportDTO)
	requestsChannel := make(chan struct{}, testRun.Concurrency)

	for i := 0; i < testRun.Requests; i++ {
		requestsChannel <- struct{}{}
		wg.Add(1)

		go func() {
			defer wg.Done()
			defer func() { <-requestsChannel }()

			status, duration, requestStart, requestEnd := MakeRequest(ctx, testRun.Url)

			// Save data if requested
			if input.ShowData {
				mu.Lock()
				data = append(data, DataOutputDTO{
					StatusCode:            status,
					DurationInMs:          duration,
					RequestStartTimestamp: FormatTimeToUTCString(requestStart),
					RequestEndTimestamp:   FormatTimeToUTCString(requestEnd),
				})
				mu.Unlock()
			}

			// Save report data
			updateReport(&mu, reportMap, strconv.Itoa(status), duration)
			updateReport(&mu, reportMap, "total", duration) // 9999 = total for all statuses
		}()
	}

	wg.Wait() // Wait for all requests to finish

	// Calculate average time
	var FinalReport []StatusReportDTO
	for _, report := range reportMap {
		report.AverageTime = float64(report.TotalTime) / float64(report.Count)
		FinalReport = append(FinalReport, *report)
	}

	// Return output
	return RunOutputDTO{
		Id:                    testRun.Id,
		Url:                   testRun.Url,
		Requests:              testRun.Requests,
		Concurrency:           testRun.Concurrency,
		TimestampStart:        FormatTimeToUTCString(testRun.Timestamp),
		TimestampEnd:          FormatTimeToUTCString(time.Now()),
		TestDurationInSeconds: int(time.Since(testRun.Timestamp).Seconds()),
		Data:                  data,
		Report:                FinalReport,
	}, nil
}

// MakeRequest makes a GET request to the given URL and returns the status code, duration, and start/end times
func MakeRequest(ctx context.Context, url string) (status, duration int, start, end time.Time) {
	start = time.Now()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return 0, 0, start, time.Now()
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return 0, 0, start, time.Now()
	}
	defer resp.Body.Close()

	end = time.Now()
	duration = int(time.Since(start).Milliseconds())
	status = resp.StatusCode
	return status, duration, start, end
}

// FormatTimeToUTCString formats a time.Time to UTC in the format "YYYY-MM-DD HH:MM:SS.sssssss"
func FormatTimeToUTCString(t time.Time) string {
	return t.UTC().Format("2006-01-02 15:04:05.0000000")
}

// updateReport updates the report map with the new data
func updateReport(mu *sync.Mutex, reportMap map[string]*StatusReportDTO, status string, duration int) {
	mu.Lock()
	defer mu.Unlock()

	report, exists := reportMap[status]
	if !exists {
		report = &StatusReportDTO{Status: status, MinTime: duration, MaxTime: duration}
		reportMap[status] = report
	}
	report.Count++
	report.TotalTime += duration
	if duration < report.MinTime {
		report.MinTime = duration
	}
	if duration > report.MaxTime {
		report.MaxTime = duration
	}
}
