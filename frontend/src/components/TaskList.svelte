<script>
  import { tasks } from "../stores/taskStore";
  import { timeLogs } from "../stores/timeLogStore";
  import { comments } from "../stores/commentStore.js";
  import { authStore } from "../stores/authStore.js";
  import TaskForm from "./TaskForm.svelte";
  import TimeLogForm from "./TimeLogForm.svelte";
  import CommentList from "./CommentList.svelte";
  import Modal from "./Modal.svelte";
  import { createEventDispatcher } from "svelte";
  import moment from "jalali-moment";

  let { project } = $props();
  const dispatch = createEventDispatcher();

  let showForm = $state(false);
  let selectedTask = $state(null);
  let showCommentsForTask = $state(null);
  let showDeleteModal = $state(false);
  let taskToDelete = $state(null);

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

  async function toggleComments(task) {
    showCommentsForTask = showCommentsForTask === task.id ? null : task.id;
    await comments.load(task.id);
  }

  async function handleTaskToggle(task) {
    await tasks.toggleComplete(task.id);
  }

  function confirmDelete(task) {
    showDeleteModal = true;
    taskToDelete = task;
  }

  async function handleDelete() {
    if (!taskToDelete) return;

    try {
      const taskId = taskToDelete.id;
      await tasks.delete(taskId);
      showDeleteModal = false;
      taskToDelete = null;
      if (selectedTask?.id === taskId) {
        selectedTask = null;
      }
    } catch (error) {
      alert(error.message);
    }
  }
</script>

<div class="space-y-6">
  <!-- Toolbar -->
  <div class="flex items-center justify-between gap-3">
    <div class="text-sm text-slate-500">
      {($tasks || []).length}
      {($tasks || []).length === 1 ? "وظیفه" : "وظیفه"}
    </div>
    <button
      onclick={toggleForm}
      class="px-4 py-3 min-h-[44px] text-sm font-medium rounded-lg transition-colors
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
        <div class="flex flex-col md:flex-row md:items-center gap-3 md:gap-4 p-3 md:p-4">
          <!-- Checkbox -->
          <button
            onclick={() => handleTaskToggle(task)}
            class="flex-shrink-0 w-11 h-11 min-w-[44px] min-h-[44px] rounded-full border-2 transition-all self-start md:self-auto
              {task.completed
              ? 'bg-emerald-500 border-emerald-500'
              : 'border-slate-300 hover:border-indigo-400'}"
          >
            {#if task.completed}
              <svg
                class="w-full h-full text-white p-2"
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
              class="font-medium text-sm md:text-base {task.completed
                ? 'line-through text-slate-400'
                : 'text-slate-900'}"
            >
              {task.title}
            </h3>
            {#if task.description}
              <p class="text-xs md:text-sm text-slate-600 mt-1 line-clamp-2 md:line-clamp-none">{task.description}</p>
            {/if}
            <div class="flex flex-wrap items-center gap-1.5 md:gap-3 mt-2 text-xs text-slate-500">
              {#if task.category}
                <span
                  class="inline-flex items-center px-2 py-0.5 md:px-2.5 rounded bg-blue-50 text-blue-700 font-medium"
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
                <div class="w-full bg-slate-200 rounded-full h-1.5 md:h-2">
                  <div
                    class="h-1.5 md:h-2 rounded-full transition-all {task.done_ratio ===
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
            class="flex-shrink-0 px-2.5 py-1 text-xs font-medium rounded-full self-start md:self-auto
            {task.priority === 'High' ? 'bg-rose-50 text-rose-700' : ''}
            {task.priority === 'Medium' ? 'bg-amber-50 text-amber-700' : ''}
            {task.priority === 'Low' ? 'bg-slate-100 text-slate-600' : ''}"
          >
            {task.priority === 'High' ? 'بالا' : task.priority === 'Medium' ? 'متوسط' : 'پایین'}
          </span>

          <!-- Actions -->
          <div
            class="flex items-center gap-1 md:gap-2 transition-opacity flex-shrink-0"
          >
            <button
              onclick={() => handleTaskSelect(task)}
              class="p-2 hover:bg-slate-100 rounded-lg transition-colors min-w-[44px] min-h-[44px] md:min-w-0 md:min-h-0"
              title="Log time"
            >
              <svg
                class="w-5 h-5 md:w-4 md:h-4 text-slate-600"
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
              onclick={() => toggleComments(task)}
              class="p-2 hover:bg-slate-100 rounded-lg transition-colors min-w-[44px] min-h-[44px] md:min-w-0 md:min-h-0"
              title="Comments"
            >
              <svg
                class="w-5 h-5 md:w-4 md:h-4 text-slate-600"
                fill="none"
                stroke="currentColor"
                stroke-width="2"
                stroke-linecap="round"
                stroke-linejoin="round"
                viewBox="0 0 24 24"
              >
                <path
                  d="M21 15a2 2 0 0 1-2 2H7l-4 4V5a2 2 0 0 1 2-2h14a2 2 0 0 1 2 2z"
                />
              </svg>
            </button>
            <button
              onclick={() => confirmDelete(task)}
              class="p-2 hover:bg-rose-50 rounded-lg transition-colors min-w-[44px] min-h-[44px] md:min-w-0 md:min-h-0"
              title="Delete task"
            >
              <svg
                class="w-5 h-5 md:w-4 md:h-4 text-slate-600 hover:text-rose-600"
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
          <div class="flex-shrink-0 text-right self-start md:self-auto">
            <div class="text-xs md:text-sm font-semibold text-slate-700">0 ساعت 0 دقیقه</div>
            <div class="text-xs text-slate-400">ثبت شده</div>
          </div>
        </div>

        {#if selectedTask?.id === task.id}
          <div class="border-t border-slate-200 bg-slate-50/50 px-3 md:px-4 py-3 md:py-4">
            <TimeLogForm {task} on:logged={() => (selectedTask = null)} />
          </div>
        {/if}

        {#if showCommentsForTask === task.id}
          <div class="border-t border-slate-200 px-3 md:px-4 py-3 md:py-4">
            <CommentList {task} authUser={$authStore.user} />
          </div>
        {/if}
      </div>
    {/each}

    {#if ($tasks || []).length === 0}
      <div class="text-center py-8 md:py-12 px-4">
        <svg
          class="w-10 h-10 md:w-12 md:h-12 mx-auto text-slate-300 mb-2 md:mb-3"
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
        <p class="text-sm md:text-base text-slate-500">هنوز وظیفه‌ای وجود ندارد</p>
        <p class="text-xs md:text-sm text-slate-400 mt-1">
          اولین وظیفه خود را برای شروع ایجاد کنید
        </p>
      </div>
    {/if}
  </div>
</div>

 <Modal show={showDeleteModal} on:close={() => { showDeleteModal = false; taskToDelete = null; }}>
    <div class="p-4 sm:p-6">
      <h3 class="text-lg font-semibold text-slate-900 mb-2">
        حذف وظیفه
      </h3>
      <p class="text-slate-600 mb-4">
        آیا مطمئن هستید که می‌خواهید این وظیفه را حذف کنید؟
      </p>
      <div class="flex flex-col sm:flex-row gap-3 justify-end sm:justify-end">
        <button
          onclick={() => { showDeleteModal = false; taskToDelete = null; }}
          class="w-full sm:w-auto px-4 py-3 min-h-[44px] sm:min-h-0 bg-slate-200 text-slate-700 rounded-lg hover:bg-slate-300 font-medium"
        >
          لغو
        </button>
        <button
          onclick={handleDelete}
          class="w-full sm:w-auto px-4 py-3 min-h-[44px] sm:min-h-0 bg-rose-600 text-white rounded-lg hover:bg-rose-700 font-medium"
        >
          حذف
        </button>
      </div>
    </div>
  </Modal>

