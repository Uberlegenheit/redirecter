package response

import (
	"encoding/json"
	"everstake-affiliate/apperrors"
	"net/http"
)

func File(w http.ResponseWriter, data []byte) {
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Write(data)
}

// Json writes to ResponseWriter a single JSON-object
func Json(w http.ResponseWriter, data interface{}) {
	js, err := json.Marshal(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

// JsonError writes to ResponseWriter error
func JsonError(w http.ResponseWriter, err error) {
	var e *apperrors.Error
	var ok bool

	if e, ok = err.(*apperrors.Error); !ok {
		e = apperrors.FromError(err)
	}

	js, _ := json.Marshal(e.ToMap())
	w.WriteHeader(e.GetHttpCode())
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}
