package server

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

// 実行するサーバのアドレスを受け取り、*http.Serverを返す
func NewHTTPServer(addr string) *http.Server {
	httpsrv := newHTTPServer()
	r := mux.NewRouter()
	r.HandleFunc("/", httpsrv.handleProduce).Methods("POST")
	r.HandleFunc("/", httpsrv.handleConsume).Methods("GET")
	return &http.Server{
		Addr:    addr,
		Handler: r,
	}
}

// サーバ構造体
type httpServer struct {
	Log *Log
}

func newHTTPServer() *httpServer {
	return &httpServer{
		Log: NewLog(),
	}
}

// ===== リクエスト構造体 =====
// APIの呼び出しもとがログに追加して欲しいレコードを含む
type ProduceRequest struct {
	Record Record `json:"record"`
}

// レスポンス構造体
// ログがどのオフセットにレコードを格納したかを伝える
type ProduceResponse struct {
	Offset uint64 `json:"offset"`
}

// APIの呼び出しもとが読み出したいレコードを指定する
type ConsumeRequest struct {
	Offset uint64 `json:"Offset"`
}

// 呼び出しもとにレコードを送り返す
type ConsumeResponse struct {
	Record Record `json:"record"`
}

// ===== サーバのハンドラ =====
func (s *httpServer) handleProduce(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	// リクエストを構造体にアンマーシャルする
	var req ProduceRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// ログにレコードを保存する
	off, err := s.Log.Append(req.Record)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// ログがレコードを保存したオフセットをレスポンスとして返す
	res := ProduceResponse{Offset: off}
	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (s *httpServer) handleConsume(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	// リクエストを構造体にアンマーシャルする
	var req ConsumeRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// ログに保存されたレコードを取得する
	record, err := s.Log.Read(req.Offset)
	if err == ErrOffsetNotFound {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// 取得したレコードをレスポンスとして返す
	res := ConsumeResponse{Record: record}
	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
