// mautrix-gvoice - A Matrix-Google Voice puppeting bridge.
// Copyright (C) 2024 Tulir Asokan
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

package connector

import (
	"context"
	"crypto/sha256"
	_ "embed"
	"encoding/base64"
	"encoding/json"
	"errors"
	"io"
	"math/rand/v2"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"sync/atomic"
	"time"

	"github.com/rs/zerolog"
	"go.mau.fi/util/exsync"
)

//go:embed puppeteer.js
var puppeteerScript string

type requestSignatureFunc func(ctx context.Context, threadID string, recipients []string, txnID int64) (string, error)

type puppeteerResponse struct {
	err  error
	resp string
}

func (gc *GVClient) runPuppeteer(ctx context.Context) {
	log := gc.UserLogin.Log.With().Str("component", "puppeteer").Logger()
	ctx = log.WithContext(ctx)
	node, _ := exec.LookPath("node")
	if node == "" {
		log.Debug().Msg("Node.js not installed")
		return
	}
	tmp, err := os.CreateTemp("", "mautrix-gvoice-puppeteer-*.js")
	if err != nil {
		log.Err(err).Msg("Failed to create temporary file")
		return
	}
	defer os.Remove(tmp.Name())
	_, err = tmp.WriteString(puppeteerScript)
	if err != nil {
		_ = tmp.Close()
		log.Err(err).Msg("Failed to write puppeteer script")
		return
	}
	_ = tmp.Close()
	cmd := exec.CommandContext(ctx, node, tmp.Name())
	// Set environment variable for debug mode if needed
	if os.Getenv("MAUTRIX_GVOICE_PUPPETEER_DEBUG") == "true" {
		cmd.Env = append(os.Environ(), "MAUTRIX_GVOICE_PUPPETEER_DEBUG=true")
	}
	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		log.Err(err).Msg("Failed to get stdout pipe")
		return
	}
	defer stdoutPipe.Close()
	stdinPipe, err := cmd.StdinPipe()
	if err != nil {
		log.Err(err).Msg("Failed to get stdin pipe")
		return
	}
	defer stdinPipe.Close()
	cmd.Stderr = os.Stderr // Redirect stderr to standard error
	gc.stopWait.Add(1)
	defer gc.stopWait.Done()
	err = cmd.Start()
	if err != nil {
		log.Err(err).Msg("Failed to start Puppeteer")
		return
	}
	kill := func() {
		if proc := cmd.Process; proc != nil {
			_ = proc.Kill()
		}
	}
	stdinJSON := json.NewEncoder(stdinPipe)
	stdoutJSON := json.NewDecoder(stdoutPipe)
	responseWaiters := exsync.NewMap[string, chan<- puppeteerResponse]()
	var program, globalName string
	var reqIDCounter atomic.Int64
	requestSignatureDirect := func(ctx context.Context, payload map[string]any) (string, error) {
		reqID := strconv.FormatInt(reqIDCounter.Add(1), 10)
		waiter := make(chan puppeteerResponse, 1)
		responseWaiters.Set(reqID, waiter)
		defer responseWaiters.Delete(reqID)
		defer close(waiter)
		zerolog.Ctx(ctx).Debug().
			Str("req_id", reqID).
			Interface("payload", payload).
			Msg("Requesting signature for message")
		err = stdinJSON.Encode(map[string]any{
			"req_id":      reqID,
			"global_name": globalName,
			"program":     program,
			"payload":     payload,
		})
		if err != nil {
			return "", err
		}
		select {
		case result := <-waiter:
			if result.err != nil {
				zerolog.Ctx(ctx).Warn().Str("req_id", reqID).Err(result.err).Msg("Received error")
				return "", result.err
			} else if result.resp == "" {
				zerolog.Ctx(ctx).Warn().Str("req_id", reqID).Msg("Received empty signature")
				return "", errors.New("empty signature")
			} else {
				zerolog.Ctx(ctx).Debug().Str("req_id", reqID).Msg("Received signature")
				return result.resp, nil
			}
		case <-time.After(10 * time.Second):
			zerolog.Ctx(ctx).Warn().Str("req_id", reqID).Msg("Timed out waiting for signature")
			return "", errors.New("request timed out")
		case <-ctx.Done():
			zerolog.Ctx(ctx).Warn().Str("req_id", reqID).Err(err).Msg("Context canceled while waiting for signature")
			return "", ctx.Err()
		}
	}
	requestSignature := func(ctx context.Context, threadID string, recipients []string, txnID int64) (string, error) {
		threadIDHash := sha256.Sum256([]byte(threadID))
		recipientsHash := sha256.Sum256([]byte(strings.Join(recipients, "|")))
		messageIDHash := sha256.Sum256([]byte(strconv.FormatInt(txnID, 10)))
		payload := map[string]any{
			"thread_id":    base64.StdEncoding.EncodeToString(threadIDHash[:]),
			"destinations": base64.StdEncoding.EncodeToString(recipientsHash[:]),
			"message_ids":  base64.StdEncoding.EncodeToString(messageIDHash[:]),
		}
		return requestSignatureDirect(ctx, payload)
	}
	funcPtr := (*requestSignatureFunc)(&requestSignature)
	gc.requestSignature.Store(funcPtr)
	log.Debug().Msg("Puppeteer started")
Loop:
	for {
		var payload map[string]string
		err = stdoutJSON.Decode(&payload)
		if errors.Is(err, io.EOF) {
			log.Debug().Msg("Puppeteer stdout closed")
			break
		} else if err != nil {
			log.Err(err).Msg("Failed to decode Puppeteer payload")
			kill()
			break
		}
		switch payload["status"] {
		case "waiting_for_init":
			log.Debug().Msg("Creating Waa payload for Puppeteer")
			waa, err := gc.Client.CreateWaa(ctx)
			if err != nil {
				log.Err(err).Msg("Failed to create Waa payload")
				kill()
				break Loop
			}
			program = waa.GetProgram()
			globalName = waa.GetGlobalName()
			hash, _ := base64.RawURLEncoding.DecodeString(waa.GetInterpreterHash())
			checksum := base64.StdEncoding.EncodeToString(hash)
			log.Debug().
				Str("script_source", waa.GetInterpreterURL().GetURL()).
				Str("global_name", globalName).
				Str("script_checksum", checksum).
				Msg("Waa payload created")
			err = stdinJSON.Encode(map[string]string{
				"script_source": waa.GetInterpreterURL().GetURL(),
				"checksum":      checksum,
			})
			if err != nil {
				log.Err(err).Msg("Failed to send Waa init payload")
				kill()
				break Loop
			}
		case "ready":
			log.Debug().Msg("Puppeteer Waa generator is ready, sending ping")
			go func() {
				pl, err := requestSignatureDirect(ctx, map[string]any{"blank_payload": true})
				if err != nil {
					log.Err(err).Msg("Failed to request signature for ping")
				} else if err = gc.Client.PingWaa(ctx, pl, rand.Int63n(2000000000)); err != nil {
					log.Err(err).Msg("Failed to send ping")
				} else {
					log.Info().Msg("Waa ping successful")
				}
			}()
		case "result":
			waiter, ok := responseWaiters.Pop(payload["req_id"])
			if !ok {
				log.Warn().Str("req_id", payload["req_id"]).Msg("Unknown response from Puppeteer")
			} else {
				select {
				case waiter <- puppeteerResponse{resp: payload["response"]}:
				default:
					log.Warn().Str("req_id", payload["req_id"]).Msg("Response channel didn't accept result")
				}
			}
		case "error":
			waiter, ok := responseWaiters.Pop(payload["req_id"])
			if !ok {
				log.Warn().Any("payload", payload).Msg("Unknown error response from Puppeteer")
			} else {
				select {
				case waiter <- puppeteerResponse{err: errors.New(payload["error"])}:
				default:
					log.Warn().Str("req_id", payload["req_id"]).Msg("Response channel didn't accept error result")
				}
			}
		default:
			log.Warn().Any("data", payload).Msg("Unknown payload from Puppeteer")
		}
	}
	_ = os.Remove(tmp.Name())
	gc.requestSignature.CompareAndSwap(funcPtr, nil)
	err = cmd.Wait()
	if err != nil && !errors.Is(err, context.Canceled) && ctx.Err() == nil {
		log.Err(err).Msg("Puppeteer exited with error")
	} else {
		log.Debug().Msg("Puppeteer exited")
	}
}
