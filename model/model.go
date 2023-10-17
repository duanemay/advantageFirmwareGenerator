package model

import "fmt"

type LayerName string

const (
	Base   LayerName = "base"
	Keypad           = "keypad"
	Fn               = "fn"
	Mod              = "mod"
)

var LayerOrder = []LayerName{Base, Keypad, Fn, Mod}

type LayerNameArray []LayerName

func (a LayerNameArray) Len() int {
	return len(a)
}
func (a LayerNameArray) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}
func (a LayerNameArray) Less(i, j int) bool {
	return a[i].ToInt() < a[j].ToInt()
}

func (n LayerName) ToInt() int {
	switch n {
	case Base:
		return 0
	case Keypad:
		return 1
	case Fn:
		return 2
	case Mod:
		return 3
	}
	return -1
}

type KeyMapFile struct {
	Keyboard   string      `json:"keyboard"`
	Keymap     string      `json:"keymap"`
	Layout     string      `json:"layout"`
	LayerNames []LayerName `json:"layer_names"`
	Layers     [][]string  `json:"layers"`
}

func NewKeyMapFile() KeyMapFile {
	return KeyMapFile{
		Keyboard:   "adv360",
		Keymap:     "default",
		Layout:     "LAYOUT",
		LayerNames: LayerOrder,
	}
}

type LayerStrings struct {
	keyStr []string
}

type KeysFile struct {
	KeyIds KeyIdArray `json:"keyIds"`
	Layers Layers     `json:"layers"`
}

type KeyGroup string

const (
	Left       KeyGroup = "left"
	Right               = "right"
	LeftThumb           = "leftThumb"
	RightThumb          = "rightThumb"
)

type KeyId struct {
	Comment string   `json:"comment"`
	KeyId   string   `json:"id"`
	Group   KeyGroup `json:"group"`
	Row     int      `json:"row"`
	Column  int      `json:"column"`
}

type KeyIdArray []KeyId

func (id KeyId) String() string {
	if id.Comment != "" {
		return fmt.Sprintf("{%s}", id.Comment)
	}
	return fmt.Sprintf("{%s, %s, %d, %d}", id.KeyId, id.Group, id.Row, id.Column)
}

func (a KeyIdArray) Len() int {
	return len(a)
}
func (a KeyIdArray) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}
func (a KeyIdArray) Less(i, j int) bool {
	iRow, iColumn := a[i].translateCoords()
	jRow, jColumn := a[j].translateCoords()

	if iRow != jRow {
		return iRow < jRow
	}
	return iColumn < jColumn
}

func (id KeyId) translateCoords() (int, int) {
	row, col := id.Row, id.Column
	switch id.Group {
	case Right:
		col += 300

	case LeftThumb:
		row += 2
		col += 100

	case RightThumb:
		row += 2
		col += 200
	}
	if id.Group == "right" {
		col += 20
	}
	return row, col
}

type Layers struct {
	Base   Layer `json:"base"`
	Keypad Layer `json:"keypad"`
	Fn     Layer `json:"fn"`
	Mod    Layer `json:"mod"`
}

func (l Layers) GetLayer(name LayerName) (Layer, error) {
	switch name {
	case Base:
		return l.Base, nil
	case Keypad:
		return l.Keypad, nil
	case Fn:
		return l.Fn, nil
	case Mod:
		return l.Mod, nil
	}
	return Layer{}, fmt.Errorf("unknown layer name: %s", name)
}

type Layer struct {
	Keys []Key `json:"keys"`
}

type Key struct {
	Id     string `json:"id"`
	Action string `json:"action"`
	Value  string `json:"value"`
}
