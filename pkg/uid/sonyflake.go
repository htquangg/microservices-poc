package uid

import (
	"os"
	"strconv"

	"github.com/sony/sonyflake"
)

type Sonyflake struct {
	sf       *sonyflake.Sonyflake
	serverID uint16
}

func NewSonyflake() *Sonyflake {
	var serverID uint16

	if val, err := strconv.ParseInt(os.Getenv("SERVER_ID"), 10, 16); err == nil {
		serverID = uint16(val)
	}

	if serverID > 0 {
		return &Sonyflake{
			sf: sonyflake.NewSonyflake(sonyflake.Settings{
				MachineID: func() (uint16, error) {
					return serverID, nil
				},
			}),
			serverID: serverID,
		}
	}

	return &Sonyflake{
		sf: sonyflake.NewSonyflake(sonyflake.Settings{}),
	}
}

func (s *Sonyflake) ID() string {
	id, err := s.sf.NextID()
	if err != nil {
		return s.ID()
	}

	return strconv.FormatUint(id, 10)
}

func (s *Sonyflake) ServerID() uint16 {
	return s.serverID
}
