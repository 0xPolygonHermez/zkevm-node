package metrics

import (
	"net/http"
	"sync"

	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	storageMutex  sync.RWMutex
	registerer    prometheus.Registerer
	gauges        map[string]prometheus.Gauge
	counters      map[string]prometheus.Counter
	counterVecs   map[string]*prometheus.CounterVec
	histograms    map[string]prometheus.Histogram
	histogramVecs map[string]*prometheus.HistogramVec
	summaries     map[string]prometheus.Summary
	initialized   bool
	initOnce      sync.Once
)

// CounterVecOpts holds options for the CounterVec type.
type CounterVecOpts struct {
	prometheus.CounterOpts
	Labels []string
}

// HistogramVecOpts holds options for the HistogramVec type.
type HistogramVecOpts struct {
	prometheus.HistogramOpts
	Labels []string
}

// Init initializes the package variables.
func Init() {
	initOnce.Do(func() {
		storageMutex = sync.RWMutex{}
		registerer = prometheus.DefaultRegisterer
		gauges = make(map[string]prometheus.Gauge)
		counters = make(map[string]prometheus.Counter)
		counterVecs = make(map[string]*prometheus.CounterVec)
		histograms = make(map[string]prometheus.Histogram)
		histogramVecs = make(map[string]*prometheus.HistogramVec)
		summaries = make(map[string]prometheus.Summary)
		initialized = true
	})
}

// Handler returns the Prometheus http handler.
func Handler() http.Handler {
	return promhttp.Handler()
}

// RegisterGauges registers the provided gauge metrics to the Prometheus
// registerer.
func RegisterGauges(opts ...prometheus.GaugeOpts) {
	if !initialized {
		return
	}

	storageMutex.Lock()
	defer storageMutex.Unlock()

	for _, options := range opts {
		registerGaugeIfNotExists(options)
	}
}

// UnregisterGauges unregisters the provided gauge metrics from the Prometheus
// registerer.
func UnregisterGauges(names ...string) {
	if !initialized {
		return
	}

	storageMutex.Lock()
	defer storageMutex.Unlock()

	for _, name := range names {
		unregisterGaugeIfExists(name)
	}
}

// Gauge retrieves gauge metric by name
func Gauge(name string) (gauge prometheus.Gauge, exist bool) {
	if !initialized {
		return
	}

	storageMutex.RLock()
	defer storageMutex.RUnlock()

	if gauge, exist = gauges[name]; !exist {
		return nil, exist
	}

	return gauge, exist
}

// GaugeSet sets the value for gauge with the given name.
func GaugeSet(name string, value float64) {
	if !initialized {
		return
	}

	if c, ok := Gauge(name); ok {
		c.Set(value)
	}
}

// RegisterCounters registers the provided counter metrics to the Prometheus
// registerer.
func RegisterCounters(opts ...prometheus.CounterOpts) {
	if !initialized {
		return
	}

	storageMutex.Lock()
	defer storageMutex.Unlock()

	for _, options := range opts {
		registerCounterIfNotExists(options)
	}
}

// Counter retrieves counter metric by name
func Counter(name string) (counter prometheus.Counter, exist bool) {
	if !initialized {
		return
	}

	storageMutex.RLock()
	defer storageMutex.RUnlock()

	if counter, exist = counters[name]; !exist {
		return nil, exist
	}

	return counter, exist
}

// CounterInc increments the counter with the given name.
func CounterInc(name string) {
	if !initialized {
		return
	}

	if c, ok := Counter(name); ok {
		c.Inc()
	}
}

// CounterAdd increments the counter with the given name.
func CounterAdd(name string, value float64) {
	if !initialized {
		return
	}

	if c, ok := Counter(name); ok {
		c.Add(value)
	}
}

// UnregisterCounters unregisters the provided counter metrics from the
// Prometheus registerer.
func UnregisterCounters(names ...string) {
	if !initialized {
		return
	}

	storageMutex.Lock()
	defer storageMutex.Unlock()

	for _, name := range names {
		unregisterCounterIfExists(name)
	}
}

// RegisterCounterVecs registers the provided counter vec metrics to the
// Prometheus registerer.
func RegisterCounterVecs(opts ...CounterVecOpts) {
	if !initialized {
		return
	}

	storageMutex.Lock()
	defer storageMutex.Unlock()

	for _, options := range opts {
		registerCounterVecIfNotExists(options)
	}
}

// CounterVec retrieves counter ver metric by name
func CounterVec(name string) (counterVec *prometheus.CounterVec, exist bool) {
	if !initialized {
		return
	}

	storageMutex.RLock()
	defer storageMutex.RUnlock()

	counterVec, exist = counterVecs[name]

	return counterVec, exist
}

// CounterVecInc increments the counter vec with the given name and label.
func CounterVecInc(name string, label string) {
	if !initialized {
		return
	}

	if cv, ok := CounterVec(name); ok {
		cv.WithLabelValues(label).Inc()
	}
}

// CounterVecAdd increments the counter vec by the given value, with the given
// name and label.
func CounterVecAdd(name string, label string, value float64) {
	if !initialized {
		return
	}

	if cv, ok := CounterVec(name); ok {
		cv.WithLabelValues(label).Add(value)
	}
}

// UnregisterCounterVecs unregisters the provided counter vec metrics from the
// Prometheus registerer.
func UnregisterCounterVecs(names ...string) {
	if !initialized {
		return
	}

	storageMutex.Lock()
	defer storageMutex.Unlock()

	for _, name := range names {
		unregisterCounterVecIfExists(name)
	}
}

// RegisterHistograms registers the provided histogram metrics to the
// Prometheus registerer.
func RegisterHistograms(opts ...prometheus.HistogramOpts) {
	if !initialized {
		return
	}

	storageMutex.Lock()
	defer storageMutex.Unlock()

	for _, options := range opts {
		registerHistogramIfNotExists(options)
	}
}

// Histogram retrieves histogram metric by name
func Histogram(name string) (histogram prometheus.Histogram, exist bool) {
	if !initialized {
		return
	}

	storageMutex.RLock()
	defer storageMutex.RUnlock()

	if histogram, exist = histograms[name]; !exist {
		return nil, exist
	}

	return histogram, exist
}

// HistogramObserve observes the histogram from the given start time.
func HistogramObserve(name string, value float64) {
	if !initialized {
		return
	}

	if histo, ok := Histogram(name); ok {
		histo.Observe(value)
	}
}

// UnregisterHistogram unregisters the provided histogram metrics from the
// Prometheus registerer.
func UnregisterHistogram(names ...string) {
	if !initialized {
		return
	}

	storageMutex.Lock()
	defer storageMutex.Unlock()

	for _, name := range names {
		unregisterHistogramIfExists(name)
	}
}

// RegisterHistogramVecs registers the provided histogram vec metrics to the
// Prometheus registerer.
func RegisterHistogramVecs(opts ...HistogramVecOpts) {
	if !initialized {
		return
	}

	storageMutex.Lock()
	defer storageMutex.Unlock()

	for _, options := range opts {
		registerHistogramVecIfNotExists(options)
	}
}

// HistogramVec retrieves histogram ver metric by name
func HistogramVec(name string) (histgramVec *prometheus.HistogramVec, exist bool) {
	if !initialized {
		return
	}

	storageMutex.RLock()
	defer storageMutex.RUnlock()

	histgramVec, exist = histogramVecs[name]

	return histgramVec, exist
}

// HistogramVecObserve observes the histogram vec with the given name, label and value.
func HistogramVecObserve(name string, label string, value float64) {
	if !initialized {
		return
	}

	if cv, ok := HistogramVec(name); ok {
		cv.WithLabelValues(label).Observe(value)
	}
}

// UnregisterHistogramVecs unregisters the provided histogram vec metrics from the
// Prometheus registerer.
func UnregisterHistogramVecs(names ...string) {
	if !initialized {
		return
	}

	storageMutex.Lock()
	defer storageMutex.Unlock()

	for _, name := range names {
		unregisterHistogramVecIfExists(name)
	}
}

// RegisterSummaries registers the provided summary metrics to the Prometheus
// registerer.
func RegisterSummaries(opts ...prometheus.SummaryOpts) {
	if !initialized {
		return
	}

	storageMutex.Lock()
	defer storageMutex.Unlock()

	for _, options := range opts {
		registerSummaryIfNotExists(options)
	}
}

// Summary retrieves summary metric by name
func Summary(name string) (summary prometheus.Summary, exist bool) {
	if !initialized {
		return
	}

	storageMutex.RLock()
	defer storageMutex.RUnlock()

	if summary, exist = summaries[name]; !exist {
		return nil, exist
	}

	return summary, exist
}

// UnregisterSummaries unregisters the provided summary metrics from the
// Prometheus registerer.
func UnregisterSummaries(names ...string) {
	if !initialized {
		return
	}

	storageMutex.Lock()
	defer storageMutex.Unlock()

	for _, name := range names {
		unregisterSummaryIfExists(name)
	}
}

// registerGaugeIfNotExists registers single gauge metric if not exists
func registerGaugeIfNotExists(opts prometheus.GaugeOpts) {
	if _, exist := gauges[opts.Name]; exist {
		log.Warnf("Gauge metric '%v' already exists.", opts.Name)
		return
	}

	log.Infof("Creating Gauge Metric '%v' ...", opts.Name)
	gauge := prometheus.NewGauge(opts)
	log.Infof("Gauge Metric '%v' successfully created! Labels: %p", opts.Name, opts.ConstLabels)

	log.Infof("Registering Gauge Metric '%v' ...", opts.Name)
	registerer.MustRegister(gauge)
	log.Infof("Gauge Metric '%v' successfully registered!", opts.Name)

	gauges[opts.Name] = gauge
}

// unregisterGaugeIfExists unregisters single gauge metric if exists
func unregisterGaugeIfExists(name string) {
	var (
		gauge prometheus.Gauge
		ok    bool
	)

	if gauge, ok = gauges[name]; !ok {
		log.Warnf("Trying to delete non-existing Gauge gauge '%v'.", name)
		return
	}

	log.Infof("Unregistering Gauge Metric '%v' ...", name)
	ok = registerer.Unregister(gauge)
	if !ok {
		log.Errorf("Failed to unregister Gauge Metric '%v'.", name)
		return
	}
	delete(gauges, name)
	log.Infof("Gauge Metric '%v' successfully unregistered!", name)
}

// registerCounterIfNotExists registers single counter metric if not exists
func registerCounterIfNotExists(opts prometheus.CounterOpts) {
	if _, exist := counters[opts.Name]; exist {
		log.Warnf("Counter metric '%v' already exists.", opts.Name)
		return
	}

	log.Infof("Creating Counter Metric '%v' ...", opts.Name)
	counter := prometheus.NewCounter(opts)
	log.Infof("Counter Metric '%v' successfully created! Labels: %p", opts.Name, opts.ConstLabels)

	log.Infof("Registering Counter Metric '%v' ...", opts.Name)
	registerer.MustRegister(counter)
	log.Infof("Counter Metric '%v' successfully registered!", opts.Name)

	counters[opts.Name] = counter
}

// unregisterCounterIfExists unregisters single counter metric if exists
func unregisterCounterIfExists(name string) {
	var (
		counter prometheus.Counter
		ok      bool
	)

	if counter, ok = counters[name]; !ok {
		log.Warnf("Trying to delete non-existing Counter counter '%v'.", name)
		return
	}

	log.Infof("Unregistering Counter Metric '%v' ...", name)
	ok = registerer.Unregister(counter)
	if !ok {
		log.Errorf("Failed to unregister Counter Metric '%v'.", name)
		return
	}
	delete(counters, name)
	log.Infof("Counter Metric '%v' successfully unregistered!", name)
}

// registerCounterVecIfNotExists registers single counter vec metric if not exists
func registerCounterVecIfNotExists(opts CounterVecOpts) {
	if _, exist := counterVecs[opts.Name]; exist {
		log.Warnf("Counter vec metric '%v' already exists.", opts.Name)
		return
	}

	log.Infof("Creating Counter Vec Metric '%v' ...", opts.Name)
	counterVec := prometheus.NewCounterVec(opts.CounterOpts, opts.Labels)
	log.Infof("Counter Vec Metric '%v' successfully created! Labels: %p", opts.Name, opts.ConstLabels)

	log.Infof("Registering Counter Vec Metric '%v' ...", opts.Name)
	registerer.MustRegister(counterVec)
	log.Infof("Counter Vec Metric '%v' successfully registered!", opts.Name)

	counterVecs[opts.Name] = counterVec
}

// unregisterCounterVecIfExists unregisters single counter metric if exists
func unregisterCounterVecIfExists(name string) {
	var (
		counterVec *prometheus.CounterVec
		ok         bool
	)

	if counterVec, ok = counterVecs[name]; !ok {
		log.Warnf("Trying to delete non-existing Counter Vec counter '%v'.", name)
		return
	}

	log.Infof("Unregistering Counter Vec Metric '%v' ...", name)
	ok = registerer.Unregister(counterVec)
	if !ok {
		log.Errorf("Failed to unregister Counter Vec Metric '%v'.", name)
		return
	}
	delete(counterVecs, name)
	log.Infof("Counter Vec Metric '%v' successfully unregistered!", name)
}

// registerHistogramIfNotExists registers single histogram metric if not exists
func registerHistogramIfNotExists(opts prometheus.HistogramOpts) {
	if _, exist := histograms[opts.Name]; exist {
		log.Warnf("Histogram metric '%v' already exists.", opts.Name)
		return
	}

	log.Infof("Creating Histogram Metric '%v' ...", opts.Name)
	histogram := prometheus.NewHistogram(opts)
	log.Infof("Histogram Metric '%v' successfully created! Labels: %p", opts.Name, opts.ConstLabels)

	log.Infof("Registering Histogram Metric '%v' ...", opts.Name)
	registerer.MustRegister(histogram)
	log.Infof("Histogram Metric '%v' successfully registered!", opts.Name)

	histograms[opts.Name] = histogram
}

// unregisterHistogramIfExists unregisters single histogram metric if exists
func unregisterHistogramIfExists(name string) {
	var (
		histogram prometheus.Histogram
		ok        bool
	)

	if histogram, ok = histograms[name]; !ok {
		log.Warnf("Trying to delete non-existing Histogram histogram '%v'.", name)
		return
	}

	log.Infof("Unregistering Histogram Metric '%v' ...", name)
	ok = registerer.Unregister(histogram)
	if !ok {
		log.Errorf("Failed to unregister Histogram Metric '%v'.", name)
		return
	}
	delete(histograms, name)
	log.Infof("Histogram Metric '%v' successfully unregistered!", name)
}

// registerHistogramVecIfNotExists unregisters single counter metric if exists
func registerHistogramVecIfNotExists(opts HistogramVecOpts) {
	if _, exist := histogramVecs[opts.Name]; exist {
		log.Warnf("Histogram vec metric '%v' already exists.", opts.Name)
		return
	}

	log.Infof("Creating Histogram Vec Metric '%v' ...", opts.Name)
	histogramVec := prometheus.NewHistogramVec(opts.HistogramOpts, opts.Labels)
	log.Infof("Histogram Vec Metric '%v' successfully created! Labels: %p", opts.Name, opts.ConstLabels)

	log.Infof("Registering Histogram Vec Metric '%v' ...", opts.Name)
	registerer.MustRegister(histogramVec)
	log.Infof("Histogram Vec Metric '%v' successfully registered!", opts.Name)

	histogramVecs[opts.Name] = histogramVec
}

// unregisterHistogramVecIfExists unregisters single histogram metric if exists
func unregisterHistogramVecIfExists(name string) {
	var (
		histogramVec *prometheus.HistogramVec
		ok           bool
	)

	if histogramVec, ok = histogramVecs[name]; !ok {
		log.Warnf("Trying to delete non-existing Histogram Vec '%v'.", name)
		return
	}

	log.Infof("Unregistering Histogram Vec Metric '%v' ...", name)
	ok = registerer.Unregister(histogramVec)
	if !ok {
		log.Errorf("Failed to unregister Histogram Vec Metric '%v'.", name)
		return
	}
	delete(histogramVecs, name)
	log.Infof("Histogram Vec Metric '%v' successfully unregistered!", name)
}

// registerSummaryIfNotExists registers single summary metric if not exists
func registerSummaryIfNotExists(opts prometheus.SummaryOpts) {
	if _, exist := summaries[opts.Name]; exist {
		log.Warnf("Summary metric '%v' already exists.", opts.Name)
		return
	}

	log.Infof("Creating Summary Metric '%v' ...", opts.Name)
	summary := prometheus.NewSummary(opts)
	log.Infof("Summary Metric '%v' successfully created! Labels: %p", opts.Name, opts.ConstLabels)

	log.Infof("Registering Summary Metric '%v' ...", opts.Name)
	registerer.MustRegister(summary)
	log.Infof("Summary Metric '%v' successfully registered!", opts.Name)

	summaries[opts.Name] = summary
}

// unregisterSummaryIfExists unregisters single summary metric if exists
func unregisterSummaryIfExists(name string) {
	var (
		summary prometheus.Summary
		ok      bool
	)

	if summary, ok = summaries[name]; !ok {
		log.Warnf("Trying to delete non-existing Summary summary '%v'.", name)
		return
	}

	log.Infof("Unregistering Summary Metric '%v' ...", name)
	ok = registerer.Unregister(summary)
	if !ok {
		log.Errorf("Failed to unregister Summary Metric '%v'.", name)
		return
	}
	delete(summaries, name)
	log.Infof("Summary Metric '%v' successfully unregistered!", name)
}
