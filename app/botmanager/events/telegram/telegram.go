package telegram

import (
	"errors"

	"github.com/Galdoba/ffstuff/app/botmanager/clients/telegram"
	"github.com/Galdoba/ffstuff/app/botmanager/events"
	"github.com/Galdoba/ffstuff/app/botmanager/storage"
	"github.com/Galdoba/ffstuff/pkg/gerror"
)

var ErrUnknownEvent = errors.New("unknown event type")
var ErrMetaUnknown = errors.New("unknown meta type")

type ProcessorUnit struct {
	tg      *telegram.Client
	offset  int
	storage storage.Storage
}

func New(client *telegram.Client, storage storage.Storage) *ProcessorUnit {
	return &ProcessorUnit{
		tg:      client,
		storage: storage,
	}
}

type Meta struct {
	UserName string
	ChatID   int
}

func (p *ProcessorUnit) Fetch(limit int) ([]events.Event, error) {
	updates, err := p.tg.Update(p.offset, limit)
	if err != nil {
		return nil, gerror.Wrap("can't get events", err)
	}
	if len(updates) == 0 {
		return nil, nil //fmt.Errorf("no updates found")
	}
	res := make([]events.Event, 0, len(updates))

	for _, u := range updates {
		res = append(res, event(u))
	}
	p.offset = updates[len(updates)-1].ID + 1
	return res, nil
}

func (p *ProcessorUnit) Process(event events.Event) error {
	switch event.Type {
	case events.Message:
		return p.processMessage(event)
	default:
		return gerror.Wrap("can't process message", ErrUnknownEvent)
	}
}

func (p *ProcessorUnit) processMessage(event events.Event) error {
	meta, err := meta(event)
	if err != nil {
		return gerror.Wrap("can't process message", err)
	}
	if err := p.doCmd(event.Text, meta.ChatID, meta.UserName); err != nil {
		return gerror.Wrap("can't process message", err)
	}
	return nil
}

func meta(e events.Event) (Meta, error) {
	res, ok := e.Meta.(Meta)
	if !ok {
		return Meta{}, gerror.Wrap("can't get meta", ErrMetaUnknown)
	}
	return res, nil
}

func event(upd telegram.Update) events.Event {
	updType := fetchType(upd)
	res := events.Event{
		Type: updType,
		Text: fetchText(upd),
	}
	if updType == events.Message {
		res.Meta = Meta{
			ChatID:   upd.Message.Chat.ID,
			UserName: upd.Message.From.UserName,
		}
	}
	return res
}

func fetchType(u telegram.Update) events.Type {
	if u.Message == nil {
		return events.Unknown
	}
	return events.Message
}

func fetchText(u telegram.Update) string {
	if u.Message == nil {
		return ""
	}
	return u.Message.Text
}
