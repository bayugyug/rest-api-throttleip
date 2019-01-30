package controllers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/bayugyug/rest-api-throttleip/config"
	"github.com/bayugyug/rest-api-throttleip/models"
	"github.com/bayugyug/rest-api-throttleip/utils"
	"github.com/go-chi/render"
)

type APIResponse struct {
	Code   int
	Status string
}

type ApiHandler struct {
}

func (api *ApiHandler) IndexPage(w http.ResponseWriter, r *http.Request) {

	//check ip details
	tracker := models.NewTrackerIP()
	trkInfo := tracker.GetIPInfo(ApiInstance.Context, r)

	//206
	if trkInfo == nil {
		api.ReplyErrContent(w, r, http.StatusPartialContent, http.StatusText(http.StatusPartialContent))
		return
	}

	//check
	if !api.CheckIPInfo(w, r, trkInfo) {
		return
	}

	//good
	render.JSON(w, r, APIResponse{
		Code:   200,
		Status: "Welcome!",
	})
}

func (api *ApiHandler) DummyReqGet(w http.ResponseWriter, r *http.Request) {

	//check ip details
	tracker := models.NewTrackerIP()
	trkInfo := tracker.GetIPInfo(ApiInstance.Context, r)

	//206
	if trkInfo == nil {
		api.ReplyErrContent(w, r, http.StatusPartialContent, http.StatusText(http.StatusPartialContent))
		return
	}

	//check
	if !api.CheckIPInfo(w, r, trkInfo) {
		return
	}

	//good
	render.JSON(w, r, APIResponse{
		Code:   200,
		Status: "DummyReqGet::Welcome",
	})
}

func (api *ApiHandler) DummyReqPost(w http.ResponseWriter, r *http.Request) {

	//check ip details
	tracker := models.NewTrackerIP()
	trkInfo := tracker.GetIPInfo(ApiInstance.Context, r)

	//206
	if trkInfo == nil {
		api.ReplyErrContent(w, r, http.StatusPartialContent, http.StatusText(http.StatusPartialContent))
		return
	}

	//check
	if !api.CheckIPInfo(w, r, trkInfo) {
		trkInfo.Status = "Denied"
		utils.Dumper(trkInfo)
		return
	}

	//good
	render.JSON(w, r, APIResponse{
		Code:   200,
		Status: "DummyReqPost::Welcome",
	})
}

func (api *ApiHandler) DummyReqPut(w http.ResponseWriter, r *http.Request) {

	//check ip details
	tracker := models.NewTrackerIP()
	trkInfo := tracker.GetIPInfo(ApiInstance.Context, r)

	//206
	if trkInfo == nil {
		api.ReplyErrContent(w, r, http.StatusPartialContent, http.StatusText(http.StatusPartialContent))
		return
	}

	//check
	if !api.CheckIPInfo(w, r, trkInfo) {
		return
	}

	//good
	render.JSON(w, r, APIResponse{
		Code:   200,
		Status: "DummyReqPut::Welcome",
	})
}

func (api *ApiHandler) DummyReqDelete(w http.ResponseWriter, r *http.Request) {

	//check ip details
	tracker := models.NewTrackerIP()
	trkInfo := tracker.GetIPInfo(ApiInstance.Context, r)

	//206
	if trkInfo == nil {
		api.ReplyErrContent(w, r, http.StatusPartialContent, http.StatusText(http.StatusPartialContent))
		return
	}

	//check
	if !api.CheckIPInfo(w, r, trkInfo) {
		return
	}

	//good
	render.JSON(w, r, APIResponse{
		Code:   200,
		Status: "DummyReqDelete::Welcome",
	})
}

//ReplyErrContent send 204 msg
//
//  http.StatusNoContent
//  http.StatusText(http.StatusNoContent)
func (api ApiHandler) ReplyErrContent(w http.ResponseWriter, r *http.Request, code int, msg string) {
	render.JSON(w, r, APIResponse{
		Code:   code,
		Status: msg,
	})
}

//CheckIPInfo check history and send error message
func (api *ApiHandler) CheckIPInfo(w http.ResponseWriter, r *http.Request, trk *models.TrackerIP) bool {

	//check
	tot := ApiInstance.IPHistory.GetIP(trk.IP)
	log.Println("IP Total:", trk.IP, tot)

	//check max reached
	if tot > config.RequestsPerMinute {
		trk.Status = "Denied"
		//save to logs
		api.SaveIPInfo(w, r, trk)
		//reply
		api.ReplyErrContent(w, r,
			http.StatusConflict,
			fmt.Sprintf("IP is not allowed. Already reached %d/%d per minute.", tot, config.RequestsPerMinute))
		return false
	}
	//save logs
	api.SaveIPInfo(w, r, trk)
	return true

}

func (api *ApiHandler) SaveIPInfo(w http.ResponseWriter, r *http.Request, trk *models.TrackerIP) {
	//pipe to redis
	ApiInstance.IPHistory.HistoryChannel <- trk
	utils.Dumper(trk)
}
