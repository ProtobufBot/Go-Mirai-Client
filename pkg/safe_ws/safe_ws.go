package safe_ws

import (
	"fmt"

	"github.com/ProtobufBot/Go-Mirai-Client/pkg/util"
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
)

// safe websocket
type SafeWebSocket struct {
	Conn          *websocket.Conn
	SendChannel   chan *WebSocketSendingMessage
	OnRecvMessage func(ws *SafeWebSocket, messageType int, data []byte)
	OnClose       func()
}

type WebSocketSendingMessage struct {
	MessageType int
	Data        []byte
}

func (ws *SafeWebSocket) Send(messageType int, data []byte) (e error) {
	defer func() {
		if err := recover(); err != nil { // 可能channel已被关闭，向已关闭的channel写入数据
			e = fmt.Errorf("failed to send websocket msg, %+v", err)
			log.Errorf("failed to send websocket msg, %+v", err)
			ws.Close()
		}
	}()
	ws.SendChannel <- &WebSocketSendingMessage{
		MessageType: messageType,
		Data:        data,
	}
	e = nil
	return
}

func (ws *SafeWebSocket) Close() {
	defer func() {
		_ = recover() // 可能已经关闭过channel
	}()
	_ = ws.Conn.Close()
	ws.OnClose()
	close(ws.SendChannel)
}

func NewSafeWebSocket(conn *websocket.Conn, OnRecvMessage func(ws *SafeWebSocket, messageType int, data []byte), onClose func()) *SafeWebSocket {
	ws := &SafeWebSocket{
		Conn:          conn,
		SendChannel:   make(chan *WebSocketSendingMessage, 100),
		OnRecvMessage: OnRecvMessage,
		OnClose:       onClose,
	}

	conn.SetCloseHandler(func(code int, text string) error {
		ws.Close()
		return nil
	})

	// 接受消息
	util.SafeGo(func() {
		for {
			messageType, data, err := conn.ReadMessage()
			if err != nil {
				log.Errorf("failed to read message, err: %+v", err)
				ws.Close()
				return
			}
			if messageType == websocket.PingMessage {
				if err := ws.Send(websocket.PongMessage, []byte("pong")); err != nil {
					ws.Close()
				}
				continue
			}
			util.SafeGo(func() {
				ws.OnRecvMessage(ws, messageType, data)
			})
		}
	})

	// 发送消息
	util.SafeGo(func() {
		for sendingMessage := range ws.SendChannel {
			if ws.Conn == nil {
				log.Errorf("failed to send websocket message, conn is nil")
				return
			}
			err := ws.Conn.WriteMessage(sendingMessage.MessageType, sendingMessage.Data)
			if err != nil {
				log.Errorf("failed to send websocket message, %+v", err)
				ws.Close()
				return
			}
		}
	})
	return ws
}
