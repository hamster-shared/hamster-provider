import { defHttp } from '/@/utils/http/axios';

enum Api {
  AccountInfo = '/v1/chain/account-info',
  StakingInfo = '/v1/chain/staking-info',
  StakingAmount = '/v1/chain/pledge',
  WithdrawAmount = '/v1/chain/withdraw-amount',
}

//get account info
export const getAccountInfoApi = () => {
  return defHttp.get({ url: Api.AccountInfo });
};

//get staking info
export const getStakingInfoApi = () => {
  return defHttp.get({ url: Api.StakingInfo });
};

// staking amount
export const stakingAmountApi = (price: number) => {
  return defHttp.post({ url: Api.StakingAmount, data: { price: price } });
};

//withdraw amount
export const withdrawAmountApi = (price: number) => {
  return defHttp.post({ url: Api.WithdrawAmount, data: { price: price } });
};
