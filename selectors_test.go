package navigator

import (
	"fmt"
	"testing"
)

func Test_selectors_At(t *testing.T) {
	t.Run("index at", func(t *testing.T) {
		var ss selectors
		for i := 0; i < 3; i++ {
			ss = append(ss, selector{
				Type:  cssType,
				Value: fmt.Sprintf("%d", i),
			})
		}
		for _, want := range []int{-10, -5, -1, 0, 1, 5, 10} {
			t.Run(fmt.Sprintf("At(%d)", want), func(t *testing.T) {
				got := ss.At(want)
				tail := len(got) - 1
				if !got[tail].Indexed {
					t.Errorf("want 'Indexed' field of the tail selector true, but got false, %+v", got)
				}
				if got[tail].Single {
					t.Errorf("want 'Single' field of the tail selector true, but got false, %+v", got)
				}
				if got[tail].Index != want {
					t.Errorf("want 'Index' field of the tail selector = %d, but got %d, %+v", want, got[tail].Index, got)
				}
			})
		}
	})
	t.Run("selectors is empty", func(t *testing.T) {
		var ss selectors
		got := ss.At(3)
		if got != nil {
			t.Errorf("want nil, got %+v", got)
		}
	})
}
