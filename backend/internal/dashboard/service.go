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
	Funnel           FunnelData          `json:"funnel"`
	Cohorts          []CohortItem        `json:"cohorts"`
	History          []HistoryItem       `json:"history"`
	RecentActivities []activity.Activity `json:"recent_activities"`
}

type Metrics struct {
	ActivationRate              float64 `json:"activation_rate"`
	NoShowRate                  float64 `json:"no_show_rate"`
	AbandonmentRate             float64 `json:"abandonment_rate"`
	FirstMeetingActivationRate  float64 `json:"first_meeting_activation_rate"`
	Activation30dRate           float64 `json:"activation_30d_rate"`
}

type FunnelData struct {
	Registered   int `json:"registered"`
	Participants int `json:"participants"`
	Active       int `json:"active"`
}

type CohortItem struct {
	Month          string  `json:"month"`
	Total          int     `json:"total"`
	Activated      int     `json:"activated"`
	ActivationRate float64 `json:"activation_rate"`
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

	// 3. Funnel: Registered > Participants > Active
	funnel, err := getFunnelData()
	if err != nil {
		return nil, err
	}

	// 4. Abandonment Rate
	abandonmentRate, err := getAbandonmentRate()
	if err != nil {
		abandonmentRate = 0
	}

	// 5. 30-Day Activation Rate
	activation30d, err := getActivation30dRate()
	if err != nil {
		activation30d = 0
	}

	// 6. First-Meeting Activation Rate
	firstMeetingRate, err := getFirstMeetingActivationRate()
	if err != nil {
		firstMeetingRate = 0
	}

	// 7. Cohort Data
	cohorts, err := getCohortData()
	if err != nil {
		cohorts = []CohortItem{}
	}

	// 8. Historical Monthly Deployments (last 6 months)
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

	// 9. Activity Feed
	recent, err := activity.GetRecentActivities()
	if err != nil {
		recent = []activity.Activity{}
	}

	return &DashboardData{
		Metrics: Metrics{
			ActivationRate:              activationRate,
			NoShowRate:                  noShowRate,
			AbandonmentRate:             abandonmentRate,
			FirstMeetingActivationRate:  firstMeetingRate,
			Activation30dRate:           activation30d,
		},
		Funnel:           funnel,
		Cohorts:          cohorts,
		History:          history,
		RecentActivities: recent,
	}, nil
}

// getFunnelData returns funnel counts: Registered > Participants > Active.
func getFunnelData() (FunnelData, error) {
	var f FunnelData

	err := db.DB.QueryRow(`SELECT COUNT(*) FROM projects`).Scan(&f.Registered)
	if err != nil {
		return f, err
	}

	err = db.DB.QueryRow(`SELECT COUNT(DISTINCT m.project_id) FROM meetings m WHERE m.status = 'completed'`).Scan(&f.Participants)
	if err != nil {
		return f, err
	}

	err = db.DB.QueryRow(`SELECT COUNT(*) FROM projects WHERE activated_at IS NOT NULL`).Scan(&f.Active)
	if err != nil {
		return f, err
	}

	return f, nil
}

// getAbandonmentRate returns the percentage of projects with no completed meetings.
func getAbandonmentRate() (float64, error) {
	var total, abandoned int
	err := db.DB.QueryRow(`SELECT COUNT(*) FROM projects`).Scan(&total)
	if err != nil {
		return 0, err
	}

	err = db.DB.QueryRow(`
		SELECT COUNT(*) FROM projects p
		WHERE NOT EXISTS (
			SELECT 1 FROM meetings m WHERE m.project_id = p.id AND m.status = 'completed'
		)
	`).Scan(&abandoned)
	if err != nil {
		return 0, err
	}

	if total == 0 {
		return 0, nil
	}
	rate := float64(abandoned) / float64(total) * 100
	return math.Round(rate*10) / 10, nil
}

// getActivation30dRate returns the percentage of participating clients activated within 30 calendar days.
func getActivation30dRate() (float64, error) {
	var participants, activatedIn30d int

	err := db.DB.QueryRow(`
		SELECT COUNT(DISTINCT m.project_id) FROM meetings m WHERE m.status = 'completed'
	`).Scan(&participants)
	if err != nil {
		return 0, err
	}

	err = db.DB.QueryRow(`
		SELECT COUNT(*) FROM projects p
		JOIN clients c ON c.id = p.client_id
		WHERE p.activated_at IS NOT NULL
		  AND p.activated_at <= c.created_at + INTERVAL '30 days'
		  AND EXISTS (SELECT 1 FROM meetings m WHERE m.project_id = p.id AND m.status = 'completed')
	`).Scan(&activatedIn30d)
	if err != nil {
		return 0, err
	}

	if participants == 0 {
		return 0, nil
	}
	rate := float64(activatedIn30d) / float64(participants) * 100
	return math.Round(rate*10) / 10, nil
}

// getFirstMeetingActivationRate returns the percentage of participants activated on their first meeting.
func getFirstMeetingActivationRate() (float64, error) {
	var totalParticipants, firstMeetingActivated int

	// Total participants (projects with at least 1 completed meeting)
	err := db.DB.QueryRow(`
		SELECT COUNT(DISTINCT m.project_id) FROM meetings m WHERE m.status = 'completed'
	`).Scan(&totalParticipants)
	if err != nil {
		return 0, err
	}

	// Activated on first meeting: project was activated at or before the first completed meeting time
	err = db.DB.QueryRow(`
		WITH first_meetings AS (
			SELECT DISTINCT ON (m.project_id)
				m.project_id, m.completed_at
			FROM meetings m
			WHERE m.status = 'completed'
			ORDER BY m.project_id, m.scheduled_at ASC
		)
		SELECT COUNT(*) FROM first_meetings fm
		JOIN projects p ON p.id = fm.project_id
		WHERE p.activated_at IS NOT NULL
		  AND p.activated_at <= fm.completed_at + INTERVAL '1 minute'
	`).Scan(&firstMeetingActivated)
	if err != nil {
		return 0, err
	}

	if totalParticipants == 0 {
		return 0, nil
	}
	rate := float64(firstMeetingActivated) / float64(totalParticipants) * 100
	return math.Round(rate*10) / 10, nil
}

// getCohortData returns activation data grouped by client purchase month.
func getCohortData() ([]CohortItem, error) {
	rows, err := db.DB.Query(`
		SELECT
			TO_CHAR(c.created_at, 'YYYY-MM') AS cohort,
			COUNT(p.id) AS total,
			COUNT(p.activated_at) AS activated
		FROM projects p
		JOIN clients c ON c.id = p.client_id
		GROUP BY cohort
		ORDER BY cohort DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var cohorts []CohortItem
	for rows.Next() {
		var item CohortItem
		if err := rows.Scan(&item.Month, &item.Total, &item.Activated); err != nil {
			continue
		}
		if item.Total > 0 {
			item.ActivationRate = math.Round(float64(item.Activated)/float64(item.Total)*1000) / 10
		}
		cohorts = append(cohorts, item)
	}

	if cohorts == nil {
		cohorts = []CohortItem{}
	}

	return cohorts, nil
}
