import type { AppRouteModule } from '/@/router/types';

import { LAYOUT } from '/@/router/constant';
import { t } from '/@/hooks/web/useI18n';

const dashboard: AppRouteModule = {
  path: '/provider',
  name: 'Provider',
  component: LAYOUT,
  redirect: '/provider/resource',
  meta: {
    orderNo: 10,
    icon: 'mdi:desktop-classic',
    title: t('routes.provider.provider'),
  },
  children: [
    {
      path: 'initialization',
      name: 'ProviderInitialization',
      component: () => import('/@/views/provider/initialization/index.vue'),
      meta: {
        // affix: true,
        icon: 'mdi:cog-outline',
        title: t('routes.provider.setting'),
      },
    },
    {
      path: 'boot',
      name: 'BootSetting',
      component: () => import('/@/views/provider/boot/index.vue'),
      meta: {
        icon: 'bi:play-circle',
        title: t('routes.provider.boot'),
      },
    },
    {
      path: 'resource',
      name: 'ResourceDetail',
      component: () => import('/@/views/provider/resource-detail/index.vue'),
      meta: {
        icon: 'mi:computer',
        title: t('routes.provider.resourceDetail'),
      },
    },
    {
      path: 'personal',
      name: 'PersonalCenter',
      component: () => import('/@/views/provider/personal-center/index.vue'),
      meta: {
        icon: 'ant-design:user-outlined',
        title: t('routes.provider.personalCenter'),
      },
    },
  ],
};

export default dashboard;
