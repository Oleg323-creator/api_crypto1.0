package handlers

import (
	"api_crypto1.0/internal/runners"
)

type Handler struct {
	Runner *runners.Runner
}

func NewHandler(runner *runners.Runner) *Handler {
	return &Handler{Runner: runner}
}
