package segments

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	conn "github.com/Enoch-Tadesse/goflag/db/connection"
	"github.com/google/uuid"
)

type rule struct {
	Attribute string `json:"attribute"`
	Operator  string `json:"operator"`
	Value     string `json:"value"`
}

func CreateSegment(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var body struct {
		Name  string `json:"name"`
		Rules []rule `json:"rules"`
	}

	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if body.Name == "" {
		http.Error(w, "Name is required", http.StatusBadRequest)
		return
	}

	// check if there exist atlease one rule
	if len(body.Rules) == 0 {
		http.Error(w, "A segment must have atlease one rule", http.StatusBadRequest)
		return
	}

	// check for rules integration
	for _, rule := range body.Rules {
		if rule.Attribute == "" || rule.Operator == "" || rule.Value == "" {
			http.Error(w, "All rules must have non-empty attribute, operator, and value", http.StatusBadRequest)
			return
		}
	}

	// check if a segment with similar name already exist
	row := conn.DB.QueryRowContext(ctx, `
		SELECT 1 FROM segments
		WHERE name = ?
	`, body.Name)
	var exists int
	err = row.Scan(&exists)

	// if the name exist
	if err == nil {
		http.Error(w, "A segment with the same name already exist", http.StatusBadRequest)
		return
	}

	if err != sql.ErrNoRows {
		log.Printf("CreateSegment: Failed to check for duplicate segment name %s: %v", body.Name, err)
		http.Error(w, "Failed to check segment duplication", http.StatusInternalServerError)
		return
	}

	// create a db transaction
	tx, err := conn.DB.Begin()
	if err != nil {
		log.Printf("CreateSegment: failed to begin DB transaction: %v", err)
		http.Error(w, "Failed to initiate transaction for database", http.StatusInternalServerError)
		return
	}

	type Segment struct {
		ID        uuid.UUID `json:"id"`
		Name      string    `json:"name"`
		CreatedAt time.Time `json:"created_at"`
	}
	type Rule struct {
		ID        uuid.UUID `json:"id"`
		Attribute string    `json:"attribute"`
		Operator  string    `json:"operator"`
		Value     string    `json:"value"`
		CreatedAt time.Time `json:"created_at"`
	}

	type Response struct {
		Segment Segment `json:"segment"`
		Rules   []Rule  `json:"rules"`
	}
	var response Response

	segmentID := uuid.New()
	now := time.Now()

	_, err = tx.ExecContext(ctx, `
		INSERT INTO segments(id, name, created_at)
		VALUES (?, ?, ?)
	`, segmentID, body.Name, now)
	if err != nil {
		rollbackWithError(w, tx, "Failed to insert segment", "Failed to insert segment", err)
		return
	}

	// write the segment data on response body
	response.Segment.ID = segmentID
	response.Segment.Name = body.Name
	response.Segment.CreatedAt = now

	stmt, err := tx.PrepareContext(ctx, `INSERT INTO segment_rules(id, seg_id, attribute, operator, value, created_at) VALUES (?, ?, ?, ?, ?, ?)`)
	if err != nil {
		http.Error(w, "Failed to prepare db statement", http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	rulesCreatedAt := time.Now()
	for _, rule := range body.Rules {
		var respRule Rule
		id := uuid.New()
		_, err := stmt.ExecContext(ctx, id, segmentID, rule.Attribute, rule.Operator, rule.Value, rulesCreatedAt)
		if err != nil {
			rollbackWithError(w, tx, "Failed to insert segmet rule", "Failed to create segment rule", err)
			return
		}
		respRule.ID = id
		respRule.Attribute = rule.Attribute
		respRule.Operator = rule.Operator
		respRule.Value = rule.Value
		respRule.CreatedAt = rulesCreatedAt

		response.Rules = append(response.Rules, respRule)
	}

	err = tx.Commit()
	if err != nil {
		log.Printf("CreateSegnemt: Failed to commit transaction: %v", err)
		http.Error(w, "Failed to fianalize segment creation", http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.Header().Set("Location", fmt.Sprintf("/segments/%s", body.Name))
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}
