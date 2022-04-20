import { defHttp } from '/@/utils/http/axios';
import { ComputingResource, UnitPriceParam} from '/@/api/provider/model/resourceModel';

enum Api {
  ChainResource = '/v1/chain/resource',
  ResourceExpirationTime = '/v1/chain/expiration-time',
  changePrice = '/v1/resource/modify-price',
  addDuration = '/v1/resource/add-duration',
  rentAgain = '/v1/resource/rent-again',
  receiveIncome = '/v1/resource/receive-income',
  deleteResource = '/v1/resource/delete-resource',
  judge = '/v1/resource/receive-income-judge',
  ModifyPrice = '/v1/chain/price',
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

// modify unit price
export const modifyUintPriceApi = (unitprice: UnitPriceParam) => {
  return defHttp.post({url: Api.ModifyPrice,params: unitprice})
}
