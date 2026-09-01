package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

        "github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type ProjetoKorpResponse struct {
	Nome    string `json:"nome"`
	Horario string `json:"horario"`
}

var requestsTotal = prometheus.NewCounter(
	prometheus.CounterOpts{
		Name: "projeto_korp_http_requests_total",
		Help: "Quantidade total de requisicoes recebidas pelo endpoint /projeto-korp",
	},
)

func init() {
	prometheus.MustRegister(requestsTotal)
}

func projetoKorpHandler(w http.ResponseWriter, r *http.Request) {
        requestsTotal.Inc()

	if r.Method != http.MethodGet {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}

	response := ProjetoKorpResponse{
		Nome:    "Projeto Korp",
		Horario: time.Now().UTC().Format(time.RFC3339),
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Erro ao gerar resposta", http.StatusInternalServerError)
		return
	}
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK\n"))
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
