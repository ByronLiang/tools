// Package tiff implements TIFF decoding as defined in TIFF 6.0 specification at
// http://partners.adobe.com/public/developer/en/tiff/TIFF6.pdf
package tiff

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
)

const ReadTagIdOffsetLen = 10

// ReadAtReader is used when decoding Tiff tags and directories
type ReadAtReader interface {
	io.Reader
	io.ReaderAt
	io.Seeker
}

// Tiff provides access to a decoded tiff data structure.
type Tiff struct {
	// Dirs is an ordered slice of the tiff's Image File Directories (IFDs).
	// The IFD at index 0 is IFD0.
	Dirs []*Dir
	// The tiff's byte-encoding (i.e. big/little endian).
	Order binary.ByteOrder
}

// Decode parses tiff-encoded data from r and returns a Tiff struct that
// reflects the structure and content of the tiff data. The first read from r
// should be the first byte of the tiff-encoded data and not necessarily the
// first byte of an os.File object.
func Decode(r io.Reader) (*Tiff, error) {
	data, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, errors.New("tiff: could not read data")
	}
	buf := bytes.NewReader(data)

	t := new(Tiff)

	// read byte order
	bo := make([]byte, 2)
	if _, err = io.ReadFull(buf, bo); err != nil {
		return nil, errors.New("tiff: could not read tiff byte order")
	}
	if string(bo) == "II" {
		t.Order = binary.LittleEndian
	} else if string(bo) == "MM" {
		t.Order = binary.BigEndian
	} else {
		return nil, errors.New("tiff: could not read tiff byte order")
	}

	// check for special tiff marker
	var sp int16
	err = binary.Read(buf, t.Order, &sp)
	if err != nil || 42 != sp {
		return nil, errors.New("tiff: could not find special tiff marker")
	}

	// load offset to first IFD
	var offset int32
	err = binary.Read(buf, t.Order, &offset)
	if err != nil {
		return nil, errors.New("tiff: could not read offset to first IFD")
	}
	// load IFD's
	var d *Dir
	prev := offset
	for offset != 0 {
		// seek to offset
		_, err := buf.Seek(int64(offset), 0)
		if err != nil {
			return nil, errors.New("tiff: seek to IFD failed")
		}

		if buf.Len() == 0 {
			return nil, errors.New("tiff: seek offset after EOF")
		}

		// load the dir
		d, offset, err = DecodeDir(buf, t.Order)
		if err != nil {
			return nil, err
		}

		if offset == prev {
			return nil, errors.New("tiff: recursive IFD")
		}
		prev = offset

		t.Dirs = append(t.Dirs, d)
	}
	// seek to offset

	return t, nil
}

func DecodeDefineField(r io.Reader, fieldId uint16) (*Tag, error) {
	data, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, errors.New("tiff: could not read data")
	}
	buf := bytes.NewReader(data)

	var order binary.ByteOrder
	// read byte order
	bo := make([]byte, 2)
	if _, err = io.ReadFull(buf, bo); err != nil {
		return nil, errors.New("tiff: could not read tiff byte order")
	}
	if string(bo) == "II" {
		order = binary.LittleEndian
	} else if string(bo) == "MM" {
		order = binary.BigEndian
	} else {
		return nil, errors.New("tiff: could not read tiff byte order")
	}

	// check for special tiff marker
	var sp int16
	err = binary.Read(buf, order, &sp)
	if err != nil || 42 != sp {
		return nil, errors.New("tiff: could not find special tiff marker")
	}

	// load offset to first IFD
	var offset int32
	err = binary.Read(buf, order, &offset)
	if err != nil {
		return nil, errors.New("tiff: could not read offset to first IFD")
	}
	// load IFD's
	var tag *Tag
	prev := offset
	for offset != 0 {
		// seek to offset
		_, err := buf.Seek(int64(offset), 0)
		if err != nil {
			return nil, errors.New("tiff: seek to IFD failed")
		}

		if buf.Len() == 0 {
			return nil, errors.New("tiff: seek offset after EOF")
		}

		// load the dir
		tag, offset, err = DecodeDefineTag(buf, order, fieldId)
		if err != nil {
			return nil, err
		}
		if tag != nil {
			return tag, nil
		}
		if offset == prev {
			return nil, errors.New("tiff: recursive IFD")
		}
		prev = offset
	}
	// seek to offset
	return nil, nil
}

func (tf *Tiff) String() string {
	var buf bytes.Buffer
	fmt.Fprint(&buf, "Tiff{")
	for _, d := range tf.Dirs {
		fmt.Fprintf(&buf, "%s, ", d.String())
	}
	fmt.Fprintf(&buf, "}")
	return buf.String()
}

// Dir provides access to the parsed content of a tiff Image File Directory (IFD).
type Dir struct {
	Tags []*Tag
}

// DecodeDir parses a tiff-encoded IFD from r and returns a Dir object.  offset
// is the offset to the next IFD.  The first read from r should be at the first
// byte of the IFD. ReadAt offsets should generally be relative to the
// beginning of the tiff structure (not relative to the beginning of the IFD).
func DecodeDir(r ReadAtReader, order binary.ByteOrder) (d *Dir, offset int32, err error) {
	d = new(Dir)

	// get num of tags in ifd
	var nTags int16
	// 下标[8:10] 获取标签总数
	err = binary.Read(r, order, &nTags)
	if err != nil {
		return nil, 0, errors.New("tiff: failed to read IFD tag count: " + err.Error())
	}
	// load tags
	for n := 0; n < int(nTags); n++ {
		// 首个标签下标索引为10
		t, err := DecodeTag(r, order)
		if err != nil {
			return nil, 0, err
		}
		d.Tags = append(d.Tags, t)
	}

	// get offset to next ifd
	err = binary.Read(r, order, &offset)
	if err != nil {
		return nil, 0, errors.New("tiff: falied to read offset to next IFD: " + err.Error())
	}

	return d, offset, nil
}

func DecodeDefineTag(r ReadAtReader, order binary.ByteOrder, fieldId uint16) (t *Tag, offset int32, err error) {
	t = new(Tag)
	// get num of tags in ifd
	var nTags int16
	// 下标[8:10] 获取标签总数
	err = binary.Read(r, order, &nTags)
	if err != nil {
		return nil, 0, errors.New("tiff: failed to read IFD tag count: " + err.Error())
	}
	// load tags
	for n := 0; n < int(nTags); n++ {
		// 获取当前字节组下标值
		//curIndex, err := r.Seek(0, io.SeekCurrent)
		//if err == nil {
		//	fmt.Println("tag data index: ", curIndex)
		//}
		// 首个标签下标索引为10
		t, err = DecodeDefineTagData(r, order, fieldId)
		if err != nil {
			return nil, 0, err
		}
		if t == nil {
			// 每个IFD 是12bit, 由于读取tag ID 已经偏移了2bit, 只需偏移10bit 就能读取下一个 IFD 数据
			_, err := r.Seek(ReadTagIdOffsetLen, io.SeekCurrent)
			if err != nil {
				return nil, 0, err
			}
			continue
		}
		break
	}

	// get offset to next ifd
	err = binary.Read(r, order, &offset)
	if err != nil {
		return nil, 0, errors.New("tiff: falied to read offset to next IFD: " + err.Error())
	}

	return t, offset, nil
}

func (d *Dir) String() string {
	s := "Dir{"
	for _, t := range d.Tags {
		s += t.String() + ", "
	}
	return s + "}"
}
