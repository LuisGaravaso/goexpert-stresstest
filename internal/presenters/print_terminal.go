package presenters

import (
	"fmt"
	"sort"
	"stresstest/internal/usecase/run"
	"time"

	"github.com/fatih/color"
)

func PrintReport(r run.RunOutputDTO) {
	cyan := color.New(color.FgCyan).SprintFunc()
	green := color.New(color.FgGreen).SprintFunc()
	red := color.New(color.FgRed).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()
	bold := color.New(color.Bold).SprintFunc()

	// Parse timestamps
	layout := "2006-01-02 15:04:05.9999999"
	start, _ := time.Parse(layout, r.TimestampStart)
	end, _ := time.Parse(layout, r.TimestampEnd)
	duration := end.Sub(start).Seconds()

	// ========== HEADER ==========
	fmt.Println(bold("📊 Stress Test Report"))
	fmt.Println("ID:         ", cyan(r.Id))
	fmt.Println("URL:        ", r.Url)
	fmt.Println("Requests:   ", r.Requests)
	fmt.Println("Concurrency:", r.Concurrency)
	fmt.Println("Start:      ", start.Format("02/01/2006 15:04:05"))
	fmt.Println("End:        ", end.Format("02/01/2006 15:04:05"))
	fmt.Printf("Duration:   %.2f seconds\n", duration)

	// ========== SEPARAR REPORTS ==========
	var status200 *run.StatusReportDTO
	var total *run.StatusReportDTO
	var others []run.StatusReportDTO

	for _, s := range r.Report {
		switch s.Status {
		case "total":
			total = &s
		case "200":
			status200 = &s
		default:
			others = append(others, s)
		}
	}

	// ========== TOTAL ==========
	if total != nil {
		fmt.Println()
		fmt.Println(bold("📌 Total Summary:"))
		printStatusLine(*total, bold)
	}

	// ========== STATUS 200 ==========
	if status200 != nil {
		fmt.Println()
		fmt.Println(bold("✅ Status 200"))
		printStatusLine(*status200, green)
	} else {
		fmt.Println()
		fmt.Println(bold("⚠️ Nenhuma requisição com status 200"))
	}

	// ========== OUTROS STATUS ==========
	if len(others) > 0 {
		fmt.Println()
		fmt.Println(bold("📦 Outros Status"))
		sort.Slice(others, func(i, j int) bool {
			return others[i].Status < others[j].Status
		})
		for _, s := range others {
			colorFunc := yellow
			if s.Status[0] == '5' {
				colorFunc = red
			}
			printStatusLine(s, colorFunc)
		}
	}
}

func printStatusLine(s run.StatusReportDTO, colorFunc func(a ...interface{}) string) {
	fmt.Printf("%s | Count: %s | Min: %dms | Max: %dms | Total: %dms | Avg: %.2fms\n",
		colorFunc("→ Status "+s.Status),
		colorFunc(fmt.Sprintf("%d", s.Count)),
		s.MinTime,
		s.MaxTime,
		s.TotalTime,
		s.AverageTime,
	)
}
