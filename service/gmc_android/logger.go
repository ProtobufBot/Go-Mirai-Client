package gmc_android

import (
	"bytes"
	"strings"

	log "github.com/sirupsen/logrus"
)

type AndroidLogger interface {
	Log(str string)
}

type AndroidWriter struct {
}

func (AndroidWriter) Write(p []byte) (n int, err error) {
	if logger != nil {
		logger.Log(string(p))
	}
	return 0, nil
}

type AndroidFormatter struct {
}

func (AndroidFormatter) Format(entry *log.Entry) ([]byte, error) {
	buf := bytes.Buffer{}
	buf.WriteByte('[')
	buf.WriteString(entry.Time.Format("2006-01-02 15:04:05"))
	buf.WriteString("] [")
	buf.WriteString(strings.ToUpper(entry.Level.String()))
	buf.WriteString("]: ")
	buf.WriteString(entry.Message)
	buf.WriteString(" \n")
	buf.Bytes()
	return append([]byte(nil), buf.Bytes()...), nil
}

