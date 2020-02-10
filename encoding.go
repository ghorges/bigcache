package bigcache

import (
	"encoding/binary"
)

/*
storage like this.
--------------------------------------------------
|timestampSize|keyHashSize|keySize|keyValue|entry|
--------------------------------------------------
 */

const (
	timestampSize = 8
	keyHashSize   = 8
	keySize       = 2
	headSize      = timestampSize + keyHashSize + keySize
)

func wrapEntry(timestamp int64, hash uint64, key string, entry []byte, buffer *[]byte) []byte {
	lenKey := len(key)
	// lenData := len(entry)	just use once,so don't use it.
	entryLen := headSize + lenKey + len(entry)
	if entryLen > cap(*buffer) {
		*buffer = make([]byte, entryLen)
	}

	// if dont't have this, will appear (*buffer)[],too cumbrous
	blob := *buffer

	binary.LittleEndian.PutUint64(blob, uint64(timestamp))
	binary.LittleEndian.PutUint64(blob[timestampSize:], hash)
	binary.LittleEndian.PutUint16(blob[timestampSize+keyHashSize:], uint16(lenKey))

	copy(blob[headSize:], key)
	copy(blob[headSize+lenKey:], entry)

	return blob[:entryLen]
}

func readEntry(buffer []byte) []byte {
	keyLen := binary.LittleEndian.Uint16(buffer[timestampSize+keyHashSize:])

	returnBuffer := make([]byte, len(buffer)-headSize-int(keyLen))
	copy(returnBuffer, buffer[headSize+keyLen:])

	return returnBuffer
}

func hashKeyToZero(buffer []byte) {
	binary.LittleEndian.PutUint64(buffer[timestampSize:], 0)
}

func getHashKey(buffer []byte) uint64 {
	return binary.LittleEndian.Uint64(buffer[timestampSize:])
}

func getKey(buffer []byte) []byte {
	keyLen := binary.LittleEndian.Uint16(buffer[timestampSize+keyHashSize:])

	returnBuffer := make([]byte, keyLen)
	copy(returnBuffer, buffer[headSize:])

	return returnBuffer
}

func readTimestamp(buffer []byte) uint64 {
	return binary.LittleEndian.Uint64(buffer[0:timestampSize])
}
