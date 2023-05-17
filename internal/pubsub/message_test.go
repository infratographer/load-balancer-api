package pubsub

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.infratographer.com/x/gidx"
	"go.infratographer.com/x/pubsubx"
)

func Test_validatePubsubMessage(t *testing.T) {
	tests := []struct {
		name    string
		msg     *pubsubx.ChangeMessage
		wantErr bool
	}{
		{
			name: "valid message",
			msg: &pubsubx.ChangeMessage{
				EventType: "test",
				Source:    "test",
				SubjectID: gidx.MustNewID("gidxtst"),
				ActorID:   gidx.MustNewID("actorid"),
			},
			wantErr: false,
		},
		{
			name: "missing source",
			msg: &pubsubx.ChangeMessage{
				EventType: "test",
				SubjectID: gidx.MustNewID("gidxtst"),
				ActorID:   gidx.MustNewID("actorid"),
			},
			wantErr: true,
		},
		{
			name: "missing subject urn",
			msg: &pubsubx.ChangeMessage{
				EventType: "test",
				Source:    "test",
				ActorID:   gidx.MustNewID("actorid"),
			},
			wantErr: true,
		},
		{
			name: "missing actor urn",
			msg: &pubsubx.ChangeMessage{
				EventType: "test",
				Source:    "test",
				SubjectID: gidx.MustNewID("gidxtst"),
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
