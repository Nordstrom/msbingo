package nbfx

type codec struct {
	dict        map[uint32]string
	reverseDict map[string]uint32
}

type Encoder interface {
	Encode(xml string) ([]byte, error)
}

type Decoder interface {
	Decode(bin []byte) (string, error)
}

func (c *codec) addDictionaryString(index uint32, value string) {
	if _, ok := c.dict[index]; ok {
		return
	}
	c.dict[index] = value
	c.reverseDict[value] = index
}

// for MultiByteInt31
const MASK_MBI31 = uint32(0x80) //0x80 = 128
