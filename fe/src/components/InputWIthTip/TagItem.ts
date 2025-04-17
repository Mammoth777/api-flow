export class TagItem {
  content: string;
  isSelected: boolean;
  private _nodes: HTMLSpanElement | null = null;
  constructor(content: string) {
    this.content = content;
    this.isSelected = false
  }

  nodes() {
    const emptyNode = document.createTextNode('\u200B');
    if (this._nodes) {
      return [
        emptyNode,
        this._nodes,
        emptyNode,
      ];
    }
    const text = this.content
    const spanElement = document.createElement('span');
    this._nodes = spanElement;
    spanElement.className = 'tag-item';
    spanElement.innerText = text;
    spanElement.setAttribute('data-value', text);
      spanElement.setAttribute('data-path', '$' + text);
    return [
      emptyNode,
      spanElement,
      emptyNode.cloneNode(),
    ]
  }

  insertByRange(range: Range) {
    const nodes = this.nodes();
    nodes.forEach(node => {
      range.insertNode(node);
    })
    const textNode = nodes[1].childNodes[0];
    range.setStartAfter(textNode);
    range.setEndAfter(textNode);
    range.collapse(true);
  }
}