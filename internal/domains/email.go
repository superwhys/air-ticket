package domains

import "context"

type EmailSender interface {
	SentMsg(ctx context.Context, targets []string, msg []byte) error
}
