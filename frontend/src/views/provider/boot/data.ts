// 基础设置 form
import { DescItem } from '/@/components/Description';
import { useI18n } from '/@/hooks/web/useI18n';
const { t } = useI18n();
export const vmSchemas: DescItem[] = [
  {
    field: 'cpu',
    label: 'cpu',
  },
  {
    field: 'mem',
    label: t('boot.boot.memory'),
  },
  {
    field: 'disk',
    label: t('boot.boot.disk'),
  },
  {
    field: 'system',
    label: t('boot.boot.system'),
  },
  {
    field: 'image',
    label: t('boot.boot.image'),
  },
  {
    field: 'accessPort',
    label: t('boot.boot.accessPort'),
  },
  {
    field: 'type',
    label: t('boot.boot.type'),
  },
];
