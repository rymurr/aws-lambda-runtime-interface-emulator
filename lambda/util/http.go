// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package util

import (
	"net/http"

	"github.com/rymurr/aws-lambda-runtime-interface-emulator/lambda/interop"
	"github.com/rymurr/aws-lambda-runtime-interface-emulator/lambda/rapidcore"
	log "github.com/sirupsen/logrus"
)

func StartHTTPServer(ipport string, sandbox *rapidcore.SandboxBuilder, bs interop.Bootstrap) {
	srv := &http.Server{
		Addr: ipport,
	}

	// Pass a channel
	http.HandleFunc("/2015-03-31/functions/function/invocations", func(w http.ResponseWriter, r *http.Request) {
		InvokeHandler(w, r, sandbox.LambdaInvokeAPI(), bs)
	})

	// go routine (main thread waits)
	if err := srv.ListenAndServe(); err != nil {
		log.Panic(err)
	}

	log.Warnf("Listening on %s", ipport)
}
