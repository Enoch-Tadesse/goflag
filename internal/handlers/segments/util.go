package segments

import (
	"database/sql"
	"log"
	"net/http"
	"reflect"
)

func rollbackWithError(w http.ResponseWriter, tx *sql.Tx, logMsg string, errMsg string, err error) {
	log.Printf("CreateSegment: %s: %v", logMsg, err)
	if rbErr := tx.Rollback(); rbErr != nil {
		log.Printf("CreateSegment: Rollback also failed: %v", rbErr)
	}
	http.Error(w, errMsg, http.StatusInternalServerError)
}

func zeroFields(data any, exclude ...string) []string {
	v := reflect.ValueOf(data)
	t := reflect.TypeOf(data)

	// dereference if data is a pointer
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
		t = t.Elem()
	}

	excluded := make(map[string]bool)
	for _, name := range exclude {
		excluded[name] = true
	}

	var zeroFieldNames []string

	for i := range v.NumField() {
		field := v.Field(i)
		fieldType := t.Field(i)

		// skip unexported fields
		if !field.CanInterface() {
			continue
		}

		// skip excluded fields
		if excluded[fieldType.Name] {
			continue
		}

		// built-in checker for zero values
		if field.IsZero() {
			zeroFieldNames = append(zeroFieldNames, fieldType.Name)
		}

	}
	return zeroFieldNames

}
