<template>
  <PageWrapper title="boot">
    <template #headerContent>
      <div class="flex justify-between items-center">
        <span class="flex-1">
          <a href="#" target="_blank">{{ name }}</a>
          是一个基于p2p,kvm虚拟化
          的共享计算平台，旨在搭建共享计算平台
        </span>
      </div>
    </template>

    <div class="py-8 bg-white flex flex-col justify-center items-center">

      <div class="flex justify-center">
        <Switch v-model:checked="option"
                :loading="loading"
                :checked-children="t('layout.setting.on')"
                :un-checked-children="t('layout.setting.off')"
                @change="onChange"></Switch>
        <label>{{t('routes.provider.start_or_stop')}}</label>

      </div>
    </div>

    <Description @register="register" class="mt-4" />

  </PageWrapper>
</template>

<script lang="ts" setup>
import { ref,onMounted } from 'vue';
import { useI18n } from '/@/hooks/web/useI18n';
import { PageWrapper } from '/@/components/Page';
import { Switch } from 'ant-design-vue';
import {getBootStateApi,setBootStateApi} from '/@/api/provider/boot'
import { getConfigApi } from '/@/api/provider/initialization';
import { Description, useDescription } from '/@/components/Description/index';
import {vmSchemas} from "/@/views/provider/boot/data";

const { pkg, lastBuildTime } = __APP_INFO__;
const {  name } = pkg;
const { t } = useI18n();
const option = ref(true)
const loading = ref(false)

const onChange = function (checked){
  loading.value = true
  setBootStateApi(checked).then(() => {

  }).finally(() => {
    loading.value = false
  })
}

const vmData = ref({})

onMounted(() => {
  loading.value = true
  getBootStateApi().then(data => {
    option.value = data
  }).finally(() => {
    loading.value = false
  })
  getConfigApi().then(data => {
    vmData.value = data.vm
  })
})

const [register] = useDescription({
  title: '资源信息',
  data: vmData,
  schema: vmSchemas,
});

</script>

<style lang="less" scoped>
.extra {
  float: right;
  margin-top: 10px;
  margin-right: 30px;
}
</style>
