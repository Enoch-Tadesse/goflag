package segments

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	conn "github.com/Enoch-Tadesse/goflag/db/connection"
	"github.com/gorilla/mux"
)

func DeleteSegment(w http.ResponseWriter, r *http.Request) {
	// Set a timeout for the DB operation
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Extract segment name from URL path parameter
	vars := mux.Vars(r)
	name := vars["name"]

	// Step 1: Look up the segment by name to get its ID
	row := conn.DB.QueryRowContext(ctx, `
		SELECT id
		FROM segments
		WHERE name = ?
	`, name)

	var segmentID string
	if err := row.Scan(&segmentID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// If no segment found, return 400
			http.Error(w, fmt.Sprintf("Segment with name %s does not exist", name), http.StatusBadRequest)
			return
		}
		// Other unexpected DB error
		log.Printf("DeleteSegment: Failed to fetch segment: %v", err)
		http.Error(w, "Failed to fetch segment", http.StatusInternalServerError)
		return
	}

	// Step 2: Start a transaction to safely delete both rules and the segment
	tx, err := conn.DB.BeginTx(ctx, nil)
	if err != nil {
		log.Printf("DeleteSegment: Failed to begin transaction: %v", err)
		http.Error(w, "Failed to initiate transaction for db", http.StatusInternalServerError)
		return
	}

	// Step 3: Delete all associated rules for the segment
	_, err = tx.ExecContext(ctx, `
		DELETE 
		FROM segment_rules
		WHERE seg_id = ?
	`, segmentID)

	if err != nil {
		rollbackWithError(w, tx, "DeleteSegment: Failed to delete segment rules", "Failed to delete segment rules", err)
		return
	}

	// Step 4: Delete the segment itself
	_, err = tx.ExecContext(ctx, `
		DELETE
		FROM segments
		WHERE id = ?
	`, segmentID)
	if err != nil {
		rollbackWithError(w, tx, "DeleteSegment: Failed to delete segment", "Failed to delete segment", err)
		return
	}

	// Step 5: Commit the transaction
	if err := tx.Commit(); err != nil {
		rollbackWithError(w, tx, "DeleteSegment: Failed to commit transaction", "Failed to commit transaction", err)
		return
	}

	// Step 6: Return success response
	type response struct {
		Message string `json:"message"`
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(
		response{Message: "Segment deleted successfully"},
	)
}
