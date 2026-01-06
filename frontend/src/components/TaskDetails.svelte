<script>
  import { tasks } from "../stores/taskStore.js";
  import { timeLogs } from "../stores/timeLogStore.js";
  import { comments } from "../stores/commentStore.js";
  import { authStore } from "../stores/authStore.js";
  import { createEventDispatcher } from "svelte";
  import JalaliDatePicker from "./JalaliDatePicker.svelte";
  import TimeLogForm from "./TimeLogForm.svelte";
  import CommentList from "./CommentList.svelte";
  import AttachmentManager from "./AttachmentManager.svelte";
  import moment from "jalali-moment";

  let { task, project } = $props();
  const dispatch = createEventDispatcher();

  let isEditing = $state(false);
  let showComments = $state(false);
  let showTimeLogs = $state(false);
  let showAttachments = $state(false);

  // Local state for editing
  let title = $state("");
  let description = $state("");
  let priority = $state("Medium");
  let category = $state("");
  let start_date = $state("");
  let due_date = $state("");
  let estimated_hours = $state("");
  let done_ratio = $state(0);

  // Validation errors
  let error = $state("");
  let dateError = $state("");
  let doneRatioError = $state("");
  let estimatedHoursError = $state("");

  function formatJalaliDate(dateString) {
    if (!dateString) return "";
    return moment(dateString).locale("fa").format("YYYY/MM/DD HH:mm");
  }

  function formatJalaliDateOnly(dateString) {
    if (!dateString) return "";
    return moment(dateString).locale("fa").format("YYYY/MM/DD");
  }

  function enterEditMode() {
    title = task.title;
    description = task.description || "";
    priority = task.priority;
    category = task.category || "";
    // Store dates in YYYY-MM-DD format for JalaliDatePicker to properly convert
    start_date = task.start_date ? task.start_date.split('T')[0] : "";
    due_date = task.due_date ? task.due_date.split('T')[0] : "";
    estimated_hours = task.estimated_hours ? task.estimated_hours.toString() : "";
    done_ratio = task.done_ratio;
    isEditing = true;
    error = "";
    dateError = "";
    doneRatioError = "";
    estimatedHoursError = "";
  }

  function cancelEdit() {
    isEditing = false;
    error = "";
    dateError = "";
    doneRatioError = "";
    estimatedHoursError = "";
  }

  // Validate date range
  function validateDates() {
    if (start_date && due_date && new Date(due_date) < new Date(start_date)) {
      dateError = "تاریخ مهلت باید بعد از تاریخ شروع یا برابر با آن باشد";
      return false;
    }
    dateError = "";
    return true;
  }

  // Validate done_ratio
  function validateDoneRatio() {
    const ratio = parseInt(done_ratio);
    if (isNaN(ratio) || ratio < 0 || ratio > 100) {
      doneRatioError = "درصد پیشرفت باید بین 0 تا 100 باشد";
      return false;
    }
    doneRatioError = "";
    return true;
  }

  // Validate estimated_hours
  function validateEstimatedHours() {
    if (estimated_hours !== "" && estimated_hours !== null) {
      const hours = parseFloat(estimated_hours);
      if (isNaN(hours) || hours < 0) {
        estimatedHoursError = "ساعات تخمینی باید بزرگتر یا مساوی 0 باشد";
        return false;
      }
    }
    estimatedHoursError = "";
    return true;
  }

  async function handleSave() {
    error = "";

    if (!title.trim()) {
      error = "عنوان الزامی است";
      return;
    }

    const isDatesValid = validateDates();
    const isDoneRatioValid = validateDoneRatio();
    const isEstimatedHoursValid = validateEstimatedHours();

    if (!isDatesValid || !isDoneRatioValid || !isEstimatedHoursValid) {
      return;
    }

    try {
      await tasks.update(task.id, {
        title: title.trim(),
        description: description.trim(),
        priority,
        category: category.trim() || null,
        start_date: start_date ? new Date(start_date).toISOString() : null,
        due_date: due_date ? new Date(due_date).toISOString() : null,
        estimated_hours: estimated_hours ? parseFloat(estimated_hours) : null,
        done_ratio: parseInt(done_ratio),
        completed: task.completed,
      });

      // Refresh task data
      const updatedTask = await tasks.getById(task.id);
      if (updatedTask) {
        task = updatedTask;
      }

      isEditing = false;
      dispatch("updated");
    } catch (err) {
      error = err.message || "به‌روزرسانی وظیفه با خطا مواجه شد";
    }
  }

  function toggleComments() {
    showComments = !showComments;
    if (showComments) {
      comments.load(task.id);
    }
  }

  function toggleTimeLogs() {
    showTimeLogs = !showTimeLogs;
    if (showTimeLogs) {
      timeLogs.load(task.id);
    }
  }

  function toggleAttachments() {
    showAttachments = !showAttachments;
  }
</script>

{#if isEditing}
  <!-- Edit Mode -->
  <div class="space-y-6">
    <div class="flex items-center justify-between">
      <h3 class="text-xl font-bold text-slate-900">ویرایش وظیفه</h3>
      <button onclick={cancelEdit} class="text-slate-600 hover:text-slate-900 font-medium">
        لغو
      </button>
    </div>

    {#if error}
      <div class="p-3 bg-rose-100 text-rose-700 rounded-lg text-sm">
        {error}
      </div>
    {/if}

    <div>
      <label for="task-title" class="block text-sm font-medium text-slate-700 mb-1">عنوان</label>
      <input
        type="text"
        id="task-title"
        bind:value={title}
        class="w-full px-3 py-3 min-h-[44px] border border-slate-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-indigo-500"
        placeholder="عنوان وظیفه"
        required
      />
    </div>

    <div>
      <label for="description" class="block text-sm font-medium text-slate-700 mb-1">توضیحات</label>
      <textarea
        id="description"
        bind:value={description}
        rows="4"
        class="w-full px-3 py-3 border border-slate-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-indigo-500 resize-none"
        placeholder="توضیحات وظیفه"
      ></textarea>
    </div>

    <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
      <div>
        <label for="priority" class="block text-sm font-medium text-slate-700 mb-1">اولویت</label>
        <select
          id="priority"
          bind:value={priority}
          class="w-full px-3 py-3 min-h-[44px] border border-slate-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-indigo-500"
        >
          <option value="Low">پایین</option>
          <option value="Medium">متوسط</option>
          <option value="High">بالا</option>
        </select>
      </div>

      <div>
        <label for="category" class="block text-sm font-medium text-slate-700 mb-1">دسته‌بندی</label>
        <input
          type="text"
          id="category"
          bind:value={category}
          class="w-full px-3 py-3 min-h-[44px] border border-slate-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-indigo-500"
          placeholder="بک‌اند، فرانت‌اند، ..."
        />
      </div>
    </div>

    <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
      <div>
        <label for="start_date" class="block text-sm font-medium text-slate-700 mb-1">تاریخ شروع</label>
        <JalaliDatePicker
          bind:value={start_date}
          onchange={validateDates}
          placeholder="1403/10/10"
          error={dateError}
        />
      </div>

      <div>
        <label for="due_date" class="block text-sm font-medium text-slate-700 mb-1">تاریخ مهلت</label>
        <JalaliDatePicker
          bind:value={due_date}
          onchange={validateDates}
          placeholder="1403/10/20"
          error={dateError}
        />
      </div>
    </div>
    {#if dateError}
      <p class="text-rose-500 text-xs">{dateError}</p>
    {/if}

    <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
      <div>
        <label for="estimated_hours" class="block text-sm font-medium text-slate-700 mb-1">ساعات تخمینی</label>
        <input
          type="number"
          id="estimated_hours"
          bind:value={estimated_hours}
          onblur={validateEstimatedHours}
          min="0"
          step="0.5"
          class="w-full px-3 py-3 min-h-[44px] border border-slate-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-indigo-500"
          class:border-rose-500={estimatedHoursError}
          placeholder="8.5"
        />
        {#if estimatedHoursError}
          <p class="text-rose-500 text-xs mt-1">{estimatedHoursError}</p>
        {/if}
      </div>

      <div>
        <label for="done_ratio" class="block text-sm font-medium text-slate-700 mb-1">
          پیشرفت (%) - {done_ratio}%
        </label>
        <input
          type="range"
          id="done_ratio"
          bind:value={done_ratio}
          onchange={validateDoneRatio}
          min="0"
          max="100"
          step="5"
          class="w-full h-3 bg-slate-200 rounded-lg appearance-none cursor-pointer"
          class:border-rose-500={doneRatioError}
        />
        {#if doneRatioError}
          <p class="text-rose-500 text-xs mt-1">{doneRatioError}</p>
        {/if}
      </div>
    </div>

    <div class="flex gap-3">
      <button
        onclick={handleSave}
        class="flex-1 min-h-[44px] bg-indigo-600 hover:bg-indigo-700 text-white px-4 py-3 rounded-lg transition-colors font-medium"
      >
        ذخیره تغییرات
      </button>
      <button
        onclick={cancelEdit}
        class="flex-1 min-h-[44px] bg-slate-200 hover:bg-slate-300 text-slate-700 px-4 py-3 rounded-lg transition-colors font-medium"
      >
        لغو
      </button>
    </div>
  </div>
{:else}
  <!-- View Mode -->
  <div class="space-y-6">
    <!-- Header -->
    <div class="flex items-start justify-between gap-4">
      <div class="flex-1">
        <div class="flex items-center gap-3 mb-2">
          <h2 class="text-xl md:text-2xl font-bold text-slate-900">{task.title}</h2>
          <span
            class="px-2.5 py-1 text-xs font-medium rounded-full
            {task.priority === 'High' ? 'bg-rose-50 text-rose-700' : ''}
            {task.priority === 'Medium' ? 'bg-amber-50 text-amber-700' : ''}
            {task.priority === 'Low' ? 'bg-slate-100 text-slate-600' : ''}"
          >
            {task.priority === 'High' ? 'بالا' : task.priority === 'Medium' ? 'متوسط' : 'پایین'}
          </span>
          {#if task.completed}
            <span class="px-2.5 py-1 text-xs font-medium rounded-full bg-emerald-100 text-emerald-700">
              تکمیل شده
            </span>
          {/if}
        </div>
      </div>
      <button
        onclick={enterEditMode}
        class="px-4 py-2 bg-indigo-600 hover:bg-indigo-700 text-white rounded-lg transition-colors font-medium text-sm"
      >
        ویرایش
      </button>
    </div>

    <!-- Description -->
    <div>
      <h3 class="text-sm font-medium text-slate-600 mb-2">توضیحات</h3>
      <p class="text-slate-800 whitespace-pre-wrap">{task.description || 'بدون توضیحات'}</p>
    </div>

    <!-- Details Grid -->
    <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
      <div>
        <h3 class="text-sm font-medium text-slate-600 mb-2">دسته‌بندی</h3>
        <div class="text-slate-800">{task.category || '-'}</div>
      </div>

      <div>
        <h3 class="text-sm font-medium text-slate-600 mb-2">تاریخ شروع</h3>
        <div class="text-slate-800">
          {task.start_date ? formatJalaliDateOnly(task.start_date) : '-'}
        </div>
      </div>

      <div>
        <h3 class="text-sm font-medium text-slate-600 mb-2">تاریخ مهلت</h3>
        <div
          class="text-slate-800 {new Date(task.due_date) < new Date() &&
          !task.completed
            ? 'text-rose-600 font-medium'
            : ''}"
        >
          {task.due_date ? formatJalaliDateOnly(task.due_date) : '-'}
        </div>
      </div>

      <div>
        <h3 class="text-sm font-medium text-slate-600 mb-2">ساعات تخمینی</h3>
        <div class="text-slate-800">
          {task.estimated_hours ? `${task.estimated_hours} ساعت` : '-'}
        </div>
      </div>
    </div>

    <!-- Progress -->
    {#if task.done_ratio > 0}
      <div>
        <div class="flex items-center justify-between text-sm text-slate-600 mb-1">
          <span>پیشرفت</span>
          <span class="font-medium">{task.done_ratio}%</span>
        </div>
        <div class="w-full bg-slate-200 rounded-full h-2">
          <div
            class="h-2 rounded-full transition-all {task.done_ratio === 100
              ? 'bg-emerald-500'
              : 'bg-indigo-500'}"
            style="width: {task.done_ratio}%"
          ></div>
        </div>
      </div>
    {/if}

    <!-- User Information -->
    <div class="bg-slate-50 rounded-lg p-4">
      <h4 class="font-medium text-slate-900 mb-3">اطلاعات کاربران</h4>
      <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
        <div>
          <span class="text-sm text-slate-600">سازنده:</span>
          <span class="font-medium text-slate-800 mr-2">
            {task.author_name || 'نامشخص'}
          </span>
        </div>
        <div>
          <span class="text-sm text-slate-600">مسئول:</span>
          <span class="font-medium text-slate-800 mr-2">
            {task.assignee_name || 'بدون مسئول'}
          </span>
        </div>
      </div>
    </div>

    <!-- Expandable Sections -->
    <div class="space-y-3">
      <button
        onclick={toggleComments}
        class="w-full flex items-center justify-between px-4 py-3 bg-slate-50 hover:bg-slate-100 rounded-lg transition-colors"
      >
        <span class="font-medium text-slate-700">نظرات</span>
        <svg
          class="w-5 h-5 text-slate-600 transition-transform {showComments
            ? 'rotate-180'
            : ''}"
          fill="none"
          stroke="currentColor"
          viewBox="0 0 24 24"
        >
          <path
            stroke-linecap="round"
            stroke-linejoin="round"
            stroke-width="2"
            d="M19 9l-7 7-7-7"
          />
        </svg>
      </button>

      {#if showComments}
        <div class="border-t border-slate-200 pt-4">
          <CommentList {task} authUser={$authStore.user} />
        </div>
      {/if}

      <button
        onclick={toggleTimeLogs}
        class="w-full flex items-center justify-between px-4 py-3 bg-slate-50 hover:bg-slate-100 rounded-lg transition-colors"
      >
        <span class="font-medium text-slate-700">ثبت زمان</span>
        <svg
          class="w-5 h-5 text-slate-600 transition-transform {showTimeLogs
            ? 'rotate-180'
            : ''}"
          fill="none"
          stroke="currentColor"
          viewBox="0 0 24 24"
        >
          <path
            stroke-linecap="round"
            stroke-linejoin="round"
            stroke-width="2"
            d="M19 9l-7 7-7-7"
          />
        </svg>
      </button>

      {#if showTimeLogs}
        <div class="border-t border-slate-200 pt-4">
          <TimeLogForm {task} on:logged={() => {}} />
        </div>
      {/if}

      <button
        onclick={toggleAttachments}
        class="w-full flex items-center justify-between px-4 py-3 bg-slate-50 hover:bg-slate-100 rounded-lg transition-colors"
      >
        <span class="font-medium text-slate-700">فایل‌های پیوست</span>
        <svg
          class="w-5 h-5 text-slate-600 transition-transform {showAttachments
            ? 'rotate-180'
            : ''}"
          fill="none"
          stroke="currentColor"
          viewBox="0 0 24 24"
        >
          <path
            stroke-linecap="round"
            stroke-linejoin="round"
            stroke-width="2"
            d="M19 9l-7 7-7-7"
          />
        </svg>
      </button>

      {#if showAttachments}
        <div class="border-t border-slate-200 pt-4">
          <AttachmentManager
            taskId={task.id}
            currentUser={$authStore.user}
            canUpload={true}
            canDelete={$authStore.user?.role === 'admin' || task.author_id === $authStore.user?.id || task.assignee_id === $authStore.user?.id}
          />
        </div>
      {/if}
    </div>

    <!-- Timestamps -->
    <div class="text-xs text-slate-400 border-t border-slate-200 pt-4 mt-6 space-y-1">
      <div>ایجاد: {formatJalaliDate(task.created_at)}</div>
      <div>به‌روزرسانی: {formatJalaliDate(task.updated_at)}</div>
    </div>
  </div>
{/if}
