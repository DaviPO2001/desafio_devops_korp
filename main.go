package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type ProjetoKorpResponse struct {
	Nome    string `json:"nome"`
	Horario string `json:"horario"`
}

type statusRecorder struct {
	http.ResponseWriter
	statusCode int
}

func (r *statusRecorder) WriteHeader(statusCode int) {
	r.statusCode = statusCode
	r.ResponseWriter.WriteHeader(statusCode)
}

var (
	requestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "projeto_korp_http_requests_total",
			Help: "Quantidade total de requisicoes recebidas pelo endpoint /projeto-korp.",
		},
		[]string{"method", "status"},
	)

	requestErrorsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "projeto_korp_http_errors_total",
			Help: "Quantidade total de respostas HTTP com status >= 400 no endpoint /projeto-korp.",
		},
		[]string{"method", "status"},
	)

	requestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "projeto_korp_http_request_duration_seconds",
			Help: "Duracao das requisicoes do endpoint /projeto-korp em segundos.",
			Buckets: []float64{
				0.005,
				0.010,
				0.025,
				0.050,
				0.100,
				0.250,
				0.500,
				1.000,
				2.500,
				5.000,
			},
		},
		[]string{"method", "status"},
	)

	requestsInFlight = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "projeto_korp_http_requests_in_flight",
			Help: "Quantidade de requisicoes atualmente em processamento no endpoint /projeto-korp.",
		},
	)
)

func init() {
	prometheus.MustRegister(
		requestsTotal,
		requestErrorsTotal,
		requestDuration,
		requestsInFlight,
	)
}

func projetoKorpHandler(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	requestsInFlight.Inc()

	recorder := &statusRecorder{
		ResponseWriter: w,
		statusCode:     http.StatusOK,
	}

	defer func() {
		requestsInFlight.Dec()

		status := strconv.Itoa(recorder.statusCode)

		requestsTotal.
			WithLabelValues(r.Method, status).
			Inc()

		requestDuration.
			WithLabelValues(r.Method, status).
			Observe(time.Since(start).Seconds())

		if recorder.statusCode >= http.StatusBadRequest {
			requestErrorsTotal.
				WithLabelValues(r.Method, status).
				Inc()
		}
	}()

	if r.Method != http.MethodGet {
		http.Error(
			recorder,
			"Método não permitido",
			http.StatusMethodNotAllowed,
		)
		return
	}

	response := ProjetoKorpResponse{
		Nome:    "Projeto Korp",
		Horario: time.Now().UTC().Format(time.RFC3339),
	}

	recorder.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(recorder).Encode(response); err != nil {
		log.Printf("erro ao gerar resposta JSON: %v", err)
		return
	}
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)

	if _, err := w.Write([]byte("OK\n")); err != nil {
		log.Printf("erro ao escrever health response: %v", err)
	}
}

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/projeto-korp", projetoKorpHandler)
	mux.HandleFunc("/health", healthHandler)
	mux.Handle("/metrics", promhttp.Handler())

	log.Println("http-server-projeto-korp iniciado na porta 8080")

	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}
}
