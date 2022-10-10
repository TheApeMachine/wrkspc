package datura

import (
	"bytes"
	"context"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/theapemachine/wrkspc/errnie"
	"github.com/theapemachine/wrkspc/spd"
	"github.com/theapemachine/wrkspc/twoface"
)

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
