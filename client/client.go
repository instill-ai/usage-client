package client

import (
	"context"

	"github.com/instill-ai/usage-client/reporter"

	usagePB "github.com/instill-ai/protogen-go/vdp/usage/v1alpha"
)

// StartReporter creates a usage reporter and start sending usage data to server regularly
// retrieveUsageData is a function that outputs any of the type:
//	*usagePB.SessionReport_MgmtUsageData
//	*usagePB.SessionReport_ConnectorUsageData
//	*usagePB.SessionReport_ModelUsageData
//	*usagePB.SessionReport_PipelineUsageData
func StartReporter(ctx context.Context, usageClient usagePB.UsageServiceClient, sessionService usagePB.Session_Service, url, edition, version string, retrieveUsageData func() interface{}) error {
	reporter, err := reporter.NewReporter(ctx, usageClient, sessionService, edition, version)
	if err != nil {
		return err
	}

	go reporter.Report(ctx, sessionService, edition, version, retrieveUsageData)

	return nil
}
