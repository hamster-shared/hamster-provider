<template>
  <div>
    <PageWrapper :title="t('initialization.initialization.basicConfiguration')">
      <CollapseContainer
        :title="t('initialization.initialization.chainSettings')"
        :canExpan="false"
      >
        <a-row :gutter="24">
          <a-col :span="18">
            <BasicForm @register="chainRegister" style="margin-top: 14px" />
          </a-col>
        </a-row>
        <a-row :gutter="24">
          <a-col :span="18">
            <div class="flex ml-5">
              <span class="mt-1 mr-1">Bootstraps:</span>
              <div>
                <div class="strap-style" v-for="(item, index) in bootstraps" :key="index">
                  <span class="mt-1 mr-1 break-all">{{ item }}</span>
                  <a-button
                    @click="removeBootstraps(index)"
                    type="primary"
                    size="small"
                    danger
                    class="mt-1 bootstraps-button"
                  >
                    {{ t('initialization.initialization.delete') }}
                  </a-button>
                </div>
                <a-button
                  @click="showAddModel"
                  type="primary"
                  size="small"
                  :disabled="bootstraps.length >= 5"
                  class="bootstraps-button mt-1"
                >
                  {{ t('initialization.initialization.add') }}
                </a-button>
              </div>
            </div>
          </a-col>
        </a-row>
      </CollapseContainer>
      <CollapseContainer
        class="mt-5"
        :title="t('initialization.initialization.virtualSpecificationSettings')"
        :canExpan="false"
        v-show="true"
      >
        <a-row :gutter="24">
          <a-col :span="18">
            <BasicForm @register="register" style="margin-top: 14px" />
          </a-col>
        </a-row>
      </CollapseContainer>

      <!--      <Description-->
      <!--        class="mt-4"-->
      <!--        :title="t('initialization.initialization.virtualMachineImage')"-->
      <!--        :column="3"-->
      <!--        :data="mockData"-->
      <!--        :schema="imageSchema"-->
      <!--      />-->
      <div class="text-center">
        <Button
          class="mt-8 submit-button"
          type="primary"
          @click="handleSubmit"
          :loading="submitLoading"
        >
          {{ t('initialization.initialization.updateInformation') }}
        </Button>
      </div>
    </PageWrapper>
    <a-modal
      v-model:visible="visible"
      :title="t('initialization.initialization.addGatewayNode')"
      :maskClosable="false"
      :footer="null"
      :centered="true"
      :closable="false"
    >
      <a-spin :spinning="addLoading">
        <div class="staking-content">
          <span class="title">{{ t('initialization.initialization.gatewayNode') }}</span>
          <a-textarea
            v-model:value="value"
            :placeholder="t('initialization.initialization.inputGatewayNodeTip')"
            :rows="3"
            @change="checkAddBootstrap"
          />
        </div>
        <span class="form-error-tip" v-if="addBootstrapTip">{{
          t('initialization.initialization.gatewayNodeTip')
        }}</span>
        <div class="staking-footer">
          <a-button class="staking-btn-close" @click="close">{{
            t('accountInfo.info.cancel')
          }}</a-button>
          <a-button class="staking-btn-ok" @click="ok">{{
            t('accountInfo.info.determine')
          }}</a-button>
        </div>
      </a-spin>
    </a-modal>
  </div>
</template>
<script lang="ts">
  import { computed, defineComponent, reactive, onMounted, h, ref, toRefs } from 'vue';
  import { Button, Row, Col } from 'ant-design-vue';
  import { BasicForm, useForm } from '/@/components/Form/index';
  import { CollapseContainer } from '/@/components/Container';
  import { PageWrapper } from '/@/components/Page';
  import { Description, DescItem } from '/@/components/Description';
  import { useI18n } from '/@/hooks/web/useI18n';
  import { useMessage } from '/@/hooks/web/useMessage';

  import headerImg from '/@/assets/images/header.jpg';
  import { getConfigApi, setConfigApi } from '/@/api/provider/initialization';
  import { vmSchemas, chainSchemas } from './data';
  import { useUserStore } from '/@/store/modules/user';
  import { ProviderConfig } from '/@/api/provider/model/settingModel';
  import AButton from '/@/components/Button/src/BasicButton.vue';

  export default defineComponent({
    components: {
      AButton,
      BasicForm,
      CollapseContainer,
      Button,
      ARow: Row,
      ACol: Col,
      PageWrapper,
      Description,
    },
    setup: function () {
      const { createMessage } = useMessage();
      const userStore = useUserStore();
      const { t } = useI18n();
      const submitLoading = ref(false);
      const mockData = reactive<Recordable>({
        windows: false,
        ubuntu: false,
      });
      const state = reactive({
        bootstraps: [] as string[],
        visible: false,
        addLoading: false,
        addBootstrapTip: false,
        value: '',
      });
      const imageSchema: DescItem[] = [
        {
          field: 'windows',
          label: 'windows',
          render: (curVal) => {
            if (!curVal) {
              return h(
                Button,
                {
                  type: 'primary',
                  onClick: consoleDebug,
                  class: 'ml-4',
                },
                () => t('initialization.initialization.download'),
              );
            } else {
              return t('initialization.initialization.downloadComplete');
            }
          },
        },
        {
          field: 'ubuntu',
          label: 'ubuntu',
          render: (curVal) => {
            if (!curVal) {
              return h(
                Button,
                {
                  type: 'primary',
                  onClick: consoleDebug,
                  class: 'ml-4',
                },
                () => t('initialization.initialization.download'),
              );
            } else {
              return t('initialization.initialization.downloadComplete');
            }
          },
        },
      ];

      const [
        chainRegister,
        { setFieldsValue: chainSetFieldsValue, validateFields: chainValidateFields },
      ] = useForm({
        labelWidth: 120,
        schemas: chainSchemas,
        showActionButtonGroup: false,
      });

      const [register, { validateFields, setFieldsValue }] = useForm({
        labelWidth: 120,
        schemas: vmSchemas,
        showActionButtonGroup: false,
      });

      onMounted(async () => {
        const data = await getConfigApi();
        state.bootstraps = data.bootstraps;
        await setFieldsValue(data.vm);
        await chainSetFieldsValue({ address: data.chainApi, account: data.seedOrPhrase });
      });

      const avatar = computed(() => {
        const { avatar } = userStore.getUserInfo;
        return avatar || headerImg;
      });

      const consoleDebug = function () {
        mockData.windows = true;
        mockData.ubuntu = true;
      };
      const showAddModel = function () {
        if (state.bootstraps.length >= 5) {
          createMessage.warning(t('initialization.initialization.nodeTip'));
          return;
        }
        state.visible = true;
      };
      const removeBootstraps = function (index) {
        state.bootstraps.splice(index, 1);
      };
      const checkAddBootstrap = function () {
        if (state.value === '') {
          state.addBootstrapTip = true;
          return;
        } else {
          state.addBootstrapTip = false;
        }
      };
      const close = function () {
        state.value = '';
        state.visible = false;
      };
      const ok = function () {
        checkAddBootstrap();
        if (state.value === '') {
          return;
        }
        state.bootstraps.push(state.value);
        close();
      };
      return {
        t,
        avatar,
        register,
        chainRegister,
        mockData,
        imageSchema,
        ...toRefs(state),
        removeBootstraps,
        checkAddBootstrap,
        close,
        ok,
        showAddModel,
        handleSubmit: () => {
          // if (state.bootstraps.length === 0) {
          //   createMessage.error(t('initialization.initialization.gatewayNodeTip'));
          //   return;
          // }
          submitLoading.value = true;
          Promise.all([validateFields(), chainValidateFields()])
            .then((data) => {
              let values = data[0];
              let chainValues = data[1];
              let config: ProviderConfig = {
                vm: values,
                chainApi: chainValues.address,
                seedOrPhrase: chainValues.account,
                bootstraps: state.bootstraps,
              };
              setConfigApi(config)
                .then(() => {
                  createMessage.success(t('initialization.initialization.updateSucceeded'));
                })
                .catch((err) => {
                  createMessage.error(t('initialization.initialization.updateFailed'), err);
                });
            })
            .catch((err) => {
              createMessage.error(t('initialization.initialization.verificationFailed'), err);
            })
            .finally(() => {
              submitLoading.value = false;
            });
        },
        submitLoading,
      };
    },
  });
</script>

<style lang="less" scoped>
  .strap-style {
    display: flex;
    align-items: center;
  }
  :deep(.ant-select:not(.ant-select-customize-input) .ant-select-selector),
  :deep(.ant-input-affix-wrapper) {
    border-radius: 8px !important;
  }
  :deep(.vben-collapse-container__header) {
    border-bottom: none;
    .vben-basic-title {
      font-size: 16px;
      border-bottom: 1px solid #eee;
      width: 100%;
      padding: 56px 4px 20px 4px;
      color: #070707;
      font-weight: bold;
      margin: 16px 12px 24px 12px;
    }
  }
  :deep(.vben-collapse-container__body) {
    padding: 30px 0;
  }
  .staking-content {
    display: flex;
    align-items: center;
    margin-top: 24px;
    padding: 0px 16px;
    .title {
      min-width: 40px;
      margin-right: 8px;
      color: rgba(0, 0, 0, 0.85);
    }
  }
  .form-error-tip {
    color: #f5313d;
    font-style: normal;
    font-weight: normal;
    font-size: 10px;
    line-height: 17px;
    margin-left: 80px;
  }
  .staking-footer {
    margin-top: 24px;
    display: grid;
    padding: 0px 16px 24px 16px;
    grid-template-columns: 1fr 1fr;
    gap: 12px;
    .staking-btn-close {
      width: 100%;
    }
    .staking-btn-ok {
      background-color: rgb(24, 144, 255);
      color: white;
    }
  }

  .bootstraps-button {
    border-radius: 8px;
  }
  .submit-button {
    border-radius: 16px;
  }
  :deep(.ant-input-number) {
    width: 100%;
    border-radius: 8px;
    .ant-input-number-handler-wrap {
      border-radius: 8px !important;
    }
  }
</style>
