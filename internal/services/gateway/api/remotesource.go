// Copyright 2019 Sorint.lab
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied
// See the License for the specific language governing permissions and
// limitations under the License.

package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"agola.io/agola/internal/services/gateway/action"
	"agola.io/agola/internal/services/types"
	"agola.io/agola/internal/util"
	"go.uber.org/zap"
	errors "golang.org/x/xerrors"

	"github.com/gorilla/mux"
)

type CreateRemoteSourceRequest struct {
	Name                string `json:"name"`
	APIURL              string `json:"apiurl"`
	Type                string `json:"type"`
	AuthType            string `json:"auth_type"`
	SkipVerify          bool   `json:"skip_verify"`
	Oauth2ClientID      string `json:"oauth_2_client_id"`
	Oauth2ClientSecret  string `json:"oauth_2_client_secret"`
	SSHHostKey          string `json:"ssh_host_key"`
	SkipSSHHostKeyCheck bool   `json:"skip_ssh_host_key_check"`
	RegistrationEnabled *bool  `json:"registration_enabled"`
	LoginEnabled        *bool  `json:"login_enabled"`
}

type CreateRemoteSourceHandler struct {
	log *zap.SugaredLogger
	ah  *action.ActionHandler
}

func NewCreateRemoteSourceHandler(logger *zap.Logger, ah *action.ActionHandler) *CreateRemoteSourceHandler {
	return &CreateRemoteSourceHandler{log: logger.Sugar(), ah: ah}
}

func (h *CreateRemoteSourceHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req CreateRemoteSourceRequest
	d := json.NewDecoder(r.Body)
	if err := d.Decode(&req); err != nil {
		httpError(w, util.NewErrBadRequest(err))
		return
	}

	creq := &action.CreateRemoteSourceRequest{
		Name:                req.Name,
		APIURL:              req.APIURL,
		Type:                req.Type,
		AuthType:            req.AuthType,
		SkipVerify:          req.SkipVerify,
		Oauth2ClientID:      req.Oauth2ClientID,
		Oauth2ClientSecret:  req.Oauth2ClientSecret,
		SSHHostKey:          req.SSHHostKey,
		SkipSSHHostKeyCheck: req.SkipSSHHostKeyCheck,
		RegistrationEnabled: req.RegistrationEnabled,
		LoginEnabled:        req.LoginEnabled,
	}
	rs, err := h.ah.CreateRemoteSource(ctx, creq)
	if httpError(w, err) {
		h.log.Errorf("err: %+v", err)
		return
	}

	res := createRemoteSourceResponse(rs)
	if err := httpResponse(w, http.StatusCreated, res); err != nil {
		h.log.Errorf("err: %+v", err)
	}
}

type UpdateRemoteSourceRequest struct {
	Name                *string `json:"name"`
	APIURL              *string `json:"apiurl"`
	SkipVerify          *bool   `json:"skip_verify"`
	Oauth2ClientID      *string `json:"oauth_2_client_id"`
	Oauth2ClientSecret  *string `json:"oauth_2_client_secret"`
	SSHHostKey          *string `json:"ssh_host_key"`
	SkipSSHHostKeyCheck *bool   `json:"skip_ssh_host_key_check"`
	RegistrationEnabled *bool   `json:"registration_enabled"`
	LoginEnabled        *bool   `json:"login_enabled"`
}

type UpdateRemoteSourceHandler struct {
	log *zap.SugaredLogger
	ah  *action.ActionHandler
}

func NewUpdateRemoteSourceHandler(logger *zap.Logger, ah *action.ActionHandler) *UpdateRemoteSourceHandler {
	return &UpdateRemoteSourceHandler{log: logger.Sugar(), ah: ah}
}

func (h *UpdateRemoteSourceHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	rsRef := vars["remotesourceref"]

	var req UpdateRemoteSourceRequest
	d := json.NewDecoder(r.Body)
	if err := d.Decode(&req); err != nil {
		httpError(w, util.NewErrBadRequest(err))
		return
	}

	creq := &action.UpdateRemoteSourceRequest{
		RemoteSourceRef: rsRef,

		Name:                req.Name,
		APIURL:              req.APIURL,
		SkipVerify:          req.SkipVerify,
		Oauth2ClientID:      req.Oauth2ClientID,
		Oauth2ClientSecret:  req.Oauth2ClientSecret,
		SSHHostKey:          req.SSHHostKey,
		SkipSSHHostKeyCheck: req.SkipSSHHostKeyCheck,
		RegistrationEnabled: req.RegistrationEnabled,
		LoginEnabled:        req.LoginEnabled,
	}
	rs, err := h.ah.UpdateRemoteSource(ctx, creq)
	if httpError(w, err) {
		h.log.Errorf("err: %+v", err)
		return
	}

	res := createRemoteSourceResponse(rs)
	if err := httpResponse(w, http.StatusCreated, res); err != nil {
		h.log.Errorf("err: %+v", err)
	}
}

type RemoteSourceResponse struct {
	ID                  string `json:"id"`
	Name                string `json:"name"`
	AuthType            string `json:"auth_type"`
	RegistrationEnabled bool   `json:"registration_enabled"`
	LoginEnabled        bool   `json:"login_enabled"`
}

func createRemoteSourceResponse(r *types.RemoteSource) *RemoteSourceResponse {
	rs := &RemoteSourceResponse{
		ID:                  r.ID,
		Name:                r.Name,
		AuthType:            string(r.AuthType),
		RegistrationEnabled: *r.RegistrationEnabled,
		LoginEnabled:        *r.LoginEnabled,
	}
	return rs
}

type RemoteSourceHandler struct {
	log *zap.SugaredLogger
	ah  *action.ActionHandler
}

func NewRemoteSourceHandler(logger *zap.Logger, ah *action.ActionHandler) *RemoteSourceHandler {
	return &RemoteSourceHandler{log: logger.Sugar(), ah: ah}
}

func (h *RemoteSourceHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	rsRef := vars["remotesourceref"]

	rs, err := h.ah.GetRemoteSource(ctx, rsRef)
	if httpError(w, err) {
		h.log.Errorf("err: %+v", err)
		return
	}

	res := createRemoteSourceResponse(rs)
	if err := httpResponse(w, http.StatusOK, res); err != nil {
		h.log.Errorf("err: %+v", err)
	}
}

type RemoteSourcesHandler struct {
	log *zap.SugaredLogger
	ah  *action.ActionHandler
}

func NewRemoteSourcesHandler(logger *zap.Logger, ah *action.ActionHandler) *RemoteSourcesHandler {
	return &RemoteSourcesHandler{log: logger.Sugar(), ah: ah}
}

func (h *RemoteSourcesHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	query := r.URL.Query()

	limitS := query.Get("limit")
	limit := DefaultRunsLimit
	if limitS != "" {
		var err error
		limit, err = strconv.Atoi(limitS)
		if err != nil {
			httpError(w, util.NewErrBadRequest(errors.Errorf("cannot parse limit: %w", err)))
			return
		}
	}
	if limit < 0 {
		httpError(w, util.NewErrBadRequest(errors.Errorf("limit must be greater or equal than 0")))
		return
	}
	if limit > MaxRunsLimit {
		limit = MaxRunsLimit
	}
	asc := false
	if _, ok := query["asc"]; ok {
		asc = true
	}

	start := query.Get("start")

	areq := &action.GetRemoteSourcesRequest{
		Start: start,
		Limit: limit,
		Asc:   asc,
	}
	csRemoteSources, err := h.ah.GetRemoteSources(ctx, areq)
	if httpError(w, err) {
		h.log.Errorf("err: %+v", err)
		return
	}

	remoteSources := make([]*RemoteSourceResponse, len(csRemoteSources))
	for i, rs := range csRemoteSources {
		remoteSources[i] = createRemoteSourceResponse(rs)
	}

	if err := httpResponse(w, http.StatusOK, remoteSources); err != nil {
		h.log.Errorf("err: %+v", err)
	}
}

type DeleteRemoteSourceHandler struct {
	log *zap.SugaredLogger
	ah  *action.ActionHandler
}

func NewDeleteRemoteSourceHandler(logger *zap.Logger, ah *action.ActionHandler) *DeleteRemoteSourceHandler {
	return &DeleteRemoteSourceHandler{log: logger.Sugar(), ah: ah}
}

func (h *DeleteRemoteSourceHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	rsRef := vars["remotesourceref"]

	err := h.ah.DeleteRemoteSource(ctx, rsRef)
	if httpError(w, err) {
		h.log.Errorf("err: %+v", err)
		return
	}

	if err := httpResponse(w, http.StatusNoContent, nil); err != nil {
		h.log.Errorf("err: %+v", err)
	}
}
