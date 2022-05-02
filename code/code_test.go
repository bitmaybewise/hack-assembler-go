package code

import (
	"testing"
)

func TestJump(t *testing.T) {
	for k, v := range jumps {
		t.Run(k, func(t *testing.T) {
			if got := Jump(k); got != v {
				t.Errorf("Jump() = %v, want %v", got, v)
			}
		})
	}
}

func TestDest(t *testing.T) {
	for k, v := range destinations {
		t.Run(k, func(t *testing.T) {
			if got := Dest(k); got != v {
				t.Errorf("Dest() = %v, want %v", got, v)
			}
		})
	}
}

func TestComp(t *testing.T) {
	for k, v := range computations {
		t.Run(k, func(t *testing.T) {
			if got := Comp(k); got != v {
				t.Errorf("Comp() = %v, want %v", got, v)
			}
		})
	}
}
