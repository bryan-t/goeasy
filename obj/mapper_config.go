package obj

import (
	"reflect"
)

type structMapKey struct {
	source      reflect.Type
	destination reflect.Type
}

type mapperConfig struct {
	//fieldMaps  map[structMapKey]map[string]*FieldMapConfig
	converters map[structMapKey]func(src any) (dst any, err error)
}

// AddTypeConverter allows adding a custom type converter for a given source and destination type.
func AddTypeConverter[sourceT any, destinationT any](mapper *Mapper, fn func(src any) (dst any, err error)) {
	var zeroSource sourceT
	var zeroDestination destinationT
	sourceType := reflect.TypeOf(zeroSource)
	destinationType := reflect.TypeOf(zeroDestination)
	structKey := structMapKey{
		source:      sourceType,
		destination: destinationType,
	}
	mapper.cfg.converters[structKey] = fn
}

/* Commenting out due to problems with supporting nested structs and assignable
// FieldMapConfig contains configuration on how to transform a given field of
// a struct.
type FieldMapConfig struct {
	Source              string
	Destination         string
	GetDestinationValue func(source any) (any, error)
}





// ConfigureFieldMaps allows overriding of how fields are mapped for sourceT and destinationT
func ConfigureFieldMaps[sourceT any, destinationT any](mapper *Mapper,
	fieldMapConfigs ...FieldMapConfig) error {
	var zeroSource sourceT
	var zeroDestination destinationT
	sourceType := reflect.TypeOf(zeroSource)
	destinationType := reflect.TypeOf(zeroDestination)
	if sourceType.Kind() != reflect.Struct || destinationType.Kind() != reflect.Struct {
		return fmt.Errorf("sourceT and destinationT must be structs")
	}

	structKey := structMapKey{
		source:      sourceType,
		destination: destinationType,
	}

	fieldMap := mapper.cfg.fieldMaps[structKey]
	if fieldMap == nil {
		fieldMap = make(map[string]*FieldMapConfig)
	}

	for _, cfg := range fieldMapConfigs {
		if cfg.Destination == "" {
			return fmt.Errorf("destination field names must be provided")
		}

		fieldMap[cfg.Destination] = &cfg
	}
	mapper.cfg.fieldMaps[structKey] = fieldMap
	return nil
}
*/
