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

char *hash_to_hex_string_reversed_wr(void *h, int size) {
	char *out = malloc(size * 2 + 1);
	hash_to_hex_string_reversed((uchar*)h, out, size);
	return out;
}

char *hash_to_hex_string_wr(void *h, int size) {
	char *out = malloc(size * 2 + 1);
	hash_to_hex_string((uchar*)h, out, size);
	return out;
}

uchar *phash_mem_wr(void *buf, size_t buf_len, int hash_size, int highfreq_factor, int* ret) {
	uchar *out = malloc(hash_size);
	*ret = phash_mem((uchar*)buf, out, buf_len, hash_size, highfreq_factor);
	return out;
}

uchar *ahash_mem_wr(void *buf, size_t buf_len, int hash_size, int* ret) {
	uchar *out = malloc(hash_size);
	*ret = ahash_mem((uchar*)buf, out, buf_len, hash_size);
	return out;
}

uchar *dhash_mem_wr(void *buf, size_t buf_len, int hash_size, int* ret) {
	uchar *out = malloc(hash_size);
	*ret = dhash_mem((uchar*)buf, out, buf_len, hash_size);
	return out;
}

uchar *whash_mem_wr(void *buf, size_t buf_len, int hash_size, int img_scale, int* ret) {
	uchar *out = malloc(hash_size);
	*ret = whash_mem((uchar*)buf, out, buf_len, hash_size, img_scale);
	return out;
}

multi_hash_t *multi_mem_wr(void* buf, size_t buf_len, int hash_size, int ph_highfreq_factor, int wh_img_scale, int* ret) {
	multi_hash_t *m = multi_hash_create(hash_size);
	*ret = multi_hash_mem(buf, m, buf_len, hash_size, ph_highfreq_factor, wh_img_scale);
	return m;
}
*/
import "C"

type Code int

const (
	Ok = 0
	ReadErr Code = -2
	DecodeErr Code = -3
)

type Hash struct {
	Size  int
	Bytes []byte
}

type MultiHash struct {
	PHash Hash
	AHash Hash
	DHash Hash
	WHash Hash
}

func retHash(hash *C.uchar, hashSize int, ret C.int) (*Hash, Code) {
	goHash := C.GoBytes(unsafe.Pointer(hash), C.int(hashSize))
	C.free(unsafe.Pointer(hash))

	return &Hash{
		Size:  hashSize,
		Bytes: goHash,
	}, Code(ret)
}

func retMultiHash(m *C.multi_hash_t, hashSize int, ret C.int) (*MultiHash, Code) {
	goPHash := C.GoBytes(unsafe.Pointer(m.phash), C.int(hashSize))
	goAHash := C.GoBytes(unsafe.Pointer(m.ahash), C.int(hashSize))
	goDHash := C.GoBytes(unsafe.Pointer(m.dhash), C.int(hashSize))
	goWHash := C.GoBytes(unsafe.Pointer(m.whash), C.int(hashSize))
	C.multi_hash_destroy(m)

	return &MultiHash{
		PHash: Hash{hashSize, goPHash},
		AHash: Hash{hashSize, goAHash},
		DHash: Hash{hashSize, goDHash},
		WHash: Hash{hashSize, goWHash},
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

func WHashFile(filepath string, hashSize, imgScale int) (*Hash, Code) {
	bytes, err := readAll(filepath)
	if err != nil {
		return nil, ReadErr
	}
	return WHashMem(bytes, hashSize, imgScale)
}

func WHashMem(buf []byte, hashSize, imgScale int) (*Hash, Code) {
	var ret C.int
	hash := C.whash_mem_wr(unsafe.Pointer(&buf[0]), C.size_t(len(buf)), C.int(hashSize), C.int(imgScale), &ret)
	return retHash(hash, hashSize, ret)
}

func MultiHashFile(filepath string, hashSize, phHighFreqFactor, whImgScale int) (*MultiHash, Code) {
	bytes, err := readAll(filepath)
	if err != nil {
		return nil, ReadErr
	}
	return MultiHashMem(bytes, hashSize, phHighFreqFactor, whImgScale)
}

func MultiHashMem(buf []byte, hashSize, phHighFreqFactor, whImgScale int) (*MultiHash, Code) {
	var ret C.int
	m := C.multi_mem_wr(unsafe.Pointer(&buf[0]), C.size_t(len(buf)), C.int(hashSize),
		C.int(phHighFreqFactor), C.int(whImgScale), &ret)
	return retMultiHash(m, hashSize, ret)
}
