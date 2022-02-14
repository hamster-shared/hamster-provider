<template>
  <PageWrapper title="资源详情">
    <Description @register="register" class="mt-4" />

    <Modal4 :ref="$modal4" @register="register4" />
  </PageWrapper>
</template>

<script lang="ts" setup>
import {ref, onMounted, h} from 'vue';
import { Button } from 'ant-design-vue';
import { PageWrapper } from '/@/components/Page';
import {getResourceInfoApi} from '/@/api/provider/resource'
import {DescItem, Description, useDescription} from '/@/components/Description/index';
import {
  displayResourceStatus,
} from "/@/views/provider/resource-detail/data";
import { useModal } from '/@/components/Modal';
import Modal4 from './Modal4.vue';

const resourceData = ref({})

const $modal4 = ref<InstanceType<typeof Modal4>>()

const [register4, { openModal: openModal4,closeModal: closeModal4 }] = useModal();

const resourceSchemas: DescItem[] = [
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
        h('span',{},'    '),
        h(Button,{onClick: editPrice},'edit'),
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

const [register] = useDescription({
  title: '资源信息',
  data: resourceData,
  schema: resourceSchemas,
});

const editPrice = function(){
  openModal4(true, {
    price: 100,
  });
}



onMounted(() => {
  getResourceInfoApi().then(data => {
    console.log(data)
    resourceData.value = data
  })
})
</script>

