package goda

import (
	"errors"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/shopspring/decimal"
)

// parse implements the common parse functionality, using reflection
// The type and value passed in shall represent the reflected type / value of the struct we want to parse the data into
func parse(s string, t reflect.Type, v reflect.Value) error {
	// fmt.Printf("Type %v value %v\n", t, v)
	// Loop through each of the struct fields
	for i := 0; i < t.NumField(); i++ {
		tt := t.Field(i)
		vv := v.Field(i)
		// fmt.Printf("Field %v: name %v, type %v tag %v value %v\n", i, tt.Name, tt.Type, tt.Tag, vv)
		// Get the offset and length from the annotated struct fields
		offset, err := strconv.Atoi(tt.Tag.Get("offset"))
		if err != nil {
			return err
		}
		length, err := strconv.Atoi(tt.Tag.Get("length"))
		if err != nil {
			return err
		}
		// Switch to correct parser by using the field type
		switch vv.Type() {
		case reflect.TypeOf(int(0)):
			value, err := strconv.Atoi(s[offset : offset+length])
			if err != nil {
				return err
			}
			// fmt.Printf("Value is %d\n", value)
			vv.SetInt(int64(value))
		case reflect.TypeOf(string("")):
			value := strings.TrimSpace(s[offset : offset+length])
			// fmt.Printf("Value is %s\n", value)
			vv.SetString(value)
		case reflect.TypeOf(bool(false)):
			value := s[offset:offset+length] != ""
			// fmt.Printf("Value is %t\n", value)
			vv.SetBool(value)
		case reflect.TypeOf(time.Time{}):
			value, err := time.Parse("020106", s[offset:offset+length])
			if err != nil {
				return err
			}
			// fmt.Printf("Value is %v\n", value)
			vv.Set(reflect.ValueOf(value))
		case reflect.TypeOf(decimal.Decimal{}):
			value, err := strconv.Atoi(s[offset : offset+length])
			if err != nil {
				return err
			}
			// Shift decimal 3 places
			balance := decimal.New(int64(value), -3)
			vv.Set(reflect.ValueOf(balance))
		}
	}
	return nil
}

// Parse is the generic record string line parser
// Each line passed is in is a raw CODA record line
// The line is matched to its specific record struct and the data is read into it
func Parse(line string) (r Record, err error) {
	if strings.HasPrefix(line, "0") {
		r = &InitialRecord{}
	} else if strings.HasPrefix(line, "1") {
		r = &OldBalanceRecord{}
	} else if strings.HasPrefix(line, "21") {
		r = &MovementRecord1{}
	} else if strings.HasPrefix(line, "22") {
		r = &MovementRecord2{}
	} else if strings.HasPrefix(line, "23") {
		r = &MovementRecord3{}
	} else if strings.HasPrefix(line, "31") {
		r = &InformationRecord1{}
	} else if strings.HasPrefix(line, "32") {
		r = &InformationRecord2{}
	} else if strings.HasPrefix(line, "33") {
		r = &InformationRecord3{}
	} else if strings.HasPrefix(line, "8") {
		r = &NewBalanceRecord{}
	} else if strings.HasPrefix(line, "4") {
		r = &FreeCommunicationRecord{}
	} else if strings.HasPrefix(line, "9") {
		r = &TrailerRecord{}
	} else {
		return nil, nil
	}

	err = r.Parse(line)
	return r, err
}

type Record interface {
	Parse(string) error
}

// Parse populates an initial record from a string
func (r *InitialRecord) Parse(s string) error {
	if !strings.HasPrefix(s, "0") {
		return errors.New("Wrong prefix")
	}

	err := parse(s, reflect.TypeOf(r).Elem(), reflect.ValueOf(r).Elem())
	if err != nil {
		return err
	}
	// Specific duplicate check
	r.IsDuplicate = s[16:17] == "D"
	return nil
}

// Parse populates an old balance record from a string
func (r *OldBalanceRecord) Parse(s string) error {
	if !strings.HasPrefix(s, "1") {
		return errors.New("Wrong prefix")
	}

	return parse(s, reflect.TypeOf(r).Elem(), reflect.ValueOf(r).Elem())
}

// Parse populates a movement record from a string
func (r *MovementRecord1) Parse(s string) error {
	if !strings.HasPrefix(s, "21") {
		return errors.New("Wrong prefix")
	}

	return parse(s, reflect.TypeOf(r).Elem(), reflect.ValueOf(r).Elem())
}

// Parse populates a movement record 2 from a string
func (r *MovementRecord2) Parse(s string) error {
	if !strings.HasPrefix(s, "22") {
		return errors.New("Wrong prefix")
	}

	return parse(s, reflect.TypeOf(r).Elem(), reflect.ValueOf(r).Elem())
}

// Parse populates a movement record 3 from a string
func (r *MovementRecord3) Parse(s string) error {
	if !strings.HasPrefix(s, "23") {
		return errors.New("Wrong prefix")
	}

	return parse(s, reflect.TypeOf(r).Elem(), reflect.ValueOf(r).Elem())
}

func (r *InformationRecord1) Parse(s string) error {
	if !strings.HasPrefix(s, "31") {
		return errors.New("Wrong prefix")
	}

	return parse(s, reflect.TypeOf(r).Elem(), reflect.ValueOf(r).Elem())
}

func (r *InformationRecord2) Parse(s string) error {
	if !strings.HasPrefix(s, "32") {
		return errors.New("Wrong prefix")
	}

	return parse(s, reflect.TypeOf(r).Elem(), reflect.ValueOf(r).Elem())
}

func (r *InformationRecord3) Parse(s string) error {
	if !strings.HasPrefix(s, "33") {
		return errors.New("Wrong prefix")
	}

	return parse(s, reflect.TypeOf(r).Elem(), reflect.ValueOf(r).Elem())
}

func (r *NewBalanceRecord) Parse(s string) error {
	if !strings.HasPrefix(s, "8") {
		return errors.New("Wrong prefix")
	}

	return parse(s, reflect.TypeOf(r).Elem(), reflect.ValueOf(r).Elem())
}

func (r *FreeCommunicationRecord) Parse(s string) error {
	if !strings.HasPrefix(s, "4") {
		return errors.New("Wrong prefix")
	}

	return parse(s, reflect.TypeOf(r).Elem(), reflect.ValueOf(r).Elem())
}

func (r *TrailerRecord) Parse(s string) error {
	if !strings.HasPrefix(s, "9") {
		return errors.New("Wrong prefix")
	}

	return parse(s, reflect.TypeOf(r).Elem(), reflect.ValueOf(r).Elem())
}
