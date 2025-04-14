// 定义节点输入字段的类型
export class InputField {
  field: string;
  type: 'string' | 'number' | 'object' | 'options' | 'boolean';
  desc: string;
  default: any;
  // 仅在 type 为 'options' 时使用
  options?: string[] | { label: string, value: number | string}[];
}