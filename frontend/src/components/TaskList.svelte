<script>
  import { tasks } from "../stores/taskStore";
  import { timeLogs } from "../stores/timeLogStore";
  import { comments } from "../stores/commentStore.js";
  import { authStore } from "../stores/authStore.js";
  import TaskForm from "./TaskForm.svelte";
  import TimeLogForm from "./TimeLogForm.svelte";
  import CommentList from "./CommentList.svelte";
  import { createEventDispatcher } from "svelte";
  import moment from "jalali-moment";

  let { project } = $props();
  const dispatch = createEventDispatcher();

  let showForm = $state(false);
  let selectedTask = $state(null);

  function formatJalaliDate(dateString) {
    if (!dateString) return "";
    return moment(dateString).locale("fa").format("YYYY/MM/DD");
  }

  function toggleForm() {
    showForm = !showForm;
  }

  async function handleTaskSelect(task) {
    selectedTask = task;
    await timeLogs.load(task.id);
    await comments.load(task.id);
  }

  async function handleTaskToggle(task) {
    await tasks.toggleComplete(task.id);
  }

  async function handleTaskDelete(taskId) {
    if (confirm("آیا مطمئن هستید که می‌خواهید این وظیفه را حذف کنید؟")) {
      await tasks.delete(taskId);
      if (selectedTask?.id === taskId) {
        selectedTask = null;
      }
    }
  }
</script>

<div class="space-y-6">
  <!-- Toolbar -->
  <div class="flex items-center justify-between">
    <div class="text-sm text-slate-500">
      {($tasks || []).length}
      {($tasks || []).length === 1 ? "وظیفه" : "وظیفه"}
    </div>
    <button
      onclick={toggleForm}
      class="px-4 py-2 text-sm font-medium rounded-lg transition-colors
        {showForm
        ? 'bg-slate-100 text-slate-700 hover:bg-slate-200'
        : 'bg-slate-900 text-white hover:bg-slate-800'}"
    >
      {showForm ? "لغو" : "+ افزودن وظیفه"}
    </button>
  </div>

  {#if showForm}
    <div class="bg-white rounded-xl shadow-sm border border-slate-200 p-6">
      <TaskForm
        {project}
        on:created={() => {
          showForm = false;
        }}
      />
    </div>
  {/if}

  <!-- Task List -->
  <div class="space-y-3">
    {#each $tasks || [] as task}
      <div
        class="group bg-white rounded-xl shadow-sm border border-slate-200 hover:shadow-md transition-shadow"
      >
        <!-- Task Row -->
        <div class="flex items-center gap-4 p-4">
          <!-- Checkbox -->
          <button
            onclick={() => handleTaskToggle(task)}
            class="flex-shrink-0 w-5 h-5 rounded-full border-2 transition-all
              {task.completed
              ? 'bg-emerald-500 border-emerald-500'
              : 'border-slate-300 hover:border-indigo-400'}"
          >
            {#if task.completed}
              <svg
                class="w-full h-full text-white"
                fill="none"
                stroke="currentColor"
                viewBox="0 0 24 24"
              >
                <path
                  stroke-linecap="round"
                  stroke-linejoin="round"
                  stroke-width="3"
                  d="M5 13l4 4L19 7"
                />
              </svg>
            {/if}
          </button>

          <!-- Task Info -->
          <div class="flex-1 min-w-0">
            <h3
              class="font-medium {task.completed
                ? 'line-through text-slate-400'
                : 'text-slate-900'}"
            >
              {task.title}
            </h3>
            {#if task.description}
              <p class="text-sm text-slate-600 mt-1">{task.description}</p>
            {/if}
            <div class="flex items-center gap-3 mt-2 text-xs text-slate-500">
              {#if task.category}
                <span
                  class="inline-flex items-center px-2 py-0.5 rounded bg-blue-50 text-blue-700 font-medium"
                >
                  {task.category}
                </span>
              {/if}
              {#if task.start_date}
                <span
                  >شروع: {formatJalaliDate(task.start_date)}</span
                >
              {/if}
              {#if task.due_date}
                <span
                  class="font-medium {new Date(task.due_date) < new Date() &&
                  !task.completed
                    ? 'text-rose-600'
                    : ''}"
                >
                  مهلت: {formatJalaliDate(task.due_date)}
                </span>
              {/if}
              {#if task.estimated_hours}
                <span>تخمین: {task.estimated_hours} ساعت</span>
              {/if}
            </div>
            {#if task.done_ratio > 0}
              <div class="mt-2">
                <div
                  class="flex items-center justify-between text-xs text-slate-600 mb-1"
                >
                  <span>پیشرفت</span>
                  <span class="font-medium">{task.done_ratio}%</span>
                </div>
                <div class="w-full bg-slate-200 rounded-full h-1.5">
                  <div
                    class="h-1.5 rounded-full transition-all {task.done_ratio ===
                    100
                      ? 'bg-emerald-500'
                      : 'bg-indigo-500'}"
                    style="width: {task.done_ratio}%"
                  ></div>
                </div>
              </div>
            {/if}
          </div>

          <!-- Priority Badge -->
          <span
            class="flex-shrink-0 px-2.5 py-1 text-xs font-medium rounded-full
            {task.priority === 'High' ? 'bg-rose-50 text-rose-700' : ''}
            {task.priority === 'Medium' ? 'bg-amber-50 text-amber-700' : ''}
            {task.priority === 'Low' ? 'bg-slate-100 text-slate-600' : ''}"
          >
            {task.priority === 'High' ? 'بالا' : task.priority === 'Medium' ? 'متوسط' : 'پایین'}
          </span>

          <!-- Hover Actions -->
          <div
            class="flex items-center gap-2 opacity-0 group-hover:opacity-100 transition-opacity"
          >
            <button
              onclick={() => handleTaskSelect(task)}
              class="p-1.5 hover:bg-slate-100 rounded-lg transition-colors"
              title="Log time"
            >
              <svg
                class="w-4 h-4 text-slate-600"
                fill="none"
                stroke="currentColor"
                viewBox="0 0 24 24"
              >
                <path
                  stroke-linecap="round"
                  stroke-linejoin="round"
                  stroke-width="2"
                  d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z"
                />
              </svg>
            </button>
            <button
              onclick={() => handleTaskDelete(task.id)}
              class="p-1.5 hover:bg-rose-50 rounded-lg transition-colors"
              title="Delete task"
            >
              <svg
                class="w-4 h-4 text-slate-600 hover:text-rose-600"
                fill="none"
                stroke="currentColor"
                viewBox="0 0 24 24"
              >
                <path
                  stroke-linecap="round"
                  stroke-linejoin="round"
                  stroke-width="2"
                  d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"
                />
              </svg>
            </button>
          </div>

          <!-- Total Time -->
          <div class="flex-shrink-0 text-right">
            <div class="text-sm font-semibold text-slate-700">0 ساعت 0 دقیقه</div>
            <div class="text-xs text-slate-400">ثبت شده</div>
          </div>
        </div>

        <!-- Expanded: Time Logging Panel -->
        {#if selectedTask?.id === task.id}
          <div class="border-t border-slate-200 bg-slate-50/50 px-4 py-4">
            <TimeLogForm {task} on:logged={() => (selectedTask = null)} />
          </div>
          <div class="border-t border-slate-200 px-4 py-4">
            <CommentList {task} authUser={$authStore.user} />
          </div>
        {/if}
      </div>
    {/each}

    {#if ($tasks || []).length === 0}
      <div class="text-center py-12">
        <svg
          class="w-12 h-12 mx-auto text-slate-300 mb-3"
          fill="none"
          stroke="currentColor"
          viewBox="0 0 24 24"
        >
          <path
            stroke-linecap="round"
            stroke-linejoin="round"
            stroke-width="1.5"
            d="M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2m-6 9l2 2 4-4"
          />
        </svg>
        <p class="text-slate-500">هنوز وظیفه‌ای وجود ندارد</p>
        <p class="text-slate-400 text-sm mt-1">
          اولین وظیفه خود را برای شروع ایجاد کنید
        </p>
      </div>
    {/if}
  </div>
</div>
