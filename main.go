package main

import (
	"encoding/binary"
	"fmt"
	"os"
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
