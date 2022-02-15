import { defHttp } from '/@/utils/http/axios';

enum Api {
  Boot = '/api/v1/config/boot',
}

// set boot state
export const setBootStateApi =  (option: boolean) => {
  return defHttp.post({ url: Api.Boot, params: {option: option} })
}

// get boot state
export const getBootStateApi = () => {
  return defHttp.get<boolean>({url:Api.Boot})
}
