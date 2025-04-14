const NodeTypeMap = {
  system: [
    'execInput'
  ],
  task: [
    'api',
    'text'
  ]
}

export const isSystemNode = (nodeType: string) => {
  console.log('nodeType', nodeType)
  return NodeTypeMap.system.includes(nodeType)
}