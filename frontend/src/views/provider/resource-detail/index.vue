<template>
  <PageWrapper title="资源详情">
    <a-button type="primary" icon="download" size="small" />
    <h1>test</h1>
    <Description @register="register" class="mt-4" />
  </PageWrapper>
</template>

<script lang="ts" setup>
import { ref,onMounted } from 'vue';
import { PageWrapper } from '/@/components/Page';
import {getResourceInfoApi} from '/@/api/provider/resource'
import { Description, useDescription } from '/@/components/Description/index';
import {resourceSchemas} from "/@/views/provider/resource-detail/data";

const resourceData = ref({})

const [register] = useDescription({
  title: '资源信息',
  data: resourceData,
  schema: resourceSchemas,
});

onMounted(() => {
  getResourceInfoApi().then(data => {
    console.log(data)
    resourceData.value = data
  })
})
</script>

