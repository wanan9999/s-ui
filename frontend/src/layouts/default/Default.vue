<template>
  <v-app>
    <drawer
      v-model:display-drawer="displayDrawer"
      :is-mobile="isMobile"
      :rail="rail"
    />
    <default-bar
      :is-mobile="isMobile"
      @toggle-drawer="toggleNavigation"
    />
    <default-view />
  </v-app>
</template>

<script lang="ts" setup>
import { computed, ref, watch } from 'vue'
import DefaultBar from './AppBar.vue'
import Drawer from './Drawer.vue'
import DefaultView from './View.vue'
import { useDisplay } from 'vuetify'

const { smAndDown } = useDisplay()
const displayDrawer = ref(false)
const rail = ref(false)

const toggleDrawer = () => {
  displayDrawer.value = !displayDrawer.value
}

const isMobile = computed((): boolean => smAndDown.value)
const toggleNavigation = () => {
  if (isMobile.value) toggleDrawer()
  else rail.value = !rail.value
}

// keep the drawer open on desktop and closed on mobile
watch(smAndDown, (v) => { displayDrawer.value = !v }, { immediate: true })
</script>
