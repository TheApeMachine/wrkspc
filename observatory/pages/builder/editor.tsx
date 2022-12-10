import React, { useState, useEffect } from "react"
import { NodeEditor } from "flume"
import config from "./config"

const Editor = () => {
  const [show, setShow] = useState(false)

  useEffect(() => {
    setShow(true)
  }, [])

  return (
    <>
      {show && <NodeEditor 
        circularBehavior="allow"
        nodeTypes={config.nodeTypes} 
        portTypes={config.portTypes} 
        defaultNodes={[
          {
            type: "cli",
            x: -(window.innerWidth/2-10), 
            y: -(window.innerHeight/2-10)
          }
        ]}
      />}
    </>
  )
} 

export default Editor
