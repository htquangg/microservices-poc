package uid

import (
	"os"
	"strconv"

	sf "github.com/sony/sonyflake"
)

type Sonyflake struct {
	sonyflake *sf.Sonyflake
}

func NewSonyflake() *Sonyflake {
	var serverID uint16

	if val, err := strconv.ParseInt(os.Getenv("SERVER_ID"), 10, 16); err == nil {
		serverID = uint16(val)
	}

	if serverID > 0 {
		return &Sonyflake{
			sf.NewSonyflake(sf.Settings{
				MachineID: func() (uint16, error) {
					return serverID, nil
				},
			}),
		}
	}

	return &Sonyflake{
		sf.NewSonyflake(sf.Settings{}),
	}
}

func (s *Sonyflake) ID() string {
	id, err := s.sonyflake.NextID()
	if err != nil {
		return s.ID()
	}

	return strconv.FormatUint(id, 10)
}
