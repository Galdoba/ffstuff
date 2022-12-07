package files

import (
	"encoding/gob"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"time"

	"github.com/Galdoba/ffstuff/app/botmanager/storage"
	"github.com/Galdoba/ffstuff/pkg/gerror"
)

const (
	defaultPerm = 0774
)

type StorageUnit struct {
	basePath string
}

func New(path string) *StorageUnit {
	return &StorageUnit{path}
}

func (su *StorageUnit) PickRandom(username string) (page *storage.Page, err error) {
	defer func() { err = gerror.WrapIfErr("can't pick random page: ", err) }()

	path := filepath.Join(su.basePath, username)

	files, err := ioutil.ReadDir(path)
	if err != nil {
		return nil, err
	}

	if len(files) == 0 {
		return nil, storage.ErrNoFiles
	}

	rand.Seed(time.Now().UnixNano())
	n := rand.Intn(len(files))
	file := files[n]
	return su.decodePage(filepath.Join(path, file.Name()))
}

func (su *StorageUnit) Remove(p *storage.Page) error {
	f, err := fileName(p)
	if err != nil {
		return gerror.Wrap("can't remove page: ", err)
	}
	path := filepath.Join(su.basePath, p.UserName, f)
	if err := os.Remove(path); err != nil {
		fmt.Println("----", err.Error())
		return gerror.Wrap(fmt.Sprintf("can't remove page [%v]: ", path), err)
	}
	log.Printf("%v removed", path)
	return nil
}

func (su *StorageUnit) IsExists(p *storage.Page) (bool, error) {
	f, err := fileName(p)
	if err != nil {
		return false, gerror.Wrap("can't check if file exists: ", err)
	}
	path := filepath.Join(su.basePath, p.UserName, f)

	switch _, err = os.Stat(path); {
	case errors.Is(err, os.ErrNotExist):
		return false, nil
	case err != nil:
		return false, gerror.Wrap("can't check if file exists: ", err)
	}
	return true, nil
}

func (su *StorageUnit) decodePage(filePath string) (*storage.Page, error) {
	f, err := os.Open(filePath)
	defer f.Close()

	if err != nil {
		return nil, gerror.Wrap("can't decode page: ", err)
	}
	var p storage.Page
	if err := gob.NewDecoder(f).Decode(&p); err != nil {
		if err.Error() == "EOF" {
			return &p, nil
		}
		return nil, gerror.Wrap("can't decode page: ", err)
	}
	return &p, nil
}

func (su *StorageUnit) Save(page *storage.Page) (err error) {
	defer func() { err = gerror.WrapIfErr("can't save: ", err) }()
	filePaths := filepath.Join(su.basePath, page.UserName)
	if err := os.MkdirAll(filePaths, defaultPerm); err != nil {
		return err
	}
	fName, err := fileName(page)
	if err != nil {
		return err
	}
	filePaths = filepath.Join(filePaths, fName)
	file, err := os.Create(filePaths)
	if err != nil {
		return err
	}
	defer func() { _ = closeLogged(file) }()

	if err := gob.NewEncoder(file).Encode(page); err != nil {
		return err
	}
	return nil
}

func fileName(p *storage.Page) (string, error) {
	return p.Hash()
}

func closeLogged(f *os.File) error {
	err := f.Close()
	if err != nil {
		log.Printf("log close err: %v")
	}
	return err
}
