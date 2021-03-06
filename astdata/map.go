package astdata

import (
	"fmt"
	"go/ast"
)

// MapType is the map in the go code
type MapType struct {
	embededData

	key   Definition
	value Definition
}

func (m *MapType) String() string {
	return fmt.Sprintf("map[%s]%s", m.key.String(), m.value.String())
}

// Key type definition
func (m *MapType) Key() Definition {
	return m.key
}

// Val is the definition of the value type
func (m *MapType) Val() Definition {
	return m.value
}

// Compare try to compare this to def
func (m *MapType) Compare(def Definition) bool {
	return m.String() == def.String()
}

func getMap(p *Package, f *File, t *ast.MapType) Definition {
	return &MapType{
		embededData: embededData{
			pkg:  p,
			fl:   f,
			node: t,
		},
		key:   newType(p, f, t.Key),
		value: newType(p, f, t.Value),
	}
}
