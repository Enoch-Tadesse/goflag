package segments

import (
	"database/sql"
	"log"
	"net/http"
)

func rollbackWithError(w http.ResponseWriter, tx *sql.Tx, logMsg string, errMsg string, err error) {
	log.Printf("CreateSegment: %s: %v", logMsg, err)
	if rbErr := tx.Rollback(); rbErr != nil {
		log.Printf("CreateSegment: Rollback also failed: %v", rbErr)
	}
	http.Error(w, errMsg, http.StatusInternalServerError)
}
