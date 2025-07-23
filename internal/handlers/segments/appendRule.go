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
	"github.com/Enoch-Tadesse/goflag/db/models"
	"github.com/google/uuid"
)

// AppendRule handles adding a new rule to an existing segment.
// It expects JSON input with segment Name, attribute, operator, and value.
// On success, it returns the newly created rule.
func AppendRule(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Input structure expects segment name to find the segment
	var body struct {
		SegName   string `json:"segment_name"`
		Attribute string `json:"attribute"`
		Operator  string `json:"operator"`
		Value     string `json:"value"`
	}

	// Decode JSON request
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate required fields are not empty
	zeroFieldNames := zeroFields(body)
	if len(zeroFieldNames) > 0 {
		http.Error(w, fmt.Sprintf("Invalid or missing values for: %s", strings.Join(zeroFieldNames, ", ")), http.StatusBadRequest)
		return
	}

	// trim the datas
	body.SegName = strings.TrimSpace(body.SegName)
	body.Attribute = strings.TrimSpace(body.Attribute)
	body.Operator = strings.TrimSpace(body.Operator)
	body.Value = strings.TrimSpace(body.Value)

	// Find segment by name, scan fields explicitly
	var seg models.Segment
	row := conn.DB.QueryRowContext(ctx, `SELECT id, name, created_at, updated_at FROM segments WHERE name = ?`, body.SegName)
	if err := row.Scan(&seg.ID, &seg.Name, &seg.CreatedAt, &seg.UpdatedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, fmt.Sprintf("Segment %s does not exist", body.SegName), http.StatusBadRequest)
			return
		}
		log.Printf("AppendRule: Failed to check if segment %s exists: %v", body.SegName, err)
		http.Error(w, "Failed to check if segment exists", http.StatusInternalServerError)
		return
	}

	// Prepare new rule ID and timestamp
	ruleID := uuid.New()
	now := time.Now()

	// Insert rule into segment_rules table
	_, err := conn.DB.ExecContext(ctx, `
		INSERT INTO segment_rules (id, seg_id, attribute, operator, value, created_at)
		VALUES (?, ?, ?, ?, ?, ?)
	`, ruleID, seg.ID, body.Attribute, body.Operator, body.Value, now)
	if err != nil {
		log.Printf("AppendRule: Failed to insert segment rule: %v", err)
		http.Error(w, "Failed to append segment rule", http.StatusInternalServerError)
		return
	}

	// Update segment's updated_at timestamp
	_, err = conn.DB.ExecContext(ctx, `
		UPDATE segments SET updated_at = ? WHERE id = ?
	`, now, seg.ID)
	if err != nil {
		log.Printf("AppendRule: Failed to update segment timestamp: %v", err)
		http.Error(w, "Failed to update segment timestamp", http.StatusInternalServerError)
		return
	}

	// Send response with the created rule
	response := struct {
		Message string             `json:"message"`
		Rule    models.SegmentRule `json:"rule"`
	}{
		Message: "Rule appended successfully",
		Rule: models.SegmentRule{
			ID:        ruleID,
			SegID:     seg.ID,
			Attribute: body.Attribute,
			Operator:  body.Operator,
			Value:     body.Value,
			CreatedAt: now,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}
