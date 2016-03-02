package mio

import (
	"errors"
	"os"
	"sort"
	"syscall"
	"unsafe"
)

var (
	MAX_KLEN = 0xff         // max key size 255 bytes
	MAX_HLEN = 0x04         // max header size 4 bytes
	NIL_BLCK = []byte{0x00} // empty block header byte
	MIN_MMAP = 0xffffff + 1 // 16MB
	SYS_PAGE = int64(syscall.Getpagesize())
	NIL_PAGE = make([]byte, SYS_PAGE, SYS_PAGE)

	KeyNilErr error = errors.New("Key was nil, or corrupt")
	DocNilErr error = errors.New("Doc was nil, or corrupt")
	KeyLenErr error = errors.New("Key exceeded MAX_KLEN")
	DocLenErr error = errors.New("Doc exceeded SYS_PAGE")
)

type mapper struct {
	path string
	file *os.File
	size int64
	docs int64
	data []byte
	gaps []int
}

func Map(path string) *mapper {
	fd, err := os.OpenFile(path,
		syscall.O_RDWR|syscall.O_CREAT|syscall.O_APPEND, 0644)
	if err != nil {
		panic(err)
	}
	fi, err := fd.Stat()
	if err != nil {
		panic(err)
	}
	size := fi.Size()
	if size == 0 {
		size = resizeFile(int(fd.Fd()), toPageSize(int64(MIN_MMAP)))
	}
	b, err := syscall.Mmap(int(fd.Fd()), 0, int(size),
		syscall.PROT_READ|syscall.PROT_WRITE, syscall.MAP_SHARED)
	if err != nil {
		panic(err)
	}
	m := &mapper{
		path: path,
		file: fd,
		size: size,
		data: b,
		gaps: make([]int, 0),
	}
	//m.initialize()
	//m.growAndRemap()
	//fmt.Printf("docs: %d, gaps: %d (%v)\n", m.docs, len(m.gaps), m.gaps)
	return m
}

// initialize the mapper; fill out the document count
// as well as the gap list when the mapper is started
func (m *mapper) initialize() {
	if m.size > 0 {
		lst := int(m.size - SYS_PAGE)
		for pos := 0; pos <= lst; pos += int(SYS_PAGE) {
			if m.data[lst] == NIL_BLCK[0] {
				lst -= int(SYS_PAGE)
			}
			if m.data[pos] != NIL_BLCK[0] {
				m.docs++
			}
		}
		m.gaps = append(m.gaps, 0, lst)
	}
}

// checks the gap list for potential gaps it can
// reuse. if a match is found, it is removed from
// the gap list and returned (with bool) for use
//
// NOTE:MAKE THIS ONLY RETURN AN INT
func (m *mapper) checkGaps(n int) ([]int, bool) {
	if len(m.gaps) > n {
		u := make([]int, n, n)
		for i := 0; i < len(m.gaps); i++ {
			if m.gaps[i+n-1]-m.gaps[i]+1 == n {
				copy(u, m.gaps[i:i+n])
				m.gaps = append(m.gaps[:i], m.gaps[i+n:]...)
				return u, true
			}
		}
	}
	return nil, false
}

func (m *mapper) find(pages int) (int, bool) {
	// check to see if there are consecutive gaps
	// that we can use, if not, move to next check
	if len(m.gaps) >= pages {
		for i := 0; i < len(m.gaps); i++ {
			if m.gaps[i+pages-1]-m.gaps[i]+1 == pages {
				m.gaps = append(m.gaps[:i], m.gaps[i+pages:]...)
				return i, true
			}
		}
	}
	// check to see if there are ANY gaps at the end that
	// we can use as in index, it not, move to next check
	if len(m.gaps) > 0 {
		offset := m.gaps[len(m.gaps)-1]
		if int64(m.gaps[len(m.gaps)-1]+pages)*SYS_PAGE < m.size {
			m.gaps = append(m.gaps[:len(m.gaps)-1], m.gaps[len(m.gaps):]...)
			return offset, true
		}
		return offset, false
	}
	// do linear search for next index, and see if we need to grow
	// ?? BRAIN NO WORK...

	return 0, false
}

/*
	pos, ok := m.find(pagesize)
	if !ok {
		m.grow()
	}
	write(data, pos)
*/

func (m *mapper) Find(offset, pagesize int) int {
	return -1 // TODO: STUFF BELOW

	/*
		// add case
		if offset == -1 {
			n := checkGaps(pagesize)
			if n == -1 {
				// find last index
				// NOTE: add m.pages
				// ie. m.pages * SYS_PAGE + 1
				// if
			}
		}
		// set case
		pos := offset * SYS_PAGE
		pgs := m.data[pos+1]
		// zero out pgs worth of data
		if pagesize < pgs {
			pgcnt := pgs - pagesize
			// add pgcnt to m.gaps & sort
			return offset
		}
		if pagesize == pgs {
			return offset
		}
		if pagesize > pgs {
			// add pages to gaps list and search
			// gaps list for a fit. if gaps list
			// returns -1, find last index linerally
		}
		//
	*/
}

// writes a key value pair at the block offset provided. It
// will increment the doc count, as well as grow the file and
// remap if nessicary. If the provided position happens to be
// in the current gap list it will also remove that entry.
func (m *mapper) _Set(key, val []byte, offset int) int {
	/*
		block, err := Block(key, val)
		if err != nil {
			return err
		}
	*/
	doc := Doc(key, val)
	if doc == nil {
		return -1
	}
	pos := int64(offset) * SYS_PAGE
	if m.data[pos] == NIL_BLCK[0] {
		// we are not not updating in this case
		if !sort.IntsAreSorted(m.gaps) {
			sort.Ints(m.gaps)
		}
		// if pos happens to be in gap list...
		x := sort.SearchInts(m.gaps, int(pos))
		if x != len(m.gaps) {
			// ...then remove it from gap list.
			copy(m.gaps[x:], m.gaps[x+1:])
			m.gaps[len(m.gaps)-1] = 0
			m.gaps = m.gaps[:len(m.gaps)-1]
		}
		// not updating a doc, so increment count
		m.docs++
		// we may be out of room soon so
		m.growAndRemap()
	}
	// write data to provided position
	pgs := doc.pages()
	copy(m.data[pos:(pos+int64(pgs))+SYS_PAGE], doc)
	return pgs
}

// writes a key value pair. It will increment the doc count,
// as well as grow the file and remap if nessicary. It will
// also check to see if the record can be placed in the gap
// list (if there are any entries in the gap list)
func (m *mapper) _Add(key, val []byte) error {
	block, err := Block(key, val)
	if err != nil {
		return err
	}
	pos := m.docs * SYS_PAGE
	// check gap list to see if we can use an empty block
	if len(m.gaps) > 0 {
		// gap list contains items, sort if nessicary
		if !sort.IntsAreSorted(m.gaps) {
			sort.Ints(m.gaps)
		}
		// shift out a position space closest to the front
		pos = int64(m.gaps[0]) * SYS_PAGE
		m.gaps = m.gaps[1:]
	} else {
		m.growAndRemap() // empty gap list, try to grow and remap
	}
	// write data to provided position
	copy(m.data[pos:pos+SYS_PAGE], block) // NOTE:OUT OF BOUNDS ON INIT TEST
	m.docs++
	return nil
}

// return a document based on its offset key
func (m *mapper) Get(offset int) []byte {
	// get position based on offset
	pos := int64(offset) * SYS_PAGE
	// get how many pages document is using
	pgs := int(m.data[pos] + 1)
	// get byte count we need to write based on pgs
	siz := int64(pgs) * SYS_PAGE
	return m.data[pos : pos+siz]
}

// delete a document and add to gap list
func (m *mapper) Del(offset int) int {
	// get position based on offset
	pos := int64(offset) * SYS_PAGE
	// get how many pages document is using
	pgs := int(m.data[pos] + 1)
	// get byte count we need to write based on pgs
	siz := int64(pgs) * SYS_PAGE
	// write zero page data across document size of siz
	copy(m.data[pos:pos+siz], make([]byte, siz, siz))
	// add deleted page offsets back to gap list
	for i := 0; i < pgs; i++ {
		m.gaps = append(m.gaps, offset+i)
	}
	// sort gap list now that we added to it
	if !sort.IntsAreSorted(m.gaps) {
		sort.Ints(m.gaps)
	}
	return pgs
}

// flush data to disk in an async fashion
func (m *mapper) Flush() {
	_, _, err := syscall.Syscall(syscall.SYS_MSYNC,
		uintptr(unsafe.Pointer(&m.data[0])), uintptr(len(m.data)),
		uintptr(syscall.MS_ASYNC))
	if err != 0 {
		panic(err)
	}
}

// unmap the current mapping
func (m *mapper) Unmap() {
	err := syscall.Munmap(m.data)
	if err != nil {
		panic(err)
	}
	m.data = nil
}

func (m *mapper) growAndRemap() {
	if m.docs < m.size/SYS_PAGE {
		return
	}
	m.Unmap()
	m.size = resizeFile(int(m.file.Fd()), toPageSize(m.size+int64(MIN_MMAP)))
	b, err := syscall.Mmap(int(m.file.Fd()), 0, int(m.size),
		syscall.PROT_READ|syscall.PROT_WRITE, syscall.MAP_SHARED)
	if err != nil {
		panic(err)
	}
	m.data = b
}

func Block(key, val []byte) ([]byte, error) {
	if key == nil {
		return nil, KeyNilErr
	}
	if len(key) > MAX_KLEN {
		return nil, KeyLenErr
	}
	if val == nil {
		return nil, DocNilErr
	}
	if int64(len(key)+len(val)+MAX_HLEN) > SYS_PAGE {
		return nil, DocLenErr
	}
	block := make([]byte, SYS_PAGE, SYS_PAGE)
	block[0] = 0x01
	block[1] = byte(len(key))
	block[2] = byte(len(val))
	block[3] = byte(len(val) >> 8)
	copy(block[MAX_HLEN:MAX_HLEN+len(key)], key)
	copy(block[MAX_HLEN+len(key):MAX_HLEN+len(key)+len(val)], val)
	return block, nil
}

func toPageSize(size int64) int64 {
	if size > 0 {
		return (size + SYS_PAGE) &^ (SYS_PAGE)
	}
	return SYS_PAGE
}

func resizeFile(fd int, size int64) int64 {
	err := syscall.Ftruncate(fd, size)
	if err != nil {
		panic(err)
	}
	return size
}
