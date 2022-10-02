package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gocolly/colly"

	"github.com/ghodss/yaml"
	_ "github.com/go-sql-driver/mysql"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/version"
	log "github.com/sirupsen/logrus"
)

var coinConfig CoinConfig

const (
	collector = "query_exporter"
)

func main() {
	var err error
	var configFile, bind string
	// =====================
	// Get OS parameter
	// =====================
	flag.StringVar(&configFile, "config", "config.yml", "configuration file")
	flag.StringVar(&bind, "bind", "0.0.0.0:9109", "bind")
	flag.Parse()

	// =====================
	// Load config & yaml
	// =====================
	var b []byte
	if b, err = ioutil.ReadFile(configFile); err != nil {
		log.Errorf("Failed to read config file: %s", err)
		os.Exit(1)
	}

	// Load yaml
	if err := yaml.Unmarshal(b, &coinConfig); err != nil {
		log.Errorf("Failed to load config: %s", err)
		os.Exit(1)
	}

	// ========================
	// Regist handler
	// ========================
	log.Infof("Regist version collector - %s", collector)
	prometheus.Register(version.NewCollector(collector))
	prometheus.Register(&QueryCollector{})

	// Regist http handler
	log.Infof("HTTP handler path - %s", "/metrics")
	http.HandleFunc("/metrics", func(w http.ResponseWriter, r *http.Request) {
		h := promhttp.HandlerFor(prometheus.Gatherers{
			prometheus.DefaultGatherer,
		}, promhttp.HandlerOpts{})
		h.ServeHTTP(w, r)
	})

	// start server
	log.Infof("Starting http server - %s", bind)
	if err := http.ListenAndServe(bind, nil); err != nil {
		log.Errorf("Failed to start http server: %s", err)
	}
}

// =============================
// Config config structure
// =============================
type CoinConfig struct {
	Metrics map[string]struct {
		URL         []string
		Type        string
		Description string
		Value       string
		metricDesc  *prometheus.Desc
	}
}

// =============================
// QueryCollector exporter
// =============================
type QueryCollector struct{}

// Describe prometheus describe
func (e *QueryCollector) Describe(ch chan<- *prometheus.Desc) {
	for metricName, metric := range coinConfig.Metrics {
		metric.metricDesc = prometheus.NewDesc(
			prometheus.BuildFQName(collector, "", metricName),
			metric.Description,
			[]string{"coin"}, nil,
		)
		coinConfig.Metrics[metricName] = metric
		log.Infof("metric description for \"%s\" registerd", metricName)
	}
}

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func RandomString() string {
	b := make([]byte, rand.Intn(10)+10)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

// Collect prometheus collect
func (e *QueryCollector) Collect(ch chan<- prometheus.Metric) {
	//var val float64
	for metricName, metric := range coinConfig.Metrics {
		log.Infof("metric description for \"%s\" registerd", metricName)
		data := make(map[string]string)

		for url := range metric.URL {

			c := colly.NewCollector()
			coinName := ""
			c.OnHTML("main", func(e *colly.HTMLElement) {
				coinName = ""
				e.ForEach("div.price-container", func(_ int, el *colly.HTMLElement) {
					coinName = e.ChildText("h2.coin-name")
				})

				e.ForEach("div.cion-item.coin-price-large", func(_ int, el *colly.HTMLElement) {
					coinResult := el.ChildText("span")
					coinResult = strings.ReplaceAll(coinResult, "$", "")
					result, err := strconv.ParseFloat(coinResult, 8)
					data[coinName] = fmt.Sprintf("%f", result)
					//fmt.Println(time.Now().Format("01-02-2006 15:04:05"), coinName, coinResult)
					log.Infof(fmt.Sprintf("Coin: %s, Price: %s", coinName, coinResult))
					if err != nil {
						panic(err)

					}
					ch <- prometheus.MustNewConstMetric(metric.metricDesc, prometheus.GaugeValue, result, coinName)
					//})
				})

			})

			c.OnRequest(func(r *colly.Request) {
				r.Headers.Set("User-Agent", RandomString())
			})

			c.OnError(func(r *colly.Response, err error) {
				fmt.Println("Request URL:", r.Request.URL, "failed with response:", r, "\nError:", err)
			})

			c.Visit(metric.URL[url])

		}
		log.Infof(fmt.Sprintf("------------------------------------------------------------------------------------------"))
	}
}
