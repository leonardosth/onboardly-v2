package dashboard

import (
	"math"
	"time"

	"onboardly-backend/internal/activity"
	"onboardly-backend/internal/db"
)

// DashboardData holds the aggregated metrics, historical trend, and activity feed.
type DashboardData struct {
	Metrics          Metrics             `json:"metrics"`
	History          []HistoryItem       `json:"history"`
	RecentActivities []activity.Activity `json:"recent_activities"`
}

type Metrics struct {
	ActivationRate float64 `json:"activation_rate"`
	NoShowRate     float64 `json:"no_show_rate"`
}

type HistoryItem struct {
	Month       string `json:"month"`
	Deployments int    `json:"deployments"`
}

// GetDashboardData aggregates DB metrics.
func GetDashboardData() (*DashboardData, error) {
	// 1. Calculate Project Activation Rate: (Go-Live projects / Total projects) * 100
	var totalProjects, goLiveProjects int
	err := db.DB.QueryRow(`SELECT COUNT(*) FROM projects`).Scan(&totalProjects)
	if err != nil {
		return nil, err
	}
	err = db.DB.QueryRow(`SELECT COUNT(*) FROM projects WHERE status = 'Go-Live'`).Scan(&goLiveProjects)
	if err != nil {
		return nil, err
	}

	activationRate := 0.0
	if totalProjects > 0 {
		activationRate = float64(goLiveProjects) / float64(totalProjects) * 100
		activationRate = math.Round(activationRate*10) / 10
	}

	// 2. Calculate Meeting No-Show Rate: (No-Show meetings / Total meetings) * 100
	var totalMeetings, noShowMeetings int
	err = db.DB.QueryRow(`SELECT COUNT(*) FROM meetings`).Scan(&totalMeetings)
	if err != nil {
		return nil, err
	}
	err = db.DB.QueryRow(`SELECT COUNT(*) FROM meetings WHERE no_show = TRUE`).Scan(&noShowMeetings)
	if err != nil {
		return nil, err
	}

	noShowRate := 0.0
	if totalMeetings > 0 {
		noShowRate = float64(noShowMeetings) / float64(totalMeetings) * 100
		noShowRate = math.Round(noShowRate*10) / 10
	}

	// 3. Historical Monthly Deployments (last 6 months)
	historyQuery := `
		SELECT TO_CHAR(updated_at, 'YYYY-MM') AS month, COUNT(id) AS count
		FROM projects
		WHERE status = 'Go-Live' AND updated_at >= NOW() - INTERVAL '6 months'
		GROUP BY month
		ORDER BY month ASC
	`
	rows, err := db.DB.Query(historyQuery)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var history []HistoryItem
	for rows.Next() {
		var item HistoryItem
		if err := rows.Scan(&item.Month, &item.Deployments); err == nil {
			history = append(history, item)
		}
	}

	// Populate fallback empty history items if none exist
	if len(history) == 0 {
		for i := 5; i >= 0; i-- {
			t := time.Now().AddDate(0, -i, 0)
			history = append(history, HistoryItem{
				Month:       t.Format("2006-01"),
				Deployments: 0,
			})
		}
	}

	// 4. Activity Feed
	recent, err := activity.GetRecentActivities()
	if err != nil {
		recent = []activity.Activity{}
	}

	return &DashboardData{
		Metrics: Metrics{
			ActivationRate: activationRate,
			NoShowRate:     noShowRate,
		},
		History:          history,
		RecentActivities: recent,
	}, nil
}
