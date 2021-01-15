package util

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"hash/crc32"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"os"
	"runtime"
	"runtime/debug"
	"strconv"
	"strings"
	"sync"
	"time"
	"unsafe"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

var (
	HEADER_AMR  = []byte("#!AMR")
	HEADER_SILK = []byte("\x02#!SILK_V3")
)

var httpClient = http.Client{
	Timeout: 15 * time.Second,
}

func GetBytes(url string) ([]byte, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header["User-Agent"] = []string{"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/83.0.4103.116 Safari/537.36 Edg/83.0.478.61"}
	resp, err := httpClient.Do(req)
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

func ToGlobalId(code int64, msgId int32) int32 {
	return int32(crc32.ChecksumIEEE([]byte(fmt.Sprintf("%d-%d", code, msgId))))
}

func IsAMRorSILK(b []byte) bool {
	return bytes.HasPrefix(b, HEADER_AMR) || bytes.HasPrefix(b, HEADER_SILK)
}

func SafeGo(fn func()) {
	go func() {
		defer func() {
			e := recover()
			if e != nil {
				log.Errorf("err recovered: %+v", e)
				log.Errorf("%s", debug.Stack())
			}
		}()
		fn()
	}()
}

func FatalError(err error) {
	log.Errorf(err.Error())
	buf := make([]byte, 64<<10)
	buf = buf[:runtime.Stack(buf, false)]
	sBuf := string(buf)
	log.Errorf(sBuf)
	time.Sleep(5 * time.Second)
	os.Exit(0)
}

func MustMd5(text string) string {
	w := md5.New()
	_, err := io.WriteString(w, text)
	if err != nil {
		panic(err)
	}
	return fmt.Sprintf("%x", w.Sum(nil))
}

var (
	client = &http.Client{
		Timeout: time.Second * 120,
		Transport: &http.Transport{
			Proxy: func(request *http.Request) (u *url.URL, e error) {
				if Proxy == "" {
					return http.ProxyFromEnvironment(request)
				}
				return url.Parse(Proxy)
			},
			DialContext: (&net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 30 * time.Second,
			}).DialContext,
			ForceAttemptHTTP2:     true,
			IdleConnTimeout:       90 * time.Second,
			TLSHandshakeTimeout:   10 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
			MaxConnsPerHost:       0,
			MaxIdleConns:          0,
			MaxIdleConnsPerHost:   999,
		},
	}
	Proxy string

	ErrOverSize = errors.New("oversize")

	UserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/87.0.4280.88 Safari/537.36 Edg/87.0.664.66"
)

func DownloadFile(url, path string, limit int64, headers map[string]string) error {
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	defer file.Close()
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	if headers != nil {
		for k, v := range headers {
			req.Header.Set(k, v)
		}
	}
	if _, ok := headers["User-Agent"]; ok {
		req.Header["User-Agent"] = []string{UserAgent}
	}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if limit > 0 && resp.ContentLength > limit {
		return ErrOverSize
	}
	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return err
	}
	return nil
}

func DownloadFileMultiThreading(url, path string, limit int64, threadCount int, headers map[string]string) error {
	if threadCount < 2 {
		return DownloadFile(url, path, limit, headers)
	}
	type BlockMetaData struct {
		BeginOffset    int64
		EndOffset      int64
		DownloadedSize int64
	}
	var blocks []*BlockMetaData
	var contentLength int64
	errUnsupportedMultiThreading := errors.New("unsupported multi-threading")
	// 初始化分块或直接下载
	initOrDownload := func() error {
		copyStream := func(s io.ReadCloser) error {
			file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE, 0666)
			if err != nil {
				return err
			}
			defer file.Close()
			if _, err = io.Copy(file, s); err != nil {
				return err
			}
			return errUnsupportedMultiThreading
		}
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return err
		}
		if headers != nil {
			for k, v := range headers {
				req.Header.Set(k, v)
			}
		}
		if _, ok := headers["User-Agent"]; ok {
			req.Header["User-Agent"] = []string{UserAgent}
		}
		req.Header.Set("range", "bytes=0-")
		resp, err := client.Do(req)
		if err != nil {
			return err
		}
		if resp.StatusCode < 200 || resp.StatusCode >= 300 {
			return errors.New("response status unsuccessful: " + strconv.FormatInt(int64(resp.StatusCode), 10))
		}
		if resp.StatusCode == 200 {
			if limit > 0 && resp.ContentLength > limit {
				return ErrOverSize
			}
			return copyStream(resp.Body)
		}
		if resp.StatusCode == 206 {
			contentLength = resp.ContentLength
			if limit > 0 && resp.ContentLength > limit {
				return ErrOverSize
			}
			blockSize := func() int64 {
				if contentLength > 1024*1024 {
					return (contentLength / int64(threadCount)) - 10
				} else {
					return contentLength
				}
			}()
			if blockSize == contentLength {
				return copyStream(resp.Body)
			}
			var tmp int64
			for tmp+blockSize < contentLength {
				blocks = append(blocks, &BlockMetaData{
					BeginOffset: tmp,
					EndOffset:   tmp + blockSize - 1,
				})
				tmp += blockSize
			}
			blocks = append(blocks, &BlockMetaData{
				BeginOffset: tmp,
				EndOffset:   contentLength - 1,
			})
			return nil
		}
		return errors.New("unknown status code.")
	}
	// 下载分块
	downloadBlock := func(block *BlockMetaData) error {
		req, _ := http.NewRequest("GET", url, nil)
		file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			return err
		}
		defer file.Close()
		_, _ = file.Seek(block.BeginOffset, io.SeekStart)
		writer := bufio.NewWriter(file)
		defer writer.Flush()
		if headers != nil {
			for k, v := range headers {
				req.Header.Set(k, v)
			}
		}
		if _, ok := headers["User-Agent"]; ok {
			req.Header["User-Agent"] = []string{UserAgent}
		}
		req.Header.Set("range", "bytes="+strconv.FormatInt(block.BeginOffset, 10)+"-"+strconv.FormatInt(block.EndOffset, 10))
		resp, err := client.Do(req)
		if err != nil {
			return err
		}
		defer resp.Body.Close()
		if resp.StatusCode < 200 || resp.StatusCode >= 300 {
			return errors.New("response status unsuccessful: " + strconv.FormatInt(int64(resp.StatusCode), 10))
		}
		var buffer = make([]byte, 1024)
		i, err := resp.Body.Read(buffer)
		for {
			if err != nil && err != io.EOF {
				return err
			}
			i64 := int64(len(buffer[:i]))
			needSize := block.EndOffset + 1 - block.BeginOffset
			if i64 > needSize {
				i64 = needSize
				err = io.EOF
			}
			_, e := writer.Write(buffer[:i64])
			if e != nil {
				return e
			}
			block.BeginOffset += i64
			block.DownloadedSize += i64
			if err == io.EOF || block.BeginOffset > block.EndOffset {
				break
			}
			i, err = resp.Body.Read(buffer)
		}
		return nil
	}

	if err := initOrDownload(); err != nil {
		if err == errUnsupportedMultiThreading {
			return nil
		}
		return err
	}
	wg := sync.WaitGroup{}
	wg.Add(len(blocks))
	var lastErr error
	for i := range blocks {
		go func(b *BlockMetaData) {
			defer wg.Done()
			if err := downloadBlock(b); err != nil {
				lastErr = err
			}
		}(blocks[i])
	}
	wg.Wait()
	return lastErr
}
