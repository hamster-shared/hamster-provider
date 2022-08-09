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
            {{ address }}
          </a-descriptions-item>
          <a-descriptions-item :label="t('accountInfo.info.accountBalance')">
            {{ amount }} Unit
          </a-descriptions-item>
        </a-descriptions>
        <a-divider />
        <a-descriptions :title="t('accountInfo.info.pledgeInformation')" :column="2">
          <a-descriptions-item :label="t('accountInfo.info.totalPledgeAmount')">
            {{ pledgeAmount }} Unit
          </a-descriptions-item>
        </a-descriptions>

        <a-divider />

        <a-descriptions :title="t('accountInfo.income.incomeInfo')" size="small" :column="2">
          <a-descriptions-item :label="t('accountInfo.income.income')">
            {{ reward }} Unit   <AButton style="margin-left: 1rem;" type="primary" size="small" shape="round" :loading="rewardLoading" @click="payoutReward"> {{t('accountInfo.income.withdraw')}} </AButton>
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
  import { getAccountInfoApi, getStakingInfoApi, stakingAmountApi, withdrawAmountApi,getRewardInfoApi,payoutRewardApi} from '/@/api/provider/account';
  import { PageWrapper } from '/@/components/Page';
  import BalanceInput from '../../../components/BalanceInput/index.vue';
  import { useI18n } from '/@/hooks/web/useI18n';
  import { defineComponent, getCurrentInstance, onMounted, reactive, toRefs } from 'vue';
  import BigNumber from 'bignumber.js';
  import { useMessage } from '/@/hooks/web/useMessage';
  import AButton from "/@/components/Button/src/BasicButton.vue";
  // eslint-disable-next-line vue/no-export-in-script-setup
  export default defineComponent({
    name: 'center',
    components: {
      AButton,
      BalanceInput,
      PageWrapper,
    },
    setup: function () {
      const { createMessage } = useMessage();
      const { t } = useI18n();
      const { proxy } = getCurrentInstance();
      const state = reactive({
        visible: false,
        title: '',
        stakingLoading: false,
        placeholder: '',
        stakingAmountTip: false,
        rewardLoading: false,
        activeModal: '',
        amount: '',
        address: '',
        pledgeAmount: '0.0000',
        activeAmount: '0.0000',
        lockAmount: '0.0000',
        reward: '0.0000',
      });
      onMounted(() => {
        getAccountInfo();
        getStaking();
        getRewardInfo();
      });
      function getAccountInfo() {
        getAccountInfoApi().then((data) => {
          state.address = data.Address;
          state.amount = new BigNumber(data.Amount)
            .div(new BigNumber(Math.pow(10, 12)))
            .toNumber()
            .toFixed(4);
        });
      }
      function getStaking() {
        getStakingInfoApi().then((data) => {
          state.pledgeAmount = new BigNumber(data.Amount)
            .div(new BigNumber(Math.pow(10, 12)))
            .toNumber()
            .toFixed(4);
          state.activeAmount = new BigNumber(data.ActiveAmount)
            .div(new BigNumber(Math.pow(10, 12)))
            .toNumber()
            .toFixed(4);
          state.lockAmount = new BigNumber(data.LockAmount)
            .div(new BigNumber(Math.pow(10, 12)))
            .toNumber()
            .toFixed(4);
        });
      }
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
      function stakingAmountClick() {
        let price = proxy.$refs.inputRef.getPrice();
        let amount = new BigNumber(state.amount).times(new BigNumber(Math.pow(10, 12))).toNumber();
        if (amount < price) {
          createMessage.warning(t('accountInfo.info.insufficientBalance'));
          return;
        }
        state.stakingLoading = true;
        stakingAmountApi(price)
          .then(() => {
            getStaking();
            getAccountInfo();
            createMessage.success(t('accountInfo.info.pledgeAmountSuccess'));
            state.stakingLoading = false;
            close();
          })
          .catch(() => {
            state.stakingLoading = false;
            createMessage.error(t('accountInfo.info.pledgeAmountFailed'));
          });
      }
      function withdrawAmount() {
        let price = proxy.$refs.inputRef.getPrice();
        let amount = new BigNumber(state.activeAmount)
          .times(new BigNumber(Math.pow(10, 12)))
          .toNumber();
        if (amount < price) {
          createMessage.warning(t('accountInfo.info.withdrawTip'));
          return;
        }
        state.stakingLoading = true;
        withdrawAmountApi(price)
          .then(() => {
            getStaking();
            getAccountInfo();
            createMessage.success(t('accountInfo.info.withdrawAmountSuccess'));
            state.stakingLoading = false;
            close();
          })
          .catch(() => {
            state.stakingLoading = false;
            createMessage.error(t('accountInfo.info.withdrawAmountFailed'));
          });
      }
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
      function getRewardInfo(){
        getRewardInfoApi().then(data => {
          state.reward = new BigNumber(data.TotalIncome)
            .div(new BigNumber(Math.pow(10, 12)))
            .toNumber()
            .toFixed(4);
        })
      }
      function payoutReward(){
        state.rewardLoading = true
        payoutRewardApi().then(() => {
          getRewardInfo()
        }).finally(() => {
          state.rewardLoading = false
        })
      }
      return {
        BalanceInput,
        checkStakingAmount,
        ok,
        close,
        showStakingModal,
        payoutReward,
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
