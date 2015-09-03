package gulpd

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"time"
	"github.com/megamsys/gulp/activities"
)

type Handler struct {
	Provider string
	d        *gulpd.Config
}

// NewHandler returns a new instance of handler.
func NewHandler(conf *gulpd.Config) *Handler {
	h := &Handler{
		d: conf,
	}	

	return h
}

func (h *Handler) ServeAMQP(r app.Requests) error {
	assembly, err := app.Get(r.Id)
	if err != nil {
		return err
	}

	di := app.ActionData{
		Assembly: assembly,
		Request: r,
		Config:   h.d,
	}
	
	p, err := activities.GetActivity(r.Category)	
	if err = p.Action(&di); err != nil {
			return err
		}	
	
	return nil
}
