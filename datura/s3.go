package datura

import (
	"net"
	"net/http"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	awshttp "github.com/aws/aws-sdk-go-v2/aws/transport/http"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/spf13/viper"
	"github.com/theapemachine/wrkspc/errnie"
	"github.com/theapemachine/wrkspc/spd"
	"github.com/theapemachine/wrkspc/twoface"
)

/*
Modifier allows for some extra generic operations to be performed while
retrieving objects from S3 compatible storage.
*/
type Modifier string

const (
	// ALLVERSIONS splits the prefix using all available versions.
	ALLVERSIONS = "*"
	// LATEST returns only the latest event(s) in a prefix
	LATEST = "@"
)

type S3 struct {
	ctx        twoface.Context
	client     *s3.Client
	region     string
	bucket     *string
	downloader *manager.Downloader
	uploader   *manager.Uploader
	pool       *twoface.Pool
}

func NewS3() *S3 {
	errnie.Traces()
	Raise()

	v := viper.GetViper()
	p := v.GetString("program")
	s := v.GetString(p + ".stage")
	c := v.GetStringMapString(p + ".stages." + s + ".s3")

	region := c["region"]
	bucket := c["bucket"]

	resolver := aws.EndpointResolverFunc(func(service, region string) (aws.Endpoint, error) {
		return aws.Endpoint{
			PartitionID:       "aws",
			URL:               c["endpoint"],
			SigningRegion:     region,
			HostnameImmutable: true,
		}, nil
	})

	conn := s3.NewFromConfig(aws.Config{
		Region:           region,
		Credentials:      credentials.NewStaticCredentialsProvider(c["key"], c["secret"], ""),
		EndpointResolver: resolver,
	}, func(o *s3.Options) {
		o.UsePathStyle = true
		o.HTTPClient = awshttp.NewBuildableClient().WithTransportOptions(func(tr *http.Transport) {
			*tr = http.Transport{
				Proxy: http.ProxyFromEnvironment,
				DialContext: (&net.Dialer{
					Timeout:   1 * time.Second,
					KeepAlive: 0,
				}).DialContext,
				MaxIdleConns:          0,
				MaxIdleConnsPerHost:   100,
				MaxConnsPerHost:       0,
				IdleConnTimeout:       0,
				DisableKeepAlives:     false,
				TLSHandshakeTimeout:   1 * time.Second,
				ExpectContinueTimeout: 1 * time.Second,
				ResponseHeaderTimeout: 1 * time.Second,
				ReadBufferSize:        4 << 10,
				WriteBufferSize:       4 << 10,
			}
		})
	})

	ctx := twoface.NewContext(nil)

	return &S3{
		ctx:        ctx,
		client:     conn,
		region:     region,
		bucket:     &bucket,
		downloader: manager.NewDownloader(conn),
		uploader:   manager.NewUploader(conn),
		pool:       twoface.NewPool(ctx).Run(),
	}
}

func (store *S3) Wait() {
	store.pool.Wait()
}

func (store *S3) PoolSize() int {
	return store.pool.Size()
}

/*
Read implements the io.Reader interface.
*/
func (store *S3) Read(p []byte) (n int, err error) {
	errnie.Traces()

	dg := spd.Unmarshal(p)

	filtered := store.Filter(store.List(spd.Payload(dg)))
	ctx := twoface.NewContext(nil)

	for _, key := range filtered {
		store.pool.Do(twoface.NewRetriableJob(
			ctx, DownloadJob{
				key: key,
				ctx: ctx,
				p:   p,
			},
		))
	}

	return len(p), err
}

/*
Write implements the io.Writer interface.
*/
func (store *S3) Write(p []byte) (n int, err error) {
	errnie.Traces()

	ctx := twoface.NewContext(nil)
	ctx.Write(p)

	store.pool.Do(twoface.NewRetriableJob(
		ctx, UploadJob{},
	))

	return len(p), err
}
