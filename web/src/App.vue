<template>
  <FrameView>
    <a-layout>
      <a-layout style="padding: 0 24px 24px">
        <a-modal v-model:open="open" :title="`${formState.id==0?'添加':'编辑'}`" @ok="handleOk">
          <a-form :model="formState" :label-col="labelCol" :wrapper-col="wrapperCol">
            <a-form-item label="port">
              <a-input v-model:value="formState.port" />
            </a-form-item>
            <a-form-item label="destination">
              <a-input v-model:value="formState.destination" />
            </a-form-item>
          </a-form>
        </a-modal>
        
        <a-layout-content :style="{ background: '#fff', padding: '24px', margin: 0, minHeight: '280px' }">
          <a-button class="editable-add-btn" style="margin-bottom: 8px" @click="showModal">Add</a-button>
          <a-table bordered :data-source="dataSource" :columns="columns">
            <template #bodyCell="{ column, text, record }">
              <template v-if="column.dataIndex === 'name'"> {{ text || ' ' }}
                <!-- <div class="editable-cell">
                  <div v-if="editableData[record.key]" class="editable-cell-input-wrapper">
                    <a-input v-model:value="editableData[record.key].name" @pressEnter="save(record.key)" />
                    <check-outlined class="editable-cell-icon-check" @click="save(record.key)" />
                  </div>
                  <div v-else class="editable-cell-text-wrapper">
                    {{ text || ' ' }}
                    <edit-outlined class="editable-cell-icon" @click="edit(record.key)" />
                  </div>
                </div> -->
              </template>
              <template v-else-if="column.dataIndex === 'operation'">
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
  id:0,
  port:"",
  destination:""
});

const columns = [
  {
    title: 'name',
    dataIndex: 'name',
    width: '30%',
  },
  {
    title: 'age',
    dataIndex: 'age',
  },
  {
    title: 'address',
    dataIndex: 'address',
  },
  {
    title: 'operation',
    dataIndex: 'operation',
  },
];
const dataSource = ref([
  {
    key: '0',
    name: 'Edward King 0',
    age: 32,
    address: 'London, Park Lane no. 0',
  },
  {
    key: '1',
    name: 'Edward King 1',
    age: 32,
    address: 'London, Park Lane no. 1',
  },
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
  console.log(toRaw(formState))
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