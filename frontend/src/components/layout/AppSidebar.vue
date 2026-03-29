<template>
  <aside class="app-sidebar">
    <router-link
      v-for="item in navItems"
      :key="item.path"
      :to="item.path"
      class="sidebar-item"
      :class="{ active: route.path === item.path }"
      :title="item.label"
    >
      <span class="sidebar-icon">{{ item.icon }}</span>
      <span class="sidebar-label">{{ item.label }}</span>
    </router-link>
  </aside>
</template>

<script setup lang="ts">
import { useRoute } from 'vue-router'

const route = useRoute()

const navItems = [
  { path: '/',     icon: '▦', label: 'Dashboard' },
  { path: '/lots', icon: '⊞', label: 'Lot / WIP'  },
  { path: '/spc',  icon: '∿', label: 'SPC'         },
]
</script>

<style scoped>
.app-sidebar {
  width: 64px;
  background: var(--mes-surface-deep);
  border-right: 1px solid var(--mes-border);
  display: flex;
  flex-direction: column;
  padding: 8px 0;
  gap: 4px;
  flex-shrink: 0;
}

.sidebar-item {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 4px;
  padding: 10px 0;
  color: var(--el-text-color-regular);
  text-decoration: none;
  transition: all 0.2s;
  cursor: pointer;
}

.sidebar-item:hover,
.sidebar-item.active {
  color: var(--el-color-primary);
  background: rgba(59, 130, 246, 0.1);
}

.sidebar-icon {
  font-size: 20px;
  line-height: 1;
}

.sidebar-label {
  font-size: 10px;
  text-align: center;
  line-height: 1.2;
}
</style>
