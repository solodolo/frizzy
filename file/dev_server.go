package file

import (
	"fmt"
	"io"
	"io/fs"
	"log"
	"net/http"
	"os"
	"time"
)

type DevServer struct {
	ServerRoot string
	Port       int
}

func (receiver *DevServer) ListenAndServe() error {
	rootFS := os.DirFS(receiver.ServerRoot)
	address := fmt.Sprintf(":%d", receiver.Port)
	server := http.Server{
		Addr:        address,
		ReadTimeout: time.Second * 5,
		Handler:     &devServerHandler{rootFS},
	}

	log.Printf("serving files from %s\n", receiver.ServerRoot)
	log.Printf("listening on http://localhost:%d\n", receiver.Port)

	return server.ListenAndServe()
}

type devServerHandler struct {
	rootFS fs.FS
}

func (receiver *devServerHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	requestPath := request.URL.EscapedPath()
	if f, err := receiver.rootFS.Open(requestPath[1:]); err != nil {
		writer.WriteHeader(404)
		io.WriteString(writer, fmt.Sprintf("File %s not found: %s", requestPath, err))
	} else {
		if bytes, err := io.ReadAll(f); err != nil {
			writer.WriteHeader(500)
			io.WriteString(writer, fmt.Sprintf("File %s could not be read", requestPath))
		} else {
			writer.Write(bytes)
		}
	}
}
