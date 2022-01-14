<template>
  <PageWrapper title="基础配置">
    <CollapseContainer title="链设置" :canExpan="false">
      <a-row :gutter="24">
        <a-col :span="18">
          <BasicForm @register="chainRegister" />
        </a-col>
      </a-row>
    </CollapseContainer>
    <CollapseContainer class="mt-5" title="虚拟规格设置" :canExpan="false">
      <a-row :gutter="24">
        <a-col :span="18">
          <BasicForm @register="register" />
        </a-col>
      </a-row>
    </CollapseContainer>

    <Description
      class="mt-4"
      title="虚拟机镜像"
      :column="3"
      :data="mockData"
      :schema="imageSchema"
    />

    <Button class="mt-4" type="primary" @click="handleSubmit"> 更新基本信息 </Button>
  </PageWrapper>
</template>
<script lang="ts">
  import { computed, defineComponent,reactive, onMounted,h } from 'vue';
  import { Button, Row, Col } from 'ant-design-vue';
  import { BasicForm, useForm } from '/@/components/Form/index';
  import { CollapseContainer } from '/@/components/Container';
  import { PageWrapper } from '/@/components/Page';
  import { Description,DescItem} from '/@/components/Description';

  import { useMessage } from '/@/hooks/web/useMessage';

  import headerImg from '/@/assets/images/header.jpg';
  import { getConfigApi,setConfigApi } from '/@/api/provider/initialization';
  import { vmSchemas,chainSchemas } from './data';
  import { useUserStore } from '/@/store/modules/user';
  import {ProviderConfig} from "/@/api/provider/model/settingModel";

  export default defineComponent({
    components: {
      BasicForm,
      CollapseContainer,
      Button,
      ARow: Row,
      ACol: Col,
      PageWrapper,
      Description,
    },
    setup() {
      const { createMessage } = useMessage();
      const userStore = useUserStore();

      const mockData = reactive<Recordable>({
        windows: false,
        ubuntu: false,
      });

      const imageSchema: DescItem[] = [
        {
          field: 'windows',
          label: 'windows',
          render: (curVal) => {
            if (!curVal ) {
              return h(Button,{
                type: "primary",
                onClick: consoleDebug,
                class: "ml-4",
              },() => "下载")
            } else {
              return "下载完成"
            }

          },
        },
        {
          field: 'ubuntu',
          label: 'ubuntu',
          render: (curVal) => {
            if (!curVal ) {
              return h(Button,{
                type: "primary",
                onClick: consoleDebug,
                class: "ml-4",
              },() => "下载")
            } else {
              return "下载完成"
            }
          },
        },
      ];

      const [chainRegister,{setFieldsValue : chainSetFieldsValue,validateFields: chainValidateFields}] = useForm({
        labelWidth: 120,
        schemas: chainSchemas,
        showActionButtonGroup: false,
      })

      const [register, { validateFields,setFieldsValue }] = useForm({
        labelWidth: 120,
        schemas: vmSchemas,
        showActionButtonGroup: false,
      });

      onMounted(async () => {
        const data = await getConfigApi();
        await setFieldsValue(data.vm);
        await chainSetFieldsValue({"address": data.chainApi,"account": data.seedOrPhrase});
      });

      const avatar = computed(() => {
        const { avatar } = userStore.getUserInfo;
        return avatar || headerImg;
      });

      const consoleDebug = function() {
        mockData.windows = true
        mockData.ubuntu = true
      }

      return {
        avatar,
        register,
        chainRegister,
        mockData,
        imageSchema,
        handleSubmit: () => {
          createMessage.success('更新成功！');
          Promise.all([validateFields(),chainValidateFields()]).then(data => {
            let values = data[0]
            let chainValues = data[1]
            console.log(values)
            let config: ProviderConfig =  {
              vm: values,
              chainApi : chainValues.address,
              seedOrPhrase: chainValues.account,
            }
            setConfigApi(config).then(() => {
              createMessage.success('更新成功！');
            })
          }).catch(err => {
            createMessage.error('校验失败',err)
          })
        },
      };
    },
  });
</script>

<style lang="less" scoped>

</style>
