/*
Copyright 2024 Hurricane1988 Authors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

     http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// TODO: Prometheus的client libraries中目前提供了四种核心的Metrics类型
//       Counter、Gauge、Histogram和Summary
package main

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"math/rand"
	"net/http"
	"time"
)

// 定义自定义指标
var (
	// 定义Counter指标
	httpRequestTotalCounter = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "http_request_total",
		Help: "Total number of HTTP requests.",
	})
	
	// 定义Gauge指标
	currentMemoryUsageGauge = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "current_memory_usage_bytes",
		Help: "Current memory usage in bytes.",
	})
	
	// 定义Histogram指标
	httpRequestDurationHistogram = prometheus.NewHistogram(prometheus.HistogramOpts{
		Name: "http_request_duration_seconds",
		Help: "Duration of request in seconds.",
		// 默认桶设置
		Buckets: prometheus.DefBuckets,
	})
	
	// 定义Summary 指标
	httpRequestDurationSummary = prometheus.NewSummary(prometheus.SummaryOpts{
		Name: "request_latency_seconds",
		Help: "Latency of requests in seconds",
		// 指定摘要目标
		Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
	})
)

func init() {
	// 注册自定义指标
	prometheus.MustRegister(httpRequestTotalCounter)
	prometheus.MustRegister(currentMemoryUsageGauge)
	prometheus.MustRegister(httpRequestDurationHistogram)
	prometheus.MustRegister(httpRequestDurationSummary)
}

func main() {
	// 给httpRequestTotal指标赋值。模拟请求，每秒递增计数器
	go func() {
		for {
			httpRequestTotalCounter.Inc()
			time.Sleep(time.Second)
		}
	}()
	
	// currentMemoryUsage
	go func() {
		for {
			// 模拟内存使用率
			memUsage := 20.0 + 10.0*rand.Float64()
			currentMemoryUsageGauge.Set(memUsage)
			time.Sleep(time.Second)
		}
	}()
	
	// 为httpRequestDuration 指标模拟赋值
	go func() {
		for {
			duration := time.Duration(rand.Intn(100)) * time.Millisecond
			httpRequestDurationHistogram.Observe(duration.Seconds())
			time.Sleep(time.Second)
		}
	}()
	
	// 为httpRequestDurationSummary模拟赋值
	go func() {
		for {
			duration := time.Duration(rand.Intn(100)) * time.Millisecond
			httpRequestDurationSummary.Observe(duration.Seconds())
			time.Sleep(time.Second)
		}
	}()
	
	// 暴漏metrics
	http.Handle("/metrics", promhttp.Handler())
	
	log.Println("Beginning to serve on port :8080")
	log.Fatalln(http.ListenAndServe(":8080", nil))
}
