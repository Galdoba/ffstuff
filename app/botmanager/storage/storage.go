package storage

import (
	"crypto/sha1"
	"errors"
	"fmt"
	"io"

	"github.com/Galdoba/ffstuff/pkg/gerror"
)

var ErrNoFiles = errors.New("no saved pages")

type Storage interface {
	Save(p *Page) error
	Remove(p *Page) error
	IsExists(p *Page) (bool, error)
	PickRandom(userName string) (*Page, error)
}

type Page struct {
	URL      string
	UserName string
}

func (p Page) Hash() (string, error) {
	h := sha1.New()
	if _, err := io.WriteString(h, p.URL); err != nil {
		return "", gerror.Wrap("can't hash", err)
	}
	if _, err := io.WriteString(h, p.UserName); err != nil {
		return "", gerror.Wrap("can't hash", err)
	}
	return fmt.Sprintf("%x", h.Sum(nil)), nil
}
