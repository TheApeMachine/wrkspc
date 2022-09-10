package datura

import (
	"bytes"
	"context"
	"net"
	"net/http"
	"sync"
	"time"

	"capnproto.org/go/capnp/v3"
	"github.com/aws/aws-sdk-go-v2/aws"
	awshttp "github.com/aws/aws-sdk-go-v2/aws/transport/http"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
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
	ctx        *twoface.Context
	client     *s3.Client
	region     string
	bucket     *string
	downloader *manager.Downloader
	uploader   *manager.Uploader
	pool       *twoface.Pool
}

func NewS3() *S3 {
	errnie.Traces()

	region := "us-east-2"
	bucket := "wrkspc"

	resolver := aws.EndpointResolverFunc(func(service, region string) (aws.Endpoint, error) {
		return aws.Endpoint{
			PartitionID:       "aws",
			URL:               "http://127.0.0.1:9000",
			SigningRegion:     region,
			HostnameImmutable: true,
		}, nil
	})

	conn := s3.NewFromConfig(aws.Config{
		Region:           region,
		Credentials:      credentials.NewStaticCredentialsProvider("minioadmin", "minioadmin", ""),
		EndpointResolver: resolver,
	}, func(o *s3.Options) {
		o.UsePathStyle = true
		o.HTTPClient = awshttp.NewBuildableClient().WithTransportOptions(func(tr *http.Transport) {
			*tr = http.Transport{
				Proxy: http.ProxyFromEnvironment,
				DialContext: (&net.Dialer{
					Timeout:   3 * time.Second,
					KeepAlive: 10 * time.Second,
				}).DialContext,
				MaxIdleConns:          100,
				MaxIdleConnsPerHost:   100,
				MaxConnsPerHost:       100,
				IdleConnTimeout:       10 * time.Second,
				DisableKeepAlives:     false,
				TLSHandshakeTimeout:   3 * time.Second,
				ExpectContinueTimeout: 5 * time.Second,
				ResponseHeaderTimeout: 5 * time.Second,
				ReadBufferSize:        16 * 1024 * 1024,
				WriteBufferSize:       16 * 1024 * 1024,
			}
		})
	})

	ctx := twoface.NewContext()

	// Start a new worker pool to perform operations concurrently.
	pool := twoface.NewPool(ctx)
	pool.Run()

	return &S3{
		ctx:        ctx,
		client:     conn,
		region:     region,
		bucket:     &bucket,
		downloader: manager.NewDownloader(conn),
		uploader:   manager.NewUploader(conn),
		pool:       pool,
	}
}

type DownloadJob struct {
	bucket     *string
	downloader *manager.Downloader
	key        string
	out        []byte
	wg         *sync.WaitGroup
}

func (job DownloadJob) Do() {
	errnie.Traces()
	defer job.wg.Done()

	buf := manager.NewWriteAtBuffer([]byte{})

	_, err := job.downloader.Download(
		context.Background(), buf, &s3.GetObjectInput{
			Bucket: job.bucket,
			Key:    &job.key,
		},
	)

	errnie.Handles(err).With(errnie.NOOP)
	job.out = make([]byte, len(buf.Bytes()))
	job.out = buf.Bytes()
}

/*
Read implements the io.Reader interface.
*/
func (store *S3) Read(p []byte) (n int, err error) {
	errnie.Traces()

	msg, err := capnp.Unmarshal(p)
	errnie.Handles(err).With(errnie.NOOP)

	dg, err := spd.ReadRootDatagram(msg)
	errnie.Handles(err).With(errnie.NOOP)

	filtered := store.Filter(store.List(spd.Prefix(dg)))
	jobs := make([]DownloadJob, len(filtered))

	var wg sync.WaitGroup

	for idx, key := range filtered {
		wg.Add(1)

		jobs[idx] = DownloadJob{
			bucket:     store.bucket,
			downloader: store.downloader,
			key:        key,
			wg:         &wg,
		}

		store.pool.Do(jobs[idx])
	}

	p = nil
	wg.Wait()

	for _, job := range jobs {
		p = append(p, job.out...)
	}

	return len(p), err
}

type UploadJob struct {
	p        []byte
	bucket   *string
	uploader *manager.Uploader
	ctx      *twoface.Context
}

func (job UploadJob) Do() {
	errnie.Traces()

	buf := bytes.NewBuffer(job.p)

	msg, err := capnp.Unmarshal(job.p)
	errnie.Handles(err).With(errnie.NOOP)

	dg, err := spd.ReadRootDatagram(msg)
	errnie.Handles(err).With(errnie.NOOP)

	prefix := spd.Prefix(dg)

	_, err = job.uploader.Upload(
		job.ctx, &s3.PutObjectInput{
			Bucket: job.bucket,
			Key:    aws.String(prefix),
			Body:   buf,
		},
	)

	errnie.Handles(err).With(errnie.NOOP)
}

/*
Write implements the io.Writer interface.
*/
func (store *S3) Write(p []byte) (n int, err error) {
	errnie.Traces()

	store.pool.Do(UploadJob{
		p:        p,
		bucket:   store.bucket,
		uploader: store.uploader,
		ctx:      store.ctx,
	})

	return len(p), err
}
