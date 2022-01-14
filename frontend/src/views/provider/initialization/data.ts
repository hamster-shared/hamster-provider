import { FormSchema } from '/@/components/Form';

// 基础设置 form
export const chainSchemas: FormSchema[] = [
  {
    field: 'address',
    component: 'Select',
    label: 'chainAddress',
    colProps: { span: 12 },
    componentProps: {
      options: [
        { label: 'dev(ws://127.0.0.1:9944)', value: 'ws://127.0.0.1:9944' },
        { label: 'test(ws://183.66.65.207:49944)', value: 'ws://183.66.65.207:49944' },
      ],
    },
    rules: [{ required: true }],
  },
  {
    field: 'account',
    component: 'Input',
    label: 'account',
    colProps: { span: 12 },
    componentProps: {
      placeholder: '区块链账户seed或助记词',
    },
    rules: [{ required: true }],
  },
]

// 基础设置 form
export const vmSchemas: FormSchema[] = [
  {
    field: 'cpu',
    component: 'InputNumber',
    label: 'cpu',
    colProps: { span: 12 },
    rules: [{ required: true }],
  },
  {
    field: 'mem',
    component: 'InputNumber',
    label: '内存',
    colProps: { span: 12 },
    rules: [{ required: true }],
  },
  {
    field: 'disk',
    component: 'InputNumber',
    label: '硬盘',
    colProps: { span: 12 },
    rules: [{ required: true }],
  },
  {
    field: 'system',
    component: 'Input',
    label: '操作系统',
    colProps: { span: 12 },
    rules: [{ required: true }],
  },
  {
    field: 'image',
    component: 'Input',
    label: '镜像',
    colProps: { span: 12 },
    rules: [{ required: true }],
  },
  {
    field: 'accessPort',
    component: 'InputNumber',
    label: '访问端口',
    colProps: { span: 12 },
    rules: [{ required: true }],
  },
  {
    field: 'type',
    component: 'Select',
    componentProps: {
      options: [
        { label: 'kvm', value: 'kvm' },
        { label: 'docker', value: 'docker' },
      ],
    },
    label: '虚拟化类型',
    colProps: { span: 12 },
    rules: [{ required: true }],
  },
];
