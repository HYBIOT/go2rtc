package core

import (
	"time"
)

func (c *Connection) StartBitrateWorker() {
	c.stopBitrateWorker = make(chan struct{}, 1)

	go func() {
		ticker := time.NewTicker(time.Second)
		defer ticker.Stop()

		prevSendBytes := c.Send
		for {
			select {
			case <-ticker.C:
				c.Bitrate = c.Send - prevSendBytes
				prevSendBytes = c.Send
			case <-c.stopBitrateWorker:
				return
			}
		}
	}()
}

func (c *Connection) StopBitrateWorker() {
	c.stopBitrateWorker <- struct{}{}
}
