import { defHttp } from '/@/utils/http/axios';
import {ComputingResource} from "/@/api/provider/model/resourceModel";

enum Api {
  ChainResource = '/v1/chain/resource',
}

// get chainInfo
export const getResourceInfoApi = () => {
  return defHttp.get<ComputingResource>({url:Api.ChainResource})
}
