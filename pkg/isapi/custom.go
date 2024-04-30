package isapi

import "github.com/AlexxIT/go2rtc/pkg/custom"

func (c *Client) startBitrateWorker() {
	c.stopBitrateWorker = make(chan struct{})
	go custom.StartBitrateWorker(&c.send, &c.bitrate, c.stopBitrateWorker)
}
