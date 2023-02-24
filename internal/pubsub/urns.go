package pubsub

import "fmt"

func newURN(kind, id string) string {
	return fmt.Sprintf("urn:infratographer:%s:%s", kind, id)
}
