package am

import (
	"context"
	"strings"
	"time"

	proto_am "github.com/htquangg/microservices-poc/internal/am/proto"
	"github.com/htquangg/microservices-poc/internal/ddd"
	"github.com/htquangg/microservices-poc/internal/registry"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const (
	CommandHandlerPrefix       = "COMMAND_"
	CommandNameHandler         = CommandHandlerPrefix + "NAME"
	CommandReplyChannelHandler = "CommandHandlerPrefix" + "REPLY_CHANNEL"
)

type (
	CommandMessage interface {
		MessageBase
		ddd.Command
	}

	IncomingCommandMessage interface {
		IncomingMessageBase
		ddd.Command
	}

	CommandPublisher interface {
		Publish(ctx context.Context, topicName string, cmd ddd.Command) error
	}

	commandPublisher struct {
		reg       registry.Registry
		publisher MessagePublisher
	}

	commandMessage struct {
		id         string
		name       string
		payload    ddd.CommandPayload
		occurredAt time.Time
		msg        IncomingMessageBase
	}

	commandMessageHandler struct {
		reg       registry.Registry
		publisher ReplyPublisher
		handler   ddd.CommandHander[ddd.Command]
	}
)

var (
	_ CommandPublisher = (*commandPublisher)(nil)
	_ CommandMessage   = (*commandMessage)(nil)
)

func NewCommandPublisher(reg registry.Registry, publisher MessagePublisher) CommandPublisher {
	return &commandPublisher{
		reg:       reg,
		publisher: publisher,
	}
}

func (c *commandPublisher) Publish(ctx context.Context, topicName string, cmd ddd.Command) error {
	payload, err := c.reg.Serialize(cmd.CommandName(), cmd.Payload())
	if err != nil {
		return err
	}

	data, err := proto.Marshal(&proto_am.CommandMessageData{
		Payload:    payload,
		OccurredAt: timestamppb.New(cmd.OccurredAt()),
	})
	if err != nil {
		return err
	}

	return c.publisher.Publish(ctx, topicName, message{
		id:       cmd.ID(),
		subject:  topicName,
		name:     cmd.CommandName(),
		data:     data,
		metadata: cmd.Metadata(),
		sentAt:   time.Now(),
	})
}

func (c commandMessage) ID() string {
	return c.id
}

func (c commandMessage) CommandName() string {
	return c.name
}

func (c commandMessage) Payload() ddd.CommandPayload {
	return c.payload
}

func (c commandMessage) Metadata() ddd.Metadata {
	return c.msg.Metadata()
}

func (c commandMessage) OccurredAt() time.Time {
	return c.occurredAt
}

func (c commandMessage) Subject() string {
	return c.msg.Subject()
}

func (c commandMessage) MessageName() string {
	return c.msg.MessageName()
}

func (c commandMessage) SentAt() time.Time {
	return c.msg.SentAt()
}

func (c commandMessage) ReceivedAt() time.Time {
	return c.msg.ReceivedAt()
}

func (c commandMessage) Ack() error {
	return c.msg.Ack()
}

func (c commandMessage) NAck() error {
	return c.msg.NAck()
}

func (c commandMessage) Extend() error {
	return c.msg.Extend()
}

func (c commandMessage) Kill() error {
	return c.msg.Kill()
}

func NewCommandHandler(
	reg registry.Registry,
	publisher ReplyPublisher,
	handler ddd.CommandHander[ddd.Command],
	mws ...MessageHandlerMiddleware,
) MessageHandler {
	return MessageHandlerWithMiddleware(commandMessageHandler{
		reg:       reg,
		publisher: publisher,
		handler:   handler,
	}, mws...)
}

func (h commandMessageHandler) HandleMessage(ctx context.Context, msg IncomingMessage) error {
	var commandData proto_am.CommandMessageData

	err := proto.Unmarshal(msg.Data(), &commandData)
	if err != nil {
		return err
	}

	commandName := msg.MessageName()

	payload, err := h.reg.Deserialize(commandName, commandData.GetPayload())
	if err != nil {
		return err
	}

	commandMsg := &commandMessage{
		id:         msg.ID(),
		name:       commandName,
		payload:    payload,
		occurredAt: commandData.GetOccurredAt().AsTime(),
		msg:        msg,
	}

	destination := commandMsg.Metadata().Get(CommandReplyChannelHandler).(string)

	reply, err := h.handler.HandleCommand(ctx, commandMsg)
	if err != nil {
		return h.publishReply(ctx, destination, h.failure(reply, commandMsg))
	}

	return h.publishReply(ctx, destination, h.success(reply, commandMsg))
}

func (h commandMessageHandler) publishReply(ctx context.Context, destination string, reply ddd.Reply) error {
	return h.publisher.Publish(ctx, destination, reply)
}

func (h commandMessageHandler) success(reply ddd.Reply, cmd ddd.Command) ddd.Reply {
	if reply == nil {
		reply = ddd.NewReply(SuccessReply, nil)
	}

	reply.Metadata().Set(ReplyOutcomeHandler, OutcomeSuccess)

	return h.applyCorrelationHeaders(reply, cmd)
}

func (h commandMessageHandler) failure(reply ddd.Reply, cmd ddd.Command) ddd.Reply {
	if reply == nil {
		reply = ddd.NewReply(FailureReply, nil)
	}

	reply.Metadata().Set(ReplyOutcomeHandler, OutcomeFailure)

	return h.applyCorrelationHeaders(reply, cmd)
}

func (h commandMessageHandler) applyCorrelationHeaders(reply ddd.Reply, cmd ddd.Command) ddd.Reply {
	for key, value := range cmd.Metadata() {
		if key == CommandNameHandler {
			continue
		}

		if strings.HasPrefix(key, CommandHandlerPrefix) {
			hdr := ReplyHandlerPrefix + key[len(CommandHandlerPrefix):]
			reply.Metadata().Set(hdr, value)
		}
	}

	return reply
}
