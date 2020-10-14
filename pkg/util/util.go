package util

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"hash/crc32"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"unsafe"
)

func GetBytes(url string) ([]byte, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header["User-Agent"] = []string{"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/83.0.4103.116 Safari/537.36 Edg/83.0.478.61"}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if strings.Contains(resp.Header.Get("Content-Encoding"), "gzip") {
		buffer := bytes.NewBuffer(body)
		r, _ := gzip.NewReader(buffer)
		defer r.Close()
		unCom, err := ioutil.ReadAll(r)
		return unCom, err
	}
	return body, nil
}


func PathExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil || os.IsExist(err)
}

func ReadAllText(path string) string {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return ""
	}
	return string(b)
}

func DelFile(path string) bool {
	err := os.Remove(path)
	if err != nil {
		// 删除失败
		log.Error(err)
		return false
	} else {
		// 删除成功
		log.Info(path + "删除成功")
		return true
	}
}

func MustMarshal(t interface{}) string {
	if b, err := json.Marshal(t); err != nil {
		return ""
	} else {
		return ByteSliceToString(b)
	}
}

func ByteSliceToString(bs []byte) string {
	return *(*string)(unsafe.Pointer(&bs))
}

func Check(err error) {
	if err != nil {
		log.Fatalf("遇到错误: %v", err)
	}
}

func ToGlobalId(code int64, msgId int32) int32 {
	return int32(crc32.ChecksumIEEE([]byte(fmt.Sprintf("%d-%d", code, msgId))))
}