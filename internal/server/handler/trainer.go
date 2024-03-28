package handler

import (
	"net/http"
)

func (h *handler) Read(w http.ResponseWriter, r *http.Request) {}

func (h *handler) Create(w http.ResponseWriter, r *http.Request) {}

func (h *handler) Update(w http.ResponseWriter, r *http.Request) {}

func (h *handler) Delete(w http.ResponseWriter, r *http.Request) {}
