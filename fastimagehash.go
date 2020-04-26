package fastimagehash

import "C"
import (
	"io/ioutil"
	"os"
	"unsafe"
)

/*
#cgo LDFLAGS: -L. -lfastimagehash
#include <fastimagehash.h>
#include <stdlib.h>
#include <string.h>

const char* Version = FASTIMAGEHASH_VERSION;

char *hash_to_hex_string_reversed_wr(void *h, int size) {
	char *out = malloc(size * size / 4 + 1);
	hash_to_hex_string_reversed((uchar*)h, out, size);
	return out;
}

char *hash_to_hex_string_wr(void *h, int size) {
	char *out = malloc(size * size / 4 + 1);
	hash_to_hex_string((uchar*)h, out, size);
	return out;
}

uchar *phash_mem_wr(void *buf, size_t buf_len, int hash_size, int highfreq_factor, int* ret) {
	uchar *out = malloc(hash_size * hash_size / 8);
	*ret = phash_mem((uchar*)buf, buf_len, out, hash_size, highfreq_factor);
	return out;
}

uchar *ahash_mem_wr(void *buf, size_t buf_len, int hash_size, int* ret) {
	uchar *out = malloc(hash_size * hash_size / 8);
	*ret = ahash_mem((uchar*)buf, buf_len, out, hash_size);
	return out;
}

uchar *mhash_mem_wr(void *buf, size_t buf_len, int hash_size, int* ret) {
	uchar *out = malloc(hash_size * hash_size / 8);
	*ret = mhash_mem((uchar*)buf, buf_len, out, hash_size);
	return out;
}

uchar *dhash_mem_wr(void *buf, size_t buf_len, int hash_size, int* ret) {
	uchar *out = malloc(hash_size * hash_size / 8);
	*ret = dhash_mem((uchar*)buf, buf_len, out, hash_size);
	return out;
}

uchar *whash_mem_wr(void *buf, size_t buf_len, int hash_size, int img_scale, int remove_max_ll, int* ret, _GoString_ go_wname) {
	uchar *out = malloc(hash_size * hash_size / 8);

	if (strncmp(_GoStringPtr(go_wname), "haar", 4) == 0) {
		*ret = whash_mem((uchar*)buf, buf_len, out, hash_size, img_scale, remove_max_ll, "haar");
	} else {
		*ret = whash_mem((uchar*)buf, buf_len, out, hash_size, img_scale, remove_max_ll, "db4");
	}

	return out;
}
*/
import "C"

type Code int

const (
	Ok             = 0
	ReadErr   Code = -2
	DecodeErr Code = -3
)

type Wave string

const (
	Haar       = "haar"
	Daubechies = "db4"
)

type Hash struct {
	Size  int    `json:"size"`
	Bytes []byte `json:"bytes"`
}

var LibVersion = C.GoString(C.Version)

func retHash(hash *C.uchar, hashSize int, ret C.int) (*Hash, Code) {
	if ret == Ok {
		goHash := C.GoBytes(unsafe.Pointer(hash), C.int(hashSize*hashSize/8))
		C.free(unsafe.Pointer(hash))

		return &Hash{
			Size:  hashSize,
			Bytes: goHash,
		}, Code(ret)
	}

	return &Hash{
		Size:  hashSize,
		Bytes: nil,
	}, Code(ret)
}

func readAll(filepath string) ([]byte, error) {
	f, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}

	bytes, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}
	return bytes, nil
}

func (h *Hash) ToHexStringReversed() (ret string) {
	out := C.hash_to_hex_string_reversed_wr(unsafe.Pointer(&h.Bytes[0]), C.int(h.Size))
	ret = C.GoString(out)
	C.free(unsafe.Pointer(out))
	return
}

func (h *Hash) ToHexString() (ret string) {
	out := C.hash_to_hex_string_wr(unsafe.Pointer(&h.Bytes[0]), C.int(h.Size))
	ret = C.GoString(out)
	C.free(unsafe.Pointer(out))
	return
}

func PHashFile(filepath string, hashSize, highFreqFactor int) (*Hash, Code) {
	bytes, err := readAll(filepath)
	if err != nil {
		return nil, ReadErr
	}
	return PHashMem(bytes, hashSize, highFreqFactor)
}

func PHashMem(buf []byte, hashSize, highFreqFactor int) (*Hash, Code) {
	var ret C.int
	hash := C.phash_mem_wr(unsafe.Pointer(&buf[0]), C.size_t(len(buf)), C.int(hashSize), C.int(highFreqFactor), &ret)
	return retHash(hash, hashSize, ret)
}

func AHashFile(filepath string, hashSize int) (*Hash, Code) {
	bytes, err := readAll(filepath)
	if err != nil {
		return nil, ReadErr
	}
	return AHashMem(bytes, hashSize)
}

func AHashMem(buf []byte, hashSize int) (*Hash, Code) {
	var ret C.int
	hash := C.ahash_mem_wr(unsafe.Pointer(&buf[0]), C.size_t(len(buf)), C.int(hashSize), &ret)
	return retHash(hash, hashSize, ret)
}

func MHashFile(filepath string, hashSize int) (*Hash, Code) {
	bytes, err := readAll(filepath)
	if err != nil {
		return nil, ReadErr
	}
	return MHashMem(bytes, hashSize)
}

func MHashMem(buf []byte, hashSize int) (*Hash, Code) {
	var ret C.int
	hash := C.mhash_mem_wr(unsafe.Pointer(&buf[0]), C.size_t(len(buf)), C.int(hashSize), &ret)
	return retHash(hash, hashSize, ret)
}

func DHashFile(filepath string, hashSize int) (*Hash, Code) {
	bytes, err := readAll(filepath)
	if err != nil {
		return nil, ReadErr
	}
	return DHashMem(bytes, hashSize)
}

func DHashMem(buf []byte, hashSize int) (*Hash, Code) {
	var ret C.int
	hash := C.dhash_mem_wr(unsafe.Pointer(&buf[0]), C.size_t(len(buf)), C.int(hashSize), &ret)
	return retHash(hash, hashSize, ret)
}

func WHashFile(filepath string, hashSize, imgScale int, removeMaxLL bool, wave Wave) (*Hash, Code) {
	bytes, err := readAll(filepath)
	if err != nil {
		return nil, ReadErr
	}
	return WHashMem(bytes, hashSize, imgScale, removeMaxLL, wave)
}

func WHashMem(buf []byte, hashSize, imgScale int, removeMaxLL bool, wave Wave) (*Hash, Code) {
	var ret C.int
	var remove_max_ll C.int
	if removeMaxLL {
		remove_max_ll = 1
	} else {
		remove_max_ll = 0
	}
	hash := C.whash_mem_wr(unsafe.Pointer(&buf[0]), C.size_t(len(buf)), C.int(hashSize), C.int(imgScale), remove_max_ll, &ret, string(wave))
	return retHash(hash, hashSize, ret)
}
