import Control from "./control"

const Port = (name: string, label: string, color: any, opts: any) => {
  return ({
    type: name,
    name: name,
    label: label,
    color: color,
    controls: [Control(name, label, opts)]
  })
}

export default Port
