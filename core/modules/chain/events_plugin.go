package chain

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/centrifuge/go-substrate-rpc-client/v4/scale"
	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
	"github.com/ethereum/go-ethereum/log"
	"reflect"
)

func DecodeEventRecordsWithIgnoreError(e types.EventRecordsRaw, m *types.Metadata, t interface{}) error {
	log.Info(fmt.Sprintf("will decode event records from raw hex: %#x", e))

	// ensure t is a pointer
	ttyp := reflect.TypeOf(t)
	if ttyp.Kind() != reflect.Ptr {
		return errors.New("target must be a pointer, but is " + fmt.Sprint(ttyp))
	}
	// ensure t is not a nil pointer
	tval := reflect.ValueOf(t)
	if tval.IsNil() {
		return errors.New("target is a nil pointer")
	}
	val := tval.Elem()
	typ := val.Type()
	// ensure val can be set
	if !val.CanSet() {
		return fmt.Errorf("unsettable value %v", typ)
	}
	// ensure val points to a struct
	if val.Kind() != reflect.Struct {
		return fmt.Errorf("target must point to a struct, but is " + fmt.Sprint(typ))
	}

	decoder := scale.NewDecoder(bytes.NewReader(e))

	// determine number of events
	n, err := decoder.DecodeUintCompact()
	if err != nil {
		return err
	}

	log.Info(fmt.Sprintf("found %v events", n))

	// iterate over events
	for i := uint64(0); i < n.Uint64(); i++ {
		log.Info(fmt.Sprintf("decoding event #%v", i))

		// decode Phase
		phase := types.Phase{}
		err := decoder.Decode(&phase)
		if err != nil {
			return fmt.Errorf("unable to decode Phase for event #%v: %v", i, err)
		}

		// decode EventID
		id := types.EventID{}
		err = decoder.Decode(&id)
		if err != nil {
			return fmt.Errorf("unable to decode EventID for event #%v: %v", i, err)
		}

		log.Info(fmt.Sprintf("event #%v has EventID %v", i, id))

		// ask metadata for method & event name for event
		moduleName, eventName, err := m.FindEventNamesForEventID(id)
		// moduleName, eventName, err := "System", "ExtrinsicSuccess", nil
		if err != nil {
			//return fmt.Errorf("unable to find event with EventID %v in metadata for event #%v: %s", id, i, err)
			log.Warn("unable to find event with EventID %v in metadata for event #%v: %s", id, i, err)

			return err
		}

		log.Info(fmt.Sprintf("event #%v is in module %v with event name %v", i, moduleName, eventName))

		// check whether name for eventID exists in t
		field := val.FieldByName(fmt.Sprintf("%v_%v", moduleName, eventName))
		var holder reflect.Value
		if !field.IsValid() {
			log.Info(fmt.Sprintf("unable to find field %v_%v for event #%v with EventID %v ", moduleName, eventName, i, id))
			holder, err = getNewStructWithEventInfo(m, moduleName, eventName)
			if err != nil {
				return err
			}
		} else {
			// create a pointer to with the correct type that will hold the decoded event
			holder = reflect.New(field.Type().Elem())
		}

		// ensure first field is for Phase, last field is for Topics
		numFields := holder.Elem().NumField()
		if numFields < 2 {
			return fmt.Errorf("expected event #%v with EventID %v, field %v_%v to have at least 2 fields "+
				"(for Phase and Topics), but has %v fields", i, id, moduleName, eventName, numFields)
		}
		phaseField := holder.Elem().FieldByIndex([]int{0})
		if phaseField.Type() != reflect.TypeOf(phase) {
			return fmt.Errorf("expected the first field of event #%v with EventID %v, field %v_%v to be of type "+
				"types.Phase, but got %v", i, id, moduleName, eventName, phaseField.Type())
		}
		topicsField := holder.Elem().FieldByIndex([]int{numFields - 1})
		if topicsField.Type() != reflect.TypeOf([]types.Hash{}) {
			return fmt.Errorf("expected the last field of event #%v with EventID %v, field %v_%v to be of type "+
				"[]types.Hash for Topics, but got %v", i, id, moduleName, eventName, topicsField.Type())
		}

		// set the phase we decoded earlier
		phaseField.Set(reflect.ValueOf(phase))

		// set the remaining fields
		for j := 1; j < numFields; j++ {
			err = decoder.Decode(holder.Elem().FieldByIndex([]int{j}).Addr().Interface())
			if err != nil {
				return fmt.Errorf("unable to decode field %v event #%v with EventID %v, field %v_%v: %v", j, i, id, moduleName,
					eventName, err)
			}
		}

		if field.IsValid() {
			// add the decoded event to the slice
			field.Set(reflect.Append(field, holder.Elem()))

			log.Info(fmt.Sprintf("decoded event #%v", i))
		}
	}
	return nil
}

func getNewStructWithEventInfo(m *types.Metadata, moduleName types.Text, eventName types.Text) (reflect.Value, error) {
	var ts []reflect.Type
	for _, mod := range m.AsMetadataV13.Modules {
		if mod.Name == moduleName {
			for _, event := range mod.Events {
				if event.Name == eventName {
					ts = append(ts, reflect.TypeOf(types.Phase{}))
					for _, arg := range event.Args {
						a, ok := TypeMap[string(arg)]
						if !ok {
							err := fmt.Errorf("unable to find type for arg %v of event %v_%v", arg, moduleName, eventName)
							return reflect.Value{}, err
						}
						ts = append(ts, a)
					}
					ts = append(ts, reflect.TypeOf([]types.Hash{}))
				}
			}
		}
	}
	newStruct := makeStruct(ts...)
	return newStruct, nil
}

func makeStruct(ts ...reflect.Type) reflect.Value {
	var sfs []reflect.StructField
	for i, t := range ts {
		sf := reflect.StructField{
			Name: fmt.Sprintf("F%d", i+1),
			Type: t,
		}
		sfs = append(sfs, sf)
	}
	st := reflect.StructOf(sfs)
	so := reflect.New(st)
	return so
}

var TypeMap = map[string]reflect.Type{
	"u8":           reflect.TypeOf(types.U8(0)),
	"u16":          reflect.TypeOf(types.U16(0)),
	"u32":          reflect.TypeOf(types.U32(0)),
	"u64":          reflect.TypeOf(types.U64(0)),
	"u128":         reflect.TypeOf(types.U128{}),
	"u256":         reflect.TypeOf(types.U256{}),
	"i8":           reflect.TypeOf(types.I8(0)),
	"i16":          reflect.TypeOf(types.I16(0)),
	"i32":          reflect.TypeOf(types.I32(0)),
	"i64":          reflect.TypeOf(types.I64(0)),
	"i128":         reflect.TypeOf(types.I128{}),
	"i256":         reflect.TypeOf(types.I256{}),
	"bool":         reflect.TypeOf(types.Bool(false)),
	"text":         reflect.TypeOf(types.Text("")),
	"hash":         reflect.TypeOf(types.Hash{}),
	"address":      reflect.TypeOf(types.Address{}),
	"AccountId":    reflect.TypeOf(types.AccountID{}),
	"BalanceOf<T>": reflect.TypeOf(types.U128{}),
}
