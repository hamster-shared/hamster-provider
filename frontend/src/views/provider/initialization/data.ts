import { FormSchema } from '/@/components/Form';
import { useI18n } from '/@/hooks/web/useI18n';
const { t } = useI18n();
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
      placeholder: t('initialization.initialization.seedTip'),
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
    label: t('initialization.initialization.memory'),
    colProps: { span: 12 },
    rules: [{ required: true }],
  },
  {
    field: 'disk',
    component: 'InputNumber',
    label: t('initialization.initialization.disk'),
    colProps: { span: 12 },
    rules: [{ required: true }],
  },
  {
    field: 'system',
    component: 'Input',
    label: t('initialization.initialization.system'),
    colProps: { span: 12 },
    rules: [{ required: true }],
  },
  {
    field: 'image',
    component: 'Input',
    label: t('initialization.initialization.image'),
    colProps: { span: 12 },
    rules: [{ required: true }],
  },
  {
    field: 'accessPort',
    component: 'InputNumber',
    label: t('initialization.initialization.accessPort'),
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
    label: t('initialization.initialization.type'),
    colProps: { span: 12 },
    rules: [{ required: true }],
  },
];
