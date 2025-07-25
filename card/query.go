package card

func (c *Card) GetOracleText() []string {
	var o []string
	if c.OracleText != nil {
		o = append(o, *c.OracleText)
	}
	for _, face := range c.CardFaces {
		if face.OracleText != nil {
			o = append(o, *face.OracleText)
		}
	}
	return o
}
