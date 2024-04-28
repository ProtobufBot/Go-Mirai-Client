package bot

import (
	"errors"
	"os"

	log "github.com/sirupsen/logrus"
)

//go:generate go run github.com/a8m/syncmap -o "gen_client_map.go" -pkg bot -name ClientMap "map[int64]*client.QQClient"
//go:generate go run github.com/a8m/syncmap -o "gen_token_map.go" -pkg bot -name TokenMap "map[int64][]byte"
var (
	Clients                    ClientMap
	LoginTokens                TokenMap
	EnergyCount                = 0
	EnergyStop                 = false
	SignCount                  = 0
	SignStop                   = false
	RegisterSignCount          = 0
	RegisterSignStop           = false
	SubmitRequestCallbackCount = 0
	SubmitRequestCallbackStop  = false
	RequestTokenCount          = 0
	RequestTokenStop           = false
	DestoryInstanceCount       = 0
	DestoryInstanceStop        = false
	RSR                        RequestSignResult
	GTL                        *GMCLogin
	SR                         SignRegister
	IsRequestTokenAgain        bool = false
	TTI_i                           = 30
)

type Logger struct {
}

type GMCLogin struct {
	DeviceSeed     int64
	ClientProtocol int32
	SignServer     string
	SignServerKey  string
}

type SignRegister struct {
	Uin       uint64
	AndroidId string
	Guid      string
	Qimei36   string
	Key       string
}

type RequestCallback struct {
	Cmd        string `json:"cmd,omitempty"` // trpc.o3.ecdh_access.EcdhAccess.SsoSecureA2Establish
	Body       string `json:"body,omitempty"`
	CallBackId int    `json:"callbackId,omitempty"`
}

type RequestSignData struct {
	Token           string `json:"token,omitempty"`
	Extra           string `json:"extra,omitempty"`
	Sign            string `json:"sign,omitempty"`
	O3dId           string `json:"o3did,omitempty"`
	RequestCallback []*RequestCallback
}

type RequestSignResult struct {
	Code int    `json:"code,omitempty"`
	Msg  string `json:"msg,omitempty"`
	Data *RequestSignData
}

func PathExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil || errors.Is(err, os.ErrExist)
}

func (l *Logger) Info(format string, args ...any) {
	log.Infof(format, args)
}
func (l *Logger) Warning(format string, args ...any) {
	log.Warnf(format, args)
}
func (l *Logger) Error(format string, args ...any) {
	log.Errorf(format, args)
}
func (l *Logger) Debug(format string, args ...any) {
	log.Debug(format, args)
}
func (l *Logger) Dump(dumped []byte, format string, args ...any) {
}

func IsClientExist(uin int64) bool {
	_, ok := Clients.Load(uin)
	return ok
}