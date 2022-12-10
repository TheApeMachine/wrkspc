const Node = (
  name: string, label: string, desc: string, inp: any, out: any
) => {  
  const getPorts = () => {
    let fin: any = {}

    if (inp.length > 0) {
      fin["inputs"] = (ports: any) => inp.map((i: any) => ports[i]())
    }

    if (out.length > 0) {
      fin["outputs"] = (ports: any) => out.map((o: any) => ports[o]())
    }

    return fin
  }

  return({
    type: name, label: label, description: desc, ...getPorts()
  })
}

export default Node
