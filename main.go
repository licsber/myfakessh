package main

import (
	"crypto/ed25519"
	"crypto/rand"
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
		"SSH-2.0-OpenSSH_9.2p1 Debian-2+deb12u4", // My Server Truely
		"SSH-2.0-OpenSSH_6.0p1 Debian-4+deb7u2", // cowrie default
		"SSH-2.0-OpenSSH_9.9", // Manjaro
	}
)

func main() {
	if len(os.Args) > 1 {
		dbPath = os.Args[1] + "/sqlite.db"
	}

	initDB()
	defer db.Close()

	serverConfig := &ssh.ServerConfig{
		MaxAuthTries:     3,
		PasswordCallback: passwordCallback,
		ServerVersion:    serverVersions[0],
	}

	_, privateKey, _ := ed25519.GenerateKey(rand.Reader)
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
