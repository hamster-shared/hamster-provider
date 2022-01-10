import { MockMethod } from 'vite-plugin-mock';
import { resultSuccess } from '../_util';

const vmConfig = {
  cpu: 1,
  mem: 2,
  disk: 50,
  system: 'windows 10',
  image: 'windows 10.ios',
  accessPort: 3389,
  // 虚拟化类型,docker/kvm
  type: 'kvm',
};

const providerConfig  = {
  chainApi: "ws://127.0.0.1:9944",
  seedOrPhrase:　"",
  vm : vmConfig,
}


//
// export default [
//   {
//     url: '/api/v1/config/settting',
//     timeout: 1000,
//     method: 'get',
//     response: () => {
//       return resultSuccess(providerConfig);
//     },
//   },
//   {
//     url: '/api/v1/config/settting',
//     timeout: 1000,
//     method: 'post',
//     response: () => {
//       return resultSuccess({});
//     },
//   },
// ] as MockMethod[];
