package segments

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	conn "github.com/Enoch-Tadesse/goflag/db/connection"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

func DeleteSegment(w http.ResponseWriter, r *http.Request) {
	// Set a timeout for the DB operation
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Extract segment id from URL path parameter
	vars := mux.Vars(r)
	idStr := strings.TrimSpace(vars["id"])
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "invalid uuid", http.StatusBadRequest)
		return
	}

	// Check if the segment exists
	row := conn.DB.QueryRowContext(ctx, `
		SELECT 1
		FROM segments
		WHERE id = ?
	`, id)

	var exists bool
	if err := row.Scan(&exists); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// If no segment found, return 400
			http.Error(w, fmt.Sprintf("Segment with id %s does not exist", id), http.StatusBadRequest)
			return
		}
		// Other unexpected DB error
		log.Printf("DeleteSegment: Failed to fetch segment: %v", err)
		http.Error(w, "Failed to fetch segment", http.StatusInternalServerError)
		return
	}

	// Start a transaction to safely delete both rules and the segment
	tx, err := conn.DB.BeginTx(ctx, nil)
	if err != nil {
		log.Printf("DeleteSegment: Failed to begin transaction: %v", err)
		http.Error(w, "Failed to initiate transaction for db", http.StatusInternalServerError)
		return
	}

	// Delete all associated rules for the segment
	_, err = tx.ExecContext(ctx, `
		DELETE 
		FROM segment_rules
		WHERE seg_id = ?
	`, id)

	if err != nil {
		rollbackWithError(w, tx, "DeleteSegment: Failed to delete segment rules", "Failed to delete segment rules", err)
		return
	}

	// Delete the segment itself
	_, err = tx.ExecContext(ctx, `
		DELETE
		FROM segments
		WHERE id = ?
	`, id)
	if err != nil {
		rollbackWithError(w, tx, "DeleteSegment: Failed to delete segment", "Failed to delete segment", err)
		return
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		rollbackWithError(w, tx, "DeleteSegment: Failed to commit transaction", "Failed to commit transaction", err)
		return
	}

	// Return success response
	type response struct {
		Message string `json:"message"`
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(
		response{Message: "Segment deleted successfully"},
	)
}
