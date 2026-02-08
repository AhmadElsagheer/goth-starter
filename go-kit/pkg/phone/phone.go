package phone

import (
	"database/sql/driver"
	"fmt"
	"strings"

	"github.com/nyaruka/phonenumbers"
)

const defaultRegion = "EG"

type PhoneNumber struct {
	*phonenumbers.PhoneNumber
}

func MustParse(s string) PhoneNumber {
	num, err := Parse(s)
	if err != nil {
		panic(err)
	}
	return num
}

func Parse(s string) (PhoneNumber, error) {
	num, err := phonenumbers.Parse(s, defaultRegion)
	if err != nil {
		return PhoneNumber{}, err
	}
	return PhoneNumber{PhoneNumber: num}, nil
}

func (d PhoneNumber) String() string {
	return phonenumbers.Format(d.PhoneNumber, phonenumbers.E164)
}

func (d *PhoneNumber) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), "\"")
	if s == "null" || s == "" {
		return nil
	}
	p, err := phonenumbers.Parse(s, defaultRegion)
	if err != nil {
		return err
	}
	d.PhoneNumber = p
	return nil
}

func (d PhoneNumber) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("\"%s\"", d.PhoneNumber.String())), nil
}

func (d *PhoneNumber) UnmarshalText(b []byte) error {
	return d.UnmarshalJSON(b)
}

func (d PhoneNumber) MarshalText() ([]byte, error) {
	return d.MarshalJSON()
}

func (d PhoneNumber) Value() (driver.Value, error) {
	return d.PhoneNumber.String(), nil
}

func (d *PhoneNumber) Scan(value interface{}) error {
	if t, ok := value.(string); ok {
		var err error
		d.PhoneNumber, err = phonenumbers.Parse(t, defaultRegion)
		return err
	}
	return fmt.Errorf("cannot scan type %T into PhoneNumber", value)
}
