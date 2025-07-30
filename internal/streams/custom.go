package streams

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/AlexxIT/go2rtc/pkg/rtsp"
)

// only support RTSP sources
func apiStreamsSpeed(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	streamName := query.Get("name")

	if streamName == "" {
		http.Error(w, "name required", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case "PUT":
		stream := Get(streamName)

		if stream == nil {
			http.Error(w, "", http.StatusNotFound)
			return
		}

		speedStr := query.Get("speed")
		if speedStr == "" {
			http.Error(w, "speed required", http.StatusBadRequest)
			return
		}

		_, err := strconv.ParseFloat(speedStr, 64)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		for _, producer := range stream.producers {
			if !isISAPIPlayback(producer.url) {
				continue
			}

			producer.speed = speedStr

			if conn, ok := producer.conn.(*rtsp.Conn); ok {
				conn.Connection.Speed = speedStr
				err := conn.Play()
				if err != nil {
					log.Error().Msgf("[stream] conn.Play(): %+v\n", err)
					http.Error(w, err.Error(), http.StatusInternalServerError)
				}
			}
		}

		return
	}

	http.Error(w, "", http.StatusNotFound)
}

func apiStreamsRemoveConsumers(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "", http.StatusNotFound)
		return
	}

	query := r.URL.Query()

	for _, streamName := range query["name"] {
		if streamName == "" {
			continue
		}

		stream := Get(streamName)

		if stream == nil {
			continue
		}

		for _, consumer := range stream.consumers {
			stream.RemoveConsumer(consumer)
		}
	}

	http.Error(w, "", http.StatusOK)
}

func setupStreamSpeed(stream *Stream, speedStr string, speedHeaders []string) {
	speedFloat, err := strconv.ParseFloat(speedStr, 64)
	if err != nil {
		return
	}

	speedStr3 := fmt.Sprintf("%.3f", speedFloat)

	for _, producer := range stream.producers {
		if speedStr3 == "1.000" {
			producer.speed = ""
			producer.speedHeaders = nil
		} else {
			producer.speed = speedStr
			producer.speedHeaders = speedHeaders
		}
	}
}

func isISAPIPlayback(url string) bool {
	return strings.Contains(url, "/Streaming/tracks")
}
