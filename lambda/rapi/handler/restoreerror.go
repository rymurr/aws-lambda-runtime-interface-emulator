// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package handler

import (
	"net/http"

	log "github.com/sirupsen/logrus"
	"github.com/rymurr/aws-lambda-runtime-interface-emulator/lambda/appctx"
	"github.com/rymurr/aws-lambda-runtime-interface-emulator/lambda/core"
	"github.com/rymurr/aws-lambda-runtime-interface-emulator/lambda/fatalerror"
	"github.com/rymurr/aws-lambda-runtime-interface-emulator/lambda/interop"
	"github.com/rymurr/aws-lambda-runtime-interface-emulator/lambda/rapi/rendering"
)

type restoreErrorHandler struct {
	registrationService core.RegistrationService
}

func (h *restoreErrorHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	appCtx := appctx.FromRequest(request)
	server := appctx.LoadInteropServer(appCtx)
	if server == nil {
		log.Panic("Invalid state, cannot access interop server")
	}

	errorType := fatalerror.GetValidRuntimeOrFunctionErrorType(request.Header.Get("Lambda-Runtime-Function-Error-Type"))
	fnError := interop.FunctionError{Type: errorType}

	runtime := h.registrationService.GetRuntime()

	if err := runtime.RestoreError(fnError); err != nil {
		log.Warn(err)
		rendering.RenderForbiddenWithTypeMsg(writer, request, rendering.ErrorTypeInvalidStateTransition, StateTransitionFailedForRuntimeMessageFormat,
			runtime.GetState().Name(), core.RuntimeRestoreErrorStateName, err)
		return
	}

	appctx.StoreInvokeErrorTraceData(appCtx, &interop.InvokeErrorTraceData{})

	rendering.RenderAccepted(writer, request)
}

func NewRestoreErrorHandler(registrationService core.RegistrationService) http.Handler {
	return &restoreErrorHandler{registrationService: registrationService}
}
