package lib

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/hashicorp/vault/api"
)

const obscurusContentKey = "obscurus"

var (
	plaintextContentHeader = middleware.SetHeader("Content-Type", "text/plain; charset=utf-8")
	jsonContentHeader      = middleware.SetHeader("Content-Type", "application/json")
)

// ObscurusHandler is an HTTP Handler
type ObscurusHandler struct {
	http.Handler
	Vault *api.Client
}

// NewHandler returns an HTTP Handler for Obscurus
func NewHandler(v *api.Client) http.Handler {
	r := chi.NewRouter()
	r.Use(
		middleware.RequestID,
		middleware.RealIP,
		middleware.Logger,
		middleware.Recoverer,
		middleware.Timeout(60*time.Second),
	)

	h := &ObscurusHandler{
		Handler: r,
		Vault:   v,
	}

	r.Handle("/", http.FileServer(http.Dir("./public")))
	r.Route("/api/v1", func(r chi.Router) {
		r.With(jsonContentHeader).Post("/secrets", h.apiSecretWrap)
		r.With(jsonContentHeader).Get("/secrets/{token}", h.apiSecretLookup)
		r.With(plaintextContentHeader).Get("/secrets/{token}/value", h.apiSecretUnwrap)
	})

	return h
}

func (h ObscurusHandler) apiSecretWrap(w http.ResponseWriter, r *http.Request) {
	b, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	s, err := h.Vault.Logical().Write("sys/wrapping/wrap", map[string]interface{}{
		obscurusContentKey: string(b),
	})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"token": s.WrapInfo.Token,
		"ttl":   s.WrapInfo.TTL,
	})
}

func (h ObscurusHandler) apiSecretLookup(w http.ResponseWriter, r *http.Request) {
	s, err := h.Vault.Logical().Write("sys/wrapping/lookup", map[string]interface{}{
		"token": chi.URLParam(r, "token"),
	})

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	json.NewEncoder(w).Encode(s.Data)
}

func (h ObscurusHandler) apiSecretUnwrap(w http.ResponseWriter, r *http.Request) {
	s, err := h.Vault.Logical().Unwrap(chi.URLParam(r, "token"))
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	fmt.Fprint(w, s.Data[obscurusContentKey].(string))
}
