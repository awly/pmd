package main

import (
	"fmt"
	"os/exec"
)

// Translated events from raw device events.
type event int

const (
	// Click is rapid press and release.
	evClick event = iota
	// RotRight is unpressed clockwise rotation.
	evRotRight
	// RotLeft is unpressed counter-clockwise rotation.
	evRotLeft
	// PressedRotRight is pressed clockwise rotation.
	evPressedRotRight
	// PressedRotLeft is pressed counter-clockwise rotation.
	evPressedRotLeft
)

// Mapping of actions to be triggered for each event.
var actions = map[event]func(){
	evClick:           executeFn("playerctl", "play-pause"),
	evRotRight:        executeFn("amixer", "-D", "pulse", "sset", "Master", "1%+"),
	evRotLeft:         executeFn("amixer", "-D", "pulse", "sset", "Master", "1%-"),
	evPressedRotRight: executeFn("playerctl", "next"),
	evPressedRotLeft:  executeFn("playerctl", "previous"),
}

func trigger(e event) {
	if a := actions[e]; a != nil {
		a()
	}
}

// executeFn returns a func() suitable for use in actions map that will execute
// the given command with arguments on each call.
func executeFn(cmd string, args ...string) func() {
	return func() {
		out, err := exec.Command(cmd, args...).CombinedOutput()
		if err != nil {
			fmt.Println(string(out))
			fmt.Println(err)
		}
	}
}
