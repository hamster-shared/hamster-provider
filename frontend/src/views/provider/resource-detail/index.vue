<template>
  <div>
    <a-spin :spinning="state.allLoading">
      <PageWrapper :title="t('resourceDetail.detail.resourceDetail')" contentBackground>
        <template #extra>
          <a-button type="primary" @click="changePriceModal" :disabled="state.disabled">
            {{ t('resourceDetail.detail.changePrice') }}
          </a-button>
          <a-button type="primary" @click="increaseDurationModal" :disabled="state.disabled">
            {{ t('resourceDetail.detail.increaseＤuration') }}
          </a-button>
          <a-button type="primary" @click="rentAgain" :disabled="state.disabled">
            {{ t('resourceDetail.detail.rentＡgain') }}</a-button
          >
          <a-button type="primary" @click="deleteResource" :disabled="state.disabled">
            {{ t('resourceDetail.detail.deleteResource') }}</a-button
          >
          <a-button type="primary" @click="receiveIncome" :disabled="state.receiveJudge">
            {{ t('resourceDetail.detail.receiveＢenefits') }}</a-button
          >
        </template>
        <div class="pt-4 m-4 desc-wrap">
          <a-descriptions
            :title="t('resourceDetail.detail.resourceInformation')"
            size="small"
            :column="2"
          >
            <a-descriptions-item :label="t('resourceDetail.detail.resourceID')">
              {{ resourceData.index }}
            </a-descriptions-item>
            <a-descriptions-item :label="t('resourceDetail.detail.system')">
              {{ resourceData.config ? resourceData.config.system : '' }}
            </a-descriptions-item>
            <a-descriptions-item :label="t('resourceDetail.detail.resourceＳtate')">
              {{ displayResourceStatus(resourceData.status) }}
            </a-descriptions-item>
            <a-descriptions-item :label="t('resourceDetail.detail.cpuＭodel')">
              {{ resourceData.config ? resourceData.config.cpuModel : '' }}
            </a-descriptions-item>
            <a-descriptions-item :label="t('resourceDetail.detail.unitPrice')">
              {{ resourceData.rentalInfo ? priceFormat(resourceData.rentalInfo.rentUnitPrice) : '' }}
            </a-descriptions-item>
            <a-descriptions-item :label="t('resourceDetail.detail.cpuCounts')">
              {{ resourceData.config ? resourceData.config.cpu + '核' : '' }}
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
      v-model:visible="state.changePriceVisible"
      :title="t('resourceDetail.detail.changeUnitPrice')"
      :maskClosable="false"
      :footer="null"
      :centered="true"
      :closable="false"
    >
      <a-spin :spinning="state.changePriceLoading">
        <div class="staking-content">
          <span class="title">{{ t('resourceDetail.detail.unitPriceTitle') }}</span>
          <a-input
            v-model:value="state.value"
            :placeholder="t('resourceDetail.detail.unitPriceText')"
            @change="checkPrice('changePrice')"
          >
            <template #addonAfter>
              <a-select style="width: 90px" v-model:value="state.uintPower">
                <a-select-option
                  v-for="(item, index) in state.uintOptions"
                  v-model:value="item.power"
                  :key="index"
                >
                  {{ item.text }}
                </a-select-option>
              </a-select>
            </template>
          </a-input>
        </div>
        <span class="form-error-tip" v-if="state.changePriceTip">{{
          t('resourceDetail.detail.inputUnitPrice')
        }}</span>
        <div class="staking-footer">
          <a-button class="staking-btn-close" @click="changePriceCancel">{{
            t('resourceDetail.detail.cancel')
          }}</a-button>
          <a-button class="staking-btn-ok" @click="changePriceOk">{{
            t('resourceDetail.detail.determine')
          }}</a-button>
        </div>
      </a-spin>
    </a-modal>
    <a-modal
      v-model:visible="state.increaseDurationVisible"
      :title="t('resourceDetail.detail.increaseＤuration')"
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
    changePriceVisible: false,
    changePriceLoading: false,
    changePriceTip: false,
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
    state.uintOptions = formatBalance.getOptions();
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
  function changePriceModal() {
    state.changePriceVisible = true;
    state.changePriceTip = false;
  }
  function increaseDurationModal() {
    state.increaseDurationVisible = true;
    state.increaseDurationTip = false;
  }
  function receiveIncome() {
    state.allLoading = true;
    receiveIncomeApi()
      .then(() => {
        createMessage.success(t('resourceDetail.detail.receiveIncomeSuccess'));
        state.allLoading = false;
      })
      .catch(() => {
        createMessage.error(t('resourceDetail.detail.receiveIncomeFailed'));
        state.allLoading = false;
      });
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
  function changePriceCancel() {
    state.changePriceVisible = false;
    state.value = '';
    state.uintPower = 0;
  }
  function increaseDurationCancel() {
    state.increaseDurationVisible = false;
    state.duration = '';
  }
  function changePriceOk() {
    checkPrice('changePrice');
    state.changePriceLoading = true;
    let price = getPrice();
    changePriceApi(price)
      .then(() => {
        state.changePriceLoading = false;
        createMessage.success(t('resourceDetail.detail.changPriceSuccess'));
        changePriceCancel();
        getResource();
      })
      .catch(() => {
        state.changePriceLoading = false;
        changePriceCancel();
        createMessage.error(t('resourceDetail.detail.changPriceError'));
      });
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
</style>
