package levels

import (
	"reflect"

	"github.com/inkyblackness/hacked/editor/event"
	"github.com/inkyblackness/hacked/ss1/content/archive/level"
)

type objectIDs struct {
	list []level.ObjectID
}

func (coords objectIDs) contains(pos level.ObjectID) bool {
	for _, entry := range coords.list {
		if entry == pos {
			return true
		}
	}
	return false
}

func (coords *objectIDs) registerAt(registry event.Registry) {
	var setEvent ObjectSelectionSetEvent
	registry.RegisterHandler(reflect.TypeOf(setEvent), coords.onObjectSelectionSetEvent)
	var addEvent ObjectSelectionAddEvent
	registry.RegisterHandler(reflect.TypeOf(addEvent), coords.onObjectSelectionAddEvent)
	var removeEvent ObjectSelectionRemoveEvent
	registry.RegisterHandler(reflect.TypeOf(removeEvent), coords.onObjectSelectionRemoveEvent)
}

func (coords *objectIDs) onObjectSelectionSetEvent(evt ObjectSelectionSetEvent) {
	coords.list = evt.objects
}

func (coords *objectIDs) onObjectSelectionAddEvent(evt ObjectSelectionAddEvent) {
	coords.list = append(coords.list, evt.objects...)
}

func (coords *objectIDs) onObjectSelectionRemoveEvent(evt ObjectSelectionRemoveEvent) {
	newList := make([]level.ObjectID, 0, len(coords.list))
	for _, oldEntry := range coords.list {
		keep := true
		for _, removedEntry := range evt.objects {
			if oldEntry == removedEntry {
				keep = false
			}
		}
		if keep {
			newList = append(newList, oldEntry)
		}
	}
	coords.list = newList
}