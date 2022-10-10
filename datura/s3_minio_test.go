/*
This test should be flexible enough to run in two different environments:
1. Using the embedded Minio server (local)
2. Using an external Minio Cluster (Docker compose or KinD)
Moreover, there should be integration tests that testing:
1. Official Minio client
2. Official AWS cli
3. Our internal Job.Do SDK based logic
Any deviation should be immediately noticed..
*/

package datura_test

import (
	"bytes"
	"context"
	"embed"
	"fmt"
	"os"
	"testing"
	"time"

	madmin "github.com/minio/madmin-go"
	mclient "github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	minio "github.com/minio/minio/cmd"
	"github.com/spf13/viper"
	"github.com/theapemachine/wrkspc/brazil"
	"github.com/theapemachine/wrkspc/datura"
	"github.com/theapemachine/wrkspc/errnie"

	. "github.com/smartystreets/goconvey/convey"
)

var (
	addr     string = "127.0.0.1:38744"
	embedded embed.FS
)

const (
	MINIOKEY    = "minioadmin"
	MINIOSECRET = "minioadmin"
)

/*
initConfig does the embedded config stuff and sets the entire program up for Viper
based config, which uses the embedded yaml config file a lot.
*/
func initConfig(t *testing.T) {
	pwd := brazil.Workdir()
	viper.AddConfigPath(pwd)
	viper.SetConfigType("yml")
	viper.SetConfigName(".s3_minio_test")
	viper.AutomaticEnv()

	errnie.Handles(
		viper.ReadInConfig(),
	)
}

/*
Setup an ebedded local Minio instance
*/
func SetupMinio(name string, t *testing.T) (string, func() error, error) {
	dir := t.TempDir()
	go minio.Main([]string{"minio", "server", "--quiet", "--address", addr, dir})
	// TODO replace with Readiness check
	time.Sleep(2 * time.Second)

	madm, err := madmin.New(addr, MINIOKEY, MINIOSECRET, false)
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

	mc, err := mclient.New(addr, &mclient.Options{
		Creds:  credentials.NewStaticV4(MINIOKEY, MINIOSECRET, ""),
		Secure: false,
	})

	errnie.Handles(
		mc.MakeBucket(context.Background(), name, mclient.MakeBucketOptions{}),
	)

	return addr, func() error {
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

/*
Just a preliminary connection attempt to the embedded minio server. Does not really
serve a testing purpose as we are using the AWS client rather than the Minio client
*/
func TestMinioPut(t *testing.T) {
	var cleanup func() error
	var err error
	addr, cleanup, err = SetupMinio("test", t)
	if err != nil {
		t.Fatalf("while starting embedded server: %s\n", err)
	}

	var mc *mclient.Client
	Convey("Given a new minio client", t, func() {
		mc, err = mclient.New(addr, &mclient.Options{
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
			"test",
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
Write a datagram to the embedded Minio bucket and count the objects in the bucket.
*/
func TestDaturaS3WriteDatagram(t *testing.T) {
	var cleanup func() error
	var err error
	addr, cleanup, err = SetupMinio("cnc-development-datalake", t)
	if err != nil {
		t.Fatalf("while starting embedded server: %s\n", err)
	}

	// sets the addr of the Minio endpoint
	initConfig(t)

	mc := datura.NewS3()
	t.Logf("S3 Client Bucket: %s", *mc.Bucket)
	t.Logf("S3 Client Client: %+v", *mc.Client)
	t.Logf("S3 Client Region: %s", mc.Region)
	t.Logf("S3 Client Ctx: %+v", *mc.Ctx)
	t.Logf("S3 Client Pool: %+v", *mc.Pool)

	// create a Pool
	// .. here

	// create a Context holding a test datagram
	// .. here

	// TODO: How can the datura.S3 instance access the previously created Context + Pool?
	// Convey("Given a datagram job", t, func() {
	// 	// Create job in queue

	// 	Convey("Writing an object should give no error", func() {
	// 		So(
	// 			// .. ,
	// 		)
	// 	})
	// })

	errnie.Handles(
		cleanup(),
	)

	t.Fatal("FAIL..")
}
