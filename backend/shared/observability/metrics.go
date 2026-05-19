package observability

import (
	"fmt"
	"net/http"
	"sort"
	"strings"
	"sync"
	"time"
)

type Metrics struct {
	mu        sync.RWMutex
	counters  map[string]float64
	durations map[string][]float64
}

func NewMetrics() *Metrics {
	return &Metrics{
		counters:  map[string]float64{},
		durations: map[string][]float64{},
	}
}

func (m *Metrics) Inc(name string) {
	m.Add(name, 1)
}

func (m *Metrics) Add(name string, value float64) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.counters[name] += value
}

func (m *Metrics) ObserveDuration(name string, started time.Time) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.durations[name] = append(m.durations[name], time.Since(started).Seconds())
}

func (m *Metrics) Handler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m.mu.RLock()
		defer m.mu.RUnlock()

		w.Header().Set("Content-Type", "text/plain; version=0.0.4")

		names := make([]string, 0, len(m.counters))
		for name := range m.counters {
			names = append(names, name)
		}
		sort.Strings(names)

		for _, name := range names {
			fmt.Fprintf(w, "# TYPE %s counter\n%s %g\n", sanitizeMetric(name), sanitizeMetric(name), m.counters[name])
		}

		histograms := make([]string, 0, len(m.durations))
		for name := range m.durations {
			histograms = append(histograms, name)
		}
		sort.Strings(histograms)

		for _, name := range histograms {
			values := m.durations[name]
			var sum float64
			for _, value := range values {
				sum += value
			}
			metric := sanitizeMetric(name)
			fmt.Fprintf(w, "# TYPE %s summary\n%s_count %d\n%s_sum %g\n", metric, metric, len(values), metric, sum)
		}
	}
}

func sanitizeMetric(name string) string {
	name = strings.ReplaceAll(name, "-", "_")
	name = strings.ReplaceAll(name, ".", "_")
	return name
}
