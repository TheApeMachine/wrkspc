import { FlumeConfig, Colors } from "flume"
import Port from "./config/port"
import Node from "./config/node"

const config = new FlumeConfig()

config.addPortType(
  Port("error", "Error", Colors.red, [])
).addPortType(
  Port("bytes", "Bytes", Colors.pink, [])
).addPortType(
  Port("string", "String", Colors.green, [])
).addPortType(
  Port("boolean", "True/False", Colors.orange, [])
).addPortType(
  Port("number", "Number", Colors.yellow, [])
).addPortType(
  Port("abstract", "drknow.Abstract", Colors.blue, [])
).addPortType(
  Port("datagram", "Datagram", Colors.purple, [])
).addPortType(
  Port("mimetype", "Type", Colors.grey, [
    "application/json", "application/octet-stream"
  ])
).addPortType(
  Port("roletype", "RoleType", Colors.grey, [
    "TEST", "DATAPOINT", "QUESTION"
  ])
).addPortType(
  Port("scopetype", "ScopeType", Colors.grey, [
    "UNIT", "BENCHMARK", "USER", "DATALAKE"
  ])
).addNodeType(Node(
  "string", "Text", "Outputs a string", 
  ["string", "string"], ["string"]
)).addNodeType(Node(
  "boolean", "True/False", "Outputs a boolean", 
  ["boolean"], ["boolean"]
)).addNodeType(Node(
  "number", "Number", "Outputs a number", 
  ["number"], ["number"]
)).addNodeType(Node(
  "pipe", "hefner.Pipe", "Transports data", 
  ["bytes"], ["bytes"]
)).addNodeType(Node(
  "workspace", "Workspace", "A Workspace instance", 
  ["string", "datagram"], ["datagram"]
)).addNodeType(Node(
  "datagram", "Datagram", "A Datagram instance",
  ["mimetype", "roletype", "scopetype"], ["datagram"]
)).addNodeType(Node(
  "unmarshal", "Unmarshal", "A Datagram unmarshaller",
  ["datagram"], ["mimetype", "roletype", "scopetype"]
)).addRootNodeType(Node(
  "cli", "CLI", "The input command", ["string"], ["string"]
))

export default config
