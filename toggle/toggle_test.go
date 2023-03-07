package toggle_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tiket/TIX-AFFILIATE-COMMON-GO/toggle"
)

const (
	spec = 2
)

type SampleExecutor struct{}

func NewSampleExecutor() toggle.ToggleExecutor[int64, int64, interface{}] {
	return &SampleExecutor{}
}

func (e *SampleExecutor) IsToggleOn(w interface{}) bool {
	return w != nil
}

func (e *SampleExecutor) OnToggleOn(t int64) int64 {
	return 1
}

func (e *SampleExecutor) OnToggleOff(t int64) int64 {
	return 0
}

func TestToggle_Run(t *testing.T) {
	sampleExecutor := SampleExecutor{}
	sampleHelper := toggle.ToggleHelper[int64, int64, interface{}]{
		Executor: &sampleExecutor,
	}
	testcases := []struct {
		testName string
		t        int64
		w        interface{}
		expected int64
	}{
		{
			testName: "should return 1 when IsToggleOn true",
			t:        spec,
			w:        true,
			expected: 1,
		},
		{
			testName: "should return 0 when IsToggleOn false",
			t:        spec,
			w:        nil,
			expected: 0,
		},
	}
	for _, tt := range testcases {
		t.Run(tt.testName, func(t *testing.T) {
			result := sampleHelper.Run(tt.t, tt.w)
			require.Equal(t, tt.expected, result)
		})
	}
}
