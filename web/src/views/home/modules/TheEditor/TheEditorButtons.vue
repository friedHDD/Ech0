<template>
  <div class="flex flex-row items-center justify-between px-2">
    <div class="flex flex-row items-center gap-2">
      <!-- ShowMore -->
      <div>
        <BaseButton
          :icon="currentMode === Mode.ECH0 ? Advance : Back"
          @click="handleChangeMode"
          :class="['w-8 h-8 sm:w-9 sm:h-9 rounded-md'].join(' ')"
          title="其它"
        />
      </div>
      <!-- Photo Upload -->
      <div v-if="currentMode === Mode.ECH0">
        <BaseButton
          :icon="ImageUpload"
          @click="handleAddImageMode"
          class="w-8 h-8 sm:w-9 sm:h-9 rounded-md"
          title="添加图片"
        />
      </div>
      <!-- Privacy Set -->
      <div v-if="currentMode === Mode.ECH0">
        <BaseButton
          :icon="echoToAdd.private ? Private : Public"
          @click="handlePrivate"
          class="w-8 h-8 sm:w-9 sm:h-9 rounded-md"
          title="是否私密"
        />
      </div>
    </div>

    <div class="flex flex-row items-center gap-2">
      <!-- Clear -->
      <!-- <div>
            <BaseButton
              :icon="Clear"
              @click="handleClear"
              class="w-8 h-8 rounded-md"
              title="清空输入和图片"
            />
          </div> -->
      <!-- Publish -->
      <div v-if="currentMode !== Mode.Panel && isUpdateMode === false">
        <BaseButton
          :icon="Publish"
          @click="handleAddorUpdate"
          class="w-8 h-8 sm:w-9 sm:h-9 rounded-md"
          title="发布Echo"
        />
      </div>
      <!-- Exit Update -->
      <div
        v-if="
          currentMode !== Mode.Panel &&
          currentMode !== Mode.TODO &&
          currentMode !== Mode.PlayMusic &&
          isUpdateMode === true
        "
      >
        <BaseButton
          :icon="ExitUpdate"
          @click="handleExitUpdateMode"
          class="w-8 h-8 sm:w-9 sm:h-9 rounded-md"
          title="退出更新模式"
        />
      </div>
      <!-- Update -->
      <div
        v-if="
          currentMode !== Mode.Panel &&
          currentMode !== Mode.TODO &&
          currentMode !== Mode.PlayMusic &&
          isUpdateMode === true
        "
      >
        <BaseButton
          :icon="Update"
          @click="handleAddorUpdate"
          class="w-8 h-8 sm:w-9 sm:h-9 rounded-md"
          title="更新Echo"
        />
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import Advance from '@/components/icons/advance.vue'
import ImageUpload from '@/components/icons/image.vue'
import Public from '@/components/icons/public.vue'
import Private from '@/components/icons/private.vue'
import Publish from '@/components/icons/publish.vue'
import Update from '@/components/icons/update.vue'
import ExitUpdate from '@/components/icons/exitupdate.vue'
import Back from '@/components/icons/back.vue'
import BaseButton from '@/components/common/BaseButton.vue'
import { Mode } from '@/enums/enums'
import { storeToRefs } from 'pinia'
import { useEchoStore } from '@/stores/echo'

const { currentMode, echoToAdd } = defineProps<{
  echoToAdd: App.Api.Ech0.EchoToAdd
  currentMode: Mode
}>()

const emit = defineEmits([
  'handleAddorUpdate',
  'handleChangeMode',
  'handleAddImageMode',
  'handlePrivate',
  'handleExitUpdateMode',
])

const echoStore = useEchoStore()
const { isUpdateMode } = storeToRefs(echoStore)

const handleAddorUpdate = () => {
  emit('handleAddorUpdate')
}

const handleChangeMode = () => {
  emit('handleChangeMode')
}

const handleAddImageMode = () => {
  emit('handleAddImageMode')
}

const handlePrivate = () => {
  emit('handlePrivate')
}

const handleExitUpdateMode = () => {
  emit('handleExitUpdateMode')
}
</script>

<style scoped></style>
