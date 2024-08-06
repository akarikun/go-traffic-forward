<template>
  <FrameView>
    <a-layout>
      <a-layout style="padding: 0 24px 24px">
        <a-modal v-model:open="open" :title="`${formState.id == 0 ? '添加' : '编辑'}`" @ok="handleOk" okText="确认"
          cancelText="取消">
          <a-form :model="formState" :label-col="{ span: 5 }" :wrapper-col="{ span: 16 }" autocomplete="off">
            <a-form-item label="端口" name="bind_port" :rules="[{ required: true, message: '端口' }]">
              <a-input v-model:value="formState.bind_port" />
            </a-form-item>
            <a-form-item label="转发" name="destination" :rules="[{ required: true, message: '转发' }]">
              <a-input v-model:value="formState.destination" />
            </a-form-item>
          </a-form>
        </a-modal>

        <a-layout-content :style="{ background: '#fff', padding: '24px', margin: 0, minHeight: '280px' }">
          <a-button class="editable-add-btn" style="margin-bottom: 8px" @click="showModal">Add</a-button>
          <a-table bordered :data-source="dataSource" :columns="columns" :loading="loading">
            <template #bodyCell="{ column, text, record }">
              <template v-if="column.dataIndex === 'add_date'">
                {{ fmtDate(record.add_date) }}
              </template>
              <template v-else-if="column.dataIndex === 'operation'">
                <a @click="handleEdit(record.id)">编辑</a> |
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
import FrameView from './components/FrameView.vue';
import * as $ from './utils/common'
import { reactive, ref, toRaw } from 'vue';
import { cloneDeep } from 'lodash-es';
import { useRouter } from 'vue-router'
const router = useRouter()
console.log(router)

const formState = ref({
  id: 0,
  bind_port: "",
  destination: ""
});
const loading = ref(true)
const columns = [
  {
    title: '编号',
    dataIndex: 'id'
  },
  {
    title: '用户',
    dataIndex: 'user'
  },
  {
    title: '端口',
    dataIndex: 'bind_port',
  },
  {
    title: '转发',
    dataIndex: 'destination',
  },
  {
    title: '使用量',
    dataIndex: 'use_total',
  },
  {
    title: '时间',
    dataIndex: 'add_date',
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
  const res = await $.GET($.URL.Forward)
  dataSource.value = res.data
  loading.value = false
}

load_data();
const onDelete = async key => {
  await $.POST($.URL.Forward_DEL, { id: key })
  await load_data()
};
const open = ref(false);
const showModal = () => {
  open.value = true;
};
const handleOk = async e => {
  formState.value.id = 0
  const { status, message, data } = await $.POST($.URL.Forward, formState.value)
  if (status == 0) {
    return
  }
  formState.value = {}
  open.value = false
  await load_data()
};
const handleEdit = (id) => {
  formState.value.id = id
  open.value = true
  formState.value = dataSource.value.filter(x => x.id == id)[0]
}

</script>
<style scoped>
#components-layout-demo-top-side-2 .logo {
  float: left;
  width: 120px;
  height: 31px;
  margin: 16px 24px 16px 0;
  background: rgba(255, 255, 255, 0.3);
}

.ant-row-rtl #components-layout-demo-top-side-2 .logo {
  float: right;
  margin: 16px 0 16px 24px;
}

.site-layout-background {
  background: #fff;
}

.editable-cell {
  position: relative;

  .editable-cell-input-wrapper,
  .editable-cell-text-wrapper {
    padding-right: 24px;
  }

  .editable-cell-text-wrapper {
    padding: 5px 24px 5px 5px;
  }

  .editable-cell-icon,
  .editable-cell-icon-check {
    position: absolute;
    right: 0;
    width: 20px;
    cursor: pointer;
  }

  .editable-cell-icon {
    margin-top: 4px;
    display: none;
  }

  .editable-cell-icon-check {
    line-height: 28px;
  }

  .editable-cell-icon:hover,
  .editable-cell-icon-check:hover {
    color: #108ee9;
  }

  .editable-add-btn {
    margin-bottom: 8px;
  }
}

.editable-cell:hover .editable-cell-icon {
  display: inline-block;
}
</style>