package main

import (
	"log/slog"

	_ "ml-elec/internal/config"
)

func main() {
	slog.Info("ml-elec starting")
	// TODO: Wire DI initialization will be added in Plan 05
	// TODO: Signal handling and graceful shutdown will be added in Plan 05
}
