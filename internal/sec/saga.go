package sec

const (
	isCompensating  = true
	notCompensating = false
)

type (
	SagaContext[T any] struct {
		ID           string
		Data         T
		Step         int
		Done         bool
		Compensating bool
	}

	Saga[T any] interface {
		AddStep() SagaStep[T]
		Name() string
		ReplyTopic() string
		getSteps() []SagaStep[T]
	}

	saga[T any] struct {
		name       string
		replyTopic string
		steps      []SagaStep[T]
	}
)

func NewSaga[T any](name string, replyTopic string) Saga[T] {
	return &saga[T]{
		name:       name,
		replyTopic: replyTopic,
	}
}

func (s *saga[T]) AddStep() SagaStep[T] {
	step := &sagaStep[T]{
		actions: map[bool]StepActionFn[T]{
			notCompensating: nil,
			isCompensating:  nil,
		},
		handlers: map[bool]map[string]StepReplyHandlerFn[T]{
			notCompensating: {},
			isCompensating:  {},
		},
	}

	s.steps = append(s.steps, step)

	return step
}

func (s *saga[T]) Name() string {
	return s.name
}

func (s *saga[T]) ReplyTopic() string {
	return s.replyTopic
}

func (s *saga[T]) getSteps() []SagaStep[T] {
	return s.steps
}

func (s *SagaContext[T]) complete() {
	s.Done = true
}

func (s *SagaContext[T]) compensate() {
	s.Compensating = true
}
