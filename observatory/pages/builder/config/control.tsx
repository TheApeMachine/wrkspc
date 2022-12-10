import { Controls } from "flume"

const Control = (control: string, label: string, opts: any) => {
  
  const getControl = () => {
    let ctrl = {name: control, label: label}
    switch(control) {
    case "error":
    case "bytes":
    case "string":
    case "abstract":
    case "datagram":
      return Controls.text(ctrl)
    case "boolean":
      return Controls.checkbox(ctrl)
    case "number":
      return Controls.number(ctrl)
    case "select":
    case "mimetype":
    case "roletype":
    case "scopetype":
      return Controls.select(
        {...ctrl, options: opts.map(
          (opt: string) => ({value: opt, label: opt})
        )}
      )
    }
  }

  return getControl()
}

export default Control
