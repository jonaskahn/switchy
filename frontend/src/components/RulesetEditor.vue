<script setup lang="ts">
import { ref, computed } from 'vue'
import type { Rule, Ruleset } from '@/types'

const FORCE_PICKER_BROWSER = '_Switchy'

const props = defineProps<{ modelValue: Ruleset; browsers: string[] }>()
const emit = defineEmits<{ 'update:modelValue': [Ruleset]; remove: []; change: [] }>()

const expanded = ref(false)

const ruleCountTag = computed(() => {
  const n = props.modelValue.Rules.length
  return `${n} RULE${n === 1 ? '' : 'S'}`
})

const ruleSummary = computed(() => {
  const n = props.modelValue.Rules.length
  if (n === 0) return 'No rules'
  const label = `${n} rule${n === 1 ? '' : 's'}`
  return props.modelValue.Enabled ? `${label} · active` : `${label} · disabled`
})

function update(partial: Partial<Ruleset>) {
  emit('update:modelValue', { ...props.modelValue, ...partial })
  emit('change')
}

function toggleEnabled(e: Event) {
  e.stopPropagation()
  update({ Enabled: !props.modelValue.Enabled })
}

function addRule(e?: Event) {
  e?.stopPropagation()
  update({ Rules: [...props.modelValue.Rules, { Pattern: '', Browser: FORCE_PICKER_BROWSER }] })
}

function removeRule(idx: number) {
  const rules = [...props.modelValue.Rules]
  rules.splice(idx, 1)
  update({ Rules: rules })
}

function updateRule(idx: number, partial: Partial<Rule>) {
  const rules = [...props.modelValue.Rules]
  rules[idx] = { ...rules[idx], ...partial }
  update({ Rules: rules })
}

function setName(e: Event) {
  const v = (e.target as HTMLInputElement).value
  update({ Name: v })
}
</script>

<template>
  <div class="row" :class="{ 'is-open': expanded }">
    <div class="row-head" @click="expanded = !expanded">
      <svg
        class="chev"
        viewBox="0 0 24 24"
        fill="none"
        stroke="currentColor"
        stroke-width="2.2"
        stroke-linecap="round"
        stroke-linejoin="round"
        aria-hidden="true"
      >
        <path d="M9 6l6 6-6 6" />
      </svg>

      <button
        type="button"
        class="switch no-drag"
        :class="{ on: modelValue.Enabled }"
        :aria-pressed="modelValue.Enabled"
        aria-label="Toggle ruleset enabled"
        @click="toggleEnabled"
      ></button>

      <div class="row-meta" @click.stop>
        <input
          class="ruleset-name-input"
          :value="modelValue.Name"
          placeholder="Ruleset name"
          aria-label="Ruleset name"
          @keydown.stop.enter="($event.target as HTMLInputElement).blur()"
          @click.stop
          @input="setName($event)"
        />
        <div class="row-path">{{ ruleSummary }}</div>
      </div>

      <span class="tag" @click.stop>
        {{ ruleCountTag }}
      </span>

      <div class="row-actions no-drag" @click.stop>
        <button
          type="button"
          class="icon-btn danger"
          aria-label="Delete ruleset"
          title="Delete"
          @click="emit('remove')"
        >
          <svg
            viewBox="0 0 24 24"
            fill="none"
            stroke="currentColor"
            stroke-width="1.6"
            stroke-linecap="round"
            stroke-linejoin="round"
          >
            <path
              d="M3 6h18M8 6V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2M19 6l-1 14a2 2 0 0 1-2 2H8a2 2 0 0 1-2-2L5 6"
            />
          </svg>
        </button>
      </div>
    </div>

    <div v-show="expanded" class="row-body">
      <p class="rule-help">
        Prefix <code>d$</code> domain · <code>r$</code> regex · <code>s$</code> contains · (none)
        contains. Use <code>_Switchy</code> as the target to force the picker.
      </p>

      <div v-if="modelValue.Rules.length === 0" class="rule-empty">
        No rules yet.
        <button type="button" class="link-btn" @click.prevent="addRule($event)">Add one</button>
      </div>

      <div v-for="(rule, idx) in modelValue.Rules" :key="idx" class="rule-row">
        <input
          class="s-input rule-pattern"
          :value="rule.Pattern"
          placeholder="d$example.com"
          @input="updateRule(idx, { Pattern: ($event.target as HTMLInputElement).value })"
        />
        <span class="rule-arrow" aria-hidden="true">→</span>
        <select
          class="s-input rule-browser"
          :value="rule.Browser"
          @change="updateRule(idx, { Browser: ($event.target as HTMLSelectElement).value })"
        >
          <option :value="FORCE_PICKER_BROWSER">{{ FORCE_PICKER_BROWSER }} (picker)</option>
          <option v-for="b in browsers" :key="b" :value="b">{{ b }}</option>
        </select>
        <button
          type="button"
          class="icon-btn danger"
          aria-label="Remove rule"
          @click="removeRule(idx)"
        >
          <svg
            width="12"
            height="12"
            viewBox="0 0 24 24"
            fill="none"
            stroke="currentColor"
            stroke-width="2"
            stroke-linecap="round"
          >
            <path d="M18 6 6 18M6 6l12 12" />
          </svg>
        </button>
      </div>

      <button
        v-if="modelValue.Rules.length > 0"
        type="button"
        class="link-btn link-btn-add"
        @click="addRule($event)"
      >
        + Add rule
      </button>
    </div>
  </div>
</template>

<style scoped>
.row {
  border-bottom: 1px solid var(--hairline);
  transition: background 100ms;
}
.row:last-child {
  border-bottom: 0;
}

.row-head {
  display: grid;
  grid-template-columns: 28px 36px 1fr auto auto;
  align-items: center;
  gap: 14px;
  padding: 12px 14px;
  cursor: pointer;
  transition: background 100ms;
}
.row:not(.is-open) .row-head:hover {
  background: var(--ink-10);
}

.chev {
  width: 14px;
  height: 14px;
  color: var(--ink-40);
  transition: transform 180ms ease;
  flex-shrink: 0;
}
.row.is-open .chev {
  transform: rotate(90deg);
  color: var(--accent);
}

.row-meta {
  min-width: 0;
}
.ruleset-name-input {
  display: block;
  width: 100%;
  max-width: 100%;
  font-size: 14px;
  font-weight: 500;
  letter-spacing: -0.005em;
  color: var(--ink-100);
  background: transparent;
  border: 0;
  outline: none;
  padding: 0;
  font-family: var(--font-ui);
}
.ruleset-name-input::placeholder {
  color: var(--ink-40);
}

.row-path {
  margin-top: 2px;
  font-family: var(--font-mono);
  font-size: 11px;
  color: var(--ink-40);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.row-actions {
  display: flex;
  gap: 2px;
  justify-content: flex-end;
}

.row-body {
  padding: 0 14px 16px 96px;
  border-top: 1px solid var(--hairline);
  background: rgba(0, 0, 0, 0.18);
}

.rule-help {
  font-size: 11px;
  line-height: 1.45;
  color: var(--ink-40);
  padding: 14px 0 10px;
  margin: 0;
}
.rule-help code {
  font-family: var(--font-mono);
  font-size: 10.5px;
  color: var(--ink-60);
  padding: 0 4px;
  border-radius: var(--r-tag);
  background: var(--ink-10);
  border: 1px solid var(--hairline);
}

.rule-empty {
  font-size: 11px;
  color: var(--ink-40);
  padding: 6px 0 12px;
}

.rule-row {
  display: grid;
  grid-template-columns: 1fr auto minmax(140px, 200px) 28px;
  align-items: center;
  gap: 8px;
  margin-bottom: 8px;
}
.rule-pattern {
  min-width: 0;
}

.rule-browser {
  width: 100%;
  min-width: 0;
  font-family: var(--font-ui);
  cursor: pointer;
}

.rule-arrow {
  font-size: 12px;
  color: var(--ink-40);
  justify-self: center;
}

.link-btn-add {
  display: inline-block;
  margin-top: 4px;
  margin-bottom: 4px;
}
</style>
