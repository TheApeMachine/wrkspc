import dynamic from "next/dynamic"

const FocusGraph = dynamic(() => import("./servicemap"), {
  ssr: false
})

export default FocusGraph
