package analytics_models

import (
	"github.com/google/uuid"
	"github.com/kadsin/sms-gateway/internal/container"
)

type SmsStats struct {
	Pending int `json:"pending"`
	Sent    int `json:"sent"`
	Failed  int `json:"failed"`
}

func SmsStatsByClient(clientID uuid.UUID) (*SmsStats, error) {
	var stats SmsStats

	query := `
	SELECT
		countIf(final_status = 'pending') AS pending,
		countIf(final_status = 'sent') AS sent,
		countIf(final_status = 'failed') AS failed
	FROM (
		SELECT
			id,
			argMax(status, updated_at) AS final_status
		FROM sms_messages
		WHERE sender_client_id = ?
		GROUP BY id
	)
	`

	if err := container.Analytics().Raw(query, clientID).Scan(&stats).Error; err != nil {
		return nil, err
	}

	return &stats, nil
}
