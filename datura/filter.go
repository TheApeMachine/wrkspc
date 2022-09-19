package datura

import (
	"context"
	"strings"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/theapemachine/wrkspc/errnie"
)

/*
FilterJob allows the filtering workload to be scheduled
onto the worker pool.
*/
type FilterJob struct {
	paginator *s3.ListObjectsV2Paginator
	modifier  Modifier
	listCache *sync.Map
	wg        *sync.WaitGroup
}

/*
Do implements the twoface.Job interface.
*/
func (job FilterJob) Do() {
	errnie.Traces()
	defer job.wg.Done()

	for job.paginator.HasMorePages() {
		page, err := job.paginator.NextPage(context.TODO())
		errnie.Handles(err)

		// Page can end up to be nil, so we should have an
		// escape route for this situation.
		if page == nil {
			continue
		}

		for _, obj := range page.Contents {
			key := *obj.Key

			if job.modifier == "@" {
				// We need to strip off the version and the
				// timestamp and uuid.
				split := strings.Split(key, "/")
				key = key[1 : len(split)-1]
			}

			// Try to read a datapoint from the cache for this key.
			val, ok := job.listCache.Load(key)

			if ok && obj.LastModified.After(val.(time.Time)) {
				// The datapoint is newer than then one currently sitting
				// in the cache, so we need to update the cache.
				job.listCache.Store(key, val)
			}

			if !ok {
				// We don't have a datapoint in the cache yet for this key,
				// so we need to start things off by writing the first one.
				job.listCache.Store(key, *obj.LastModified)
			}
		}
	}
}

/*
Filter an object listing based on a set of specified wildcards and other
modifiers that add another layer of usability onto the S3 API.
*/
func (store *S3) Filter(
	paginators []*s3.ListObjectsV2Paginator, modifier Modifier,
) []string {
	errnie.Traces()

	var out []string
	var wg sync.WaitGroup
	listCache := sync.Map{}

	for _, paginator := range paginators {
		wg.Add(1)

		// Start a job for each paginator so they do not have
		// to wait around for each other.
		store.pool.Do(FilterJob{
			paginator: paginator,
			modifier:  modifier,
			listCache: &listCache,
			wg:        &wg,
		})
	}

	// Wait until the filter jobs are done.
	wg.Wait()

	// We can now range over our concurrency safe map to
	// get the final prefixes we need to download.
	listCache.Range(func(key, _ any) bool {
		out = append(out, key.(string))
		return true
	})

	return out
}
