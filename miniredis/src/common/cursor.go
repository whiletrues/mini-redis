package common

type Cursor struct {
	buffer []byte
	Index  int
}

func (c *Cursor) Next() {
	c.Index++
}

func (c *Cursor) Previous() {
	c.Index--
}

func (c *Cursor) Current() (int, byte) {
	return c.Index, c.buffer[c.Index]
}

func (c *Cursor) hasNext() bool {
	return c.Index < len(c.buffer)
}

func (c *Cursor) nextLine() []byte {

	start := c.Index

	for c.hasNext() {
		index, char := c.Current()

		if char == '\r' {
			c.Next()
			c.Next()
			return c.buffer[start:index]
		}
		c.Next()
	}

	return c.buffer[start:c.Index]
}
