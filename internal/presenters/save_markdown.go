package presenters

import (
	"fmt"
	"sort"
	"stresstest/internal/usecase/run"
	"strings"
	"time"
)

func ToMarkdown(r run.RunOutputDTO) string {
	var markdown strings.Builder
	md := func(format string, a ...interface{}) {
		markdown.WriteString(fmt.Sprintf(format, a...))
		markdown.WriteString("\n")
	}

	// Parse timestamps
	layout := "2006-01-02 15:04:05.9999999"
	start, _ := time.Parse(layout, r.TimestampStart)
	end, _ := time.Parse(layout, r.TimestampEnd)
	duration := end.Sub(start).Seconds()

	// Header
	md("## 📊 Stress Test Report")
	md("**ID:** `%s`", r.Id)
	md("**URL:** %s", r.Url)
	md("**Requests:** %d", r.Requests)
	md("**Concurrency:** %d", r.Concurrency)
	md("**Start:** %s", start.Format("02/01/2006 15:04:05"))
	md("**End:** %s", end.Format("02/01/2006 15:04:05"))
	md("**Duration:** %.2f seconds", duration)

	// Separar os reports
	var status200 *run.StatusReportDTO
	var total *run.StatusReportDTO
	var others []run.StatusReportDTO

	for _, s := range r.Report {
		if s.Status == "total" {
			total = &s
		} else if s.Status == "200" {
			status200 = &s
		} else {
			others = append(others, s)
		}
	}

	// Total Summary
	md("\n### 📌 Total Summary")
	md("| Count | Min Time | Max Time | Total Time | Average Time |")
	md("|-------|----------|----------|------------|---------------|")
	if total != nil {
		md("| %d | %dms | %dms | %dms | %.2fms |", total.Count, total.MinTime, total.MaxTime, total.TotalTime, total.AverageTime)
	}

	// Status 200
	if status200 != nil {
		md("\n### ✅ Status 200")
		md("| Count | Min Time | Max Time | Total Time | Average Time |")
		md("|-------|----------|----------|------------|---------------|")
		md("| %d | %dms | %dms | %dms | %.2fms |", status200.Count, status200.MinTime, status200.MaxTime, status200.TotalTime, status200.AverageTime)
	} else {
		md("\n⚠️ Nenhuma requisição com status 200")
	}

	// Outros Status
	if len(others) > 0 {
		md("\n### 📦 Outros Status")
		md("| Status | Count | Min Time | Max Time | Total Time | Average Time |")
		md("|--------|-------|----------|----------|------------|---------------|")
		sort.Slice(others, func(i, j int) bool {
			return others[i].Status < others[j].Status
		})
		for _, s := range others {
			md("| %s | %d | %dms | %dms | %dms | %.2fms |", s.Status, s.Count, s.MinTime, s.MaxTime, s.TotalTime, s.AverageTime)
		}
	}

	return markdown.String()
}
