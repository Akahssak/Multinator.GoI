package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os/exec"
)

type ExecutionRequest struct {
	Script string `json:"script"` // The script to execute
}

// Handler to execute specified Go scripts
func executeGoFile(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Parse the request body to get the script name
	var req ExecutionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
		return
	}

	// Define allowed scripts
	allowedScripts := map[string]string{
		"script1": "D:\\projects\\go\\hackathon\\ingestion\\http\\api.go",
	}

	// Validate the script name
	scriptPath, exists := allowedScripts[req.Script]
	if !exists {
		http.Error(w, "Invalid script specified", http.StatusBadRequest)
		return
	}

	// Execute the specified Go script
	cmd := exec.Command("go", "run", scriptPath)
	output, err := cmd.CombinedOutput()

	if err != nil {
		http.Error(w, fmt.Sprintf("Error executing Go script: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	// Respond with the output of the Go script
	response := map[string]string{
		"message": string(output),
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func main() {
	http.HandleFunc("/api/run-go", executeGoFile)

	fmt.Println("Go server running at http://localhost:33060")
	http.ListenAndServe(":33060", nil)
}
