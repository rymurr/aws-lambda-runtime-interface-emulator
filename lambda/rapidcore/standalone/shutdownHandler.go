// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package standalone

import (
	"context"
	"net/http"

	"github.com/rymurr/aws-lambda-runtime-interface-emulator/lambda/interop"
	"github.com/rymurr/aws-lambda-runtime-interface-emulator/lambda/metering"
)

type shutdownAPIRequest struct {
	TimeoutMs int64 `json:"timeoutMs"`
}

func ShutdownHandler(w http.ResponseWriter, r *http.Request, s InteropServer, shutdownFunc context.CancelFunc) {
	shutdown := shutdownAPIRequest{}
	if lerr := readBodyAndUnmarshalJSON(r, &shutdown); lerr != nil {
		lerr.Send(w, r)
		return
	}

	internalState := s.Shutdown(&interop.Shutdown{
		DeadlineNs: metering.Monotime() + int64(shutdown.TimeoutMs*1000*1000),
	})

	w.Write(internalState.AsJSON())

	shutdownFunc()
}
