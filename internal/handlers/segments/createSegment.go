package segments

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	conn "github.com/Enoch-Tadesse/goflag/db/connection"
	"github.com/google/uuid"
)

// rule represents a condition that defines which users belong in a segment.
type rule struct {
	Attribute string `json:"attribute"`
	Operator  string `json:"operator"`
	Value     string `json:"value"`
}

// CreateSegment handles the creation of a new segment with rules.
func CreateSegment(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Parse and validate the request body
	var body struct {
		Name  string `json:"name"`
		Rules []rule `json:"rules"`
	}

	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	body.Name = strings.TrimSpace(body.Name)

	if body.Name == "" {
		http.Error(w, "Name is required", http.StatusBadRequest)
		return
	}

	// check if there exists atleast one rule
	if len(body.Rules) == 0 {
		http.Error(w, "A segment must have at least one rule", http.StatusBadRequest)
		return
	}

	// Validate each rule in the segment
	for _, rule := range body.Rules {
		// trim the datas
		rule.Attribute = strings.TrimSpace(rule.Attribute)
		rule.Value = strings.TrimSpace(rule.Value)
		rule.Operator = strings.TrimSpace(rule.Operator)
		if rule.Attribute == "" || rule.Operator == "" || rule.Value == "" {
			http.Error(w, "All rules must have non-empty attribute, operator, and value", http.StatusBadRequest)
			return
		}
	}

	// Check if a segment with the same name already exists
	row := conn.DB.QueryRowContext(ctx, `
		SELECT 1 FROM segments WHERE name = ?
	`, body.Name)

	var exists int
	err = row.Scan(&exists)

	if err == nil {
		http.Error(w, "A segment with the same name already exists", http.StatusBadRequest)
		return
	}

	if err != sql.ErrNoRows {
		// An unexpected error occurred when querying for segment name
		log.Printf("CreateSegment: Failed to check for duplicate segment name %s: %v", body.Name, err)
		http.Error(w, "Failed to check segment duplication", http.StatusInternalServerError)
		return
	}

	// Begin DB transaction
	tx, err := conn.DB.Begin()
	if err != nil {
		log.Printf("CreateSegment: Failed to begin DB transaction: %v", err)
		http.Error(w, "Failed to initiate transaction for database", http.StatusInternalServerError)
		return
	}

	// Define local response structures
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

	// Insert the segment into the `segments` table
	_, err = tx.ExecContext(ctx, `
		INSERT INTO segments(id, name, created_at)
		VALUES (?, ?, ?)
	`, segmentID, body.Name, now)
	if err != nil {
		rollbackWithError(w, tx, "Failed to insert segment", "Failed to insert segment", err)
		return
	}

	// Populate segment response data
	response.Segment = Segment{
		ID:        segmentID,
		Name:      body.Name,
		CreatedAt: now,
	}

	// Prepare statement for inserting multiple rules efficiently
	stmt, err := tx.PrepareContext(ctx, `
		INSERT INTO segment_rules(id, seg_id, attribute, operator, value, created_at)
		VALUES (?, ?, ?, ?, ?, ?)
	`)
	if err != nil {
		http.Error(w, "Failed to prepare DB statement", http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	rulesCreatedAt := time.Now()
	for _, rule := range body.Rules {
		var respRule Rule
		ruleID := uuid.New()

		// Insert each rule into the `segment_rules` table
		_, err := stmt.ExecContext(ctx, ruleID, segmentID, rule.Attribute, rule.Operator, rule.Value, rulesCreatedAt)
		if err != nil {
			rollbackWithError(w, tx, "Failed to insert segment rule", "Failed to create segment rule", err)
			return
		}

		// Populate rule response
		respRule = Rule{
			ID:        ruleID,
			Attribute: rule.Attribute,
			Operator:  rule.Operator,
			Value:     rule.Value,
			CreatedAt: rulesCreatedAt,
		}
		response.Rules = append(response.Rules, respRule)
	}

	// Commit the transaction
	err = tx.Commit()
	if err != nil {
		log.Printf("CreateSegment: Failed to commit transaction: %v", err)
		http.Error(w, "Failed to finalize segment creation", http.StatusInternalServerError)
		return
	}

	// Return success response with Location header
	w.Header().Add("Content-Type", "application/json")
	w.Header().Set("Location", fmt.Sprintf("/segments/%s", body.Name))
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}
