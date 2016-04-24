package main

import (
	"encoding/binary"
	"fmt"
	"os"
	"os/exec"
	"time"
)

func main() {
	for {
		err := connectAndListen("/dev/input/powermate")
		if err != nil {
			fmt.Println(err)
		}
		time.Sleep(time.Second)
	}
}

func connectAndListen(fname string) error {
	device, err := os.OpenFile(fname, os.O_RDONLY, 0666)
	if err != nil {
		return err
	}
	defer device.Close()

	s := &state{pressed: false}
	buf := make([]byte, 48)
	for {
		n, err := device.Read(buf)
		if err != nil {
			return err
		}
		event := buf[:n]

		w := event[16:20]
		typ, _ := binary.Varint(w)

		w = event[20:24]
		val, _ := binary.Varint(w)

		s.handle(int32(typ), int32(val))
	}
}

const (
	typeRot   = 1
	typePress = -1
	valLeft   = 0
	valRight  = -1
	valUp     = 0
	valDown   = -1
)

type state struct {
	pressed   bool
	pressedAt time.Time
}

func (s *state) handle(typ, val int32) {
	switch typ {
	case typeRot:
		switch val {
		case valRight:
			if s.pressed {
				execute("playerctl", "next")
			} else {
				execute("amixer", "-D", "pulse", "sset", "Master", "1%+")
			}
		case valLeft:
			if s.pressed {
				execute("playerctl", "previous")
			} else {
				execute("amixer", "-D", "pulse", "sset", "Master", "1%-")
			}
		}
	case typePress:
		switch val {
		case valDown:
			s.pressed = true
			s.pressedAt = time.Now()
		case valUp:
			s.pressed = false
			if time.Since(s.pressedAt) < time.Second {
				execute("playerctl", "play-pause")
			}
		}
	}
}

func execute(cmd string, args ...string) {
	out, err := exec.Command(cmd, args...).CombinedOutput()
	if err != nil {
		fmt.Println(string(out))
		fmt.Println(err)
	}
}
