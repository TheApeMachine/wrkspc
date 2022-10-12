package datura

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	madmin "github.com/minio/madmin-go"
	mclient "github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	minio "github.com/minio/minio/cmd"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/theapemachine/wrkspc/errnie"
	"github.com/theapemachine/wrkspc/spd"
)

const (
	MINIOKEY    = "minioadmin"
	MINIOSECRET = "minioadmin"
	MINIOADDR   = "127.0.0.1:9000"
)

func BenchmarkWrite(b *testing.B) {
	errnie.Tracing(false)
	store := NewS3()

	for i := 0; i < b.N; i++ {
		store.Write(spd.NewCached(
			"datapoint", "test", "test.wrkspc.org", "test",
		))
	}
}

func TestS3WriteRead(t *testing.T) {
	var cleanup func() error
	var err error
	_, cleanup, err = SetupMinio("wrkspc", t)
	if err != nil {
		t.Fatalf("while starting embedded server: %s\n", err)
	}

	Convey("Given an S3 connection", t, func() {
		tree := NewS3()

		Convey("And a value is written", func() {
			dg := spd.NewCached(
				"datapoint", "test", "test.wrkspc.org", "test",
			)

			tree.Write(dg)

			Convey("It should be able to retrieve the value", func() {
				for _, key := range []string{
					"v4.0.0/datapoint/test/test.wrkspc.org",
					"v4.0.0/datapoint/test/test.wr",
					"v4.0.0/datapoint/",
				} {
					q := spd.NewCached(
						"datapoint", "test", "test.wrkspc.org",
						key,
					)

					tree.Read(q)

					So(
						string(spd.Unmarshal(q).Payload()),
						ShouldEqual,
						string(spd.Unmarshal(dg).Payload()),
					)
				}
			})
		})
	})

	errnie.Handles(
		cleanup(),
	)
}

// func TestS3NewS3ClientError(t *testing.T) {
// 	Convey("Given a invalid S3 connection", t, func() {
// 		_ = NewS3()

// 		Convey("The connection attempt should fail gracefully", func() {})
// 	})
// }

/*
Just a preliminary connection attempt to the embedded minio server. Does not really
serve a testing purpose as we are using the AWS client rather than the Minio client
*/
func TestMinioPut(t *testing.T) {
	var cleanup func() error
	var err error
	_, cleanup, err = SetupMinio("wrkspc", t)
	if err != nil {
		t.Fatalf("while starting embedded server: %s\n", err)
	}

	var mc *mclient.Client
	Convey("Given a new minio client", t, func() {
		mc, err = mclient.New(MINIOADDR, &mclient.Options{
			Creds:  credentials.NewStaticV4(MINIOKEY, MINIOSECRET, ""),
			Secure: false,
		})

		Convey("There should be no error", func() {
			So(
				err,
				ShouldBeNil,
			)
		})
	})

	data := []byte("test")
	Convey("Using the minio client", t, func() {
		_, err = mc.PutObject(
			context.Background(),
			"wrkspc",
			"foo/var",
			bytes.NewReader(data),
			int64(len(data)),
			mclient.PutObjectOptions{},
		)

		Convey("Writing an object should give no error", func() {
			So(
				err,
				ShouldBeNil,
			)
		})
	})

	errnie.Handles(
		cleanup(),
	)
}

/*
Setup an ebedded local Minio instance
*/
func SetupMinio(name string, t *testing.T) (string, func() error, error) {
	dir := t.TempDir()
	go minio.Main([]string{"minio", "server", "--quiet", "--address", MINIOADDR, dir})
	// TODO replace with Readiness check
	time.Sleep(2 * time.Second)

	madm, err := madmin.New(MINIOADDR, MINIOKEY, MINIOSECRET, false)
	if err != nil {
		return "", nil, errnie.NewError(
			fmt.Errorf("while creating madmin.. %s", err),
		)
	}

	// Fetch service status.
	_, err = madm.ServerInfo(context.TODO())
	if err != nil {
		return "", nil, errnie.NewError(
			fmt.Errorf("while getting madmin info.. %s", err),
		)
	}

	mc, err := mclient.New(MINIOADDR, &mclient.Options{
		Creds:  credentials.NewStaticV4(MINIOKEY, MINIOSECRET, ""),
		Secure: false,
	})

	errnie.Handles(
		mc.MakeBucket(context.Background(), name, mclient.MakeBucketOptions{}),
	)

	return MINIOADDR, func() error {
		t.Cleanup(func() {
			errnie.Handles(
				os.RemoveAll(dir),
			)
			errnie.Handles(
				madm.ServiceStop(context.Background()),
			)
		})

		return nil
	}, nil

}
