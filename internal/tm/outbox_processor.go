package tm

import (
	"context"
	"time"

	"github.com/htquangg/microservices-poc/internal/am"
)

const (
	messageLimit    = 50
	pollingInterval = 500 * time.Millisecond
)

type OutboxProcessor interface {
	Start(ctx context.Context) error
}

type outboxProcessor struct {
	publisher am.MessagePublisher
	store     OutboxStore
}

func NewOutboxProcessor(publisher am.MessagePublisher, store OutboxStore) OutboxProcessor {
	return outboxProcessor{
		publisher: publisher,
		store:     store,
	}
}

func (p outboxProcessor) Start(ctx context.Context) error {
	errC := make(chan error)

	go func() {
		errC <- p.processMessages(ctx)
	}()

	select {
	case <-ctx.Done():
		return nil
	case err := <-errC:
		return err
	}
}

func (p outboxProcessor) processMessages(ctx context.Context) error {
	timer := time.NewTimer(0)

	// TOIMPROVE: will persistence currentOffset for recovering
	currentOffset := "0"
	for {
		msgs, err := p.store.FindUnpublished(ctx, currentOffset, messageLimit)
		if err != nil {
			return err
		}

		if len(msgs) > 0 {
			topicMessages := make(map[string][]am.Message, len(msgs))

			for _, msg := range msgs {
				// group the messages by subject
				topicMessages[msg.Subject()] = append(topicMessages[msg.Subject()], msg)
			}

			for topic, messages := range topicMessages {
				if err := p.publisher.Publish(ctx, topic, messages...); err != nil {
					return err
				}
			}

			err = p.store.MarkPublished(ctx, msgs[0].ID(), msgs[len(msgs)-1].ID())
			if err != nil {
				return err
			}

			currentOffset = msgs[len(msgs)-1].ID()

			// poll again immediately
			continue
		}

		if !timer.Stop() {
			select {
			case <-timer.C:
			default:
			}
		}

		// wait a short time before polling again
		timer.Reset(pollingInterval)

		select {
		case <-ctx.Done():
			return nil
		case <-timer.C:
		}
	}
}
