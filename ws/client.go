package ws

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/websocket"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

var upgrader = websocket.Upgrader{
	CheckOrigin:     func(r *http.Request) bool { return true },
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type token struct {
	Token string `json:"token"`
}

// connection is an middleman between the websocket connection and the hub.
type connection struct {
	h *hub
	// The websocket connection.
	ws *websocket.Conn

	// Buffered channel of outbound messages.
	send chan []byte
}

// readPump pumps messages from the websocket connection to the hub.
func (s subscription) readPump() {
	c := s.conn
	defer func() {
		c.h.Unregister <- s
		c.ws.Close()
	}()
	c.ws.SetReadLimit(maxMessageSize)
	c.ws.SetReadDeadline(time.Now().Add(pongWait))
	c.ws.SetPongHandler(func(string) error { c.ws.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, msg, err := c.ws.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
				log.Printf("error: %v", err)
			}
			break
		}
		m := Message{msg, s.room}
		c.h.Broadcast <- m
	}
}

// write writes a message with the given message type and payload.
func (c *connection) write(mt int, payload []byte) error {
	c.ws.SetWriteDeadline(time.Now().Add(writeWait))
	return c.ws.WriteMessage(mt, payload)
}

// writePump pumps messages from the hub to the websocket connection.
func (s *subscription) writePump() {
	c := s.conn
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.ws.Close()
	}()
	for {
		select {
		case message, ok := <-c.send:
			if !ok {
				c.write(websocket.CloseMessage, []byte{})
				return
			}
			if err := c.write(websocket.TextMessage, message); err != nil {
				return
			}
		case <-ticker.C:
			if err := c.write(websocket.PingMessage, []byte{}); err != nil {
				return
			}
		}
	}
}

// serveWs handles websocket requests from the peer.
func ServeWs(hub *hub, w http.ResponseWriter, r *http.Request, roomId string) {
	// token := r.URL.Query().Get("token")
	// fmt.Println(token)
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err.Error())
		return
	}
	t := token{}
	err = ws.ReadJSON(&t)
	if err != nil {
		log.Println(err)
	}

	fmt.Println(string(t.Token))

	token, err := ValidateToken(string(t.Token))
	if err != nil || !token.Valid {
		w.WriteHeader(http.StatusUnauthorized)
		ws.Close()
	}

	c := &connection{h: hub, send: make(chan []byte, 256), ws: ws}
	s := subscription{c, roomId}
	c.h.Register <- s
	go s.writePump()
	go s.readPump()
}

func ServeSMS(hub *hub, w http.ResponseWriter, r *http.Request, roomId string) {
	// token := r.URL.Query().Get("token")
	// fmt.Println(token)
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err.Error())
		return
	}
	ws.WriteMessage(1, []byte("connected"))
	t := token{}
	err = ws.ReadJSON(&t)
	if err != nil {
		log.Println(err)
	}

	fmt.Println(string(t.Token))

	token, err := ValidateToken(string(t.Token))
	if err != nil || !token.Valid {
		w.WriteHeader(http.StatusUnauthorized)
		ws.Close()
	}
	c := &connection{h: hub, send: make(chan []byte, 256), ws: ws}
	s := subscription{c, roomId}
	c.h.Register <- s
	go s.writePump()
	go s.readPump()
}

func ValidateToken(token string) (*jwt.Token, error) {
	secretKey := "secret"
	return jwt.Parse(token, func(t_ *jwt.Token) (interface{}, error) {
		if _, ok := t_.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method %v", t_.Header["alg"])
		}
		return []byte(secretKey), nil
	})
}
