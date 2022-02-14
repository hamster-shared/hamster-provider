import { defHttp } from '/@/utils/http/axios';
import {ComputingResource, UnitPriceParam} from "/@/api/provider/model/resourceModel";

enum Api {
  ChainResource = '/v1/chain/resource',
  ModifyPrice = '/v1/chain/price'
}

// get chainInfo
export const getResourceInfoApi = () => {
  return defHttp.get<ComputingResource>({url:Api.ChainResource})
}

// modify unit price
export const modifyUintPriceApi = (unitprice: UnitPriceParam) => {
  return defHttp.post({url: Api.ModifyPrice,params: unitprice})
}
