// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2018 Datadog, Inc.

package listener

import (
	"fmt"
	"log"
	"net"

	"github.com/DataDog/datadog-agent/pkg/logs/config"
	"github.com/DataDog/datadog-agent/pkg/logs/pipeline"
)

// A TCPListener listens and accepts TCP connections and delegates the work to connHandler
type TCPListener struct {
	listener    net.Listener
	connHandler *ConnectionHandler
}

// NewTCPListener returns an initialized TCPListener
func NewTCPListener(pp pipeline.Provider, source *config.IntegrationConfigLogSource) (*TCPListener, error) {
	log.Println("Starting TCP forwarder on port", source.Port)
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", source.Port))
	if err != nil {
		return nil, err
	}
	connHandler := &ConnectionHandler{
		pp:     pp,
		source: source,
	}
	return &TCPListener{
		listener:    listener,
		connHandler: connHandler,
	}, nil
}

// Start listens to TCP connections on another routine
func (tcpListener *TCPListener) Start() {
	go tcpListener.run()
}

// run accepts new TCP connections and lets connHandler handle them
func (tcpListener *TCPListener) run() {
	for {
		conn, err := tcpListener.listener.Accept()
		if err != nil {
			log.Println("Can't listen:", err)
			return
		}
		go tcpListener.connHandler.handleConnection(conn)
	}
}
