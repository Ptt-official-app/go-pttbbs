package ptttype

import "testing"

func TestMQType_String(t *testing.T) {
	tests := []struct {
		name     string
		m        MQType
		expected string
	}{
		// TODO: Add test cases.
		{
			m:        MQ_TEXT,
			expected: "text",
		},
		{
			m:        MQ_UUENCODE,
			expected: "uuencode",
		},
		{
			m:        MQ_JUSTIFY,
			expected: "mq-justify",
		},
		{
			m:        MQType(100),
			expected: "unknown",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.m.String(); got != tt.expected {
				t.Errorf("MQType.String() = %v, want %v", got, tt.expected)
			}
		})
	}
}
