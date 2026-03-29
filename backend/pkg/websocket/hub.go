package websocket

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

// Message WebSocket 廣播訊息格式
type Message struct {
	Event   string          `json:"event"`
	Payload json.RawMessage `json:"payload"`
}

// Hub WebSocket 連線廣播中心（Singleton）
type Hub struct {
	clients   map[*client]struct{}
	broadcast chan Message
	register  chan *client
	unregister chan *client
	mu        sync.RWMutex
}

type client struct {
	hub  *Hub
	conn *websocket.Conn
	send chan Message
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		// Phase 3 可加入 Origin 驗證
		return true
	},
}

// NewHub 建立 Hub 實例
func NewHub() *Hub {
	return &Hub{
		clients:    make(map[*client]struct{}),
		broadcast:  make(chan Message, 256),
		register:   make(chan *client),
		unregister: make(chan *client),
	}
}

// Run 啟動事件迴圈（goroutine）
func (h *Hub) Run() {
	for {
		select {
		case c := <-h.register:
			h.mu.Lock()
			h.clients[c] = struct{}{}
			h.mu.Unlock()
			log.Printf("[WS] 新連線，目前 %d 個客戶端", len(h.clients))

		case c := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[c]; ok {
				delete(h.clients, c)
				close(c.send)
			}
			h.mu.Unlock()
			log.Printf("[WS] 連線關閉，目前 %d 個客戶端", len(h.clients))

		case msg := <-h.broadcast:
			h.mu.RLock()
			for c := range h.clients {
				select {
				case c.send <- msg:
				default:
					// 緩衝區滿，移除該客戶端
					close(c.send)
					delete(h.clients, c)
				}
			}
			h.mu.RUnlock()
		}
	}
}

// Broadcast 向所有連線的客戶端廣播事件
func (h *Hub) Broadcast(event string, payload any) {
	b, err := json.Marshal(payload)
	if err != nil {
		log.Printf("[WS] payload 序列化失敗: %v", err)
		return
	}
	h.broadcast <- Message{Event: event, Payload: b}
}

// ServeWS HTTP Upgrade 處理器（掛到 /ws 路由）
func (h *Hub) ServeWS(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("[WS] Upgrade 失敗: %v", err)
		return
	}
	c := &client{hub: h, conn: conn, send: make(chan Message, 64)}
	h.register <- c

	go c.writePump()
	go c.readPump()
}

// readPump 讀取客戶端訊息（目前只處理 ping/close）
func (c *client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()
	for {
		if _, _, err := c.conn.ReadMessage(); err != nil {
			break
		}
	}
}

// writePump 將 send channel 的訊息寫入 WebSocket
func (c *client) writePump() {
	defer c.conn.Close()
	for msg := range c.send {
		if err := c.conn.WriteJSON(msg); err != nil {
			break
		}
	}
}
