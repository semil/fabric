package etcdraft
//
//import (
//	"go.etcd.io/etcd/raft/v3/raftpb"
//	"go.etcd.io/etcd/server/v3/wal"
//	"go.etcd.io/etcd/server/v3/wal/walpb"
//	"go.uber.org/zap"
//	"sync/atomic"
//)
//
//type WalProxy struct {
//	*wal.WAL
//	SetUnsafeNoFsyncCalled bool
//	SyncCounter            uint64
//}
//
////import (
////"bytes"
////"errors"
////"fmt"
////"hash/crc32"
////"io"
////"os"
////"path/filepath"
////"strings"
////"sync"
////"time"
////
////"go.etcd.io/etcd/client/pkg/v3/fileutil"
////"go.etcd.io/etcd/pkg/v3/pbutil"
////"go.etcd.io/etcd/raft/v3"
////"go.etcd.io/etcd/raft/v3/raftpb"
////"go.etcd.io/etcd/server/v3/wal/walpb"
////
////"go.uber.org/zap"
////)
//
////func CreateWalProxy(lg *zap.Logger, dirpath string, metadata []byte) (*WalProxy, error) {
////	w, err := wal.Create(lg, dirpath, metadata)
////	if err != nil {
////		return nil, err
////	}
////	return &WalProxy{
////		WAL:                    w,
////		SetUnsafeNoFsyncCalled: false,
////	}, nil
////}
//
//func CreateWalProxy(w *wal.WAL) *WalProxy {
//	return &WalProxy{
//		WAL:                    w,
//		SetUnsafeNoFsyncCalled: false,
//	}
//}
//
//func OpenWalProxy(lg *zap.Logger, dirpath string, snap walpb.Snapshot) (*WalProxy, error) {
//	w, err := wal.Open(lg, dirpath, snap)
//	if err != nil {
//		return nil, err
//	}
//	return CreateWalProxy(w), nil
//}
//
//func (w *WalProxy) SetUnsafeNoFsync() {
//	w.WAL.SetUnsafeNoFsync()
//	w.SetUnsafeNoFsyncCalled = true
//}
//
//func (w *WalProxy) ReadAll() (metadata []byte, state raftpb.HardState, ents []raftpb.Entry, err error) {
//	return w.WAL.ReadAll()
//}
//
////func (w *WalProxy) Sync() error {
////	err := w.WAL.Sync()
////	atomic.AddUint64(&w.SyncCounter, 1)
////	return err
////}
//
//func (w *WalProxy) sync() error {
//	err := w.WAL.Sync()
//	atomic.AddUint64(&w.SyncCounter, 1)
//	return err
//}
