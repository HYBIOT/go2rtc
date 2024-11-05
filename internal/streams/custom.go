package streams

import (
	"net/http"
	"strconv"

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
		streamsMu.RLock()
		stream := Get(streamName)
		defer streamsMu.RUnlock()

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

		streamsMu.RLock()
		stream := Get(streamName)
		streamsMu.RUnlock()

		if stream == nil {
			continue
		}

		for _, consumer := range stream.consumers {
			stream.RemoveConsumer(consumer)
		}
	}

	http.Error(w, "", http.StatusOK)
}
