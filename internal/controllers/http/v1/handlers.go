package v1

import (
	"log/slog"
	"net/http"
	"strconv"
)

func Auth(log *slog.Logger, userAuther UserAuther) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const fn = "http.v1.api.Auth"

		log = log.With(
			slog.String("fn", fn),
		)

		// Get Qury params
		token := r.URL.Query().Get("token")

		messengerType := r.URL.Query().Get("mtype")
		if messengerType == "" {
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}

		messengerID, err := strconv.Atoi(r.URL.Query().Get("mid"))
		if err != nil {
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}

		chatID, err := strconv.Atoi(r.URL.Query().Get("chatid"))
		if err != nil {
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}

		redirectURL := r.URL.Query().Get("redirect")
		if redirectURL == "" {
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}

		// Activate Account
		err = userAuther.AccountActivation(messengerType, int64(messengerID), int64(chatID), token)
		if err != nil {
			http.Error(w, "Bad Request", http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, redirectURL, http.StatusSeeOther)
	}
}
