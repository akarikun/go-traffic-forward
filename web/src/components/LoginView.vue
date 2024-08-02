<template>
  <a-card title="后台登录">
    <a-form :model="formState" name="basic" :label-col="{ span: 8 }" :wrapper-col="{ span: 16 }" autocomplete="off"
      @finish="onFinish" @finishFailed="onFinishFailed">
      <a-form-item label="用户名" name="username" :rules="[{ required: true, message: '请输入用户名' }]">
        <a-input v-model:value="formState.username" />
      </a-form-item>

      <a-form-item label="密码" name="password" :rules="[{ required: true, message: '请输入密码' }]">
        <a-input-password v-model:value="formState.password" />
      </a-form-item>

      <a-form-item :wrapper-col="{ offset: 8, span: 16 }">
        <a-button type="primary" html-type="submit">登录</a-button>
      </a-form-item>
    </a-form>
  </a-card>

</template>
<script setup>
import { reactive } from 'vue';

import { useRouter } from 'vue-router'

import * as $ from './../utils/common.js'

const router = useRouter()
const formState = reactive({
  username: '',
  password: '',
});
const onFinish = async values => {
  const { status, data, message } = await $.POST($.URL.Login, { username: values.username, password: values.password })
  if (status == 0) {
    return;
  }
  router.push({ name: 'main' });
};
const onFinishFailed = errorInfo => {
  console.log('Failed:', errorInfo);
};
</script>
<style scoped>
.ant-card {
  width: 400px;
  height: 250px;
  position: absolute;
  top: 50%;
  left: 50%;
  margin-left: -200px;
  margin-top: -125px
}
</style>