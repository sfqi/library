package service

import (
	"fmt"
	"testing"
)

func TestGenerator_Do(t *testing.T) {
	for i := 0; i < 10; i++ {
		t.Run(fmt.Sprintf("test %d", i), func(t *testing.T) {
			g := &Generator{}
			got, err := g.Do()
			if err != nil {
				t.Errorf("Do() error = %v", err)
				return
			}
			if len(got) != 36 {
				t.Errorf("uuid is not allowed lenght")
			}
		})
	}
}
