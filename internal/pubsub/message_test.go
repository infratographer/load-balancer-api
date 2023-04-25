package pubsub

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.infratographer.com/x/pubsubx"
)

func Test_validatePubsubMessage(t *testing.T) {
	tests := []struct {
		name    string
		msg     *pubsubx.Message
		wantErr bool
	}{
		{
			name: "valid message",
			msg: &pubsubx.Message{
				EventType:  "test",
				Source:     "test",
				SubjectURN: "foo",
				ActorURN:   "bar",
			},
			wantErr: false,
		},
		{
			name: "missing source",
			msg: &pubsubx.Message{
				EventType:  "test",
				SubjectURN: "foo",
				ActorURN:   "bar",
			},
			wantErr: true,
		},
		{
			name: "missing subject urn",
			msg: &pubsubx.Message{
				EventType: "test",
				Source:    "test",
				ActorURN:  "bar",
			},
			wantErr: true,
		},
		{
			name: "missing actor urn",
			msg: &pubsubx.Message{
				EventType:  "test",
				Source:     "test",
				SubjectURN: "foo",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validatePubsubMessage(tt.msg)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
		})
	}
}
