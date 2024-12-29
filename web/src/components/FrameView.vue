<template>
  <a-layout>
    <a-layout-header class="header">
      <div class="logo" />
      <!-- <a-menu v-model:selectedKeys="selectedKeys1" theme="dark" mode="horizontal" :style="{ lineHeight: '64px' }">
        <a-menu-item key="1">nav 1</a-menu-item>
        <a-menu-item key="2">nav 2</a-menu-item>
        <a-menu-item key="3">nav 3</a-menu-item>
      </a-menu> -->
    </a-layout-header>
    <a-layout>
      <a-layout-sider width="200" style="background: #fff">
        <a-menu v-model:selectedKeys="selectedKeys2" v-model:openKeys="openKeys" mode="inline"
          :style="{ height: '100%', borderRight: 0 }">
          <a-sub-menu key="sub1">
            <template #title>
              <span>
                后台管理
              </span>
            </template>
            <a-menu-item v-for="item in menuItems" :key="item.key">
              <router-link :to="item.key">{{ item.title }}</router-link>
            </a-menu-item>
          </a-sub-menu>
        </a-menu>
      </a-layout-sider>
      <a-layout style="padding: 0 24px 24px">
        <a-breadcrumb style="margin: 16px 0">
          <a-breadcrumb-item>后台管理</a-breadcrumb-item>
          <a-breadcrumb-item>{{ breadcrumb }}</a-breadcrumb-item>
        </a-breadcrumb>
        <slot></slot>
      </a-layout>
    </a-layout>
  </a-layout>
</template>
<script setup>
import { useRouter } from 'vue-router'
import { ref, onMounted } from 'vue';
const selectedKeys1 = ref(['1']);

const openKeys = ref(['sub1']);
const selectedKeys2 = ref(['main']);
const breadcrumb = ref('')
const menuItems = ref([
  { title: '转发配置', key: 'main' },
  { title: '用户管理', key: 'user' },
  { title: 'WAF(ufw)', key: 'waf' }
])
const router = useRouter()
onMounted(() => {
  const name = router.currentRoute.value.name;
  selectedKeys2.value = [name]
  breadcrumb.value = menuItems.value.find(x=>x.key == name).title
})

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