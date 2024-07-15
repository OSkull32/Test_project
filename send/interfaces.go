package send

import "context"

type interfaceSend interface {
	PublishMessage(ctx context.Context, message string) error
}
