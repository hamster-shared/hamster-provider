import { defHttp } from '/@/utils/http/axios';
import { ComputingResource } from '/@/api/provider/model/resourceModel';

enum Api {
  ChainResource = '/api/v1/chain/resource',
  ResourceExpirationTime = '/api/v1/chain/expiration-time',
  changePrice = '/api/v1/resource/modify-price',
  addDuration = '/api/v1/resource/add-duration',
  rentAgain = '/api/v1/resource/rent-again',
  receiveIncome = '/api/v1/resource/receive-income',
  deleteResource = '/api/v1/resource/delete-resource',
  judge = '/api/v1/resource/receive-income-judge',
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

//add duration
export const addDurationAPi = (duration: number) => {
  return defHttp.post({ url: Api.addDuration, data: { duration: duration } });
};

//rent again
export const rentAgainApi = () => {
  return defHttp.post({ url: Api.rentAgain });
};

//receive income
export const receiveIncomeApi = () => {
  return defHttp.post({ url: Api.receiveIncome });
};

// delete resource
export const deleteResourceApi = () => {
  return defHttp.post({ url: Api.deleteResource });
};

export const judgeReceiveIncomeApi = () => {
  return defHttp.get({ url: Api.judge });
};
