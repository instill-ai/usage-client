package usage

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"runtime"
	"time"

	"github.com/catalinc/hashcash"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/instill-ai/usage-client/internal/logger"

	usagePB "github.com/instill-ai/protogen-go/vdp/usage/v1alpha"
)

const (
	// HashBits is Number of zero bits
	HashBits = 20
	// SaltLen is Random salt length
	SaltLen = 40
	// DefaultExtension Extension to add to the minted stamp
	DefaultExtension = ""
	// timeout
	timeout = 15 * time.Second
)

var (
	hasher          = hashcash.New(HashBits, SaltLen, DefaultExtension)
	reportFrequency = 1 * time.Minute
)

// Session is the alias of session
type Session usagePB.Session

// Hash converts the usage data into base64 encoded checksum
func (s *Session) Hash() (string, error) {
	buf := new(bytes.Buffer)
	err := json.NewEncoder(buf).Encode(s)
	if err != nil {
		return "", err
	}
	sum := sha256.Sum256(buf.Bytes())
	return base64.URLEncoding.EncodeToString(sum[:]), nil
}

// Expired checks whether the usage data is expired
func (s *Session) Expired() bool {
	return time.Since(s.ReportTime.AsTime()) > time.Minute
}

// Minter a interface that creates Hashcash stamp
type Minter interface {
	Mint(string) (string, error)
}

// Reporter interface
type Reporter interface {
	// SingleReport represents send one report to the usage server
	// Types that are assignable to usageData:
	//	*usagePB.SessionReport_MgmtUsageData
	//	*usagePB.SessionReport_ConnectorUsageData
	//	*usagePB.SessionReport_ModelUsageData
	//	*usagePB.SessionReport_PipelineUsageData
	SingleReport(ctx context.Context, service usagePB.Session_Service, edition, version string, usageData interface{}) error
	// Report sends report to the server regularly based on the report frequency
	// retrieveUsageData is a function that outputs any of the type:
	//	*usagePB.SessionReport_MgmtUsageData
	//	*usagePB.SessionReport_ConnectorUsageData
	//	*usagePB.SessionReport_ModelUsageData
	//	*usagePB.SessionReport_PipelineUsageData
	Report(ctx context.Context, service usagePB.Session_Service, edition, version string, retrieveUsageData func() interface{})
}

// reporter represents a reporter that sends usage data to the server on a regular basis
type reporter struct {
	client     usagePB.UsageServiceClient
	sessionUID string
	start      time.Time
	minter     Minter
	token      string
}

// NewReporter creates a new usage reporter
func NewReporter(ctx context.Context, client usagePB.UsageServiceClient, service usagePB.Session_Service, edition, version string) (Reporter, error) {

	// Create the session
	resp, err := client.CreateSession(ctx,
		&usagePB.CreateSessionRequest{
			Session: &usagePB.Session{
				Service:    service,
				Edition:    edition,
				Version:    version,
				Arch:       runtime.GOARCH,
				Os:         runtime.GOOS,
				Uptime:     0,
				ReportTime: timestamppb.New(time.Now()),
			},
		})
	if err != nil {
		return nil, err
	}

	// Validation: token
	token := resp.GetSession().GetToken()
	if token == "" {
		return nil, errors.New("invalid empty token. New session creation failed, no token")
	}

	r := &reporter{
		client:     client,
		sessionUID: resp.GetSession().GetUid(),
		start:      time.Now(),
		minter:     hasher,
		token:      token,
	}

	return r, nil
}

// SingleReport represents send one report to the usage server
func (r *reporter) SingleReport(ctx context.Context, service usagePB.Session_Service, edition, version string, usageData interface{}) error {
	s := Session{
		Service:    service,
		Edition:    edition,
		Version:    version,
		Arch:       runtime.GOARCH,
		Os:         runtime.GOOS,
		Uptime:     int64(time.Since(r.start).Truncate(time.Second).Seconds()),
		ReportTime: timestamppb.New(time.Now()),
	}

	// Generate Proof-of-Word (PoW) with token + hash of session data
	base, err := s.Hash()
	if err != nil {
		return err
	}
	pow, err := r.minter.Mint(r.token + base)
	if err != nil {
		return err
	}

	report := usagePB.SessionReport{
		SessionUid: r.sessionUID,
		Token:      r.token,
		Pow:        pow,
	}

	// Session report
	pbSessionData := usagePB.Session(s)
	report.Session = &pbSessionData
	// Usage data
	invalidUsageDataErr := errors.New("invalid usage data type")
	switch service {
	case usagePB.Session_SERVICE_MGMT:
		if ud, ok := usageData.(*usagePB.SessionReport_MgmtUsageData); ok {
			pbUsageData := usagePB.SessionReport_MgmtUsageData(*ud)
			report.UsageData = &pbUsageData
		} else {
			return fmt.Errorf("[mgmt-backend] %v", invalidUsageDataErr)
		}
	case usagePB.Session_SERVICE_CONNECTOR:
		if ud, ok := usageData.(*usagePB.SessionReport_ConnectorUsageData); ok {
			pbUsageData := usagePB.SessionReport_ConnectorUsageData(*ud)
			report.UsageData = &pbUsageData
		} else {
			return fmt.Errorf("[connector-backend] %v", invalidUsageDataErr)
		}
	case usagePB.Session_SERVICE_MODEL:
		if ud, ok := usageData.(*usagePB.SessionReport_ModelUsageData); ok {
			pbUsageData := usagePB.SessionReport_ModelUsageData(*ud)
			report.UsageData = &pbUsageData
		} else {
			return fmt.Errorf("[model-backend] %v", invalidUsageDataErr)
		}
	case usagePB.Session_SERVICE_PIPELINE:
		if ud, ok := usageData.(*usagePB.SessionReport_PipelineUsageData); ok {
			pbUsageData := usagePB.SessionReport_PipelineUsageData(*ud)
			report.UsageData = &pbUsageData
		} else {
			return fmt.Errorf("[pipeline-backend] %v", invalidUsageDataErr)
		}
	default:
		return invalidUsageDataErr
	}

	if _, err = r.client.SendSessionReport(ctx, &usagePB.SendSessionReportRequest{
		Report: &report,
	}); err != nil {
		return err
	}

	return nil
}

// Report sends report to the server regularly based on the report frequency
// retrieveUsageData is a function that outputs any of the type:
//	*usagePB.SessionReport_MgmtUsageData
//	*usagePB.SessionReport_ConnectorUsageData
//	*usagePB.SessionReport_ModelUsageData
//	*usagePB.SessionReport_PipelineUsageData
func (r *reporter) Report(ctx context.Context, service usagePB.Session_Service, edition, version string, retrieveUsageData func() interface{}) {

	logger, _ := logger.GetZapLogger()
	defer logger.Sync() //nolint

	for {
		localCtx, _ := context.WithTimeout(ctx, timeout)
		if err := r.SingleReport(localCtx, service, edition, version, retrieveUsageData()); err != nil {
			logger.Error(err.Error())
		}
		select {
		case <-ctx.Done():
			return
		case <-time.After(reportFrequency):
		}
	}
}

// StartReporter creates a usage reporter and start sending usage data to server regularly
// retrieveUsageData is a function that outputs any of the type:
//	*usagePB.SessionReport_MgmtUsageData
//	*usagePB.SessionReport_ConnectorUsageData
//	*usagePB.SessionReport_ModelUsageData
//	*usagePB.SessionReport_PipelineUsageData
func StartReporter(ctx context.Context, usageClient usagePB.UsageServiceClient, sessionService usagePB.Session_Service, edition, version string, retrieveUsageData func() interface{}) error {
	// Delay a short period time to start collecting data
	usageDelay := 5 * time.Second
	time.Sleep(usageDelay)

	reporter, err := NewReporter(ctx, usageClient, sessionService, edition, version)
	if err != nil {
		return err
	}

	go reporter.Report(ctx, sessionService, edition, version, retrieveUsageData)

	return nil
}
