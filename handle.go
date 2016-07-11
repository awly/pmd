package main

import "time"

const (
	// Raw event types received from the device.
	typeRot   = 1
	typePress = -1
	valLeft   = 0
	valRight  = -1
	valUp     = 0
	valDown   = -1

	// Threshold to trigger pressed rotation actions.
	// Pressed rotation actions have more impact and need to be less sensitive.
	pressedRotTicks = 2
)

// Current state of the device.
type state struct {
	pressed   bool
	pressedAt time.Time
	rot       int
}

// Handle incoming parsed event.
func (s *state) handle(typ, val int32) {
	switch typ {
	case typeRot:
		switch val {
		case valRight:
			s.rot++
			if s.pressed {
				if s.rot%pressedRotTicks == 0 {
					trigger(evPressedRotRight)
				}
			} else {
				trigger(evRotRight)
			}
		case valLeft:
			s.rot--
			if s.pressed {
				if s.rot%pressedRotTicks == 0 {
					trigger(evPressedRotLeft)
				}
			} else {
				trigger(evRotLeft)
			}
		}
	case typePress:
		s.rot = 0
		switch val {
		case valDown:
			s.pressed = true
			s.pressedAt = time.Now()
		case valUp:
			s.pressed = false
			if time.Since(s.pressedAt) < time.Second {
				trigger(evClick)
			}
		}
	}
}
