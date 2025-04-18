export class TagItem {
  content: string;
  isSelected: boolean;
  private parent?: HTMLElement | null;
  constructor(content: string, parent?: HTMLElement | null) {
    this.content = content;
    this.isSelected = false
    this.parent = parent;
  }

  nodes() {
    if (this.parent) {
      const textNode = document.createTextNode(this.content)
      return [textNode]
    } else {
      const emptyNode = document.createTextNode('\u200B');
      const text = this.content
      const spanElement = document.createElement('span');
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
  }

  insertByRange(range: Range) {
    const nodes = this.nodes();
    nodes.forEach(node => {
      range.insertNode(node);
    })
    const textNode = this.parent ? nodes[0] : nodes[1].childNodes[0];
    range.setStartAfter(textNode);
    range.setEndAfter(textNode);
    range.collapse(true);
  }
}