<template>
  <FrameView>
    <a-layout>
      <a-layout style="padding: 0 24px 24px">
        <a-modal v-model:open="open" title="添加" @ok="handleOk" okText="确认" cancelText="取消">
          <a-form :model="formState" :label-col="{ span: 5 }" :wrapper-col="{ span: 16 }" autocomplete="off">
            <a-form-item label="编号(可选)" name="id" :rules="[{ message: '编号' }]">
              <a-input v-model:value="formState.id" />
            </a-form-item>
            <a-form-item label="port(可选)" name="port" :rules="[{ message: 'port' }]">
              <a-input v-model:value="formState.port" />
            </a-form-item>
            <a-form-item label="from(可选)" name="from" :rules="[{ message: 'from' }]">
              <a-input v-model:value="formState.from" />
            </a-form-item>

            <a-form-item label="action" name="action">
              <a-select ref="select" v-model:value="action_value" style="width: 120px" @change="handleChange">
                <a-select-option value="deny">deny</a-select-option>
                <a-select-option value="allow">allow</a-select-option>
              </a-select>
            </a-form-item>

            <a-form-item label="预览">
              <label>{{ preview }}</label>
            </a-form-item>

            <!-- <a-form-item label="转发" name="action" :rules="[{ required: true, message: '转发' }]">
              <a-input v-model:value="formState.action" />
            </a-form-item> -->
          </a-form>
        </a-modal>

        <a-layout-content :style="{ background: '#fff', padding: '24px', margin: 0, minHeight: '280px' }">
          <a-button class="editable-add-btn" style="margin-bottom: 8px" @click="showModal">添加</a-button>
          <a-table bordered :data-source="dataSource" :columns="columns" :loading="loading">
            <template #bodyCell="{ column, text, record }">
              <template v-if="column.dataIndex === 'add_date'">
                {{ fmtDate(record.add_date) }}
              </template>
              <template v-else-if="column.dataIndex === 'operation'">
                <a-popconfirm v-if="dataSource.length" title="是否删除,该操作不可恢复" @confirm="onDelete(record.id)">
                  <a>删除</a>
                </a-popconfirm>
              </template>
              <template v-else-if="column.dataIndex === 'use_total'">
                {{ $.formatBytes(text) }}
              </template>
              <template v-else>
                {{ text }}
              </template>
            </template>
          </a-table>
        </a-layout-content>
      </a-layout>
    </a-layout>
  </FrameView>
</template>
<script setup>
import FrameView from './../components/FrameView.vue';
import * as $ from '../utils/common'
import { reactive, ref, toRaw, watch } from 'vue';
import { cloneDeep } from 'lodash-es';
import { useRouter } from 'vue-router'
const router = useRouter()

const action_value = ref('deny')
const preview = ref('')
const formState = ref({
  id: "",
  port: "",
  from: "",
  action: ""
});

watch(formState, (val) => {
  update_prev();
}, { deep: true })

const update_prev = () => {
  const { id, port, from, action } = formState.value;
  let str = ['ufw']
  if (id) {
    str.push(`insert ${id}`);
  }
  str.push(`${action_value.value}`)
  if (from) {
    str.push(`from ${from}`)
    if (port) {
      str.push(`to any port ${port}`)
    }
  } else if (port) {
    str.push(`${port}`)
  }
  // console.log(str);
  preview.value = str.join(' ')
}

const loading = ref(true)
const columns = [
  {
    title: '编号',
    dataIndex: 'id'
  },
  {
    title: 'To',
    dataIndex: 'to'
  },
  {
    title: 'Action',
    dataIndex: 'action',
  },
  {
    title: 'From',
    dataIndex: 'from',
  },
  {
    title: '操作',
    dataIndex: 'operation',
  },
];
const dataSource = ref([]);
const fmtDate = (add_date) => {
  return new Date(add_date).toLocaleString()
}

const load_data = async () => {
  loading.value = true
  var res = await $.GET($.URL.WAF)
  // console.log(res);
  if (res.status == 1) {
    var { status, msg, data } = await $.GET($.URL.WAF_STATUS)
    // console.log(data);
    // console.log({ status, msg, data });
    dataSource.value = format_data(data);
  }
  loading.value = false
}

const format_data = (str) => {
  const regex = /\[\s*\d+\].*/g;
  const matches = str.match(regex);

  if (matches) {
    const result = matches.map((line) =>
      line.split(/\s{4,}/)
    ).map(x => {
      const to = x[0].split('] ');
      return { id: to[0].replace('[', ''), to: to[1], action: x[1], from: x[2] }
    })
    // console.log(result);
    return result;
  }
  return [];
}

const handleChange = key => {
  update_prev();
}


load_data();
const onDelete = async key => {
  await $.POST($.URL.WAF_DELETE, { id: key })
  await load_data()
};
const open = ref(false);
const showModal = () => {
  formState.value = {}
  open.value = true;
};
const handleOk = async e => {
  const res = await $.POST($.URL.WAF_UPDATE, { cmd: preview.value.replace('ufw ', '') })
  open.value = false;
  formState.value = {}
  await load_data()
};

</script>