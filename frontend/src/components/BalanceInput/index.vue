<template>
  <div class="input-container">
    <a-input v-model:value="value" :placeholder="placeholder" @change="changeClick">
      <template #addonAfter>
        <a-select style="width: 90px" v-model:value="uintPower">
          <a-select-option
            v-for="(item, index) in uintOptions"
            v-model:value="item.power"
            :key="index"
          >
            {{ item.text }}
          </a-select-option>
        </a-select>
      </template>
    </a-input>
  </div>
</template>

<script lang="ts">
  import { formatBalance } from '@polkadot/util';
  import BigNumber from 'bignumber.js';
  import { defineComponent, onMounted, reactive, toRefs } from 'vue';
  export default defineComponent({
    name: 'BalanceInput',
    props: {
      changeClick: {
        type: Function,
        default: () => {},
      },
      placeholder: {
        type: String,
        default: '',
      },
    },
    setup: function () {
      const state = reactive({
        value: '',
        uintPower: 0,
        uintOptions: [],
      });
      onMounted(() => {
        getUintOptions();
      });

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
      return {
        ...toRefs(state),
        getPrice,
      };
    },
  });
</script>

<style lang="less" scoped>
  .input-container {
    width: 100%;
  }
</style>
