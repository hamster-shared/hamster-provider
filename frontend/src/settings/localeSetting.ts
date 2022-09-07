import type { DropMenu } from '../components/Dropdown';
import type { LocaleSetting, LocaleType } from '/#/config';

export const LOCALE: { [key: string]: LocaleType } = {
  EN_US: 'en',
  ZH_CN: 'zh_CN',
};

export const localeSetting: LocaleSetting = {
  showPicker: false,
  // Locale
  locale: LOCALE.EN_US,
  // Default locale
  fallback: LOCALE.EN_US,
  // available Locales
  availableLocales: [LOCALE.EN_US,LOCALE.ZH_CN],
};

// locale list
export const localeList: DropMenu[] = [
  {
    text: 'English',
    event: LOCALE.EN_US,
  },
  {
    text: '简体中文',
    event: LOCALE.ZH_CN,
  },
];
