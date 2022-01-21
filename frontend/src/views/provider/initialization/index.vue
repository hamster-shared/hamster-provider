<template>
  <PageWrapper :title="t('initialization.initialization.basicConfiguration')">
    <CollapseContainer :title="t('initialization.initialization.chainSettings')" :canExpan="false">
      <a-row :gutter="24">
        <a-col :span="18">
          <BasicForm @register="chainRegister" />
        </a-col>
      </a-row>
    </CollapseContainer>
    <CollapseContainer class="mt-5" :title="t('initialization.initialization.virtualSpecificationSettings')" :canExpan="false">
      <a-row :gutter="24">
        <a-col :span="18">
          <BasicForm @register="register" />
        </a-col>
      </a-row>
    </CollapseContainer>

    <Description
      class="mt-4"
      :title="t('initialization.initialization.virtualMachineImage')"
      :column="3"
      :data="mockData"
      :schema="imageSchema"
    />

    <Button class="mt-4" type="primary" @click="handleSubmit"> {{ t('initialization.initialization.updateInformation') }} </Button>
  </PageWrapper>
</template>
<script lang="ts">
  import { computed, defineComponent,reactive, onMounted,h } from 'vue';
  import { Button, Row, Col } from 'ant-design-vue';
  import { BasicForm, useForm } from '/@/components/Form/index';
  import { CollapseContainer } from '/@/components/Container';
  import { PageWrapper } from '/@/components/Page';
  import { Description,DescItem} from '/@/components/Description';
  import { useI18n } from '/@/hooks/web/useI18n';
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
      const { t } = useI18n();
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
              },() => t('initialization.initialization.download'))
            } else {
              return t('initialization.initialization.downloadComplete');
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
              },() => t('initialization.initialization.download'))
            } else {
              return t('initialization.initialization.downloadComplete');
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
        t,
        avatar,
        register,
        chainRegister,
        mockData,
        imageSchema,
        handleSubmit: () => {
          createMessage.success(t('initialization.initialization.updateSucceeded'));
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
              createMessage.success(t('initialization.initialization.updateSucceeded'));
            })
          }).catch(err => {
            createMessage.error(t('initialization.initialization.verificationFailed'),err)
          })
        },
      };
    },
  });
</script>

<style lang="less" scoped>

</style>
