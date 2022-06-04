package usage

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"runtime"
	"time"

	"github.com/catalinc/hashcash"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"

	usagev1alpha "github.com/instill-ai/protogen-go/vdp/usage/v1alpha"
)

const (
	// Default server
	defaultURL = "https://usage.instill.tech"
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
type Session usagev1alpha.Session

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
	// 	*usagev1alpha.SessionReport_MgmtUsageData
	//	*usagev1alpha.SessionReport_ConnectorUsageData
	//	*usagev1alpha.SessionReport_ModelUsageData
	//	*usagev1alpha.SessionReport_PipelineUsageData
	SingleReport(ctx context.Context, service usagev1alpha.Session_Service, env, version string, usageData interface{}) error
	// Report sends report to the server regularly based on the report frequency
	Report(ctx context.Context, db *gorm.DB, service usagev1alpha.Session_Service, env, version string, retrieveUsageData func(*gorm.DB) (usageData interface{}, err error))
}

// reporter represents a reporter that sends usage data to the server on a regular basis
type reporter struct {
	client     usagev1alpha.UsageServiceClient
	sessionUID string
	start      time.Time
	minter     Minter
	token      string
	url        string
}

// NewReporter creates a new usage reporter
func NewReporter(ctx context.Context, conn *grpc.ClientConn, service usagev1alpha.Session_Service, url, env, version string) (Reporter, error) {
	if url == "" {
		url = defaultURL
	}

	// Initialize client
	client := usagev1alpha.NewUsageServiceClient(conn)

	// Create the session
	resp, err := client.CreateSession(ctx,
		&usagev1alpha.CreateSessionRequest{
			Session: &usagev1alpha.Session{
				Service:    service,
				Env:        env,
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
		url:        url,
	}

	return r, nil
}

// SingleReport represents send one report to the usage server
func (r *reporter) SingleReport(ctx context.Context, service usagev1alpha.Session_Service, env, version string, usageData interface{}) error {
	s := Session{
		Service:    service,
		Env:        env,
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

	report := usagev1alpha.SessionReport{
		SessionUid: r.sessionUID,
		Token:      r.token,
		Pow:        pow,
	}

	// Session report
	pbSessionData := usagev1alpha.Session(s)
	report.Session = &pbSessionData
	// Usage data
	invalidUsageDataErr := status.Error(codes.InvalidArgument, "invalid usage data")
	switch service {
	case usagev1alpha.Session_SERVICE_MGMT:
		if ud, ok := usageData.(*usagev1alpha.SessionReport_MgmtUsageData); ok {
			pbUsageData := usagev1alpha.SessionReport_MgmtUsageData(*ud)
			report.UsageData = &pbUsageData
		} else {
			return invalidUsageDataErr
		}
	case usagev1alpha.Session_SERVICE_CONNECTOR:
		if ud, ok := usageData.(*usagev1alpha.SessionReport_ConnectorUsageData); ok {
			pbUsageData := usagev1alpha.SessionReport_ConnectorUsageData(*ud)
			report.UsageData = &pbUsageData
		} else {
			return invalidUsageDataErr
		}
	case usagev1alpha.Session_SERVICE_MODEL:
		if ud, ok := usageData.(*usagev1alpha.SessionReport_ModelUsageData); ok {
			pbUsageData := usagev1alpha.SessionReport_ModelUsageData(*ud)
			report.UsageData = &pbUsageData
		} else {
			return invalidUsageDataErr
		}
	case usagev1alpha.Session_SERVICE_PIPELINE:
		if ud, ok := usageData.(*usagev1alpha.SessionReport_PipelineUsageData); ok {
			pbUsageData := usagev1alpha.SessionReport_PipelineUsageData(*ud)
			report.UsageData = &pbUsageData
		} else {
			return invalidUsageDataErr
		}
	default:
		return invalidUsageDataErr
	}

	_, err = r.client.SendSessionReport(ctx, &usagev1alpha.SendSessionReportRequest{
		Report: &report,
	})
	return err
}

// Report sends report to the server regularly based on the report frequency
// retrieveUsageData is a function that takes *gorm.DB as the input, and outputs (usagev1alpha.isSessionReport_UsageData, error)
func (r *reporter) Report(ctx context.Context, db *gorm.DB, service usagev1alpha.Session_Service, env, version string, retrieveUsageData func(*gorm.DB) (interface{}, error)) {

	for {
		usageData, _ := retrieveUsageData(db)
		localCtx, _ := context.WithTimeout(ctx, timeout)
		r.SingleReport(localCtx, service, env, version, usageData)
		select {
		case <-ctx.Done():
			return
		case <-time.After(reportFrequency):
		}
	}
}

// StartReporter creates a usage reporter and start sending usage data to server regularly
// retrieveUsageData is a function that takes *gorm.DB as the input, and outputs (usagev1alpha.isSessionReport_UsageData, error)
func StartReporter(ctx context.Context, db *gorm.DB, conn *grpc.ClientConn, service usagev1alpha.Session_Service, url, env, version string, retrieveUsageData func(*gorm.DB) (interface{}, error)) error {
	// Delay a short period time to start collecting data
	usageDelay := 5 * time.Second
	time.Sleep(usageDelay)

	reporter, err := NewReporter(ctx, conn, service, url, env, version)
	if err != nil {
		return err
	}

	go reporter.Report(ctx, db, service, env, version, retrieveUsageData)

	return nil
}
