<template>
  <v-dialog

    width="800"
  >
    <v-card class="rounded-lg">
      <v-card-title>
        {{ title }}
      </v-card-title>
      <v-divider />
      <v-card-text>
        <div class="code-editor">
          <div
            ref="lineNumbers"
            class="line-numbers"
          >
            <span
              v-for="n in lineCount"
              :key="n"
            >{{ n }}</span>
          </div>
          <v-textarea
            ref="textareaRef"
            v-model="content"
            hide-details
            variant="outlined"
            bg-color="background"
            wrap="off"
            :spellcheck="false"
            no-resize
            auto-grow
            @scroll.capture="syncScroll"
          />
        </div>
      </v-card-text>
      <v-card-actions>
        <v-spacer />
        <v-btn
          color="primary"
          variant="outlined"
          @click="closeModal"
        >
          {{ $t('actions.close') }}
        </v-btn>
        <v-btn
          color="primary"
          variant="tonal"
          @click="saveChanges"
        >
          {{ $t('actions.save') }}
        </v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>
</template>

<script lang="ts">
import { useTheme } from 'vuetify'

export default {
  props: {
    visible: { type: Boolean, required: true },
    data: { type: String, required: true },
    title: { type: String, required: true }
  },
  emits: ['close', 'save'],
  data() {
    return {
      content: this.$props.data,
      theme: useTheme()
    }
  },
  computed: {
    lineCount() {
      return this.content?.split('\n').length
    }
  },
  watch: {
    visible(v) {
      if (v) {
        this.content = this.$props.data
      }
    }
  },
  methods: {
    syncScroll() {
      const field = this.$refs.textareaRef as { $el: HTMLElement } | undefined
      const textarea = field?.$el.querySelector('textarea')
      const lineNumbers = this.$refs.lineNumbers as HTMLElement | undefined
      if (lineNumbers && textarea) {
        lineNumbers.scrollTop = textarea.scrollTop
      }
    },
    closeModal() {
      this.$emit('close')
    },
    saveChanges() {
      this.$emit('save', this.content)
    }
  }
}
</script>

<style scoped>
.code-editor {
  direction: ltr;
  display: flex;
  border: 1px solid v-bind('theme.current.colors["outline"]');
  border-radius: 4px;
  overflow: hidden;
  font-size: 14px; /* Consistent font size */
}

.line-numbers {
  flex: 0 0 48px;
  background: v-bind('theme.current.colors["surface"]');
  text-align: right;
  padding: 12px 8px 12px 4px; /* Match textarea padding */
  line-height: 1.5; /* Match textarea line height */
  overflow-y: hidden; /* Prevent independent scrolling */
  user-select: none;
  display: flex;
  flex-direction: column;
}

.line-numbers span {
  display: block;
  line-height: 1.5; /* Match textarea line height */
  height: 1.5em; /* Ensure consistent height per line */
  font-family: monospace; /* Match textarea font */
}

/* Override Vuetify textarea styles for alignment */
:deep(.v-textarea .v-field__input) {
  padding: 12px 8px;
  line-height: 1.5;
  font-family: monospace;
  white-space: pre;
  mask-image: inherit;
  font-size: 14px;
}

:deep(.v-textarea) {
  min-width: 0;
}
</style>
