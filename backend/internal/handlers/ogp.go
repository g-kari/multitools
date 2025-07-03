package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"ogp-verification-service/internal/models"
	"ogp-verification-service/internal/services"
)

type OGPHandler struct {
	service *services.OGPService
	limiter *RateLimiter
}

type RateLimiter struct {
	clients map[string]*ClientInfo
	mu      sync.RWMutex
}

type ClientInfo struct {
	requests  int
	lastReset time.Time
}

func NewOGPHandler() *OGPHandler {
	return &OGPHandler{
		service: services.NewOGPService(),
		limiter: &RateLimiter{
			clients: make(map[string]*ClientInfo),
		},
	}
}

func (h *OGPHandler) VerifyOGP(w http.ResponseWriter, r *http.Request) {
	// Handle CORS preflight requests
	if r.Method == http.MethodOptions {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	clientIP := h.getClientIP(r)
	if !h.limiter.Allow(clientIP) {
		http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
		return
	}

	var req models.OGPRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if req.URL == "" {
		http.Error(w, "URL is required", http.StatusBadRequest)
		return
	}

	response, err := h.service.FetchOGPData(req.URL)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error fetching OGP data: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		return
	}
}

func (h *OGPHandler) getClientIP(r *http.Request) string {
	xff := r.Header.Get("X-Forwarded-For")
	if xff != "" {
		return xff
	}
	// Extract IP from host:port format
	if idx := strings.LastIndex(r.RemoteAddr, ":"); idx != -1 {
		return r.RemoteAddr[:idx]
	}
	return r.RemoteAddr
}

func (rl *RateLimiter) Allow(clientIP string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	client, exists := rl.clients[clientIP]
	
	if !exists {
		rl.clients[clientIP] = &ClientInfo{
			requests:  1,
			lastReset: now,
		}
		return true
	}

	if now.Sub(client.lastReset) >= time.Minute {
		client.requests = 1
		client.lastReset = now
		return true
	}

	if client.requests >= 10 {
		return false
	}

	client.requests++
	return true
}