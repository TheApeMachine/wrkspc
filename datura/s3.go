package datura

import (
	"bytes"
	"context"
	"net"
	"net/http"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	awshttp "github.com/aws/aws-sdk-go-v2/aws/transport/http"
	"github.com/aws/aws-sdk-go-v2/config"
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
	Raise()

	v := viper.GetViper()
	p := v.GetString("program")
	s := v.GetString(p + ".stage")
	c := v.GetStringMapString(p + ".stages." + s + ".s3")

	region := c["region"]
	bucket := c["bucket"]

	resolver := aws.EndpointResolverWithOptionsFunc(
		func(service, region string, options ...interface{}) (aws.Endpoint, error) {
			return aws.Endpoint{
				PartitionID:       "aws",
				URL:               c["endpoint"],
				SigningRegion:     region,
				HostnameImmutable: true,
			}, &aws.EndpointNotFoundError{}
		},
	)

	conn := s3.NewFromConfig(aws.Config{
		Region:           region,
		Credentials:      credentials.NewStaticCredentialsProvider(c["key"], c["secret"], ""),
		EndpointResolver: config.WithEndpointResolverWithOptions(resolver),
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

func (store *S3) Wait() {
	store.pool.Wait()
}

func (store *S3) PoolSize() int {
	return store.pool.Size()
}

type DownloadJob struct {
	bucket     *string
	downloader *manager.Downloader
	key        string
	out        []byte
}

func (job DownloadJob) Do() {
	errnie.Traces()

	buf := manager.NewWriteAtBuffer([]byte{})

	_, err := job.downloader.Download(
		context.Background(), buf, &s3.GetObjectInput{
			Bucket: job.bucket,
			Key:    &job.key,
		},
	)

	errnie.Handles(err)
	job.out = make([]byte, len(buf.Bytes()))
	job.out = buf.Bytes()
}

/*
Read implements the io.Reader interface.
*/
func (store *S3) Read(p []byte) (n int, err error) {
	errnie.Traces()

	dg := spd.Unmarshal(p)

	filtered := store.Filter(store.List(spd.Payload(dg)))
	jobs := make([]DownloadJob, len(filtered))

	for idx, key := range filtered {
		jobs[idx] = DownloadJob{
			bucket:     store.bucket,
			downloader: store.downloader,
			key:        key,
		}

		store.pool.Do(jobs[idx])
	}

	p = nil

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
	delay    int64
}

func (job UploadJob) Do() {
	errnie.Traces()

	buf := bytes.NewBuffer(job.p)

	dg := spd.Unmarshal(job.p)

	_, err := job.uploader.Upload(
		job.ctx, &s3.PutObjectInput{
			Bucket: job.bucket,
			Key:    aws.String(dg.Prefix()),
			Body:   buf,
		},
	)

	if err != nil {
		job.delay += 5
		time.Sleep(time.Duration(job.delay))
		job.Do()
	}
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
