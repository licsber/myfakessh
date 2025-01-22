package main

import (
	"crypto/rand"
	"crypto/rsa"
	"errors"
	"log"
	"net"
	"os"
	"time"

	"golang.org/x/crypto/ssh"
)

var (
	errBadPassword = errors.New("permission denied")
	serverVersions = []string{
		"SSH-2.0-OpenSSH_9.9",
	}
)

func main() {
	if len(os.Args) > 1 {
		dbPath = os.Args[1] + "/sqlite.db"
	}

	initDB()
	defer db.Close()

	serverConfig := &ssh.ServerConfig{
		MaxAuthTries:     6,
		PasswordCallback: passwordCallback,
		ServerVersion:    serverVersions[0],
	}

	privateKey, _ := rsa.GenerateKey(rand.Reader, 2048)
	signer, _ := ssh.NewSignerFromSigner(privateKey)
	serverConfig.AddHostKey(signer)

	listener, err := net.Listen("tcp", ":22")
	if err != nil {
		log.Fatal("Failed to listen:", err)
		return
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal("Failed to accept:", err)
			break
		}
		go handleConn(conn, serverConfig)
	}
}

func passwordCallback(conn ssh.ConnMetadata, password []byte) (*ssh.Permissions, error) {
	log.Println(conn.RemoteAddr(), string(conn.ClientVersion()), conn.User(), string(password))
	logLogin(conn, password)
	time.Sleep(100 * time.Millisecond)
	return nil, errBadPassword
}

func handleConn(conn net.Conn, serverConfig *ssh.ServerConfig) {
	defer conn.Close()
	log.Println(conn.RemoteAddr())
	logConnect(conn)
	ssh.NewServerConn(conn, serverConfig)
}
