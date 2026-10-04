<template>
  <v-navigation-drawer
    v-model="showDrawer"
    :temporary="isMobile"
    :rail="!isMobile && rail"
    :permanent="!isMobile"
    :width="232"
    class="app-navigation"
  >
    <v-list-item
      height="72"
      title="S-UI"
      class="app-brand"
    >
      <template #prepend>
        <v-avatar
          color="primary"
          variant="tonal"
          rounded="lg"
          size="36"
        >
          <v-icon
            icon="mdi-lan"
            size="22"
          />
        </v-avatar>
      </template>
      <template
        v-if="isMobile"
        #append
      >
        <v-btn
          icon="mdi-close"
          variant="text"
          size="small"
          :aria-label="$t('actions.close')"
          @click="$emit('update:displayDrawer', false)"
        />
      </template>
    </v-list-item>

    <v-divider />

    <v-list
      density="compact"
      nav
      class="navigation-list"
    >
      <v-list-item
        v-for="item in menu"
        :key="item.title"
        link
        :to="item.path"
        :active="router.currentRoute.value.path == item.path"
        color="primary"
        rounded="lg"
        @click="isMobile ? $emit('update:displayDrawer', false) : null"
      >
        <template #prepend>
          <v-icon
            :icon="item.icon"
            size="20"
          />
        </template>
        <v-list-item-title>{{ $t(item.title) }}</v-list-item-title>
      </v-list-item>
    </v-list>
    <template #append>
      <v-list-item
        prepend-icon="mdi-logout"
        :title="$t('menu.logout')"
        @click="Logout"
      />
    </template>
  </v-navigation-drawer>
</template>

<script lang="ts" setup>
import { computed } from 'vue'
import router from '@/router'
import { logout } from '@/plugins/httputil'

const props = defineProps<{ isMobile: boolean, displayDrawer: boolean, rail: boolean }>()
const emit = defineEmits<{ 'update:displayDrawer': [visible: boolean] }>()

const showDrawer = computed({
  get: () => props.displayDrawer,
  set: (visible: boolean) => {
    emit('update:displayDrawer', visible)
  },
})

const menu = [
  { title: 'pages.home', icon: 'mdi-home',  path: '/' },
  { title: 'pages.inbounds', icon: 'mdi-cloud-download',  path: '/inbounds' },
  { title: 'pages.clients', icon: 'mdi-account-multiple',  path: '/clients' },
  { title: 'pages.outbounds', icon: 'mdi-cloud-upload',  path: '/outbounds' },
  { title: 'pages.endpoints', icon: 'mdi-cloud-tags',  path: '/endpoints' },
  { title: 'pages.services', icon: 'mdi-server',  path: '/services' },
  { title: 'pages.tls', icon: 'mdi-certificate',  path: '/tls' },
  { title: 'pages.basics', icon: 'mdi-application-cog',  path: '/basics' },
  { title: 'pages.rules', icon: 'mdi-routes',  path: '/rules' },
  { title: 'pages.dns', icon: 'mdi-dns',  path: '/dns' },
  { title: 'pages.admins', icon: 'mdi-account-tie',  path: '/admins' },
  { title: 'pages.settings', icon: 'mdi-cog',  path: '/settings' },
]

const Logout = async () => {
  logout()
}
</script>
