package nanokontrol2

import (
	"log"
  "time"
)

import "github.com/rakyll/portmidi"


 /*
   * MIDI Mappings
   *
   * Sliders  0 - 7
   * Knobs  16 - 23
   * S Button 32 - 39
   * Play    41
   * Stop   42
   * FF   43
   * RV   44
   * REC   45
   * Cycle  46
   * M Button 48 - 55
   * Track Back  58
   * Track FWD  59
   * Marker Set 60
   * Marker Back  61
   * Marker FWD 62
   * R Button 64 - 71
   */
  
const REFRESH_RATE = 50

type Nanokontrol struct{
  state []float32
  in *portmidi.Stream
}


func (n *Nanokontrol) Get(index uint8) float32{
	return n.state[index]
}

func (n *Nanokontrol) Poll(){
	for{
        if n.in != nil{
            result, err := n.in.Poll()
            if err != nil {
                log.Fatal(err)
            }

            if result {
                msg, err := n.in.Read(1024)
                if err != nil {
                    log.Fatal(err)
                }
                for b := range msg {
                    event := msg[b]
                    n.state[event.Data1] = float32(event.Data2) / 127
                }
            }
        }
        time.Sleep(REFRESH_RATE)
    }
}

func midiStream() *portmidi.Stream{
    var stream *portmidi.Stream
    portmidi.Initialize()
    if portmidi.CountDevices() == 0{
        log.Printf("No MIDI controller found")
    }else{
        log.Printf("midi found")
        midistream, err := portmidi.NewInputStream(portmidi.DefaultInputDeviceID(), 1024)
        if err != nil {
            log.Fatal(err)
        }
        stream = midistream
    }
    return stream
}

func Initialize() *Nanokontrol{
  n := Nanokontrol{make([]float32, 71), midiStream()}
  go n.Poll()
  return &n
}