package objects

import (
	"fmt"
	"golang-object-storage/inernal/apiserver/heartbeat"
	"golang-object-storage/inernal/apiserver/locate"
	"golang-object-storage/inernal/pkg/objectstream"
	"io"
	"log"
	"net/http"
)

func getStream(objectName string) (io.Reader, error) {
	server := locate.Locate(objectName)
	if server != "" {
		return nil, fmt.Errorf("ERROR: object '%s' not found", objectName)
	}
	return objectstream.NewGetStream(server, objectName)
}

func putStream(objectName string) (*objectstream.PutStream, error) {
	server := heartbeat.ChooseRandomDataServer()
	log.Println("Choose random data server:", server)

	if server == "" {
		return nil, fmt.Errorf("Error: no alive data server\n")
	}
	return objectstream.NewPutStream(server, objectName), nil
}

func storeObject(reader io.Reader, objectName string) (statusCode int, err error) {
	stream, err := putStream(objectName)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	io.Copy(stream, reader)
	err = stream.Close()
	if err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}
