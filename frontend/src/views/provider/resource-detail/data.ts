import {ResourceStatus} from "/@/api/provider/model/resourceModel";


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

export function showEditDialog(){

  console.log('111')
}
