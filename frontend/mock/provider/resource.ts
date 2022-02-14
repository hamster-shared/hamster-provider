import {resultSuccess} from "../_util";
import {MockMethod} from "vite-plugin-mock";

const resource = {
  "index": 2,
  "accountId": [
    214,
    187,
    86,
    189,
    103,
    202,
    168,
    219,
    156,
    120,
    169,
    24,
    111,
    135,
    225,
    107,
    85,
    59,
    237,
    189,
    0,
    71,
    27,
    211,
    146,
    68,
    165,
    247,
    243,
    132,
    106,
    28
  ],
  "peerId": "QmYybVsLX5jewh7P7obSAgPnjF9nuAtzNwWyRZcP2TskjM",
  "config": {
    "cpu": 1,
    "memory": 1,
    "system": "Ubuntu 18",
    "cpuModel": "Intel(R) Core(TM) i5-1038NG7 CPU @ 2.00GHz"
  },
  "rentalStatistics": {
    "rentalCount": 0,
    "rentalDuration": 0,
    "faultCount": 0,
    "faultDuration": 0
  },
  "rentalInfo": {
    "rentUnitPrice": 123,
    "rentDuration": 143400,
    "endOfRent": 144782
  },
  "status": {
    "isInuse": false,
    "isLocked": false,
    "isUnused": true,
    "isOffline": false
  }
}

export default [
  // {
  //   url: '/api/v1/chain/resource',
  //   timeout: 1000,
  //   method: 'get',
  //   response: () => {
  //     return resultSuccess(resource)
  //   },
  // },
  // {
  //   url: '/api/v1/chain/price',
  //   timeout: 2000,
  //   method: 'post',
  //   response: ()=> {
  //     return resultSuccess("")
  //   }
  // }

] as MockMethod[];
