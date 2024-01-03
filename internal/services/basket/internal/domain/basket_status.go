package domain

type BasketStatus string

const (
	BasketUnknown      BasketStatus = ""
	BasketIsOpen       BasketStatus = "open"
	BasketIsCancelled  BasketStatus = "cancelled"
	BasketIsCheckedOut BasketStatus = "checked_out"
)

func (s BasketStatus) String() string {
	switch s {
	case BasketIsOpen, BasketIsCancelled, BasketIsCheckedOut:
		return string(s)
	default:
		return ""
	}
}

func ToBasketStatus(status string) BasketStatus {
	switch status {
	case BasketIsOpen.String():
		return BasketIsOpen
	case BasketIsCancelled.String():
		return BasketIsCancelled
	case BasketIsCheckedOut.String():
		return BasketIsCheckedOut
	default:
		return BasketUnknown
	}
}
