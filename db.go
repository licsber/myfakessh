package main

import (
	"database/sql"
	"fmt"
	"log"
	"net"
	"strconv"
	"time"

	"golang.org/x/crypto/ssh"
	_ "modernc.org/sqlite"
)

var (
	db     *sql.DB
	dbPath = "./sqlite.db"
)

func initDB() {
	_db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		log.Fatal(err)
	}
	db = _db

	createTableSQL := `
    CREATE TABLE IF NOT EXISTS connect_attempts (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        time TEXT NOT NULL,
        ip TEXT NOT NULL,
        port INTEGER NOT NULL
    );
    CREATE TABLE IF NOT EXISTS login_attempts (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        time TEXT NOT NULL,
        ip TEXT NOT NULL,
        port INTEGER NOT NULL,
        client_version TEXT NOT NULL,
        username TEXT NOT NULL,
        password TEXT NOT NULL
    );`
	_, err = db.Exec(createTableSQL)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Table created successfully")
}

func logConnect(conn net.Conn) {
	currentTime := time.Now().Format("2006-01-02 15:04:05")
	remoteAddr := conn.RemoteAddr().String()
	host, port, _ := net.SplitHostPort(remoteAddr)
	portInt, _ := strconv.Atoi(port)

	insertSQL := `INSERT INTO connect_attempts (time, ip, port) VALUES (?, ?, ?)`
	_, err := db.Exec(insertSQL, currentTime, host, portInt)
	if err != nil {
		log.Fatal("Failed to insert connect attempt:", err)
		return
	}
}

func logLogin(conn ssh.ConnMetadata, password []byte) {
	currentTime := time.Now().Format("2006-01-02 15:04:05")
	remoteAddr := conn.RemoteAddr().String()
	host, port, _ := net.SplitHostPort(remoteAddr)
	portInt, _ := strconv.Atoi(port)
	clientVersion := string(conn.ClientVersion())
	username := conn.User()
	passwordStr := string(password)

	insertSQL := `INSERT INTO login_attempts (time, ip, port, client_version, username, password) VALUES (?, ?, ?, ?, ?, ?)`
	_, err := db.Exec(insertSQL, currentTime, host, portInt, clientVersion, username, passwordStr)
	if err != nil {
		log.Fatal("Failed to insert login attempt:", err)
		return
	}
}
