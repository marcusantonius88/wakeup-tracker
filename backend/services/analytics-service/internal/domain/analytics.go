package domain

import "wakeup-tracker/backend/shared/events"

type BehaviorMetrics struct {
	UserID                 string  `json:"user_id"`
	Confirmed              int     `json:"confirmed"`
	Failed                 int     `json:"failed"`
	MorningIntents         int     `json:"morning_intents"`
	ConsistencyScore       float64 `json:"consistency_score"`
	LastRegressionDetected bool    `json:"last_regression_detected"`
}

func Apply(metric BehaviorMetrics, event events.Envelope) (BehaviorMetrics, []events.Envelope) {
	userID, _ := event.Payload["user_id"].(string)
	if userID != "" {
		metric.UserID = userID
	}
	published := []events.Envelope{}

	switch event.EventType {
	case events.WakeUpConfirmed:
		metric.Confirmed++
	case events.WakeUpFailed:
		metric.Failed++
	case events.MorningIntentSubmitted:
		metric.MorningIntents++
	}

	total := metric.Confirmed + metric.Failed
	if total > 0 {
		metric.ConsistencyScore = float64(metric.Confirmed) / float64(total)
	}
	if metric.ConsistencyScore >= 0.8 && event.EventType == events.WakeUpConfirmed {
		published = append(published, events.NewEnvelope(events.ConsistencyImproved, metric.UserID, event.CorrelationID, map[string]any{"score": metric.ConsistencyScore}))
	}
	if metric.Failed >= 2 && event.EventType == events.WakeUpFailed {
		metric.LastRegressionDetected = true
		published = append(published, events.NewEnvelope(events.WakeUpRegressionDetected, metric.UserID, event.CorrelationID, map[string]any{"failures": metric.Failed}))
	}
	return metric, published
}
