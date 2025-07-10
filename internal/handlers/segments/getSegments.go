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
	"github.com/Enoch-Tadesse/goflag/db/models"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

func GetAllSegments(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// get the segments first
	rows, err := conn.DB.QueryContext(ctx, `
		SELECT id, name, created_at
		FROM segments
	`)

	if err != nil {
		log.Printf("GetAllSegments: Failed to fetch all segments. %v", err)
		http.Error(w, "Failed to fetch segments", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	type Segment struct {
		ID        string    `json:"id"`
		Name      string    `json:"name"`
		CreatedAt time.Time `json:"created_at"`
	}

	var segments []Segment

	// iterate over all segments and collect
	for rows.Next() {
		var segment Segment
		if err := rows.Scan(&segment.ID, &segment.Name, &segment.CreatedAt); err != nil {
			log.Printf("GetAllSegments: Error scanning segment row: %v", err)
			continue
		}
		segments = append(segments, segment)
	}

	if len(segments) == 0 {
		w.Header().Add("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"message": "No segments currently"})
		return
	}

	type rule struct {
		ID        string `json:"id"`
		Attribute string `json:"attribute"`
		Operator  string `json:"operator"`
		Value     string `json:"value"`
		CreatedAt string `json:"created_at"`
	}

	type Result struct {
		Segment Segment `json:"segment"`
		Rules   []rule  `json:"rules"`
	}

	var results struct {
		AllSegments []Result `json:"all_segments"`
	}

	for _, seg := range segments {
		result := Result{
			Segment: seg,
			Rules:   []rule{},
		}
		id := seg.ID
		// get all segment rules specific to a single segment
		innerRows, err := conn.DB.QueryContext(ctx, `
			SELECT id, attribute, operator, value, created_at
			FROM segment_rules
			WHERE seg_id = ?
		`, id)
		if err != nil {
			log.Printf("GetAllSegments: Failed to fetch rules for segment id: %s, name: %s: %v", seg.ID, seg.Name, err)
			http.Error(w, fmt.Sprintf("Failed to fetch rules for %s segment", seg.Name), http.StatusInternalServerError)
			return
		}
		// Immediate Invoked Function Expression for defer innerRows resource wastage
		func() {
			defer innerRows.Close()
			result.Segment = seg

			for innerRows.Next() {
				var rule rule
				if err := innerRows.Scan(&rule.ID, &rule.Attribute, &rule.Operator, &rule.Value, &rule.CreatedAt); err != nil {
					log.Printf("GetAllSegments: Failed to scan innerRow: %v", err)
					continue
				}
				result.Rules = append(result.Rules, rule)
			}
		}()
		results.AllSegments = append(results.AllSegments, result)
	}
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(results)
}

func GetSegmentByName(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	vars := mux.Vars(r)
	name := vars["name"]

	// get the segment by name
	row := conn.DB.QueryRowContext(ctx, `
		SELECT id, name, created_at
		FROM segments
		WHERE name = ?
		LIMIT 1
	`, name)

	var segment models.Segment

	err := row.Scan(&segment.ID, &segment.Name, &segment.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, fmt.Sprintf("Segment with name %s does not exist", name), http.StatusBadRequest)
			return
		}
		log.Printf("GetSegmentByName: Failed to check if segment %s exists: %v", name, row.Err())
		http.Error(w, "Failed to check if segment exist", http.StatusInternalServerError)
		return
	}

	// get the associalted segment rules
	rows, err := conn.DB.QueryContext(ctx, `
		SELECT id, attribute, operator, value, created_at
		FROM segment_rules
		WHERE seg_id = ?
	`, segment.ID)

	if err != nil {
		log.Printf("GetSegmentByName: failed to fetch rows: %v", err)
		http.Error(w, "Failed to fetch segment rules", http.StatusInternalServerError)
		return
	}

	defer rows.Close()

	type segmentRule struct {
		ID        uuid.UUID `json:"id"`
		Attribute string    `json:"attribute"`
		Operator  string    `json:"operator"`
		Value     string    `json:"value"`
		CreatedAt time.Time `json:"created_at"`
	}

	var rules []segmentRule

	// scan all the segment rules
	for rows.Next() {
		var rule segmentRule
		if err := rows.Scan(&rule.ID, &rule.Attribute, &rule.Operator, &rule.Value, &rule.CreatedAt); err != nil {
			if err == sql.ErrNoRows {
				log.Print("some error")
				break
			}
			log.Printf("GetSegmentByName: Failed to scan result rows: %v", err)
			http.Error(w, "Failed to fetch segment rules", http.StatusInternalServerError)
			return
		}
		rules = append(rules, rule)
	}

	var response struct {
		Segment models.Segment `json:"segment"`
		Rules   []segmentRule  `json:"rules"`
	}

	response.Segment = segment
	response.Rules = rules

	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)

}
