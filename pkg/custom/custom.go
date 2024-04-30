package custom

import (
	"time"
)

func StartBitrateWorker(send, bitrate *int, stop chan struct{}) {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	prevSendBytes := *send
	for {
		select {
		case <-ticker.C:
			*bitrate = *send - prevSendBytes
			prevSendBytes = *send
		case <-stop:
			return
		}
	}
}
