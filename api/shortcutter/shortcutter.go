package shortcutter

import (
	"bytes"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"net/http"
)

const RealIpHeader = "X-Forwarded-For"

func (api *API) DefaultRedirect(w http.ResponseWriter, r *http.Request) {
	log.Debug("DefaultRedirect")
	//Use permanent redirect only for root path
	//Browser save redirects and never make this request again
	http.Redirect(w, r, api.config.DefaultRedirectURL, http.StatusMovedPermanently)
}

func (api *API) GetShortCutInfo(w http.ResponseWriter, r *http.Request) {

	var bt bytes.Buffer
	var err error
	urlID := mux.Vars(r)["url_id"]
	log.Debugf("ShortCut url_id %s Header %+v", urlID, r.Header)

	bt, err = api.service.GetShortCutInfo(urlID)
	if err != nil {
		log.Errorf("GetShortCutInfo error: %s", err.Error())
	}

	_, err = w.Write(bt.Bytes())
	if err != nil {
		log.Errorf("Write error: %s", err.Error())
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	//Cache for 30 days
	w.Header().Set("Cache-Control", "max-age:2592000, public")
}
