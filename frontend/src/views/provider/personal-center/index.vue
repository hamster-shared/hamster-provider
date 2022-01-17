<template>
  <div>
    <PageWrapper :title="t('accountInfo.info.accountInfo')" contentBackground>
      <template #extra>
        <a-button type="primary" @click="showStakingModal('staking')">
          {{ t('accountInfo.info.pledge') }}
        </a-button>
        <a-button type="primary" @click="showStakingModal('withdraw')">
          {{ t('accountInfo.info.reclaimPledge') }}
        </a-button>
      </template>
      <div class="pt-4 m-4 desc-wrap">
        <a-descriptions :title="t('accountInfo.info.accountInfo')" size="small" :column="2">
          <a-descriptions-item :label="t('accountInfo.info.accountAddress')">
            fafdafaafa
          </a-descriptions-item>
          <a-descriptions-item :label="t('accountInfo.info.accountBalance')">
            fafdafaf
          </a-descriptions-item>
        </a-descriptions>
        <a-divider />
        <a-descriptions :title="t('accountInfo.info.pledgeInformation')" :column="2">
          <a-descriptions-item :label="t('accountInfo.info.totalPledgeAmount')">
            111
          </a-descriptions-item>
          <a-descriptions-item :label="t('accountInfo.info.activePledgeAmount')">
            2017-08-08
          </a-descriptions-item>
          <a-descriptions-item :label="t('accountInfo.info.lockedPledgeAmount')">
            725
          </a-descriptions-item>
        </a-descriptions>
      </div>
    </PageWrapper>
    <a-modal
      v-model:visible="visible"
      :title="title"
      :maskClosable="false"
      :footer="null"
      :centered="true"
      :closable="false"
    >
      <a-spin :spinning="stakingLoading">
        <div class="staking-content">
          <span class="title">{{ t('accountInfo.info.pledgeAmount') }}</span>
          <BalanceInput
            :placeholder="placeholder"
            :changeClick="checkStakingAmount"
            ref="inputRef"
          />
        </div>
        <span class="form-error-tip" v-if="stakingAmountTip">{{
          t('accountInfo.info.stakingAmountTip')
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
  import { PageWrapper } from '/@/components/Page';
  import BalanceInput from '../../../components/BalanceInput/index.vue';
  import { useI18n } from '/@/hooks/web/useI18n';
  import { defineComponent, getCurrentInstance, reactive, toRefs } from 'vue';
  // eslint-disable-next-line vue/no-export-in-script-setup
  export default defineComponent({
    name: 'center',
    components: {
      BalanceInput,
      PageWrapper,
    },
    setup: function () {
      const { t } = useI18n();
      const { proxy } = getCurrentInstance();
      const state = reactive({
        visible: false,
        title: '',
        stakingLoading: false,
        placeholder: '',
        stakingAmountTip: false,
        activeModal: '',
      });
      function checkStakingAmount() {
        if (proxy.$refs.inputRef.value === '') {
          state.stakingAmountTip = true;
          return;
        } else {
          state.stakingAmountTip = false;
        }
        proxy.$refs.inputRef.value = proxy.$refs.inputRef.value.replace(/\D/g, '');
      }
      function ok() {
        if (state.activeModal === 'staking') {
          stakingAmountClick();
        } else {
          withdrawAmount();
        }
      }
      function close() {
        state.visible = false;
        proxy.$refs.inputRef.value = '';
        proxy.$refs.inputRef.uintPower = 0;
      }
      //质押
      function stakingAmountClick() {
        let price = proxy.$refs.inputRef.getPrice();
        console.log(price);
      }
      //取回质押
      function withdrawAmount() {}
      function showStakingModal(params: string) {
        state.stakingLoading = false;
        state.visible = true;
        state.stakingAmountTip = false;
        if (params === 'staking') {
          state.title = t('accountInfo.info.pledgeAmountModal');
          state.placeholder = t('accountInfo.info.inputPledgeAmount');
          state.activeModal = 'staking';
        } else {
          state.title = t('accountInfo.info.withdrawPledgeAmount');
          state.placeholder = t('accountInfo.info.inputWithdrawTip');
          state.activeModal = 'withdraw';
        }
      }
      return {
        BalanceInput,
        checkStakingAmount,
        ok,
        close,
        showStakingModal,
        t,
        ...toRefs(state),
      };
    },
  });
</script>

<style lang="less" scoped>
  .staking-content {
    display: flex;
    align-items: center;
    margin-top: 24px;
    padding: 0px 16px;
    .title {
      width: 90px;
      color: rgba(0, 0, 0, 0.85);
    }
  }
  .form-error-tip {
    color: #f5313d;
    font-style: normal;
    font-weight: normal;
    font-size: 10px;
    line-height: 17px;
    margin-left: 90px;
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
