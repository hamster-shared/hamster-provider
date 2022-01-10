import type { AppRouteModule } from '/@/router/types';

import { LAYOUT } from '/@/router/constant';
import { t } from '/@/hooks/web/useI18n';

const dashboard: AppRouteModule = {
  path: '/provider',
  name: 'Provider',
  component: LAYOUT,
  redirect: '/provider/setting',
  meta: {
    orderNo: 10,
    icon: 'mdi:desktop-classic',
    title: t('routes.provider.provider'),
  },
  children: [
    {
      path: 'setting',
      name: 'ProviderSetting',
      component: () => import('/@/views/provider/setting/index.vue'),
      meta: {
        // affix: true,
        icon: 'mdi:cog-outline',
        title: t('routes.provider.setting'),
      },
    },
  ],
};

export default dashboard;
