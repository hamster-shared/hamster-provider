import {DescItem} from "/@/components/Description";
import {ResourceStatus} from "/@/api/provider/model/resourceModel";
import {h} from 'vue'
import { Button } from 'ant-design-vue';

export const resourceSchemas: DescItem[] = [
  {
    field: 'index',
    label: '机器id',
  },
  {
    field: 'config.cpuModel',
    label: 'cpu型号',
  },
  {
    field: 'config.cpu',
    label: 'cpu核数',
  },
  {
    field: 'config.memory',
    label: '内存',
  },
  {
    field: 'config.system',
    label: '系统',
  },
  {
    field: 'rentalInfo.rentUnitPrice',
    label: '单价/h',
    render: function (v,data){
      return h("div",{},[
        v,
        h(Button,{type: 'primary',shape:'circle',icon:'search'}),

      ])
    }
  },
  {
    field: 'rentalInfo.endOfRent',
    label: '租用到期时间',
  },
  {
    label: '状态',
    render: function (v,data){
      return displayResourceStatus(data.status)
    }
  }
];

export function displayResourceStatus(status: ResourceStatus) {
  if (status == undefined){
    return ''
  }
  if (status.isInuse) {
    return '正在出租'
  }else if(status.isLocked){
    return '已锁定'
  }else if(status.isOffline){
    return '已离线'
  }else if(status.isUnused) {
    return '未使用'
  }
}
