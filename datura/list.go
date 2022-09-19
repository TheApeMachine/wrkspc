package datura

import (
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/theapemachine/wrkspc/errnie"
)

/*
List all the objects under the given prefix, which allows for some
basic, but rapid inspection of the objects. Specifically useful when
you can make decisions based on the LastModified timestamp.
*/
func (store *S3) List(prefix []byte) ([]*s3.ListObjectsV2Paginator, Modifier) {
	errnie.Traces()
	str := string(prefix)
	split := strings.Split(str, "/")

	var prefixes []string
	var modifier Modifier
	var m Modifier
	var out []*s3.ListObjectsV2Paginator

	if len(split) < 2 {
		return out, modifier
	}

	for idx, mod := range split[:2] {
		prefixes, m = store.modify(idx, str, mod)

		if m != "" {
			modifier = m
		}
	}

	for _, prfx := range prefixes {
		out = append(out, s3.NewListObjectsV2Paginator(
			store.client, &s3.ListObjectsV2Input{
				Bucket: store.bucket,
				Prefix: aws.String(prfx),
			},
		))
	}

	return out, modifier
}

func (store *S3) modify(idx int, prefix, mod string) ([]string, Modifier) {
	errnie.Traces()

	out := []string{prefix}
	var modifier Modifier

	if mod == ALLVERSIONS {
		out = []string{}

		for i := 0; i < 2; i++ {
			out = append(
				out,
				prefix[0:idx]+fmt.Sprintf("v%d.0.0", i+2)+prefix[idx+1:],
			)
		}
	}

	if mod == LATEST {
		found := strings.Index(prefix, mod)

		if found != -1 {
			out[idx] = prefix[0:idx-1] + prefix[idx+1:]
			modifier = LATEST
		}
	}

	return out, modifier
}
