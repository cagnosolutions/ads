package main

import "github.com/cagnosolutions/ads/mio"

func main() {
	m := mio.Map("m.db")
	m.Set([]byte(`0`), 0)
	m.Set([]byte(`1`), 1)
	m.Set([]byte(`2`), 2)
	m.Set([]byte(`3`), 3)
	m.Set([]byte(`4`), 4)
	m.Set([]byte(`5`), 5)
	m.Set([]byte(`6`), 6)
	m.Set([]byte(`7`), 7)
	m.Set([]byte(`8`), 8)
	m.Set([]byte(`9`), 9)
	m.Unmap()
}
