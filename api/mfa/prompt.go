/*
Copyright 2023 Gravitational, Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package mfa

import (
	"context"
	"fmt"

	"github.com/gravitational/teleport/api/client/proto"
)

// Prompt is an MFA prompt.
type Prompt interface {
	// Run prompts the user to complete an MFA authentication challenge.
	Run(ctx context.Context, chal *proto.MFAAuthenticateChallenge) (*proto.MFAAuthenticateResponse, error)
}

// PromptFunc is a function wrapper that implements the Prompt interface.
type PromptFunc func(ctx context.Context, chal *proto.MFAAuthenticateChallenge) (*proto.MFAAuthenticateResponse, error)

// Run prompts the user to complete an MFA authentication challenge.
func (f PromptFunc) Run(ctx context.Context, chal *proto.MFAAuthenticateChallenge) (*proto.MFAAuthenticateResponse, error) {
	return f(ctx, chal)
}

// PromptConstructor is a function that creates a new MFA prompt.
type PromptConstructor func(...PromptOpt) Prompt

// PromptConfig contains common mfa prompt config options.
type PromptConfig struct {
	// PromptReason is an optional message to share with the user before an MFA Prompt.
	// It is intended to provide context about why the user is being prompted where it may
	// not be obvious, such as for admin actions or per-session MFA.
	PromptReason string
	// DeviceType is an optional device description to emphasize during the prompt.
	DeviceType DeviceDescriptor
	// Quiet suppresses users prompts.
	Quiet bool
}

// DeviceDescriptor is a descriptor for a device, such as "registered".
type DeviceDescriptor string

// DeviceDescriptorRegistered is a registered device.
const DeviceDescriptorRegistered = "registered"

// PromptOpt applies configuration options to a prompt.
type PromptOpt func(*PromptConfig)

// WithQuiet sets the prompt's Quiet field.
func WithQuiet() PromptOpt {
	return func(cfg *PromptConfig) {
		cfg.Quiet = true
	}
}

// WithPromptReason sets the prompt's PromptReason field.
func WithPromptReason(hint string) PromptOpt {
	return func(cfg *PromptConfig) {
		cfg.PromptReason = hint
	}
}

// WithPromptReasonAdminAction sets the prompt's PromptReason field to a standard admin action message.
func WithPromptReasonAdminAction(actionName string) PromptOpt {
	adminMFAPromptReason := fmt.Sprintf("MFA is required for admin-level API request: %q", actionName)
	return WithPromptReason(adminMFAPromptReason)
}

// WithPromptReasonSessionMFA sets the prompt's PromptReason field to a standard session mfa message.
func WithPromptReasonSessionMFA(serviceType, serviceName string) PromptOpt {
	return WithPromptReason(fmt.Sprintf("MFA is required to access %s %q", serviceType, serviceName))
}

// WithPromptDeviceType sets the prompt's DeviceType field.
func WithPromptDeviceType(deviceType DeviceDescriptor) PromptOpt {
	return func(cfg *PromptConfig) {
		cfg.DeviceType = deviceType
	}
}
