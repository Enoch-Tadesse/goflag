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
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	vars := mux.Vars(r)
	name := vars["name"]

	row := conn.DB.QueryRowContext(ctx, `
		SELECT id
		FROM segments
		WHERE name = ?
	`, name)

	var segmentID string
	if err := row.Scan(&segmentID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, fmt.Sprintf("Segment with name %s does not exist", name), http.StatusBadRequest)
			return
		}
		log.Printf("DeleteSegment: Failed to fetch segment: %v", err)
		http.Error(w, "Failed to fetch segment", http.StatusInternalServerError)
		return
	}

	tx, err := conn.DB.BeginTx(ctx, nil)
	if err != nil {
		log.Printf("DeleteSegment: Failed to begin transaction: %v", err)
		http.Error(w, "Failed to initiate transaction for db", http.StatusInternalServerError)
		return
	}

	_, err = tx.ExecContext(ctx, `
		DELETE 
		FROM segment_rules
		WHERE seg_id = ?
	`, segmentID)

	if err != nil {
		rollbackWithError(w, tx, "DeleteSegment: Failed to delete segment rules", "Failed to delete segment rules", err)
		return
	}
	_, err = tx.ExecContext(ctx, `
		DELETE
		FROM segments
		WHERE id = ?
	`, segmentID)
	if err != nil {
		rollbackWithError(w, tx, "DeleteSegment: Failed to delete segment", "Failed to delete segment", err)
		return
	}
	if err := tx.Commit(); err != nil {
		rollbackWithError(w, tx, "DeleteSegment: Failed to commit transaction", "Failed to commit transaction", err)
		return
	}

	type response struct {
		Message string
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(
		response{Message: "Segment delete successfully"},
	)
}
