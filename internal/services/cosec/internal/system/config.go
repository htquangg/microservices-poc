package system

import (
	"fmt"

	"github.com/htquangg/microservices-poc/internal/services/cosec/config"
	"github.com/htquangg/microservices-poc/pkg/uid"
)

type Config struct {
	*config.Config
	webID   string
	webName string
	rpcID   string
	rpcName string
}

func (s *System) webID() string {
	if s.cfg.webID == "" {
		id := uid.GetManager().ID()
		s.cfg.webID = fmt.Sprintf("http-%s-svc-%s", s.cfg.Name, id)
	}
	return s.cfg.webID
}

func (s *System) webName() string {
	if s.cfg.webName == "" {
		s.cfg.webName = fmt.Sprintf("http-%s-svc", s.cfg.Name)
	}
	return s.cfg.webName
}

func (s *System) rpcID() string {
	if s.cfg.rpcID == "" {
		id := uid.GetManager().ID()
		s.cfg.rpcID = fmt.Sprintf("rpc-%s-svc-%s", s.cfg.Name, id)
	}
	return s.cfg.rpcID
}

func (s *System) rpcName() string {
	if s.cfg.rpcName == "" {
		s.cfg.rpcName = fmt.Sprintf("rpc-%s-svc", s.cfg.Name)
	}
	return s.cfg.rpcName
}
