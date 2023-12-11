package kafka

import (
	"context"
	"sync"

	"github.com/htquangg/microservices-poc/pkg/logger"

	"github.com/segmentio/kafka-go"
)

type consumeFn func(msg *kafka.Message) error

type subscription struct {
	r           *kafka.Reader
	ctx         context.Context
	log         logger.Logger
	messageCh   chan *kafka.Message
	quit        chan struct{}
	cancelFn    context.CancelFunc
	wg          sync.WaitGroup
	concurrency int
	once        sync.Once
	consumeFn   func(msg *kafka.Message) error
}

func (s *subscription) Unsubscribe() error {
	var err error
	s.once.Do(func() {
		s.cancelFn()
		close(s.quit)
		close(s.messageCh)
		s.wg.Wait()
		err = s.r.Close()
	})
	return err
}

func (s *subscription) consume() {
	s.wg.Add(1)

	go func() {
		defer s.wg.Done()

		for {
			select {
			case <-s.quit:
				return
			default:
				msg, err := s.r.ReadMessage(s.ctx)
				if err != nil {
					if s.ctx.Err() != nil {
						continue
					}
					s.log.Errorf("message could not read %s", err.Error())
					continue
				}

				s.messageCh <- &msg
			}
		}
	}()

	for i := 0; i < s.concurrency; i++ {
		s.wg.Add(1)
		go func() {
			defer s.wg.Done()
			for msg := range s.messageCh {
				s.process(msg)
			}
		}()
	}
}

func (s *subscription) process(msg *kafka.Message) {
	consumeErr := s.consumeFn(msg)

	if consumeErr != nil {
		s.log.Warnf("consume function err %s", consumeErr.Error())
		// try to process same message again
		if consumeErr = s.consumeFn(msg); consumeErr != nil {
			s.log.Warnf("consume function again err %s", consumeErr.Error())
		}
	}
}
