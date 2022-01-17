<template>
  <div>
    <PageWrapper :title="t('resourceDetail.detail.resourceDetail')" contentBackground>
      <template #extra>
        <a-button type="primary" @click="changePriceModal">
          {{ t('resourceDetail.detail.changePrice') }}
        </a-button>
        <a-button type="primary" @click="increaseDurationModal">
          {{ t('resourceDetail.detail.increaseＤuration') }}
        </a-button>
        <a-button type="primary" @click="receiveIncome">
          {{ t('resourceDetail.detail.receiveＢenefits') }}</a-button
        >
        <a-button type="primary" @click="rentAgain">
          {{ t('resourceDetail.detail.rentＡgain') }}</a-button
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
            {{ resourceData.rentalInfo ? resourceData.rentalInfo.rentUnitPrice : '' }}
          </a-descriptions-item>
          <a-descriptions-item :label="t('resourceDetail.detail.cpuCounts')">
            {{ resourceData.config ? resourceData.config.cpu + '核' : '' }}
          </a-descriptions-item>
          <a-descriptions-item :label="t('resourceDetail.detail.expireDate')">
            {{ resourceData.rentalInfo ? resourceData.rentalInfo.endOfRent : '' }}
          </a-descriptions-item>
          <a-descriptions-item :label="t('resourceDetail.detail.memory')">
            {{ resourceData.config ? resourceData.config.memory + 'Ｇ' : '' }}
          </a-descriptions-item>
        </a-descriptions>
      </div>
    </PageWrapper>
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
          <BalanceInput
            :placeholder="t('resourceDetail.detail.unitPriceText')"
            :changeClick="checkPrice('changePrice')"
            ref="inputRef"
          />
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
      :title="t('resourceDetail.detail.changeUnitPrice')"
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
            @change="checkPrice()"
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
  import { getResourceInfoApi } from '/@/api/provider/resource';
  import { PageWrapper } from '/@/components/Page';
  import { ResourceStatus } from '/@/api/provider/model/resourceModel';
  import BalanceInput from '/@/components/BalanceInput/index.vue';
  import * as cluster from "cluster";
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
  onMounted(() => {
    getResourceInfoApi().then((data) => {
      resourceData.value = data;
    });
  });
  function changePriceModal() {
    state.changePriceVisible = true;
    state.changePriceTip = false;
  }
  function increaseDurationModal() {
    state.increaseDurationVisible = true;
    state.increaseDurationTip = false;
  }
  function receiveIncome() {}
  function rentAgain() {}
  function changePriceCancel() {
    state.changePriceVisible = false;
    proxy.$refs.inputRef.value = '';
    proxy.$refs.inputRef.uintPower = 0;
  }
  function increaseDurationCancel() {
    state.increaseDurationVisible = false;
    state.duration = '';
  }
  function changePriceOk() {
    state.changePriceVisible = false;
  }
  function increaseDurationOk() {
    state.increaseDurationVisible = false;
  }
  function checkPrice(params: string) {
    if (params === 'changePrice') {
      if (proxy.$refs.inputRef) {
        if (proxy.$refs.inputRef.value === '') {
          state.changePriceTip = true;
          return;
        } else {
          state.changePriceTip = false;
        }
        proxy.$refs.inputRef.value = proxy.$refs.inputRef.value.replace(/\D/g, '');
      }
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
