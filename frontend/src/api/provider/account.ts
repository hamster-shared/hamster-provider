import { defHttp } from '/@/utils/http/axios';

enum Api {
  AccountInfo = '/api/v1/chain/account-info',
  StakingInfo = '/api/v1/chain/staking-info',
  StakingAmount = '/api/v1/chain/pledge',
  WithdrawAmount = '/api/v1/chain/withdraw-amount',
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
