package iface

import (
	"io"
	"os"
)

// run this with:
// go generate ./...

//go:generate mockgen -destination=../mocks/mainmock.go -package=mocks github.com/FredHutch/sftp_downloader/iface Sftper,Filer,Walker

//go:generate mockgen -destination=../mocks/fileinfo.go -package=mocks os FileInfo

//go:generate mockgen -destination=../mocks/ioreader.go -package=mocks io Reader

// Sftper helps make things testable
type Sftper interface {
	// return interface
	Walk(root string) Walker

	// return interface
	Create(path string) (Filer, error)

	// return concrete type since no methods on FileInfo are used by doTheWork
	Lstat(p string) (os.FileInfo, error)

	Close() error

	ReadDir(p string) ([]os.FileInfo, error)

	Open(path string) (io.Reader, error)
}

// Filer interface has methods used by doTheWork
type Filer interface {
	Write([]byte) (int, error)
	Close() error
}

// Walker helps make things testable
type Walker interface {
	Step() bool
	Err() error
	Path() string
}
