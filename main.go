package main

import (
	"embed"
	"encoding/json"
	"io/fs"
	"log"
	"net/http"
	"os"
	"os/exec"
	"sync"

	"github.com/creack/pty"
	"github.com/gorilla/websocket"
)

//go:embed static
var staticFiles embed.FS

// ResizeMsg 前端调整窗口大小指令
type ResizeMsg struct {
	Type string `json:"type"`
	Cols int    `json:"cols"`
	Rows int    `json:"rows"`
}

// Session 代表一个运行中的终端进程
type Session struct {
	ID      string
	PTY     *os.File
	History []byte                   // 简单的内存历史缓冲区
	Clients map[*websocket.Conn]bool // 支持多个标签页同时看一个终端
	mu      sync.Mutex               // 修复：将字段名从 Lock 改为 mu，避免命名冲突
}

// SessionManager 管理所有会话
var (
	sessions = make(map[string]*Session)
	mu       sync.Mutex // 全局锁
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true },
	}
)

func getOrCreateSession(id string) (*Session, error) {
	mu.Lock()
	defer mu.Unlock()

	if sess, ok := sessions[id]; ok {
		return sess, nil
	}

	// 创建新 Shell
	shell := os.Getenv("SHELL")
	if shell == "" {
		shell = "/bin/bash"
		// Windows 兼容
		if _, err := os.Stat("C:\\Windows\\System32\\WindowsPowerShell\\v1.0\\powershell.exe"); err == nil {
			shell = "powershell.exe"
		}
	}

	cmd := exec.Command(shell)
	cmd.Env = append(os.Environ(), "TERM=xterm-256color")

	ptmx, err := pty.Start(cmd)
	if err != nil {
		return nil, err
	}

	sess := &Session{
		ID:      id,
		PTY:     ptmx,
		History: make([]byte, 0),
		Clients: make(map[*websocket.Conn]bool),
	}
	sessions[id] = sess

	// 启动一个协程，专门读取 PTY 数据并广播给所有连接的 WebSocket
	go func() {
		buf := make([]byte, 4096)
		for {
			n, err := ptmx.Read(buf)
			if err != nil {
				// 进程结束（如输入 exit），清理会话
				mu.Lock()
				delete(sessions, id)
				mu.Unlock()

				// 关闭所有客户端连接
				sess.mu.Lock() // 修复：使用 sess.mu.Lock()
				for conn := range sess.Clients {
					conn.Close()
				}
				sess.mu.Unlock() // 修复：使用 sess.mu.Unlock()
				return
			}

			data := make([]byte, n)
			copy(data, buf[:n])

			sess.mu.Lock() // 修复：使用 sess.mu.Lock()
			// 1. 写入历史记录
			if len(sess.History) < 1024*1024 {
				sess.History = append(sess.History, data...)
			} else {
				// 简单的截断策略：保留后半部分
				sess.History = append(sess.History[512*1024:], data...)
			}

			// 2. 广播给所有客户端
			for conn := range sess.Clients {
				err := conn.WriteMessage(websocket.BinaryMessage, data)
				if err != nil {
					conn.Close()
					delete(sess.Clients, conn)
				}
			}
			sess.mu.Unlock() // 修复：使用 sess.mu.Unlock()
		}
	}()

	return sess, nil
}

func handleTerminal(w http.ResponseWriter, r *http.Request) {
	// 获取终端 ID
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "Missing session ID", http.StatusBadRequest)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade error:", err)
		return
	}

	sess, err := getOrCreateSession(id)
	if err != nil {
		conn.WriteMessage(websocket.TextMessage, []byte("Error: "+err.Error()))
		conn.Close()
		return
	}

	// 注册客户端
	sess.mu.Lock() // 修复：使用 sess.mu.Lock()
	sess.Clients[conn] = true
	// 连接建立时，立即发送历史记录回放
	if len(sess.History) > 0 {
		conn.WriteMessage(websocket.BinaryMessage, sess.History)
	}
	sess.mu.Unlock() // 修复：使用 sess.mu.Unlock()

	// 退出时清理
	defer func() {
		sess.mu.Lock() // 修复：使用 sess.mu.Lock()
		delete(sess.Clients, conn)
		sess.mu.Unlock() // 修复：使用 sess.mu.Unlock()
		conn.Close()
	}()

	// 处理 WebSocket 输入
	for {
		msgType, message, err := conn.ReadMessage()
		if err != nil {
			break
		}

		if msgType == websocket.BinaryMessage {
			// 直接写入 PTY
			sess.PTY.Write(message)
		} else {
			var msg ResizeMsg
			if err := json.Unmarshal(message, &msg); err == nil {
				if msg.Type == "resize" {
					pty.Setsize(sess.PTY, &pty.Winsize{
						Rows: uint16(msg.Rows),
						Cols: uint16(msg.Cols),
					})
				}
			}
		}
	}
}

func main() {
	assets, err := fs.Sub(staticFiles, "static")
	if err != nil {
		log.Fatal(err)
	}

	http.Handle("/", http.FileServer(http.FS(assets)))
	http.HandleFunc("/ws", handleTerminal)

	log.Println("🚀 Persistent ssho running at: http://localhost:8080")
	if err := http.ListenAndServe("0.0.0.0:8080", nil); err != nil {
		log.Fatal(err)
	}
}
