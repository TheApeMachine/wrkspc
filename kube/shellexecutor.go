package kube

// import (
// 	"time"

// 	"github.com/theapemachine/wrkspc/shell"
// )

// type ShellExecutor struct {
// 	base   Base
// 	client *shell.PowerShell
// 	line   string
// }

// func NewShellExecutor(client *shell.PowerShell, line string) MigratableKind {
// 	return &ShellExecutor{
// 		client: client,
// 		line:   line,
// 	}
// }

// func (kind *ShellExecutor) Up() error {
// 	kind.client.Execute(kind.line)
// 	kind.base = NewBase(kind)
// 	kind.base.waiter(true)

// 	return nil
// }

// func (kind *ShellExecutor) Check() bool {
// 	// Obvs this needs a better solution.
// 	time.Sleep(5 * time.Second)
// 	return true
// }

// func (kind *ShellExecutor) Down() error {
// 	return kind.base.teardown()
// }

// func (kind *ShellExecutor) Delete() error {
// 	return nil
// }

// func (kind *ShellExecutor) Name() string {
// 	return kind.line
// }
