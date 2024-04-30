package mp4

import "github.com/AlexxIT/go2rtc/pkg/custom"

func (c *Consumer) startBitrateWorker() {
	c.stopBitrateWorker = make(chan struct{})
	go custom.StartBitrateWorker(&c.SuperConsumer.Send, &c.SuperConsumer.Bitrate, c.stopBitrateWorker)
}
