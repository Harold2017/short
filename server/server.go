package server

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"short/short"
	"short/utils"
	"strings"
)

type Req struct {
	ReqURL string `json:"req_url"`
}

type Resp struct {
	RespURL string `json:"resp_url"`
}

type ErrResp struct {
	Msg string `json:"msg"`
}

func Short(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println("short handler request error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		errMsg, _ := json.Marshal(ErrResp{Msg: http.StatusText(http.StatusInternalServerError)})
		w.Write(errMsg)
	}

	var req Req
	err = json.Unmarshal(body, &req)
	if err != nil {
		log.Println("short handler parse error:", err)
		w.WriteHeader(http.StatusBadRequest)
		errMsg, _ := json.Marshal(ErrResp{Msg: http.StatusText(http.StatusBadRequest)})
		w.Write(errMsg)
		return
	}

	longURL, err := url.Parse(req.ReqURL)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		errMsg, _ := json.Marshal(ErrResp{Msg: "request url is invalid"})
		w.Write(errMsg)
		return
	}
	if longURL.Host == utils.Conf.Host {
		w.WriteHeader(http.StatusBadRequest)
		errMsg, _ := json.Marshal(ErrResp{Msg: "request url is already shortened"})
		w.Write(errMsg)
		return
	}

	shortURL, err := short.DefaultShorter.Short(req.ReqURL)
	if err != nil {
		log.Println("short handler short error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		errMsg, _ := json.Marshal(ErrResp{Msg: http.StatusText(http.StatusInternalServerError)})
		w.Write(errMsg)
		return
	}
	shortURL = (&url.URL{
		Scheme: "http",
		Host:   utils.Conf.Host,
		Path:   shortURL,
	}).String()
	resp, _ := json.Marshal(Resp{RespURL: shortURL})
	w.Write(resp)
}

func Long(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println("long handler request error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		errMsg, _ := json.Marshal(ErrResp{Msg: http.StatusText(http.StatusInternalServerError)})
		w.Write(errMsg)
	}

	var req Req
	err = json.Unmarshal(body, &req)
	if err != nil {
		log.Println("long handler parse error:", err)
		w.WriteHeader(http.StatusBadRequest)
		errMsg, _ := json.Marshal(ErrResp{Msg: http.StatusText(http.StatusBadRequest)})
		w.Write(errMsg)
	}

	shortURL, err := url.Parse(req.ReqURL)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		errMsg, _ := json.Marshal(ErrResp{Msg: "request url is invalid"})
		w.Write(errMsg)
		return
	}

	longURL, err := short.DefaultShorter.Long(strings.TrimLeft(shortURL.Path, "/"))
	if err != nil {
		log.Println("long handler long error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		errMsg, _ := json.Marshal(ErrResp{Msg: http.StatusText(http.StatusInternalServerError)})
		w.Write(errMsg)
		return
	}
	resp, _ := json.Marshal(Resp{RespURL: longURL})
	w.Write(resp)
}

func Redirect(w http.ResponseWriter, r *http.Request) {
	shortURL := r.URL.Query().Get("shortURL")
	longURL, err := short.DefaultShorter.Long(shortURL)
	if err != nil {
		log.Println("redirect error: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(http.StatusText(http.StatusInternalServerError)))
		return
	}
	if len(longURL) != 0 {
		w.Header().Set("Location", longURL)
		w.WriteHeader(http.StatusTemporaryRedirect)  // 302
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
