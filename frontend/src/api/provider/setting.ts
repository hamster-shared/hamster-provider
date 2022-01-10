import { defHttp } from '/@/utils/http/axios';
import {ProviderConfig,} from '/@/api/provider/model/settingModel';

enum Api {
  Setting = '/api/v1/config/settting',
}

//获取系统配置
export const getConfigApi =  () => {
  return defHttp.get<ProviderConfig>({ url: Api.Setting });
}

// 修改配置
export const setConfigApi =  (config: ProviderConfig) => {
  return defHttp.post({ url: Api.Setting, params: config })
}
