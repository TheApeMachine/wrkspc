package datura

import (
	"bytes"

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
	ctx        twoface.Context
	p          []byte
}

func (job DownloadJob) Do() errnie.Error {
	errnie.Traces()

	buf := manager.NewWriteAtBuffer([]byte{})

	_, err := job.downloader.Download(
		job.ctx.Get(), buf, &s3.GetObjectInput{
			Bucket: job.bucket,
			Key:    &job.key,
		},
	)

	if e := errnie.Handles(err); e.Type == errnie.NIL {
		return e
	}

	job.p = append(job.p, buf.Bytes()...)
	return errnie.NewError(nil)
}

type UploadJob struct {
	p        []byte
	bucket   *string
	uploader *manager.Uploader
	ctx      twoface.Context
}

func (job UploadJob) Do() errnie.Error {
	errnie.Traces()

	buf := bytes.NewBuffer(job.p)
	dg := spd.Unmarshal(job.p)

	_, err := job.uploader.Upload(
		job.ctx.Get(), &s3.PutObjectInput{
			Bucket: job.bucket,
			Key:    aws.String(dg.Prefix()),
			Body:   buf,
		},
	)

	return errnie.Handles(err)
}
