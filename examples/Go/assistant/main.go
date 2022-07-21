// This example streams the microphone thru Snowboy to listen for the hotword,
// by using the PortAudio interface.
//
// HOW TO USE:
// 	go run examples/Go/listen/main.go [path to snowboy resource file] [path to snowboy hotword file]
//
package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"os"
	"io/ioutil"
	"io"
	"sync"
	"time"

	"github.com/brentnd/go-snowboy"
	"github.com/byebyebruce/aggrsdk/pkg/hear"
	"github.com/gordonklaus/portaudio"
)

// Sound represents a sound stream implementing the io.Reader interface
// that provides the microphone data.
type Sound struct {
	stream *portaudio.Stream
	data   []int16
}

// Init initializes the Sound's PortAudio stream.
func (s *Sound) Init() error {
	inputChannels := 1
	outputChannels := 0
	sampleRate := 16000
	s.data = make([]int16, 1024)

	// initialize the audio recording interface
	err := portaudio.Initialize()
	if err != nil {
		return fmt.Errorf("Error initialize audio interface: %s", err)

	}

	// open the sound input stream for the microphone
	stream, err := portaudio.OpenDefaultStream(inputChannels, outputChannels, float64(sampleRate), len(s.data), s.data)
	if err != nil {
		return fmt.Errorf("Error open default audio stream: %s", err)
	}

	err = stream.Start()
	if err != nil {
		return fmt.Errorf("Error on stream start: %s", err)
	}

	s.stream = stream
	return nil
}

// Close closes down the Sound's PortAudio connection.
func (s *Sound) Close() {
	s.stream.Close()
	portaudio.Terminate()
}

// Read is the Sound's implementation of the io.Reader interface.
func (s *Sound) Read(p []byte) (int, error) {
	if err := s.stream.Read(); err != nil {
		return 0, err
	}

	buf := &bytes.Buffer{}
	for _, v := range s.data {
		binary.Write(buf, binary.LittleEndian, v)
	}

	copy(p, buf.Bytes())
	return len(p), nil
}

func DetectHotWord(stopChan <-chan struct{}) (bool, error) {
	said := false
	// open the mic
	mic := &Sound{}
	if err := mic.Init(); err != nil {
		return false, err
	}
	once := sync.Once{}
	stop := func() {
		once.Do(func() {
			mic.Close()
		})
	}
	defer stop()

	go func() {
		<-stopChan
		stop()
	}()

	// open the snowboy detector
	d := snowboy.NewDetector(os.Args[1])
	defer d.Close()

	// set the handlers
	d.HandleFunc(snowboy.NewHotword(os.Args[2], 0.5), func(string) {
		said = true
		stop()
		fmt.Println("You said the hotword!")
	})

	d.HandleSilenceFunc(1*time.Second, func(string) {
		fmt.Println("Silence detected.")
	})

	// display the detector's expected audio format
	sr, nc, bd := d.AudioFormat()
	fmt.Printf("sample rate=%d, num channels=%d, bit depth=%d\n", sr, nc, bd)

	// start detecting using the microphone
	 err := d.ReadAndDetect(mic); 
	 if err != nil && err==io.EOF {
		err = nil
	}
	return said, nil
}

func main() {
	stopChan := make(chan struct{})
	for i:=0; i<10; i++ {
		fmt.Println("time:",i)
		
		ok, err := DetectHotWord(stopChan)
		fmt.Println("result", ok,err)
		if err!=nil || !ok {
			continue
		}
		b,err :=hear.Hear2Wav(time.Second*10,time.Second*4)
		if err!=nil {
			continue
		}
		fmt.Println("listen ok")
	 _ = ioutil.WriteFile(fmt.Sprintf("result_%d.wav",i), b, os.ModePerm)
	 	time.Sleep(time.Second)
	}
}
