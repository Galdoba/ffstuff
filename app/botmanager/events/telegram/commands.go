package telegram

import (
	"errors"
	"log"
	"net/url"
	"strings"

	"github.com/Galdoba/ffstuff/app/botmanager/storage"
	"github.com/Galdoba/ffstuff/pkg/gerror"
)

const (
	RndCmd   = "/rnd"
	HelpCmd  = "/help"
	StartCmd = "/start"
)

func (p *ProcessorUnit) doCmd(text string, chatID int, username string) error {
	text = strings.TrimSpace(text)
	log.Printf("got new command '%s' from '%s'", text, username)

	if isAddCmd(text) {
		return p.savePage(chatID, text, username)
	}

	//add page: http://...
	//rnd page: /rnd
	//help: /help
	//start: /start: hi + help
	switch text {
	case RndCmd:
		return p.sendRandom(chatID, username)
	case HelpCmd:
		return p.sendHelp(chatID)
	case StartCmd:
		return p.sendHello(chatID)
	default:
		return p.tg.SendMessage(chatID, msgUnknownCommand)
	}
}

func (p *ProcessorUnit) savePage(chatID int, pageURL string, username string) (err error) {
	defer func() { err = gerror.WrapIfErr("can't do cmd 'save page'", err) }()
	page := &storage.Page{
		URL:      pageURL,
		UserName: username,
	}
	isExists, err := p.storage.IsExists(page)
	if err != nil {
		return err
	}
	if isExists {
		return p.tg.SendMessage(chatID, msgAlreadyExists)
	}
	if err := p.storage.Save(page); err != nil {
		return err
	}
	if err := p.tg.SendMessage(chatID, msgSaved); err != nil {
		return err
	}
	return nil
}

func (p *ProcessorUnit) sendRandom(chatID int, username string) (err error) {
	defer func() { err = gerror.WrapIfErr("can't do cmd 'send random'", err) }()
	page, err := p.storage.PickRandom(username)
	if err != nil && !errors.Is(err, storage.ErrNoFiles) {
		return err
	}
	if errors.Is(err, storage.ErrNoFiles) {
		p.tg.SendMessage(chatID, msgNoSavedPages)
	}
	if page == nil {
		return nil
	}
	if err := p.tg.SendMessage(chatID, page.URL); err != nil {
		return err
	}
	return p.storage.Remove(page)
}

func (p *ProcessorUnit) sendHelp(chatID int) (err error) {
	return p.tg.SendMessage(chatID, msgHelp)
}

func (p *ProcessorUnit) sendHello(chatID int) (err error) {
	return p.tg.SendMessage(chatID, msgHello)
}

func isAddCmd(text string) bool {
	return isURL(text)
}

func isURL(t string) bool {
	u, err := url.Parse(t)
	return err == nil && u.Host != ""
}
