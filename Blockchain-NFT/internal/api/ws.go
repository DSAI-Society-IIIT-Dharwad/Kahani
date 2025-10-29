package api

import (
	"net/http"
)

func (a *API) handleEvents(w http.ResponseWriter, r *http.Request) {
	if a.observer == nil {
		writeError(w, http.StatusServiceUnavailable, "event streaming not configured")
		return
	}

	conn, err := a.upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	defer conn.Close()

	id, ch := a.observer.Subscribe(32)
	if id == "" {
		return
	}
	defer a.observer.Unsubscribe(id)

	// Notify client that connection is established.
	_ = conn.WriteJSON(map[string]string{"type": "connection.ready"})

	for {
		select {
		case <-r.Context().Done():
			return
		case event, ok := <-ch:
			if !ok {
				return
			}

			if err := conn.WriteJSON(event); err != nil {
				return
			}
		}
	}
}
