package d2i18n

import (
	"io"
)

// I18n is the interface for accessing i18n texts for Dofus 2
type I18n interface {
	GetUndiacriticalText(int32) (string, bool, error)
	GetText(int32) (string, bool, error)
	GetNamedText(string) (string, bool, error)
}

type i18n struct {
	r Reader

	indexes              map[int32]int32
	undiacriticalIndexes map[int32]int32
	textIndexes          map[string]int32
	textSortIndex        map[int32]int
}

// Parse parses a I18n files.
// Since it uses Seek() to retrieve the position, ReadSeekers like bytes.Buffer might not work
func Parse(r io.ReadSeeker) (I18n, error) {
	reader := NewReader(r)
	i18n := &i18n{
		reader,
		map[int32]int32{}, map[int32]int32{},
		map[string]int32{}, map[int32]int{},
	}
	if err := i18n.parseIndexes(); err != nil {
		return nil, err
	}
	return i18n, nil
}

func (i *i18n) parseIndexes() error {
	tablePosition, err := i.r.ReadInt32()
	if err != nil {
		return err
	}
	if err := i.r.Goto(int64(tablePosition)); err != nil {
		return err
	}
	if err := i.parseNumIndexes(); err != nil {
		return err
	}
	if err := i.parseTextIndexes(); err != nil {
		return err
	}
	if err := i.parseSortIndexes(); err != nil {
		return err
	}
	return nil
}

func (i *i18n) parseNumIndexes() error {
	numIndexesLength, err := i.r.ReadInt32()
	if err != nil {
		return err
	}
	for position := int64(0); position < int64(numIndexesLength); position += 9 {
		index, err := i.r.ReadInt32()
		if err != nil {
			return err
		}
		hasUndiacritical, err := i.r.ReadBoolean()
		if err != nil {
			return err
		}
		value, err := i.r.ReadInt32()
		if err != nil {
			return err
		}
		i.indexes[index] = value
		if hasUndiacritical {
			undiacriticalValue, err := i.r.ReadInt32()
			if err != nil {
				return err
			}
			position += 4
			i.undiacriticalIndexes[index] = undiacriticalValue
		} else {
			i.undiacriticalIndexes[index] = value
		}

	}
	return nil
}

func (i *i18n) parseTextIndexes() error {
	length, err := i.r.ReadInt32()
	if err != nil {
		return err
	}
	for length > 0 {
		begin, err := i.r.Position()
		if err != nil {
			return err
		}

		index, err := i.r.ReadString()
		if err != nil {
			return err
		}
		value, err := i.r.ReadInt32()
		if err != nil {
			return err
		}
		i.textIndexes[index] = value

		end, err := i.r.Position()
		if err != nil {
			return err
		}
		length -= int32(end - begin)
	}
	return nil
}

func (i *i18n) parseSortIndexes() error {
	length, err := i.r.ReadInt32()
	if err != nil {
		return err
	}
	count := 1
	for length > 0 {
		begin, err := i.r.Position()
		if err != nil {
			return err
		}

		index, err := i.r.ReadInt32()
		if err != nil {
			return err
		}

		i.textSortIndex[index] = count
		count++

		end, err := i.r.Position()
		if err != nil {
			return err
		}
		length -= int32(end - begin)
	}
	return nil
}

func (i *i18n) GetUndiacriticalText(index int32) (string, bool, error) {
	offset, found := i.undiacriticalIndexes[index]
	if !found {
		return "", false, nil
	}
	if err := i.r.Goto(int64(offset)); err != nil {
		return "", true, err
	}
	str, err := i.r.ReadString()
	return str, false, err
}

func (i *i18n) GetText(index int32) (string, bool, error) {
	offset, found := i.indexes[index]
	if !found {
		return "", false, nil
	}
	if err := i.r.Goto(int64(offset)); err != nil {
		return "", true, err
	}
	str, err := i.r.ReadString()
	return str, false, err
}

func (i *i18n) GetNamedText(name string) (string, bool, error) {
	offset, found := i.textIndexes[name]
	if !found {
		return "", false, nil
	}
	if err := i.r.Goto(int64(offset)); err != nil {
		return "", true, err
	}
	str, err := i.r.ReadString()
	return str, true, err
}
