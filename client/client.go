package client

import (
	"context"

	"github.com/instill-ai/usage-client/reporter"

	usagePB "github.com/instill-ai/protogen-go/base/usage/v1alpha"
)

// InitReporter creates a usage reporter
func InitReporter(ctx context.Context, usageClient usagePB.UsageServiceClient, sessionService usagePB.Session_Service, edition, version string, defaultOwnerUid string) (reporter.Reporter, error) {
	reporter, err := reporter.NewReporter(ctx, usageClient, sessionService, edition, version, defaultOwnerUid)
	if err != nil {
		return nil, err
	}
	return reporter, nil
}

// StartReporter uses a usage reporter to start sending usage data to server regularly
// retrieveUsageData is a function that outputs any of the type:
//
//	*usagePB.SessionReport_MgmtUsageData
//	*usagePB.SessionReport_ConnectorUsageData
//	*usagePB.SessionReport_ModelUsageData
//	*usagePB.SessionReport_PipelineUsageData
func StartReporter(ctx context.Context, reporter reporter.Reporter, sessionService usagePB.Session_Service, edition, version string, ownerUid string, retrieveUsageData func() interface{}) error {
	go reporter.Report(ctx, sessionService, edition, version, ownerUid, retrieveUsageData)

	return nil
}

// SingleReporter uses a usage reporter and sends one-time usage data to server
func SingleReporter(ctx context.Context, reporter reporter.Reporter, sessionService usagePB.Session_Service, edition, version string, ownerUid string, usageData interface{}) error {
	err := reporter.SingleReport(ctx, sessionService, edition, version, ownerUid, usageData)
	if err != nil {
		return err
	}
	return nil
}
