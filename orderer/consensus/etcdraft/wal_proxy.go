package etcdraft

import (
	"errors"
	"go.etcd.io/etcd/client/pkg/v3/fileutil"
	"go.etcd.io/etcd/raft/v3/raftpb"
	"go.etcd.io/etcd/server/v3/wal"
	"go.etcd.io/etcd/server/v3/wal/walpb"
	"go.uber.org/zap"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync/atomic"
)

//import (
//"bufio"
//"encoding/binary"
//"hash"
//"io"
//"sync"
//
//"go.etcd.io/etcd/pkg/v3/crc"
//"go.etcd.io/etcd/pkg/v3/pbutil"
//"go.etcd.io/etcd/raft/v3/raftpb"
//"go.etcd.io/etcd/server/v3/wal/walpb"
//)

//const (
//	metadataType int64 = iota + 1
//	entryType
//	stateType
//	crcType
//	snapshotType
//
//	// warnSyncDuration is the amount of time allotted to an fsync before
//	// logging a warning
//	warnSyncDuration = time.Second
//)
//
//var crcTable = crc32.MakeTable(crc32.Castagnoli)
//
//const minSectorSize = 512
//
//// frameSizeBytes is frame size in bytes, including record size and padding size.
//const frameSizeBytes = 8
//
//type Decoder struct {
//	mu  sync.Mutex
//	brs []*bufio.Reader
//
//	// lastValidOff file offset following the last valid decoded record
//	lastValidOff int64
//	crc          hash.Hash32
//}
//
//func NewDecoder(r ...io.Reader) *Decoder {
//	readers := make([]*bufio.Reader, len(r))
//	for i := range r {
//		readers[i] = bufio.NewReader(r[i])
//	}
//	return &Decoder{
//		brs: readers,
//		crc: crc.New(0, crcTable),
//	}
//}
//
//func (d *Decoder) Decode(rec *walpb.Record) error {
//	rec.Reset()
//	d.mu.Lock()
//	defer d.mu.Unlock()
//	return d.decodeRecord(rec)
//}
//
//// raft max message size is set to 1 MB in etcd server
//// assume projects set reasonable message size limit,
//// thus entry size should never exceed 10 MB
//const maxWALEntrySizeLimit = int64(10 * 1024 * 1024)
//
//func (d *Decoder) decodeRecord(rec *walpb.Record) error {
//	if len(d.brs) == 0 {
//		return io.EOF
//	}
//
//	l, err := readInt64(d.brs[0])
//	if err == io.EOF || (err == nil && l == 0) {
//		// hit end of file or preallocated space
//		d.brs = d.brs[1:]
//		if len(d.brs) == 0 {
//			return io.EOF
//		}
//		d.lastValidOff = 0
//		return d.decodeRecord(rec)
//	}
//	if err != nil {
//		return err
//	}
//
//	recBytes, padBytes := decodeFrameSize(l)
//	if recBytes >= maxWALEntrySizeLimit-padBytes {
//		return wal.ErrMaxWALEntrySizeLimitExceeded
//	}
//
//	data := make([]byte, recBytes+padBytes)
//	if _, err = io.ReadFull(d.brs[0], data); err != nil {
//		// ReadFull returns io.EOF only if no bytes were read
//		// the Decoder should treat this as an ErrUnexpectedEOF instead.
//		if err == io.EOF {
//			err = io.ErrUnexpectedEOF
//		}
//		return err
//	}
//	if err := rec.Unmarshal(data[:recBytes]); err != nil {
//		if d.isTornEntry(data) {
//			return io.ErrUnexpectedEOF
//		}
//		return err
//	}
//
//	// skip crc checking if the record type is crcType
//	if rec.Type != crcType {
//		d.crc.Write(rec.Data)
//		if err := rec.Validate(d.crc.Sum32()); err != nil {
//			if d.isTornEntry(data) {
//				return io.ErrUnexpectedEOF
//			}
//			return err
//		}
//	}
//	// record decoded as valid; point last valid offset to end of record
//	d.lastValidOff += frameSizeBytes + recBytes + padBytes
//	return nil
//}
//
//func decodeFrameSize(lenField int64) (recBytes int64, padBytes int64) {
//	// the record size is stored in the lower 56 bits of the 64-bit length
//	recBytes = int64(uint64(lenField) & ^(uint64(0xff) << 56))
//	// non-zero padding is indicated by set MSb / a negative length
//	if lenField < 0 {
//		// padding is stored in lower 3 bits of length MSB
//		padBytes = int64((uint64(lenField) >> 56) & 0x7)
//	}
//	return recBytes, padBytes
//}
//
//// isTornEntry determines whether the last entry of the WAL was partially written
//// and corrupted because of a torn write.
//func (d *Decoder) isTornEntry(data []byte) bool {
//	if len(d.brs) != 1 {
//		return false
//	}
//
//	fileOff := d.lastValidOff + frameSizeBytes
//	curOff := 0
//	chunks := [][]byte{}
//	// split data on sector boundaries
//	for curOff < len(data) {
//		chunkLen := int(minSectorSize - (fileOff % minSectorSize))
//		if chunkLen > len(data)-curOff {
//			chunkLen = len(data) - curOff
//		}
//		chunks = append(chunks, data[curOff:curOff+chunkLen])
//		fileOff += int64(chunkLen)
//		curOff += chunkLen
//	}
//
//	// if any data for a sector chunk is all 0, it's a torn write
//	for _, sect := range chunks {
//		isZero := true
//		for _, v := range sect {
//			if v != 0 {
//				isZero = false
//				break
//			}
//		}
//		if isZero {
//			return true
//		}
//	}
//	return false
//}
//
//func (d *Decoder) updateCRC(prevCrc uint32) {
//	d.crc = crc.New(prevCrc, crcTable)
//}
//
//func (d *Decoder) lastCRC() uint32 {
//	return d.crc.Sum32()
//}
//
//func (d *Decoder) lastOffset() int64 { return d.lastValidOff }
//
//func mustUnmarshalEntry(d []byte) raftpb.Entry {
//	var e raftpb.Entry
//	pbutil.MustUnmarshal(&e, d)
//	return e
//}
//
//func mustUnmarshalState(d []byte) raftpb.HardState {
//	var s raftpb.HardState
//	pbutil.MustUnmarshal(&s, d)
//	return s
//}
//
//func readInt64(r io.Reader) (int64, error) {
//	var n int64
//	err := binary.Read(r, binary.LittleEndian, &n)
//	return n, err
//}

type WalProxy struct {
	*wal.WAL
	SetUnsafeNoFsyncCalled bool
	SyncCounter            uint64
}

func OpenWALFiles(lg *zap.Logger, dirpath string, names []string, nameIndex int, write bool) ([]io.Reader, []*fileutil.LockedFile, func() error, error) {
	rcs := make([]io.ReadCloser, 0)
	rs := make([]io.Reader, 0)
	ls := make([]*fileutil.LockedFile, 0)
	for _, name := range names[nameIndex:] {
		p := filepath.Join(dirpath, name)
		if write {
			l, err := fileutil.TryLockFile(p, os.O_RDWR, fileutil.PrivateFileMode)
			if err != nil {
				closeAll(lg, rcs...)
				return nil, nil, nil, err
			}
			ls = append(ls, l)
			rcs = append(rcs, l)
		} else {
			rf, err := os.OpenFile(p, os.O_RDONLY, fileutil.PrivateFileMode)
			if err != nil {
				closeAll(lg, rcs...)
				return nil, nil, nil, err
			}
			ls = append(ls, nil)
			rcs = append(rcs, rf)
		}
		rs = append(rs, rcs[len(rcs)-1])
	}

	closer := func() error { return closeAll(lg, rcs...) }

	return rs, ls, closer, nil
}

func closeAll(lg *zap.Logger, rcs ...io.ReadCloser) error {
	stringArr := make([]string, 0)
	for _, f := range rcs {
		if err := f.Close(); err != nil {
			lg.Warn("failed to close: ", zap.Error(err))
			stringArr = append(stringArr, err.Error())
		}
	}
	if len(stringArr) == 0 {
		return nil
	}
	return errors.New(strings.Join(stringArr, ", "))
}

func CreateWalProxy(w *wal.WAL) *WalProxy {
	return &WalProxy{
		WAL:                    w,
		SetUnsafeNoFsyncCalled: false,
	}
}

func OpenWalProxy(lg *zap.Logger, dirpath string, snap walpb.Snapshot) (*WalProxy, error) {
	w, err := wal.Open(lg, dirpath, snap)
	if err != nil {
		return nil, err
	}
	return CreateWalProxy(w), nil
}

func (w *WalProxy) SetUnsafeNoFsync() {
	w.WAL.SetUnsafeNoFsync()
	w.SetUnsafeNoFsyncCalled = true
}

//func (w *WalProxy) GetDecoder() {
//	NewDecoder(w)
//	w.WAL.decoder
//}

func (w *WalProxy) ReadAll() (metadata []byte, state raftpb.HardState, ents []raftpb.Entry, err error) {
	return w.WAL.ReadAll()
}

//func (w *WalProxy) Sync() error {
//	err := w.WAL.Sync()
//	atomic.AddUint64(&w.SyncCounter, 1)
//	return err
//}

func (w *WalProxy) sync() error {
	err := w.WAL.Sync()
	atomic.AddUint64(&w.SyncCounter, 1)
	return err
}
