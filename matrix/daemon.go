package matrix

import (
	"os"

	"github.com/containerd/containerd/cmd/containerd/command"
	"github.com/containerd/containerd/pkg/seed"
	"github.com/theapemachine/wrkspc/errnie"
)

func init() {
	seed.WithTimeAndRand()
}

/*
NewDaemon starts a new ContainerD daemon.
*/
func NewDaemon() {
	errnie.Traces()
	go func() {
		errnie.Handles(
			command.App().Run(os.Args),
		).With(errnie.KILL)
	}()
	// time.Sleep(3 * time.Second)

	// client, err := containerd.New("/run/containerd/containerd.sock")
	// if err != nil {
	// 	log.Println(err)
	// 	return
	// }
	// defer client.Close()

	// ctx := namespaces.WithNamespace(context.Background(), "apeterm2")
	// image, err := client.Pull(ctx, "docker.io/theapemachine/term:latest", containerd.WithPullUnpack)
	// if err != nil {
	// 	log.Println(err)
	// 	return
	// }
	// log.Printf("Successfully pulled %s image\n", image.Name())

	// container, err := client.NewContainer(
	// 	ctx,
	// 	"apeterm2",
	// 	containerd.WithNewSnapshot("apeterm2-snapshot", image),
	// 	containerd.WithNewSpec(oci.WithImageConfig(image)),
	// )
	// if err != nil {
	// 	log.Println(err)
	// 	return
	// }
	// defer container.Delete(ctx, containerd.WithSnapshotCleanup)

	// spec, err := container.Spec(ctx)
	// if err != nil {
	// 	log.Println(err)
	// 	return
	// }
	// var (
	// 	con console.Console
	// 	tty = spec.Process.Terminal
	// )
	// if tty {
	// 	con = console.Current()
	// 	defer con.Reset()
	// 	if err := con.SetRaw(); err != nil {
	// 		log.Println(err)
	// 		return
	// 	}
	// }

	// // create a task from the container
	// task, err := container.NewTask(ctx, cio.NewCreator(cio.WithStdio))

	// if err != nil {
	// 	log.Println(err)
	// 	return
	// }

	// defer task.Delete(ctx)

	// // make sure we wait before calling start
	// exitStatusC, err := task.Wait(ctx)
	// if err != nil {
	// 	fmt.Println(err)
	// }

	// // call start on the task to execute the redis server
	// if err := task.Start(ctx); err != nil {
	// 	log.Println(err)
	// 	return
	// }

	// if tty {
	// 	if err := tasks.HandleConsoleResize(ctx, task, con); err != nil {
	// 		logrus.WithError(err).Error("console resize")
	// 	}
	// } else {
	// 	sigc := commands.ForwardAllSignals(ctx, task)
	// 	defer commands.StopCatch(sigc)
	// }

	// // kill the process and get the exit status
	// if err := task.Kill(ctx, syscall.SIGTERM); err != nil {
	// 	log.Println(err)
	// 	return
	// }

	// // wait for the process to fully exit and print out the exit status

	// status := <-exitStatusC
	// code, _, err := status.Result()
	// if err != nil {
	// 	log.Println(err)
	// 	return
	// }
	// fmt.Printf("redis-server exited with status: %d\n", code)
}
