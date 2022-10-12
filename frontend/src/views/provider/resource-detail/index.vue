<template>
  <div>
    <a-spin :spinning="state.allLoading">
      <PageWrapper :title="t('resourceDetail.detail.resourceDetail')" contentBackground>
        <template #extra>
          <a-button type="primary" @click="increaseDurationModal" :disabled="state.disabled">
            {{ t('resourceDetail.detail.increaseDuration') }}
          </a-button>
          <a-button type="primary" @click="rentAgain" :disabled="state.disabled">
            {{ t('resourceDetail.detail.rentAgain') }}
          </a-button>
          <a-button type="primary" @click="deleteResource" :disabled="state.disabled">
            {{ t('resourceDetail.detail.deleteResource') }}
          </a-button>
        </template>
        <div class="pt-4 m-4 desc-wrap">
          <a-descriptions
            :title="t('resourceDetail.detail.resourceInformation')"
            size="small"
            layout="vertical"
          >
            <a-descriptions-item :label="t('resourceDetail.detail.resourceID')">
              {{ resourceData.index }}
            </a-descriptions-item>
            <a-descriptions-item :label="t('resourceDetail.detail.system')">
              {{ resourceData.config ? resourceData.config.system : '' }}
            </a-descriptions-item>
            <a-descriptions-item :label="t('resourceDetail.detail.resourceState')">
              {{ displayResourceStatus(resourceData.status) }}
            </a-descriptions-item>
            <a-descriptions-item :label="t('resourceDetail.detail.cpuModel')">
              {{ resourceData.config ? resourceData.config.cpuModel : '' }}
            </a-descriptions-item>
            <a-descriptions-item :label="t('resourceDetail.detail.cpuCounts')">
              {{ resourceData.config ? resourceData.config.cpu + 'Core' : '' }}
            </a-descriptions-item>
            <a-descriptions-item :label="t('resourceDetail.detail.expireDate')">
              {{ resourceData.expirationTime ? resourceData.expirationTime : '' }}
            </a-descriptions-item>
            <a-descriptions-item :label="t('resourceDetail.detail.memory')">
              {{ resourceData.config ? resourceData.config.memory + 'Ｇ' : '' }}
            </a-descriptions-item>
          </a-descriptions>
        </div>
      </PageWrapper>
    </a-spin>
    <a-modal
      v-model:visible="state.increaseDurationVisible"
      :title="t('resourceDetail.detail.increaseDuration')"
      :maskClosable="false"
      :footer="null"
      :centered="true"
      :closable="false"
    >
      <a-spin :spinning="state.increaseDurationLoading">
        <div class="staking-content">
          <span class="title">{{ t('resourceDetail.detail.unitPriceTitle') }}</span>
          <a-input
            v-model:value="state.duration"
            @change="checkPrice('changeDuration')"
            :placeholder="t('resourceDetail.detail.unitPriceText')"
          >
            <template #addonAfter>
              <span>/h</span>
            </template>
          </a-input>
        </div>
        <span class="form-error-tip" v-if="state.increaseDurationTip">{{
          t('resourceDetail.detail.inputUnitPrice')
        }}</span>
        <div class="staking-footer">
          <a-button class="staking-btn-close" @click="increaseDurationCancel">{{
            t('resourceDetail.detail.cancel')
          }}</a-button>
          <a-button class="staking-btn-ok" @click="increaseDurationOk">{{
            t('resourceDetail.detail.determine')
          }}</a-button>
        </div>
      </a-spin>
    </a-modal>
  </div>
</template>

<script lang="ts" setup>
  import { useMessage } from '/@/hooks/web/useMessage';
  import { ref, onMounted, reactive, getCurrentInstance } from 'vue';
  import { useI18n } from '/@/hooks/web/useI18n';
  import {
    getResourceInfoApi,
    getExpirationTimeApi,
    changePriceApi,
    addDurationAPi,
    rentAgainApi,
    receiveIncomeApi,
    deleteResourceApi,
    judgeReceiveIncomeApi,
  } from '/@/api/provider/resource';
  import { PageWrapper } from '/@/components/Page';
  import { ResourceStatus } from '/@/api/provider/model/resourceModel';
  import BigNumber from 'bignumber.js';
  import { formatToDateTime } from '/@/utils/dateUtil';
  import * as cluster from 'cluster';
  import { formatBalance } from '@polkadot/util';
  const { t } = useI18n();
  const resourceData = ref({});
  const state = reactive({
    increaseDurationVisible: false,
    increaseDurationLoading: false,
    increaseDurationTip: false,
    duration: '',
    value: '',
    uintPower: 0,
    uintOptions: [],
    allLoading: false,
    disabled: false,
    receiveJudge: false,
  });
  const { createMessage } = useMessage();
  const { proxy } = getCurrentInstance();
  function displayResourceStatus(status: ResourceStatus) {
    if (status === undefined || !status) {
      return '';
    }
    if (status.isInuse) {
      return t('resourceDetail.detail.inuse');
    } else if (status.isLocked) {
      return t('resourceDetail.detail.locked');
    } else if (status.isOffline) {
      return t('resourceDetail.detail.offline');
    } else if (status.isUnused) {
      return t('resourceDetail.detail.unused');
    }
  }
  function priceFormat(params: number) {
    return new BigNumber(params).div(new BigNumber(Math.pow(10, 12))).toNumber() + ' Uint';
  }
  onMounted(() => {
    getUintOptions();
    getResource();
    judge();
  });
  function getResource() {
    state.disabled = true;
    getResourceInfoApi()
      .then((data) => {
        console.log(data);
        getExpirationTimeApi(data.rentalInfo.endOfRent).then((res) => {
          data['expirationTime'] = formatToDateTime(new Date(new Date().getTime() + res));
          resourceData.value = data;
          state.disabled = false;
        });
      })
      .catch(() => {
        state.disabled = true;
        resourceData.value = '';
      });
  }
  function getUintOptions() {
    state.uintOptions.unshift({ power: -3, text: 'milli', value: '-' });
    state.uintOptions.unshift({ power: -6, text: 'micro', value: '-' });
    state.uintOptions.unshift({ power: -9, text: 'nano', value: '-' });
    state.uintOptions.unshift({ power: -12, text: 'pico', value: '-' });
  }
  function getPrice() {
    return new BigNumber(state.value)
      .times(new BigNumber(Math.pow(10, state.uintPower)))
      .times(new BigNumber(Math.pow(10, 12)))
      .toNumber();
  }
  function increaseDurationModal() {
    state.increaseDurationVisible = true;
    state.increaseDurationTip = false;
  }
  function judge() {
    judgeReceiveIncomeApi().then((data) => {
      state.receiveJudge = !data;
    });
  }
  function rentAgain() {
    state.allLoading = true;
    if (resourceData.value) {
      if (resourceData.value['status'].isOffline) {
        rentAgainApi()
          .then(() => {
            getResource();
            createMessage.success(t('resourceDetail.detail.rentAgainSuccess'));
            state.allLoading = false;
          })
          .catch(() => {
            createMessage.error(t('resourceDetail.detail.rentAgainFailed'));
            state.allLoading = false;
          });
      } else {
        createMessage.warning(t('resourceDetail.detail.rentAgainWarn'));
        state.allLoading = false;
      }
    } else {
      state.allLoading = false;
    }
  }
  function deleteResource() {
    state.allLoading = true;
    if (resourceData.value) {
      if (resourceData.value['status'].isInuse || resourceData.value['status'].isLocked) {
        createMessage.warning(t('resourceDetail.detail.deleteResourceWarn'));
        state.allLoading = false;
        return;
      }
      deleteResourceApi()
        .then(() => {
          createMessage.success(t('resourceDetail.detail.deleteResourceSuccess'));
          getResource();
          state.allLoading = false;
        })
        .catch(() => {
          createMessage.error(t('resourceDetail.detail.deleteResourceFailed'));
          state.allLoading = false;
        });
    } else {
      state.allLoading = false;
    }
  }
  function increaseDurationCancel() {
    state.increaseDurationVisible = false;
    state.duration = '';
  }
  function increaseDurationOk() {
    checkPrice('changeDuration');
    state.increaseDurationLoading = true;
    addDurationAPi(parseInt(state.duration))
      .then(() => {
        state.increaseDurationLoading = false;
        createMessage.success(t('resourceDetail.detail.increaseDurationSucceeded'));
        getResource();
        increaseDurationCancel();
      })
      .catch(() => {
        state.increaseDurationLoading = false;
        createMessage.error(t('resourceDetail.detail.increaseDurationFailed'));
      });
    // state.increaseDurationVisible = false;
  }
  function checkPrice(params: string) {
    if (params === 'changePrice') {
      if (state.value === '') {
        state.changePriceTip = true;
        return;
      } else {
        state.changePriceTip = false;
      }
      state.value = state.value.replace(/\D/g, '');
    } else {
      if (state.duration === '') {
        state.increaseDurationTip = true;
        return;
      } else {
        state.increaseDurationTip = false;
      }
      state.duration = state.duration.replace(/\D/g, '');
    }
  }
</script>
<style lang="less" scoped>
  :deep(.vben-page-wrapper .vben-page-wrapper-content) {
    border-radius: 6px;
  }
  .desc-wrap {
    padding: 16px;
    background-color: @component-background;
  }
  .staking-content {
    display: flex;
    align-items: center;
    margin-top: 24px;
    padding: 0px 16px;
    .title {
      min-width: 40px;
      color: rgba(0, 0, 0, 0.85);
    }
  }
  .form-error-tip {
    color: #f5313d;
    font-style: normal;
    font-weight: normal;
    font-size: 10px;
    line-height: 17px;
    margin-left: 55px;
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

  :deep(.ant-descriptions-header) {
    border-bottom: 1px solid #eee;
    padding-bottom: 20px;
    margin-bottom: 10px;
    .ant-descriptions-title {
      margin-top: -12px;
    }
  }
  :deep(.ant-descriptions-item-container) {
    .ant-descriptions-item-label {
      color: #2e3c43;
      margin-top: 28px;
      &::after {
        content: '';
      }
    }
    .ant-descriptions-item-content {
      color: #222;
      line-height: 1;
      font-size: 16px;
      font-weight: 500;
    }
  }
</style>
