<template>
  <div
    class="bg-white rounded-lg ring-1 ring-gray-200 ring-inset mx-auto shadow-sm hover:shadow-md"
  >
    <div class="mx-auto w-full px-3 py-4">
      <!-- The Title && Nav -->
      <TheTitleAndNav />

      <!-- The Editor -->
      <div class="rounded-lg p-2 sm:p-3 mb-1">
        <TheMdEditor
          v-model="echoToAdd.content"
          class="rounded-lg"
          v-if="currentMode === Mode.ECH0"
        />

        <!-- ImageMode : TheImageEditor -->
        <div v-if="currentMode === Mode.Image">
          <h2 class="text-gray-500 font-bold my-2">插入图片（支持直链、本地上传）</h2>
          <div class="flex items-center justify-between mb-3">
            <div class="flex items-center gap-2">
              <span class="text-gray-500">选择添加方式：</span>
              <!-- 直链 -->
              <BaseButton
                :icon="Url"
                class="w-7 h-7 sm:w-7 sm:h-7 rounded-md"
                @click="imageToAdd.image_source = ImageSource.URL"
                title="插入图片链接"
              />
              <!-- 上传本地 -->
              <BaseButton
                :icon="Upload"
                class="w-7 h-7 sm:w-7 sm:h-7 rounded-md"
                @click="imageToAdd.image_source = ImageSource.LOCAL"
                title="上传本地图片"
              />
            </div>
            <div>
              <BaseButton
                v-if="imageToAdd.image_url != ''"
                :icon="Addmore"
                class="w-7 h-7 sm:w-7 sm:h-7 rounded-md"
                @click="handleAddMoreImage"
                title="添加更多图片"
              />
            </div>
          </div>
          <div class="my-1">
            <!-- 图片上传本地 -->
            <input
              id="file-input"
              class="hidden"
              type="file"
              accept="image/*"
              ref="fileInput"
              @change="handleUploadImage"
            />
            <BaseButton
              v-if="imageToAdd.image_source === ImageSource.LOCAL"
              @click="handleTriggerUpload"
              class="rounded-md"
              title="上传本地图片"
            >
              <span class="text-gray-400">点击上传</span>
            </BaseButton>
            <!-- 图片直链 -->
            <BaseInput
              v-if="imageToAdd.image_source === ImageSource.URL"
              v-model="imageToAdd.image_url"
              class="rounded-lg h-auto w-full"
              placeholder="请输入图片链接..."
            />
          </div>
        </div>

        <!-- TodoMode : TheTodoModeEditor -->
        <TheTodoModeEditor :current-mode="currentMode" v-model:content="todoToAdd.content" />

        <!-- MusicMode : TheMusicModeEditor -->
        <TheMusicModeEditor @refresh-audio="handleRefreshAudio" :current-mode="currentMode" />

        <!-- The Mode Panel -->
        <TheModePanel
          v-if="currentMode === Mode.Panel"
          v-model:mode="currentMode"
          v-model:extension-type="currentExtensionType"
          v-model:extension-to-add="extensionToAdd"
        />

        <!-- ExtensionMode: TheExtensionEditor -->
        <TheExtensionEditor
          v-model:current-mode="currentMode"
          v-model:current-extension-type="currentExtensionType"
          v-model:extension-to-add="extensionToAdd"
          v-model:video-u-r-l="videoURL"
          v-model:website-to-add="websiteToAdd"
        />
      </div>

      <!-- Editor Buttons -->
      <TheEditorButtons
        :echo-to-add="echoToAdd"
        :current-mode="currentMode"
        @handle-addor-update="handleAddorUpdate"
        @handle-change-mode="handleChangeMode"
        @handle-add-image-mode="handleAddImageMode"
        @handle-exit-update-mode="handleExitUpdateMode"
        @handle-private="handlePrivate"
      />

      <!-- Editor Image -->
      <TheEditorImage
        :imagesToAdd="imagesToAdd"
        :current-mode="currentMode"
        @handleAddorUpdateEcho="handleAddorUpdateEcho"
      />
    </div>
  </div>
</template>

<script setup lang="ts">
import Upload from '@/components/icons/upload.vue'
import Url from '@/components/icons/url.vue'
import Addmore from '@/components/icons/addmore.vue'

import BaseButton from '@/components/common/BaseButton.vue'
import BaseInput from '@/components/common/BaseInput.vue'

import TheMdEditor from '@/components/advanced/TheMdEditor.vue'
import TheModePanel from './TheEditor/TheModePanel.vue'
import TheTitleAndNav from './TheEditor/TheTitleAndNav.vue'
import TheEditorImage from './TheEditor/TheEditorImage.vue'
import TheEditorButtons from './TheEditor/TheEditorButtons.vue'
import TheTodoModeEditor from './TheEditor/TheTodoModeEditor.vue'
import TheMusicModeEditor from './TheEditor/TheMusicModeEditor.vue'
import TheExtensionEditor from './TheEditor/TheExtensionEditor.vue'

import { theToast } from '@/utils/toast'
import { ref, watch } from 'vue'
import { fetchUploadImage, fetchAddEcho, fetchAddTodo, fetchUpdateEcho } from '@/service/api'
import { useEchoStore } from '@/stores/echo'
import { useTodoStore } from '@/stores/todo'
import { storeToRefs } from 'pinia'
import { Mode, ExtensionType, ImageSource } from '@/enums/enums'

const emit = defineEmits(['refreshAudio'])

const echoStore = useEchoStore()
const todoStore = useTodoStore()

const { setTodoMode, getTodos } = todoStore

const { todoMode } = storeToRefs(todoStore)
const { echoToUpdate, isUpdateMode } = storeToRefs(echoStore)

const currentMode = ref<Mode>(Mode.ECH0)
const currentExtensionType = ref<ExtensionType>()

const handleChangeMode = () => {
  if (currentMode.value === Mode.ECH0) {
    currentMode.value = Mode.Panel
  } else if (
    currentMode.value === Mode.TODO ||
    currentMode.value === Mode.PlayMusic ||
    currentMode.value === Mode.EXTEN
  ) {
    currentMode.value = Mode.Panel
    setTodoMode(false)

    if (!echoToAdd.value.image_url || echoToAdd.value.image_url.length === 0) {
      echoToAdd.value.image_source = null
    }
  } else {
    currentMode.value = Mode.ECH0
  }
}

const websiteToAdd = ref<{
  title: string
  site: string
}>({
  title: '',
  site: '',
}) // 临时网站链接变量
const videoURL = ref<string>('') // 临时Bilibili链接变量
const extensionToAdd = ref({
  extension: '',
  extension_type: '',
}) // 临时扩展变量
const imageIndex = ref<number>(0) // 临时图片索引变量
const imageSourceMemory = ref<string>() // 临时图片来源变量
const imageToAdd = ref<App.Api.Ech0.ImageToAdd>({
  image_url: '',
  image_source: '',
}) // 临时图片添加变量
const imagesToAdd = ref<App.Api.Ech0.ImageToAdd[]>([])
const echoToAdd = ref<App.Api.Ech0.EchoToAdd>({
  content: '',
  image_url: null,
  image_source: null,
  images: [],
  private: false,
  extension: null,
  extension_type: null,
})

const todoToAdd = ref<App.Api.Todo.TodoToAdd>({
  content: '',
})

const handleAddMoreImage = () => {
  imagesToAdd.value.push({
    image_url: String(imageToAdd.value.image_url),
    image_source: String(imageToAdd.value.image_source),
  })

  imageToAdd.value.image_url = ''
  imageToAdd.value.image_source = ''
}

const handleAddImageMode = () => {
  currentMode.value = Mode.Image
}
const fileInput = ref<HTMLInputElement | null>(null)
const handleTriggerUpload = () => {
  imageSourceMemory.value = imageToAdd.value.image_source
  if (fileInput.value) {
    fileInput.value.click()
  }
}
const handleUploadImage = async (event: Event) => {
  const target = event.target as HTMLInputElement
  if (target.files && target.files.length > 0) {
    const file = target.files[0]
    fetchUploadImage(file)
      .then((res) => {
        if (res.code === 1) {
          // 改成新版的图片数组
          // echoToAdd.value.image_url = res.data
          imageToAdd.value.image_url = String(res.data)
          imageToAdd.value.image_source = ImageSource.LOCAL
          theToast.success('图片上传成功！')

          handleAddMoreImage()

          // 如果当前处于Echo更新模式，则需要立马执行更新（图片上传操作不可逆，需要立马更新确保后端数据同步）
          if (isUpdateMode.value && echoToUpdate.value) {
            handleAddorUpdateEcho(true)
          }

          // if (currentMode.value === Mode.Image) {
          //   currentMode.value = Mode.ECH0
          // }
        }
      })
      .finally(() => {
        // 重置文件输入
        if (fileInput.value) {
          fileInput.value.value = ''
        }
      })
  }
}

const handleRefreshAudio = () => {
  emit('refreshAudio')
}

const handlePrivate = () => {
  echoToAdd.value.private = !echoToAdd.value.private
}

const handleClear = () => {
  echoToAdd.value.content = ''
  echoToAdd.value.image_url = null
  echoToAdd.value.image_source = null
  echoToAdd.value.images = []
  echoToAdd.value.private = false
  echoToAdd.value.extension = null
  echoToAdd.value.extension_type = null
  extensionToAdd.value.extension = ''
  extensionToAdd.value.extension_type = ''
  videoURL.value = ''
  imagesToAdd.value = []
  imageToAdd.value.image_url = ''
  imageToAdd.value.image_source = ''
  imageIndex.value = 0
}

const handleAddorUpdateEcho = (justSyncImages: boolean) => {
  echoToAdd.value.images = imagesToAdd.value // 将图片数组添加到Echo中

  // 检查是否有外部链接分享
  if (extensionToAdd.value.extension_type === ExtensionType.WEBSITE) {
    // 检查是否存在网站链接
    if (websiteToAdd.value.title.length > 0 && websiteToAdd.value.site.length === 0) {
      theToast.error('网站链接不能为空！')
      return
    }

    // 检查是否存在网站标题
    if (websiteToAdd.value.title.length === 0 && websiteToAdd.value.site.length > 0) {
      websiteToAdd.value.title = '外部链接'
    }

    // 检查网站标题和链接是否都存在
    if (websiteToAdd.value.title.length > 0 && websiteToAdd.value.site.length > 0) {
      // 将网站标题和链接添加到扩展中 (序列化为json)
      extensionToAdd.value.extension = JSON.stringify({
        title: websiteToAdd.value.title,
        site: websiteToAdd.value.site,
      })
    } else {
      extensionToAdd.value.extension = ''
      extensionToAdd.value.extension_type = ''
    }
  }

  // 检查最终的Extension模块是否有内容
  if (extensionToAdd.value.extension.length > 0 && extensionToAdd.value.extension_type.length > 0) {
    echoToAdd.value.extension = extensionToAdd.value.extension
    echoToAdd.value.extension_type = extensionToAdd.value.extension_type
  } else {
    echoToAdd.value.extension = null
    echoToAdd.value.extension_type = null
  }

  // 检查Echo是否为空
  if (
    !echoToAdd.value.content &&
    (!echoToAdd.value.images || echoToAdd.value.images.length === 0) &&
    !echoToAdd.value.extension &&
    !echoToAdd.value.extension_type
  ) {
    if (isUpdateMode.value) {
      theToast.error('待更新的Echo不能为空！')
      return
    } else {
      theToast.error('待添加的Echo不能为空！')
      return
    }
  }

  // if (!echoToAdd.value.image_url || echoToAdd.value.image_url.length === 0) {
  //   echoToAdd.value.image_source = null
  // }

  // 检查是否处于更新模式
  if (isUpdateMode.value) {
    // 处于更新模式，执行更新操作
    if (!echoToUpdate.value) {
      theToast.error('没有待更新的Echo！')
      return
    }

    // 回填 echoToUpdate
    echoToUpdate.value.content = echoToAdd.value.content
    echoToUpdate.value.private = echoToAdd.value.private
    echoToUpdate.value.images = echoToAdd.value.images
    echoToUpdate.value.extension = echoToAdd.value.extension
    echoToUpdate.value.extension_type = echoToAdd.value.extension_type

    // 更新Echo
    fetchUpdateEcho(echoToUpdate.value).then((res) => {
      if (res.code === 1 && !justSyncImages) {
        theToast.success('更新成功！')
        handleClear()
        echoStore.refreshEchos()
        isUpdateMode.value = false
        echoToUpdate.value = null
        currentMode.value = Mode.ECH0
      } else if (res.code === 1 && justSyncImages) {
        theToast.success('发现图片更改，已自动更新同步Echo！')
      } else {
        theToast.error('更新失败，请稍后再试！')
      }
    })
    return
  }

  // 不是Echo更新模式，执行添加操作
  fetchAddEcho(echoToAdd.value).then((res) => {
    if (res.code === 1) {
      theToast.success('发布成功！')
      handleClear()
      echoStore.refreshEchos()
      currentMode.value = Mode.ECH0
    }
  })
}

const handleAddTodo = () => {
  if (todoToAdd.value.content === '') {
    theToast.error('内容不能为空！')
    return
  }

  fetchAddTodo(todoToAdd.value).then((res) => {
    if (res.code === 1) {
      theToast.success('添加成功！')
      todoToAdd.value.content = ''
      getTodos()
    }
  })
}

const handleAddorUpdate = () => {
  if (todoMode.value) {
    handleAddTodo()
  } else {
    handleAddorUpdateEcho(false)
  }
}

const handleExitUpdateMode = () => {
  isUpdateMode.value = false
  echoToUpdate.value = null
  handleClear()
  theToast.info('已退出更新模式！')
}

// 监听用户输入
watch(
  () => videoURL.value,
  (newVal) => {
    if (newVal.length > 0) {
      const bvRegex = /(BV[0-9A-Za-z]{10})/
      const ytRegex =
        /(?:https?:\/\/(?:www\.)?)?(?:youtu\.be\/|youtube\.com\/(?:(?:watch)?\?(?:.*&)?v(?:i)?=|(?:embed)\/))([\w-]+)/
      let match = newVal.match(bvRegex)
      if (match) {
        extensionToAdd.value.extension = match[0] //bilibili
      } else {
        match = newVal.match(ytRegex)
        if (match) {
          extensionToAdd.value.extension = match[1] //youtube
        } else {
          theToast.error('请输入正确的B站/YT分享链接！')
        }
      }
    }
  },
)

// 监听是否处于更新模式
watch(
  () => isUpdateMode.value,
  (newVal) => {
    if (newVal) {
      // 处于更新模式（将待更新的数据填充到编辑器）
      // 0. 清空编辑器
      handleClear()

      // 1. 填充本文
      echoToAdd.value.content = echoToUpdate.value?.content || ''
      echoToAdd.value.private = echoToUpdate.value?.private || false

      // 2. 填充图片
      if (echoToUpdate.value?.images && echoToUpdate.value.images.length > 0) {
        imagesToAdd.value = echoToUpdate.value.images.map((img) => ({
          image_url: img.image_url || '',
          image_source: img.image_source || '',
        }))
      } else {
        imagesToAdd.value = []
      }

      // 3. 填充扩展
      if (echoToUpdate.value?.extension && echoToUpdate.value.extension_type) {
        currentExtensionType.value = echoToUpdate.value.extension_type as ExtensionType
        extensionToAdd.value.extension = echoToUpdate.value.extension
        extensionToAdd.value.extension_type = echoToUpdate.value.extension_type
        // 根据扩展类型填充
        switch (echoToUpdate.value.extension_type) {
          case ExtensionType.MUSIC:
            break

          case ExtensionType.VIDEO:
            videoURL.value = echoToUpdate.value.extension // 直接使用extension填充B站链接
            break

          case ExtensionType.GITHUBPROJ:
            break

          case ExtensionType.WEBSITE:
            // 反序列化网站链接
            const websiteData = JSON.parse(echoToUpdate.value.extension) as {
              title?: string
              site?: string
            }
            websiteToAdd.value.title = websiteData.title || ''
            websiteToAdd.value.site = websiteData.site || ''
            break
        }
      }

      // 4. 回到页面顶部（触发BackToTop）
      window.scrollTo({ top: 0, behavior: 'smooth' })

      // 5. 弹出通知，提示可以编辑了
      theToast.info('已进入更新模式，请编辑内容后点击更新按钮！')
    } else {
      // 退出更新模式
      echoToUpdate.value = null
    }
  },
)
</script>
