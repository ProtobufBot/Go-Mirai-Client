package handler

import (
	"errors"
	"io/ioutil"

	"github.com/gin-gonic/gin"
	"github.com/golang/protobuf/proto"
)

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func Bind(c *gin.Context, req any) error {
	buf, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		return err
	}
	if r, ok := req.(proto.Message); ok {
		if err := proto.Unmarshal(buf, r); err != nil {
			return err
		}
	} else {
		return errors.New("obj is not ProtoMessage")
	}
	return nil
}
