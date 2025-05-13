package obj

import (
	"fmt"
	"math"
	"reflect"
	"testing"
	"time"
	"unsafe"

	"github.com/stretchr/testify/assert"
)

type testAllTypes struct {
	Bool       bool
	Int        int
	Int8       int8
	Int16      int16
	Int32      int32
	Int64      int64
	Uint       uint
	Uint8      uint8
	Uint16     uint16
	Uint32     uint32
	Uint64     uint64
	Float32    float32
	Float64    float64
	Complex64  complex64
	Complex128 complex128
	Pointer    *int
	String     string
}

func TestMap(t *testing.T) {
	i := 1
	src := struct {
		AllTypes testAllTypes
	}{
		AllTypes: testAllTypes{
			Bool:       true,
			Int:        math.MaxInt,
			Int8:       math.MaxInt8,
			Int16:      math.MaxInt16,
			Int32:      math.MaxInt32,
			Int64:      math.MaxInt64,
			Uint:       math.MaxUint,
			Uint8:      math.MaxUint8,
			Uint16:     math.MaxUint16,
			Uint32:     math.MaxUint32,
			Uint64:     math.MaxUint64,
			Float32:    math.MaxFloat32,
			Float64:    math.MaxFloat64,
			Complex64:  complex(math.MaxFloat32, math.MaxFloat32),
			Complex128: complex(math.MaxFloat64, math.MaxFloat64),
			Pointer:    &i,
			String:     "test",
		},
	}

	mapper := NewMapper()
	dst := struct {
		AllTypes testAllTypes
	}{
		AllTypes: testAllTypes{},
	}
	err := mapper.Map(src, &dst)

	assert.Nil(t, err, "Map returned an error")
	assert.Equal(t, src, dst, "Not equal")
}

func TestMapNoEquivalentField(t *testing.T) {
	mapper := NewMapper()
	i := 1
	src := testAllTypes{
		Bool:       true,
		Int:        math.MaxInt,
		Int8:       math.MaxInt8,
		Int16:      math.MaxInt16,
		Int32:      math.MaxInt32,
		Int64:      math.MaxInt64,
		Uint:       math.MaxUint,
		Uint8:      math.MaxUint8,
		Uint16:     math.MaxUint16,
		Uint32:     math.MaxUint32,
		Uint64:     math.MaxUint64,
		Float32:    math.MaxFloat32,
		Float64:    math.MaxFloat64,
		Complex64:  complex(math.MaxFloat32, math.MaxFloat32),
		Complex128: complex(math.MaxFloat64, math.MaxFloat64),
		Pointer:    &i,
		String:     "test",
	}
	dst := struct {
		Int  int
		Test string
	}{}

	err := mapper.Map(src, &dst)
	assert.Nil(t, err, "Map returned an error")
	assert.Equal(t, src.Int, dst.Int, "Int are not equal")
	assert.Empty(t, dst.Test)
}

func TestMapArrayDst(t *testing.T) {
	type IntStruct struct {
		Int int
	}
	tests := []struct {
		name string
		src  interface{}
		dst  [3]IntStruct
		err  error
	}{
		{
			name: "Equal length array",
			src:  [3]IntStruct{{1}, {2}, {3}},
		},
		{
			name: "Longer array",
			src:  [4]IntStruct{{1}, {2}, {3}, {4}},
			err:  ErrInsufficientCapacity,
		},
		{
			name: "Shorter array",
			src:  [3]IntStruct{{1}, {2}},
		},
		{
			name: "Equal length slice",
			src:  []IntStruct{{1}, {2}, {3}},
		},
		{
			name: "Longer slice",
			src:  []IntStruct{{1}, {2}, {3}, {4}},
			err:  ErrInsufficientCapacity,
		},
		{
			name: "Shorter slice",
			src:  make([]IntStruct, 2, 5),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mapper := NewMapper()
			err := mapper.Map(test.src, &test.dst)
			if err != nil || test.err != nil {
				assert.Equal(t, test.err, err, "Error not equal")
				return
			}
			srcV := reflect.ValueOf(test.src)
			for i := range srcV.Len() {
				v := srcV.Index(i).Field(0)
				assert.Equal(t, int(v.Int()), test.dst[i].Int, "Has unequal value")
			}
		})
	}
}

func TestMapMismatchType(t *testing.T) {
	i := 1
	src := struct {
		Int    int
		String string
		Ptr    *int
	}{Ptr: &i}
	tests := []struct {
		name string
		dst  interface{}
	}{
		{
			name: "Bool",
			dst: &struct {
				Int bool
			}{},
		},
		{
			name: "Int",
			dst: &struct {
				String int
			}{},
		},
		{
			name: "Int8",
			dst: &struct {
				Int int8
			}{},
		},
		{
			name: "Int16",
			dst: &struct {
				Int int8
			}{},
		},
		{
			name: "Int32",
			dst: &struct {
				Int int8
			}{},
		},
		{
			name: "Int64",
			dst: &struct {
				Int int64
			}{},
		},
		{
			name: "Uint",
			dst: &struct {
				Int uint
			}{},
		},
		{
			name: "Uint8",
			dst: &struct {
				Int uint8
			}{},
		},
		{
			name: "Uint16",
			dst: &struct {
				Int uint16
			}{},
		},
		{
			name: "Uint32",
			dst: &struct {
				Int uint32
			}{},
		},
		{
			name: "Uint64",
			dst: &struct {
				Int uint64
			}{},
		},
		{
			name: "Float32",
			dst: &struct {
				Int float32
			}{},
		},
		{
			name: "Float64",
			dst: &struct {
				Int float64
			}{},
		},
		{
			name: "Complex64",
			dst: &struct {
				Int complex64
			}{},
		},
		{
			name: "Complex128",
			dst: &struct {
				Int complex128
			}{},
		},
		{
			name: "Pointer",
			dst: &struct {
				Ptr *int8
			}{},
		},
		{
			name: "String",
			dst: &struct {
				Int string
			}{},
		},
	}

	for _, test := range tests {
		//t.Run(test.name, func(t *testing.T) {
		mapper := NewMapper()
		err := mapper.Map(src, test.dst)
		assert.Equal(t, ErrMismatchType, err)
		//})
	}
}

func TestMapSliceDst(t *testing.T) {
	type IntStruct struct {
		Int int
	}
	tests := []struct {
		name string
		src  interface{}
		dst  []IntStruct
		err  error
	}{
		{
			name: "Array source",
			src:  [3]IntStruct{{1}, {2}, {3}},
		},
		{
			name: "Wrong source type",
			src:  1,
			err:  ErrMismatchType,
		},
		{
			name: "Slice source",
			src:  [3]IntStruct{{1}, {2}},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mapper := NewMapper()
			err := mapper.Map(test.src, &test.dst)
			if err != nil || test.err != nil {
				assert.Equal(t, test.err, err, "Error not equal")
				return
			}
			srcV := reflect.ValueOf(test.src)
			for i := range srcV.Len() {
				v := srcV.Index(i).Field(0)
				assert.Equal(t, int(v.Int()), test.dst[i].Int, "Has unequal value")
			}
		})
	}
}

func TestMapMapDst(t *testing.T) {
	type IntStruct struct {
		Int int
	}
	tests := []struct {
		name string
		src  interface{}
		dst  map[int]IntStruct
		err  error
	}{
		{
			name: "Success",
			src:  map[int]IntStruct{1: {1}, 2: {2}},
		},
		{
			name: "Wrong source type",
			src:  1,
			err:  ErrMismatchType,
		},
		{
			name: "Wrong key type",
			src:  map[string]IntStruct{"1": {1}},
			err:  ErrMismatchType,
		},
		{
			name: "Wrong value type",
			src:  map[int]int{1: 1},
			err:  ErrMismatchType,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mapper := NewMapper()
			err := mapper.Map(test.src, &test.dst)
			if err != nil || test.err != nil {
				assert.Equal(t, test.err, err, "Error not equal")
				return
			}
			assert.Equal(t, test.src, test.dst)
		})
	}
}

// AI generated code start
type testUserDTO struct {
	ID             int
	withGetterName string
}

type testUser struct {
	ID   int
	Name string
}

func (u testUserDTO) GetName() string {
	return "Mr. " + u.withGetterName
}

func TestMapWithGetter(t *testing.T) {
	dto := testUserDTO{
		ID:             1,
		withGetterName: "John",
	}

	user := testUser{}
	mapper := NewMapper()
	err := mapper.Map(dto, &user)

	assert.Nil(t, err, "Map returned an error")
	assert.Equal(t, dto.ID, user.ID, "ID not equal")
	assert.Equal(t, "Mr. John", user.Name, "Name not equal")
}

type testUserWithSetter struct {
	withSetterID   int
	withSetterName string
}

func (u *testUserWithSetter) SetID(id int) {
	u.withSetterID = id
}

func (u *testUserWithSetter) SetName(name string) {
	u.withSetterName = name
}

func TestMapWithSetter(t *testing.T) {
	dto := testUserDTO{
		ID:             1,
		withGetterName: "John",
	}

	user := testUserWithSetter{}
	mapper := NewMapper()
	err := mapper.Map(dto, &user)

	assert.Nil(t, err, "Map returned an error")
	assert.Equal(t, dto.ID, user.withSetterID, "ID not equal")
	assert.Equal(t, dto.GetName(), user.withSetterName, "Name not equal")
}

type testUserWithDifferentSetter struct {
	withSetterID   int
	withSetterName string
}

func (u *testUserWithDifferentSetter) SetID(id int64) {
	u.withSetterID = int(id)
}

func (u *testUserWithDifferentSetter) SetName(name string) {
	u.withSetterName = name
}

func TestMapWithDifferentSetter(t *testing.T) {
	dto := testUserDTO{
		ID:             1,
		withGetterName: "John",
	}

	user := testUserWithDifferentSetter{}
	mapper := NewMapper()
	err := mapper.Map(dto, &user)

	assert.Equal(t, ErrMismatchType, err, "Error not equal")
	assert.Equal(t, 0, user.withSetterID, "ID not equal")
	assert.Equal(t, "", user.withSetterName, "Name not equal")
}

type testUserWithGetterAndSetter struct {
	withGetterSetterID   int
	withGetterSetterName string
}

func (u *testUserWithGetterAndSetter) SetID(id int) {
	u.withGetterSetterID = id
}

func (u *testUserWithGetterAndSetter) SetName(name string) {
	u.withGetterSetterName = name
}

func (u testUserDTO) GetID() int {
	return u.ID
}

func TestMapWithGetterAndSetter(t *testing.T) {
	dto := testUserDTO{
		ID:             1,
		withGetterName: "John",
	}

	user := testUserWithGetterAndSetter{}
	mapper := NewMapper()
	err := mapper.Map(dto, &user)

	assert.Nil(t, err, "Map returned an error")
	assert.Equal(t, dto.GetID(), user.withGetterSetterID, "ID not equal")
	assert.Equal(t, dto.GetName(), user.withGetterSetterName, "Name not equal")
}

type testUserWithSetterNoField struct {
	withSetterID int
}

func (u *testUserWithSetterNoField) SetID(id int) {
	u.withSetterID = id
}

func TestMapWithSetterNoField(t *testing.T) {
	dto := testUserDTO{
		ID:             1,
		withGetterName: "John",
	}

	user := testUserWithSetterNoField{}
	mapper := NewMapper()
	err := mapper.Map(dto, &user)

	assert.Nil(t, err, "Map returned an error")
	assert.Equal(t, dto.ID, user.withSetterID, "ID not equal")
}

// AI generated code start
func TestMapWithInterfaceField(t *testing.T) {
	type StructWithInterface struct {
		Data interface{}
	}

	tests := []struct {
		name     string
		src      StructWithInterface
		dst      StructWithInterface
		expected interface{}
		err      error
	}{
		{
			name:     "String interface",
			src:      StructWithInterface{Data: "test string"},
			dst:      StructWithInterface{},
			expected: "test string",
		},
		{
			name:     "Int interface",
			src:      StructWithInterface{Data: 42},
			dst:      StructWithInterface{},
			expected: 42,
		},
		{
			name:     "Struct interface",
			src:      StructWithInterface{Data: testAllTypes{Int: 1, String: "test"}},
			dst:      StructWithInterface{},
			expected: testAllTypes{Int: 1, String: "test"},
		},
		{
			name:     "Nil interface",
			src:      StructWithInterface{Data: nil},
			dst:      StructWithInterface{},
			expected: nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mapper := NewMapper()
			err := mapper.Map(test.src, &test.dst)
			assert.Equal(t, test.err, err, "Error not equal")
			assert.Equal(t, test.expected, test.dst.Data, "Data not equal")
		})
	}
}

// AI generated code end

// AI generated code start
func TestMapWithNestedInterfaceField(t *testing.T) {
	type InnerStruct struct {
		Data interface{}
	}

	type OuterStruct struct {
		Inner interface{}
	}

	tests := []struct {
		name     string
		src      OuterStruct
		dst      OuterStruct
		expected interface{}
		err      error
	}{
		{
			name: "Nested string interface",
			src:  OuterStruct{Inner: InnerStruct{Data: "test string"}},
			dst:  OuterStruct{},
			expected: InnerStruct{
				Data: "test string",
			},
		},
		{
			name: "Nested int interface",
			src:  OuterStruct{Inner: InnerStruct{Data: 42}},
			dst:  OuterStruct{},
			expected: InnerStruct{
				Data: 42,
			},
		},
		{
			name: "Nested struct interface",
			src:  OuterStruct{Inner: InnerStruct{Data: testAllTypes{Int: 1, String: "test"}}},
			dst:  OuterStruct{},
			expected: InnerStruct{
				Data: testAllTypes{Int: 1, String: "test"},
			},
		},
		{
			name: "Nested nil interface",
			src:  OuterStruct{Inner: InnerStruct{Data: nil}},
			dst:  OuterStruct{},
			expected: InnerStruct{
				Data: nil,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mapper := NewMapper()
			err := mapper.Map(test.src, &test.dst)
			assert.Equal(t, test.err, err, "Error not equal")
			assert.Equal(t, test.expected, test.dst.Inner, "Inner not equal")
		})
	}
}

// AI generated code end

// AI generated code start
func TestMapNotAddressable(t *testing.T) {
	type TestStruct struct {
		Field int
	}

	src := TestStruct{Field: 42}
	dst := TestStruct{}

	mapper := NewMapper()
	err := mapper.Map(src, dst) // Passing non-addressable value as destination

	assert.Equal(t, ErrNotAddresable, err, "Error not equal")
}

// AI generated code end

// AI generated code start
func TestMapWithTypeConverter(t *testing.T) {
	type Source struct {
		Value string
	}
	type Destination struct {
		Value int
	}

	mapper := NewMapper()
	AddTypeConverter[string, int](mapper, func(src any) (any, error) {
		str, ok := src.(string)
		if !ok {
			return nil, ErrMismatchType
		}
		switch str {
		case "one":
			return 1, nil
		case "two":
			return 2, nil
		default:
			return 0, nil
		}
	})

	src := Source{Value: "one"}
	dst := Destination{}
	err := mapper.Map(src, &dst)
	assert.Nil(t, err, "Map returned an error")
	assert.Equal(t, 1, dst.Value, "Type converter did not convert value correctly")

	src = Source{Value: "two"}
	dst = Destination{}
	err = mapper.Map(src, &dst)
	assert.Nil(t, err, "Map returned an error")
	assert.Equal(t, 2, dst.Value, "Type converter did not convert value correctly")

	src = Source{Value: "unknown"}
	dst = Destination{}
	err = mapper.Map(src, &dst)
	assert.Nil(t, err, "Map returned an error")
	assert.Equal(t, 0, dst.Value, "Type converter did not convert value correctly")
}

func TestMapWithTypeConverterError(t *testing.T) {
	type Source struct {
		Value string
	}
	type Destination struct {
		Value int
	}

	mapper := NewMapper()
	AddTypeConverter[string, int](mapper, func(src any) (any, error) {
		return nil, fmt.Errorf("conversion error")
	})

	src := Source{Value: "fail"}
	dst := Destination{}
	err := mapper.Map(src, &dst)
	assert.EqualError(t, err, "conversion error")
}

func TestMapWithTypeConverterMismatchType(t *testing.T) {
	type Source struct {
		Value string
	}
	type Destination struct {
		Value int
	}

	mapper := NewMapper()
	AddTypeConverter[string, int](mapper, func(src any) (any, error) {
		return "not an int", nil
	})

	src := Source{Value: "fail"}
	dst := Destination{}
	err := mapper.Map(src, &dst)
	assert.Equal(t, ErrMismatchType, err)
}

//
// AI generated code end
//
// AI generated code start

func TestMapAssignableStruct_timeTime(t *testing.T) {
	type Src struct {
		CreatedAt time.Time
	}
	type Dst struct {
		CreatedAt time.Time
	}

	now := time.Now()
	src := Src{CreatedAt: now}
	dst := Dst{}

	mapper := NewMapper()
	err := mapper.Map(src, &dst)
	assert.Nil(t, err, "Map returned an error")
	assert.Equal(t, src.CreatedAt, dst.CreatedAt, "time.Time values not equal")
}

func TestMapAssignableStruct_timeTimePointer(t *testing.T) {
	type Src struct {
		CreatedAt *time.Time
	}
	type Dst struct {
		CreatedAt *time.Time
	}

	now := time.Now()
	src := Src{CreatedAt: &now}
	dst := Dst{}

	mapper := NewMapper()
	err := mapper.Map(src, &dst)
	assert.Nil(t, err, "Map returned an error")
	assert.NotNil(t, dst.CreatedAt, "Destination pointer is nil")
	assert.Equal(t, *src.CreatedAt, *dst.CreatedAt, "time.Time pointer values not equal")
}

func TestMapAssignableStruct_timeTimeSlice(t *testing.T) {
	type Src struct {
		Times []time.Time
	}
	type Dst struct {
		Times []time.Time
	}

	now := time.Now()
	src := Src{Times: []time.Time{now, now.Add(time.Hour)}}
	dst := Dst{}

	mapper := NewMapper()
	err := mapper.Map(src, &dst)
	assert.Nil(t, err, "Map returned an error")
	assert.Equal(t, src.Times, dst.Times, "time.Time slices not equal")
}

func TestMapAssignableStruct_timeTimeArray(t *testing.T) {
	type Src struct {
		Times [2]time.Time
	}
	type Dst struct {
		Times [2]time.Time
	}

	now := time.Now()
	src := Src{Times: [2]time.Time{now, now.Add(time.Hour)}}
	dst := Dst{}

	mapper := NewMapper()
	err := mapper.Map(src, &dst)
	assert.Nil(t, err, "Map returned an error")
	assert.Equal(t, src.Times, dst.Times, "time.Time arrays not equal")
}

func TestMapAssignableStruct_timeTimeMap(t *testing.T) {
	type Src struct {
		Times map[string]time.Time
	}
	type Dst struct {
		Times map[string]time.Time
	}

	now := time.Now()
	src := Src{Times: map[string]time.Time{"a": now, "b": now.Add(time.Hour)}}
	dst := Dst{}

	mapper := NewMapper()
	err := mapper.Map(src, &dst)
	assert.Nil(t, err, "Map returned an error")
	assert.Equal(t, src.Times, dst.Times, "time.Time maps not equal")
}

// AI generated code end

// AI generated code start
// --- Deeply Nested Structs ---
func TestMap_DeeplyNestedStructs(t *testing.T) {
	type Inner struct{ Value int }
	type Middle struct{ Inner Inner }
	type Outer struct{ Middle Middle }

	src := Outer{Middle: Middle{Inner: Inner{Value: 42}}}
	var dst Outer

	mapper := NewMapper()
	err := mapper.Map(src, &dst)
	assert.Nil(t, err)
	assert.Equal(t, src, dst)
}

// --- Pointer-to-Value and Value-to-Pointer Mapping ---
func TestMap_PointerToValueAndViceVersa(t *testing.T) {
	type S struct{ V *int }
	type D struct{ V int }
	val := 99
	src := S{V: &val}
	var dst D

	mapper := NewMapper()
	err := mapper.Map(src, &dst)
	assert.Nil(t, err)
	assert.Equal(t, val, dst.V)

	// Value to pointer
	type S2 struct{ V int }
	type D2 struct{ V *int }
	src2 := S2{V: 123}
	var dst2 D2
	err = mapper.Map(src2, &dst2)
	assert.Nil(t, err)
	assert.NotNil(t, dst2.V)
	assert.Equal(t, src2.V, *dst2.V)
}

// --- Nil Handling ---
func TestMap_NilPointerSliceMap(t *testing.T) {
	type S struct {
		P *int
		S []int
		M map[string]int
	}
	type D = S

	src := S{}
	var dst D
	mapper := NewMapper()
	err := mapper.Map(src, &dst)
	assert.Nil(t, err)
	assert.Nil(t, dst.P)
	assert.Nil(t, dst.S)
	assert.Nil(t, dst.M)
}

// --- Unexported Fields ---
func TestMap_UnexportedFieldsIgnored(t *testing.T) {
	type S struct {
		Exported   int
		unexported int
	}
	type D struct {
		Exported   int
		unexported int
	}
	src := S{Exported: 1, unexported: 2}
	var dst D
	mapper := NewMapper()
	err := mapper.Map(src, &dst)
	assert.Nil(t, err)
	assert.Equal(t, 1, dst.Exported)
	// Unexported field should remain zero
	assert.Equal(t, 0, dst.unexported)
}

// --- Partial Overlap ---
func TestMap_PartialFieldOverlap(t *testing.T) {
	type S struct{ A, B int }
	type D struct{ B, C int }
	src := S{A: 1, B: 2}
	var dst D
	mapper := NewMapper()
	err := mapper.Map(src, &dst)
	assert.Nil(t, err)
	assert.Equal(t, 2, dst.B)
	assert.Equal(t, 0, dst.C)
}

// --- Custom Setter/Getter with Different Types ---
type DstWithSetter struct {
	val int
}

func (d *DstWithSetter) SetVal(v string) { d.val = len(v) }

type SrcWithGetter struct{}

func (s SrcWithGetter) GetVal() int { return 42 }

func TestMap_SetterGetterTypeMismatch(t *testing.T) {
	src := SrcWithGetter{}
	dst := DstWithSetter{}
	mapper := NewMapper()
	err := mapper.Map(src, &dst)
	assert.Equal(t, ErrMismatchType, err)
}

// --- Map Key/Value Conversion ---
func TestMap_MapKeyValueTypeConversion(t *testing.T) {
	type S struct{ M map[string]int }
	type D struct{ M map[string]string }
	mapper := NewMapper()
	AddTypeConverter[int, string](mapper, func(src any) (any, error) {
		return fmt.Sprintf("%d", src.(int)), nil
	})
	src := S{M: map[string]int{"a": 1, "b": 2}}
	var dst D
	err := mapper.Map(src, &dst)
	assert.Nil(t, err)
	assert.Equal(t, map[string]string{"a": "1", "b": "2"}, dst.M)
}

// --- Empty/Zero Value Handling ---
func TestMap_ZeroValues(t *testing.T) {
	type S struct {
		A int
		B []int
	}
	type D struct {
		A int
		B []int
	}
	src := S{}
	var dst D
	mapper := NewMapper()
	err := mapper.Map(src, &dst)
	assert.Nil(t, err)
	assert.Equal(t, 0, dst.A)
	assert.Nil(t, dst.B)
}

// --- Unsupported Types ---
func TestMap_UnsupportedTypes(t *testing.T) {
	type S struct {
		F func()
		C chan int
		U unsafe.Pointer
	}
	type D S
	src := S{}
	var dst D
	mapper := NewMapper()
	err := mapper.Map(src, &dst)
	assert.Nil(t, err) // Should ignore unsupported types, not error
}

// --- Multiple Type Converters ---
func TestMap_MultipleTypeConverters(t *testing.T) {
	type S struct{ V int }
	type D struct{ V string }
	mapper := NewMapper()
	AddTypeConverter[int, string](mapper, func(src any) (any, error) {
		return fmt.Sprintf("int:%d", src.(int)), nil
	})
	AddTypeConverter[int, string](mapper, func(src any) (any, error) {
		return fmt.Sprintf("override:%d", src.(int)), nil
	})
	src := S{V: 5}
	var dst D
	err := mapper.Map(src, &dst)
	assert.Nil(t, err)
	assert.Equal(t, "override:5", dst.V)
}

// --- AddTypeConverter Edge Cases ---
func TestMap_AddTypeConverterEdgeCases(t *testing.T) {
	type S struct{ V int }
	type D struct{ V string }
	mapper := NewMapper()
	// Add a converter for types that are not assignable
	AddTypeConverter[float64, string](mapper, func(src any) (any, error) {
		return "should not be used", nil
	})
	src := S{V: 1}
	var dst D
	err := mapper.Map(src, &dst)
	assert.Equal(t, ErrMismatchType, err)
	// Overwrite existing converter
	AddTypeConverter[int, string](mapper, func(src any) (any, error) {
		return "new", nil
	})
	src = S{V: 2}
	dst = D{}
	err = mapper.Map(src, &dst)
	assert.Nil(t, err)
	assert.Equal(t, "new", dst.V)
}

// AI generated code end
