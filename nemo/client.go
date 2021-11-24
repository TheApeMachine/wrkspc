package nemo

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/theapemachine/wrkspc/errnie"
	"github.com/theapemachine/wrkspc/spdg"
)

/*
Client is a wrapper around the AWS S3 client library.
*/
type Client struct {
	session *session.Session
	conn    *s3.S3
	region  *string
	bucket  *string
}

/*
NewClient constructs a Client and returns a reference pointer to it.
*/
func NewClient(key, secret, region, bucket string) Client {
	errnie.Traces()

	session, err := session.NewSession(&aws.Config{
		Region:      &region,
		Credentials: credentials.NewStaticCredentials(key, secret, ""),
	})

	errnie.Handles(err).With(errnie.NOOP)
	errnie.Logs(
		"setting up S3 client in", region, "for bucket", bucket,
	).With(errnie.DEBUG)

	return Client{
		session: session,
		conn:    s3.New(session),
		region:  &region,
		bucket:  &bucket,
	}
}

/*
Peek is used to retrieve data from the S3 store.
*/
func (client Client) Peek(dg *spdg.Datagram) chan *spdg.Datagram {
	errnie.Traces()

	out := make(chan *spdg.Datagram)
	return out
}

/*
Poke is used to save data in the S3 Store.
*/
func (client Client) Poke(dg *spdg.Datagram) {
	errnie.Traces()

	uploader := s3manager.NewUploaderWithClient(client.conn)

	upParams := &s3manager.UploadInput{
		Bucket: client.bucket,
		Key:    dg.Context.Prefix(),
		Body:   dg.Encode(),
	}

	result, err := uploader.Upload(upParams)
	errnie.Handles(err).With(errnie.NOOP)
	errnie.Logs(result).With(errnie.DEBUG)
}
