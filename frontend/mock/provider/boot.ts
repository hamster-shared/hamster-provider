import {resultSuccess} from "../_util";
import {MockMethod} from "vite-plugin-mock";

export default [
  {
    url: '/api/v1/config/boot',
    timeout: 1000,
    method: 'post',
    response: () => {
      return resultSuccess('')
    },
  },
  {
    url: '/api/v1/config/boot',
    timeout: 1000,
    method: 'get',
    response: () => {
      return resultSuccess(false)
    },
  },
] as MockMethod[];
