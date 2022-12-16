package infra

import (
	"bytes"
	"io"
	"strings"

	"github.com/theapemachine/wrkspc/brazil"
	"github.com/theapemachine/wrkspc/container"
	"github.com/theapemachine/wrkspc/tweaker"
	"github.com/theapemachine/wrkspc/twoface"
)

type PXEBoot struct {
	ctx *twoface.Context
}

func (pxe *PXEBoot) Read(p []byte) (n int, err error) {
	brazil.NewFile(
		"~/tmp/wrkspc/quay.io/coreos/butane", "butane.yml",
		bytes.NewBuffer([]byte{}),
	)

	container.NewDocker(
		pxe.ctx, "quay.io/coreos", "butane", "latest",
	).Pull().Create(
		nil, &[]string{"butane.yml", ">", "ignition.json"},
	).Start()

	container.NewDocker(
		pxe.ctx, "pixiecore", "pixiecore", "",
	).Pull().Create(
		nil, &[]string{"boot", strings.Join([]string{
			tweaker.GetString("pxe.channel"),
			tweaker.GetString("pxe.kernel"),
		}, "/"), strings.Join([]string{
			tweaker.GetString("pxe.channel"),
			tweaker.GetString("pxe.image"),
		}, "/")},
	).Start()

	ip := tweaker.GetString("pxe.machines.deepthought.ip")
	user := tweaker.GetString("pxe.machines.deepthought.username")
	pass := tweaker.GetString("pxe.machines.deepthought.password")

	dt := NewIDRAC(ip, user, pass)
	dt.Reboot()

	return len(p), io.EOF
}
