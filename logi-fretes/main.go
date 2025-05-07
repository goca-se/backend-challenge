package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"
)

const (
	expectedTokenPrefix = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9"
	failureRate         = 0.3 // 30% failure rate
	minResponseTime     = 500 * time.Millisecond
	maxResponseTime     = 2000 * time.Millisecond
)

// Request represents the structure of the incoming JSON payload
type Request struct {
	Origem struct {
		CEP string `json:"cep"`
	} `json:"origem"`
	Destino struct {
		CEP string `json:"cep"`
	} `json:"destino"`
	Pacote struct {
		Peso      float64 `json:"peso"`
		Dimensoes struct {
			Altura      float64 `json:"altura"`
			Largura     float64 `json:"largura"`
			Comprimento float64 `json:"comprimento"`
		} `json:"dimensoes"`
		Valor float64 `json:"valor"`
	} `json:"pacote"`
	Servicos []string `json:"servicos"`
}

// SuccessResponse represents the structure of a successful response
type SuccessResponse struct {
	Status    string    `json:"status"`
	RequestID string    `json:"request_id"`
	Cotacoes  []Cotacao `json:"cotacoes"`
	Meta      struct {
		ProcessadoEm string `json:"processado_em"`
	} `json:"meta"`
}

// Cotacao represents a shipping quote option
type Cotacao struct {
	Tipo             string  `json:"tipo"`
	Codigo           string  `json:"codigo"`
	Valor            float64 `json:"valor"`
	PrazoDias        int     `json:"prazo_dias"`
	RegiaoDisponivel bool    `json:"regiao_disponivel"`
}

// ErrorResponse represents the structure of an error response
type ErrorResponse struct {
	Status    string   `json:"status"`
	ErrorCode string   `json:"error_code"`
	Message   string   `json:"message"`
	Details   []Detail `json:"details,omitempty"`
	RequestID string   `json:"request_id"`
}

// Detail represents validation error details
type Detail struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// ServiceUnavailableResponse represents a service unavailable error
type ServiceUnavailableResponse struct {
	Status    string `json:"status"`
	ErrorCode string `json:"error_code"`
	Message   string `json:"message"`
	RequestID string `json:"request_id"`
}

// Random number generator instance
var rng = rand.New(rand.NewSource(time.Now().UnixNano()))

// generateRequestID generates a unique request ID
func generateRequestID() string {
	return "b3f5a7d9-e201-4b7f-91d5-" + randomString(12)
}

// randomString generates a random string of specified length
func randomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyz0123456789"
	result := make([]byte, length)
	for i := range result {
		result[i] = charset[rng.Intn(len(charset))]
	}
	return string(result)
}

// validateToken checks if the provided token starts with the expected prefix
func validateToken(token string) bool {
	return strings.HasPrefix(token, expectedTokenPrefix)
}

// simulateProcessingTime simulates API processing time between 0.5-2 seconds
func simulateProcessingTime() time.Duration {
	duration := minResponseTime + time.Duration(rng.Float64()*float64(maxResponseTime-minResponseTime))
	return duration
}

// shouldFail returns true with a probability equal to the failure rate
func shouldFail() bool {
	return rng.Float64() < failureRate
}

// validateRequest validates the request payload
func validateRequest(req Request) (bool, []Detail) {
	var details []Detail

	// Check for valid CEP
	if req.Origem.CEP == "" {
		details = append(details, Detail{
			Field:   "origem.cep",
			Message: "CEP de origem é obrigatório",
		})
	}

	if req.Destino.CEP == "" {
		details = append(details, Detail{
			Field:   "destino.cep",
			Message: "CEP de destino é obrigatório",
		})
	}

	// Check for valid weight
	if req.Pacote.Peso <= 0 {
		details = append(details, Detail{
			Field:   "pacote.peso",
			Message: "Peso deve ser maior que zero",
		})
	}

	// Check for valid dimensions
	if req.Pacote.Dimensoes.Altura <= 0 {
		details = append(details, Detail{
			Field:   "pacote.dimensoes.altura",
			Message: "Altura deve ser maior que zero",
		})
	}

	if req.Pacote.Dimensoes.Largura <= 0 {
		details = append(details, Detail{
			Field:   "pacote.dimensoes.largura",
			Message: "Largura deve ser maior que zero",
		})
	}

	if req.Pacote.Dimensoes.Comprimento <= 0 {
		details = append(details, Detail{
			Field:   "pacote.dimensoes.comprimento",
			Message: "Comprimento deve ser maior que zero",
		})
	}

	return len(details) == 0, details
}

// handleQuote handles the shipping quote request
func handleQuote(w http.ResponseWriter, r *http.Request) {
	// Check HTTP method
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	// Validate authentication
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(ErrorResponse{
			Status:    "error",
			ErrorCode: "UNAUTHORIZED",
			Message:   "Token de autenticação não fornecido",
			RequestID: generateRequestID(),
		})
		return
	}

	token := strings.TrimPrefix(authHeader, "Bearer ")
	if !validateToken(token) {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(ErrorResponse{
			Status:    "error",
			ErrorCode: "INVALID_TOKEN",
			Message:   "Token de autenticação inválido",
			RequestID: generateRequestID(),
		})
		return
	}

	// Parse the request body
	var req Request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{
			Status:    "error",
			ErrorCode: "INVALID_JSON",
			Message:   "Formato de JSON inválido",
			RequestID: generateRequestID(),
		})
		return
	}

	// Validate request
	valid, details := validateRequest(req)
	if !valid {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{
			Status:    "error",
			ErrorCode: "VALIDATION_ERROR",
			Message:   "Erro de validação nos dados informados",
			Details:   details,
			RequestID: generateRequestID(),
		})
		return
	}

	// Generate a unique request ID
	requestID := generateRequestID()

	// Simulate processing time
	processingTime := simulateProcessingTime()
	time.Sleep(processingTime)

	// Simulate random failures (30% chance)
	if shouldFail() {
		w.WriteHeader(http.StatusServiceUnavailable)
		json.NewEncoder(w).Encode(ServiceUnavailableResponse{
			Status:    "error",
			ErrorCode: "SERVICE_UNAVAILABLE",
			Message:   "Serviço temporariamente indisponível. Tente novamente mais tarde.",
			RequestID: requestID,
		})
		return
	}

	// Create a successful response
	var response SuccessResponse
	response.Status = "success"
	response.RequestID = requestID

	// Add quote options based on requested services
	for _, service := range req.Servicos {
		switch service {
		case "standard":
			response.Cotacoes = append(response.Cotacoes, Cotacao{
				Tipo:             "standard",
				Codigo:           "STD",
				Valor:            28.50,
				PrazoDias:        3,
				RegiaoDisponivel: true,
			})
		case "express":
			response.Cotacoes = append(response.Cotacoes, Cotacao{
				Tipo:             "express",
				Codigo:           "EXP",
				Valor:            42.75,
				PrazoDias:        1,
				RegiaoDisponivel: true,
			})
		case "economic":
			response.Cotacoes = append(response.Cotacoes, Cotacao{
				Tipo:             "economic",
				Codigo:           "ECO",
				Valor:            19.99,
				PrazoDias:        5,
				RegiaoDisponivel: true,
			})
		}
	}

	// If no specific services were requested, return all available options
	if len(req.Servicos) == 0 || (len(req.Servicos) == 1 && req.Servicos[0] == "todos") {
		response.Cotacoes = []Cotacao{
			{
				Tipo:             "standard",
				Codigo:           "STD",
				Valor:            28.50,
				PrazoDias:        3,
				RegiaoDisponivel: true,
			},
			{
				Tipo:             "express",
				Codigo:           "EXP",
				Valor:            42.75,
				PrazoDias:        1,
				RegiaoDisponivel: true,
			},
			{
				Tipo:             "economic",
				Codigo:           "ECO",
				Valor:            19.99,
				PrazoDias:        5,
				RegiaoDisponivel: true,
			},
		}
	}

	// Set the processing time in the response
	response.Meta.ProcessadoEm = fmt.Sprintf("%.4fs", processingTime.Seconds())

	// Send the response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func main() {
	// Set up routes
	http.HandleFunc("/api/cotacoes", handleQuote)

	// Start the server
	port := "8080"
	if envPort := os.Getenv("PORT"); envPort != "" {
		port = envPort
	}

	log.Printf("LogiFretes API server starting on port %s...\n", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
} 