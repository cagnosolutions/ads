package mio

import (
	"fmt"
	"log"
	"os"
	"sort"
	"syscall"
	"unsafe"
)

var (
	NIL_PAGE = []byte{0x00}
	MIN_MMAP = 0xffffff + 1 // 16MB
	SYS_PAGE = int64(syscall.Getpagesize())
)

func (m *mapper) String() string {
	str := "mapper info\n" +
		"===========\n" +
		"path: " + m.path + "\n" +
		"file: " + m.file.Name() + "\n" +
		"size: " + fmt.Sprintf("%d\n", m.size) +
		"pags: " + fmt.Sprintf("%d\n", m.pags) +
		"data: " + fmt.Sprintf("%d\n", len(m.data)) +
		"gaps: " + fmt.Sprintf("%v\n", m.gaps)
	return str
}

// mapper struct
type mapper struct {
	path string
	file *os.File
	size int64
	pags int64
	data []byte
	gaps []int
}

func (m *mapper) Pages() int {
	return int(m.pags)
}

func (m *mapper) Gaps() int {
	return len(m.gaps)
}

// open a backing file at the provided path, and
// grow the file to minimum mmap size if it is a
// new file. attempt to memory map backing file
// and then initialize and return a mapper.
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
	var fresh bool
	if size == 0 {
		size = resize(fd.Fd(), align(int64(MIN_MMAP)))
		fresh = true
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
		pags: 0,
		data: b,
		gaps: make([]int, 0),
	}
	m.initialize(fresh)
	return m
}

// initialize the mapper; fill out the document count
// as well as the gap list when the mapper is started
func (m *mapper) initialize(fresh bool) {
	// if the file was just created we don't
	// need to load anything into the mapper
	if !fresh { // otherwise...
		// file is not fresh; it contains data
		// that we need to add to the mapper...
		lst, pos := int(m.size-SYS_PAGE), 0
		for m.data[lst] == NIL_PAGE[0] {
			lst--
		}
		for pos <= lst {
			if m.data[pos] != NIL_PAGE[0] {
				pgs := m.data[pos+1]
				m.pags += int64(pgs)
				pos += int(pgs) * int(SYS_PAGE)
			} else {
				m.gaps = append(m.gaps, pos/int(SYS_PAGE))
				pos += int(SYS_PAGE)
			}
		}
	}
}

// writes a binary document at the block offset provided. It
// will increment the doc count, as well as grow the file and
// remap if nessicary. If the provided position happens to be
// in the current gap list it will also remove that entry.
func (m *mapper) Set(dat []byte, offset int) bool {
	// create new document; sanitize and prepend header based
	// on the raw dat supplied. if !ok, print log and return
	doc, ok := document(dat)
	if !ok {
		log.Println("data exceeded maximum page limit")
		return false
	}
	// read header to get docs page count
	pgs := int(doc[1])
	// calculate proper position to set based on
	// doc's page count and provided offset
	pos, pct := m.positionToSet(offset, pgs)
	// write data to provided position
	copy(m.data[pos:pos+(pgs*int(SYS_PAGE))], doc)
	m.pags += int64(pct) // page count to incr or decr
	return true
}

// return a usable offset based on the offset
// provided and the number of pages the data
// requires in order to perform a set / update.
func (m *mapper) positionToSet(offset, pages int) (int, int) {
	// get the position of the offset in bytes
	// and check header to see if need to wipe
	pos, pgs := offset*int(SYS_PAGE), 0
	if m.data[pos] == 0x01 && m.data[pos+1] != 0x00 {
		// wipe all pages of the document located at
		// the offset provided so we are left with a
		// clean slate. return the newly wiped pages.
		pgs = m.del(offset)
	}
	// if the number of pages provided is smaller
	// than then freshly wiped pages...
	if pages < pgs {
		// add difference back into gap list
		m.addOffset(offset+pages, pgs-pages)
		// re-use same position in offset, it fits
		return offset * int(SYS_PAGE), pages - pgs
	}
	// if the number of pages provided is exactly the
	// same number as the freshly wiped pages...
	if pages == pgs {
		// re-use same position in offset, it fits
		return offset * int(SYS_PAGE), 0
	}
	// ...otherwise, the number of pages provided is
	// greater than the freshly wiped pages; so add
	// freshly wiped pages into gaps list and search
	// gaps list for a different offset position that
	// fits...
	m.addOffset(offset, pgs)
	offset, ok := m.getOffset(pages)
	// if no position fits, grow file and return last page
	if !ok {
		m.grow()
	}
	return offset * int(SYS_PAGE), pages - pgs
}

// writes a blob of data. It will increment the page count,
// as well as grow the file and remap if nessicary. It will
// also check to see if the record can be placed in the gap
// list (if there are any entries in the gap list)
func (m *mapper) Add(dat []byte) bool {
	// create new document; sanitize and prepend header based
	// on the raw dat supplied. if !ok, print log and return
	doc, ok := document(dat)
	if !ok {
		log.Println("data exceeded maximum page limit")
		return false
	}
	pgs := int(doc[1])
	// calculate proper position to add based on doc's page count
	pos := m.positionToAdd(pgs)
	// write data to provided position
	copy(m.data[pos:pos+(pgs*int(SYS_PAGE))], doc)
	m.pags += int64(pgs)
	return true
}

// calculate proper position for document to be added
// based on documents page count; grow file if needed
func (m *mapper) positionToAdd(pages int) int {
	// return a usable offset based on the
	// number of pages the data requires
	offset, ok := m.getOffset(pages)
	// grow file if there isn't enough room
	if !ok {
		m.grow()
	}
	// return position of offset in mapping
	return offset * int(SYS_PAGE)
}

// return a document based on its offset key
func (m *mapper) Get(offset int) ([]byte, bool) {
	// get position based on offset
	pos := int64(offset) * SYS_PAGE
	// check for valid header before returning
	if m.data[pos] == 0x00 || m.data[pos+1] == 0x00 {
		return nil, false
	}
	// get how many pages document is using
	pgs := int(m.data[pos+1])
	// get byte count we need to read based on pgs
	siz := int64(pgs) * SYS_PAGE
	// return document slice
	return m.data[pos : pos+siz], true
}

// delete a document and add re-claimed pages to the
// to m.gaps list. decrement m.pags (ie. pages in use)
func (m *mapper) Del(offset int) bool {
	pages := m.del(offset)
	m.addOffset(offset, pages)
	m.pags -= int64(pages)
	return pages > 0
}

// delete, ie. write nil bytes, starting at offset
// across document size (page count from header) and
// return then number of pages that were deleted
func (m *mapper) del(offset int) int {
	// get position based on offset
	pos := int64(offset) * SYS_PAGE
	// get how many pages document is using
	pgs := int(m.data[pos+1])
	// get byte count we need to write based on pgs
	siz := int64(pgs) * SYS_PAGE
	// write zero page data across document size of siz
	copy(m.data[pos:pos+siz], make([]byte, siz, siz))
	return pgs
}

// adds number of sequential pages starting from the
// offset to the gaps list, then ensures gaps are sorted
func (m *mapper) addOffset(offset, pages int) {
	// add pages to gap list
	for i := 0; i < pages; i++ {
		m.gaps = append(m.gaps, offset+i)
	}
	// sort gap list now that we added to it
	if !sort.IntsAreSorted(m.gaps) {
		sort.Ints(m.gaps)
	}
}

// checks the gap list for potential gaps it can
// reuse. if a match is found, it is removed from
// the gap list and returned (with bool) for use
func (m *mapper) getOffset(pages int) (int, bool) {
	if len(m.gaps) > 0 {
		if pages == 1 {
			i := m.gaps[0]
			m.gaps = m.gaps[1:]
			return i, true
		}
		// check to see if there are consecutive gaps
		// that we can use, if not, get last offset...
		if len(m.gaps) >= pages {
			for i := 0; i < len(m.gaps); i++ {
				if (m.gaps[i+pages-1]-m.gaps[i])+1 == pages {
					m.gaps = append(m.gaps[:i], m.gaps[i+pages:]...)
					return i, true
				}
			}
		}
	}
	// ...get last page offset, then from the offset
	// compare the supplied page count we need to ensure
	// that it will 1) fit in current file, or 2) signal
	// that we have to grow the file in order to add it
	lst := int(m.pags) + len(m.gaps)
	pos := int64(lst+pages) * SYS_PAGE
	return lst, pos < m.size
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

// grow the underlying file and re-map if nessicary
func (m *mapper) grow() {
	// NOTE: CAN POSSIBLY REMOVE THIS CHECK...
	if m.pags+int64(len(m.gaps)) < (m.size/SYS_PAGE)-1 {
		return
	}
	m.Unmap()
	m.size = resize(m.file.Fd(), align(m.size+int64(MIN_MMAP)))
	b, err := syscall.Mmap(int(m.file.Fd()), 0, int(m.size),
		syscall.PROT_READ|syscall.PROT_WRITE, syscall.MAP_SHARED)
	if err != nil {
		panic(err)
	}
	m.data = b
}

// create a new document from supplied data;
// sanitize/prepend 2 byte header (opt, pgs)
func document(data []byte) ([]byte, bool) {
	size := len(data)
	if size < 1 || size+2 > 0xff*int(SYS_PAGE) {
		return nil, false
	}
	pgs := (size + 2 + int(SYS_PAGE) - 1) &^ (int(SYS_PAGE) - 1) / int(SYS_PAGE)
	doc := make([]byte, pgs*int(SYS_PAGE), pgs*int(SYS_PAGE))
	doc[0], doc[1] = 0x01, byte(pgs)
	copy(doc[2:], data)
	return doc, true
}

// verify that the data provided is contains
// a proper header; return data and bool
func verify(data []byte) ([]byte, bool) {
	if data != nil && len(data) > 1 && data[0] != 0x00 && data[1] <= 0xff {
		return data, true
	}
	return nil, false
}

// align a given size to the nearest
// system page... always rounding up
func align(size int64) int64 {
	if size > 0 {
		return (size + SYS_PAGE) &^ (SYS_PAGE)
	}
	return SYS_PAGE
}

// resize/truncate a file based on the
// supplied file pointer to size, size
func resize(fd uintptr, size int64) int64 {
	err := syscall.Ftruncate(int(fd), size)
	if err != nil {
		panic(err)
	}
	return size
}
