package segments

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	conn "github.com/Enoch-Tadesse/goflag/db/connection"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

// DeleteRule deletes a segment rule by its UUID.
// The rule ID is extracted from the URL path parameter `id`.
func DeleteRule(w http.ResponseWriter, r *http.Request) {
	// Extract rule ID from URL path parameter
	vars := mux.Vars(r)
	idStr := strings.TrimSpace(vars["id"])
	ruleID, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "Invalid UUID format", http.StatusBadRequest)
		return
	}

	// Execute DELETE statement
	result, err := conn.DB.Exec(`
		DELETE FROM segment_rules WHERE id = ?
	`, ruleID)
	if err != nil {
		log.Printf("DeleteRule: Failed to delete rule with id %s: %v", ruleID, err)
		http.Error(w, "Failed to delete rule", http.StatusInternalServerError)
		return
	}

	// Check if a row was actually deleted
	affected, err := result.RowsAffected()
	if err != nil {
		log.Printf("DeleteRule: Failed to check affected rows for id %s: %v", ruleID, err)
		http.Error(w, "Failed to confirm deletion", http.StatusInternalServerError)
		return
	}

	if affected == 0 {
		http.Error(w, "Rule not found", http.StatusNotFound)
		return
	}

	// Return success response
	response := struct {
		Message string `json:"message"`
	}{
		Message: "Rule deleted successfully",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
