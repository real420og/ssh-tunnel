package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"

	"golang.org/x/crypto/ssh"
)

type Endpoint struct {
	Host string
	Port int
}

func (endpoint *Endpoint) String() string {
	return fmt.Sprintf("%s:%d", endpoint.Host, endpoint.Port)
}

type SSHtunnel struct {
	Local  *Endpoint
	Server *Endpoint
	Remote *Endpoint

	Config *ssh.ClientConfig
}

func (tunnel *SSHtunnel) Start() error {
	listener, err := net.Listen("tcp", tunnel.Local.String())
	if err != nil {
		return err
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			return err
		}
		go tunnel.forward(conn)
	}
}

func (tunnel *SSHtunnel) forward(localConn net.Conn) {
	serverConn, err := ssh.Dial("tcp", tunnel.Server.String(), tunnel.Config)
	if err != nil {
		log.Printf("Server dial error: %s\n", err)
		return
	}

	remoteConn, err := serverConn.Dial("tcp", tunnel.Remote.String())
	if err != nil {
		log.Printf("Remote dial error: %s\n", err)
		return
	}

	copyConn := func(writer, reader net.Conn) {
		defer writer.Close()
		defer reader.Close()

		_, err := io.Copy(writer, reader)
		if err != nil {
			log.Printf("io.Copy error: %s", err)
		}
	}

	go copyConn(localConn, remoteConn)
	go copyConn(remoteConn, localConn)
}

func SSHAgent() ssh.AuthMethod {
	signer, err := ssh.ParsePrivateKey([]byte(os.Getenv("SSH_AUTH_KEY")))
	if err != nil {
		log.Fatalf("parse key failed: %v", err)
	}

	return ssh.PublicKeys(signer)
}

func main() {
	localPort := flag.Int("port", 8080, "get local port")
	remoteHost := flag.String("host", "", "get remote host")
	serverHost := flag.String("server", "", "get server host")

	flag.Parse()

	localEndpoint := &Endpoint{
		Host: "0.0.0.0",
		Port: *localPort,
	}

	remoteEndpoint := &Endpoint{
		Host: *remoteHost,
		Port: 443,
	}

	serverEndpoint := &Endpoint{
		Host: *serverHost,
		Port: 22,
	}

	sshConfig := &ssh.ClientConfig{
		User: os.Getenv("SSH_AUTH_LOGIN"),
		Auth: []ssh.AuthMethod{
			SSHAgent(),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	tunnel := &SSHtunnel{
		Config: sshConfig,
		Local:  localEndpoint,
		Server: serverEndpoint,
		Remote: remoteEndpoint,
	}

	tunnel.Start()
}
