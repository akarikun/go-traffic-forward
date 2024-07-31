<template>
  <FrameView>
    <a-layout>
      <a-layout style="padding: 0 24px 24px">
        <a-modal v-model:open="open" :title="`${formState.id == 0 ? '添加' : '编辑'}`" @ok="handleOk" okText="确认"
          cancelText="取消">
          <a-form :model="formState" name="basic" :label-col="{ span: 5 }" :wrapper-col="{ span: 16 }"
            autocomplete="off">
            <a-form-item label="端口" name="port" :rules="[{ required: true, message: '端口' }]">
              <a-input v-model:value="formState.port" />
            </a-form-item>
            <a-form-item label="转发" name="destination" :rules="[{ required: true, message: '转发' }]">
              <a-input v-model:value="formState.destination" />
            </a-form-item>
          </a-form>
        </a-modal>

        <a-layout-content :style="{ background: '#fff', padding: '24px', margin: 0, minHeight: '280px' }">
          <a-button class="editable-add-btn" style="margin-bottom: 8px" @click="showModal">Add</a-button>
          <a-table bordered :data-source="dataSource" :columns="columns">
            <template #bodyCell="{ column, text, record }">
              {{ text }}
              <template v-if="column.dataIndex === 'operation'">
                <a-popconfirm v-if="dataSource.length" title="Sure to delete?" @confirm="onDelete(record.key)">
                  <a>Edit</a> |
                </a-popconfirm>
                <a-popconfirm v-if="dataSource.length" title="Sure to delete?" @confirm="onDelete(record.key)">
                  <a>Delete</a>
                </a-popconfirm>
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
import { reactive, ref, toRaw } from 'vue';
import { cloneDeep } from 'lodash-es';

const formState = reactive({
  id: 0,
  port: "",
  destination: ""
});

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
    dataIndex: 'port',
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
const dataSource = ref([
  { id: 1, user:'admin', port: "8088", port_status:"1", destination: "localhost:59992", use_total: "100MB", add_date: '2022-07-13' }
]);
//const count = computed(() => dataSource.value.length + 1);
const editableData = reactive({});
const edit = key => {
  editableData[key] = cloneDeep(dataSource.value.filter(item => key === item.key)[0]);
};
const save = key => {
  Object.assign(dataSource.value.filter(item => key === item.key)[0], editableData[key]);
  delete editableData[key];
};
const onDelete = key => {
  dataSource.value = dataSource.value.filter(item => item.key !== key);
};
// const handleAdd = () => {
//   const newData = {
//     key: `${count.value}`,
//     name: `Edward King ${count.value}`,
//     age: 32,
//     address: `London, Park Lane no. ${count.value}`,
//   };
//   dataSource.value.push(newData);
// };
const open = ref(false);
const showModal = () => {
  open.value = true;
};
const handleOk = e => {
  const data = toRaw(formState)
  console.log(data)
  // console.log(e);
  // open.value = false;
};

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