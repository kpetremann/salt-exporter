package listener

import (
	"context"
	"os"
	"path"
	"strings"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/kpetremann/salt-exporter/pkg/event"
	"github.com/rs/zerolog/log"
)

const DefaultPKIDirpath = "/etc/salt/pki/master"

type PKIWatcher struct {
	ctx        context.Context
	pkiDirPath string
	watcher    *fsnotify.Watcher
	eventChan  chan<- event.WatchEvent
	lock       sync.RWMutex
}

func NewPKIWatcher(ctx context.Context, pkiDirPath string, eventChan chan event.WatchEvent) (*PKIWatcher, error) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}

	w := &PKIWatcher{
		ctx:        ctx,
		pkiDirPath: pkiDirPath,
		watcher:    watcher,
		eventChan:  eventChan,
		lock:       sync.RWMutex{},
	}

	return w, nil
}

// SetPKIDirectory sets the filepath to the salt-master pki directory
//
// The directory must be readable by the user running the exporter (usually salt).
//
// Default: /etc/salt/pki.
func (w *PKIWatcher) SetPKIDirectory(filepath string) {
	w.pkiDirPath = filepath
}

func (w *PKIWatcher) open() {
	for {
		select {
		case <-w.ctx.Done():
			return
		default:
		}

		minionsDir := path.Join(w.pkiDirPath, "minions")

		log.Info().Msg("loading currently accepted minions")
		entries, err := os.ReadDir(minionsDir)
		if err != nil {
			log.Error().Str("error", err.Error()).Msg("failed to list PKI directory")
			time.Sleep(5 * time.Second)
		} else {
			for _, e := range entries {
				if !e.IsDir() {
					w.eventChan <- event.WatchEvent{
						MinionName: e.Name(),
						Op:         event.Accepted,
					}
					log.Info().Msgf("minion %s loaded", e.Name())
				}
			}

			// Add a path.
			err = w.watcher.Add(minionsDir)
			if err != nil {
				log.Error().Str("error", err.Error()).Msg("failed to watch PKI directory")
				time.Sleep(time.Second * 5)
			} else {
				return
			}
		}
	}
}

func (w *PKIWatcher) StartWatching() {
	w.open()

	for {
		select {
		case <-w.ctx.Done():
			w.Stop()
			return
		case evt := <-w.watcher.Events:
			minionName := path.Base(evt.Name)
			if minionName == ".key_cache" || strings.HasPrefix(minionName, ".___atomic_write") {
				continue
			}
			if evt.Op == fsnotify.Create {
				w.eventChan <- event.WatchEvent{
					MinionName: minionName,
					Op:         event.Accepted,
				}
				log.Info().Msgf("minion %s accepted by master", minionName)
			}
			if evt.Op == fsnotify.Remove {
				w.eventChan <- event.WatchEvent{
					MinionName: minionName,
					Op:         event.Removed,
				}
				log.Info().Msgf("minion %s removed from master", minionName)
			}
		case err := <-w.watcher.Errors:
			log.Error().Str("error", err.Error()).Msg("fail processing watch event")
		}
	}
}

func (w *PKIWatcher) Stop() {
	log.Info().Msg("stop listening for PKI changes")
	w.watcher.Close()
}
