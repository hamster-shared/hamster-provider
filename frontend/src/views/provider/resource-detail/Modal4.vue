<template>
  <BasicModal
    title="Modify Price"
    v-bind="$attrs"
    @ok="handlerOk"
    @register="register"
    @visible-change="handleVisibleChange"
  >
    <div class="pt-3px pr-3px">
      <BasicForm :model="model" @register="registerForm" />
    </div>
  </BasicModal>
</template>
<script lang="ts">
  import { defineComponent, ref, nextTick } from 'vue';
  import { BasicModal, useModalInner } from '/@/components/Modal';
  import { BasicForm, FormSchema, useForm } from '/@/components/Form/index';
  import {modifyUintPriceApi} from '/@/api/provider/resource'
  const schemas: FormSchema[] = [
    {
      field: 'price',
      component: 'Input',
      label: 'price per hour',
      colProps: {
        span: 24,
      },
      defaultValue: '100',
    },
  ];
  export default defineComponent({
    components: { BasicModal, BasicForm },
    props: {
      userData: { type: Object },
    },
    setup(props) {
      const modelRef = ref({});
      const [
        registerForm,
        {
          // setFieldsValue,
          // setProps
        },
      ] = useForm({
        labelWidth: 120,
        schemas,
        showActionButtonGroup: false,
        actionColOptions: {
          span: 24,
        },
      });

      const [register,{ setModalProps,closeModal }] = useModalInner((data) => {
        data && onDataReceive(data);
      });

      const handlerOk = function(){
        // send new data to api
        setModalProps({ loading: true, confirmLoading: true });
        modifyUintPriceApi({unitPrice: modelRef.value.price}).then(() => {
          closeModal()
        }).finally(() => {
          setModalProps({ loading: false, confirmLoading: false });
        })
      }

      function onDataReceive(data) {
        console.log('Data Received', data);
        // 方式1;
        // setFieldsValue({
        //   field2: data.data,
        //   field1: data.info,
        // });

        // // 方式2
        modelRef.value = { price: data.price };

        // setProps({
        //   model:{ field2: data.data, field1: data.info }
        // })
      }

      function handleVisibleChange(v) {
        v && props.userData && nextTick(() => onDataReceive(props.userData));
      }

      return { register, schemas, registerForm, model: modelRef, handleVisibleChange,handlerOk };
    },
  });
</script>
