<template>
  <div class="px-11">
    <!-- 列出所有连接（列出每个连接的头像） -->
    <div
      class="rounded-md shadow-sm hover:shadow-md ring-1 ring-gray-200 ring-inset bg-white p-4 mb-4"
    >
      <h2 class="text-gray-600 font-bold text-lg mb-2 flex items-center">
        <Connect class="mr-2" />连接到的Ech0:
      </h2>
      <div v-if="!loading">
        <div v-if="!connectsInfo.length" class="text-gray-500 text-sm mb-2">当前暂无连接</div>
        <div v-else class="flex flex-wrap gap-4">
          <div
            v-for="(connect, index) in connectsInfo"
            :key="index"
            class="relative flex flex-col items-center justify-center w-8 h-8 border-2 border-gray-200 shadow-sm rounded-full hover:shadow-md transition duration-200 ease-in-out group"
          >
            <a :href="connect.server_url" target="_blank">
              <img :src="connect.logo" alt="avatar" class="w-8 h-8 rounded-full object-cover" />
              <!-- 热力圆点 -->
              <span
                class="absolute top-0 right-0 w-2.5 h-2.5 border-2 border-white rounded-full"
                :style="{
                  transform: 'translate(35%, -35%)',
                  backgroundColor: getColor(connect.today_echos || 0),
                }"
              ></span>
            </a>
            <!-- Tooltip -->
            <div
              class="absolute z-10 left-1/2 -translate-x-1/2 top-10 min-w-max bg-gray-800 text-white text-xs rounded px-3 py-2 opacity-0 group-hover:opacity-100 pointer-events-none transition-opacity duration-200 shadow-lg"
            >
              <div class="font-bold mb-1">{{ connect.server_name }}</div>
              <div v-if="connect.sys_username">管理员: {{ connect.sys_username }}</div>
              <div v-if="connect.total_echos">共有: {{ connect.total_echos }}</div>
              <div v-if="connect.today_echos">今日: {{ connect.today_echos }}</div>
            </div>
          </div>
        </div>
      </div>
      <div v-else>
        <div class="text-gray-500 text-sm mb-2">加载中...</div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import Connect from '@/components/icons/connect.vue'
import { useConnectStore } from '@/stores/connect'
import { storeToRefs } from 'pinia'
import { onMounted } from 'vue'

const connectStore = useConnectStore()
const { getConnectInfo } = connectStore
const { loading, connectsInfo } = storeToRefs(connectStore)

const getColor = (count: number): string => {
  if (count >= 4) return '#196127'
  if (count >= 3) return '#239a3b'
  if (count >= 2) return '#7bc96f'
  if (count >= 1) return '#c6e48b'
  return '#b7bbb7' // 默认颜色
}

onMounted(() => {
  getConnectInfo()
})
</script>

<style scoped></style>
