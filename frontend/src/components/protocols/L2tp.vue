<template>
  <v-card>
    <v-alert
      type="info"
      variant="tonal"
      class="my-3"
    >
      仅开放 UDP 500 和 4500，无需开放裸 UDP 1701。用户页管理账号；DNS、分流和流量统计由 sing-box 处理。暂支持 IPv4 TCP/UDP。
    </v-alert>
    <v-text-field
      v-model="data.listen"
      label="监听 IPv4"
    />
    <v-text-field
      v-model.number="data.listen_port"
      label="IKE 端口"
      type="number"
    />
    <v-text-field
      v-model="data.public_ip"
      label="服务器公网 IPv4"
    />
    <v-text-field
      v-model="data.psk"
      label="IPsec 预共享密钥（PSK）"
    />
    <v-text-field
      v-model="data.pool"
      label="客户端地址池（IPv4 CIDR）"
    />
  </v-card>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import type { Inbound, L2TP } from '@/types/inbounds'
const model = defineModel<Inbound>('data', { required: true })
const data = computed(() => model.value as L2TP)
</script>
