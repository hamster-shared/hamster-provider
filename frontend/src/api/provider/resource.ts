import { defHttp } from '/@/utils/http/axios';
import {ComputingResource} from "/@/api/provider/model/resourceModel";

enum Api {
  ChainResource = '/api/v1/chain/resource',
  ResourceExpirationTime = '/api/v1/chain/expiration-time',
  changePrice = '/api/v1/resource/modify-price',
}

// get chainInfo
export const getResourceInfoApi = () => {
  return defHttp.get<ComputingResource>({ url:Api.ChainResource})
};

//get Resource Expiration time
export const getExpirationTimeApi = (expireBlock: number) => {
  return defHttp.get({ url: Api.ResourceExpirationTime, params: { expireBlock: expireBlock } });
};

//change price
export const changePriceApi = (price: number) => {
  return defHttp.post({ url: Api.changePrice, data: { price: price } });
};
