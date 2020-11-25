package alerts

import (
	"log"
	"os"
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
	"github.com/sprucewillis/nvidia-finder/internal/webscraper/inventory"
)

// SetupAudioAlerts(chan inventory.Item) play an audio alert when something is in stock
func SetupAudioAlerts(c chan inventory.Item) error {
	// read audio file and set up streamer
	fileName := "src/github.com/sprucewillis/nvidia-finder/internal/util/alerts/objection.mp3"
	f, err := os.Open(fileName)
	if err != nil {
		log.Println("ERROR: unable to open audio alert file", err)
		return err
	}
	streamer, format, err := mp3.Decode(f)
	if err != nil {
		log.Println("ERROR: unable to decode audio alert file", err)
		return err
	}
	speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))
	buffer := beep.NewBuffer(format)
	buffer.Append(streamer)
	streamer.Close()
	for range c {
		ping := buffer.Streamer(0, buffer.Len()*2/3) // short enough to play in its entirety
		speaker.Play(ping)
	}
	return nil
}
