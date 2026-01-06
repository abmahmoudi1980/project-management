<script>
  import { tasks } from "../stores/taskStore";
  import { api } from "../lib/api.js";
  import { createEventDispatcher } from "svelte";
  import JalaliDatePicker from "./JalaliDatePicker.svelte";
  import AttachmentFormUploader from "./AttachmentFormUploader.svelte";

  let { project } = $props();
  const dispatch = createEventDispatcher();

  let title = $state("");
  let description = $state("");
  let priority = $state("Medium");
  let category = $state("");
  let start_date = $state("");
  let due_date = $state("");
  let estimated_hours = $state("");
  let done_ratio = $state(0);
  let attachmentFiles = $state([]);
  let error = $state("");
  let dateError = $state("");
  let doneRatioError = $state("");
  let estimatedHoursError = $state("");
  let attachmentError = $state("");

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
        estimatedHoursError =
          "ساعات تخمینی باید بزرگتر یا مساوی 0 باشد";
        return false;
      }
    }
    estimatedHoursError = "";
    return true;
  }

  async function handleSubmit() {
    error = "";
    attachmentError = "";

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
      // Create the task first
      const newTask = await tasks.create(project.id, {
        title: title.trim(),
        description: description.trim(),
        priority,
        category: category.trim() || null,
        start_date: start_date ? new Date(start_date).toISOString() : null,
        due_date: due_date ? new Date(due_date).toISOString() : null,
        estimated_hours: estimated_hours ? parseFloat(estimated_hours) : null,
        done_ratio: parseInt(done_ratio),
      });

      // Upload attachments if any
      if (attachmentFiles.length > 0) {
        try {
          await api.attachments.upload(newTask.id, attachmentFiles);
        } catch (attachmentErr) {
          console.error('Attachment upload error:', attachmentErr);
          attachmentError = "وظیفه ایجاد شد اما آپلود فایل‌ها با خطا مواجه شد: " + (attachmentErr.message || "خطای نامشخص");
          // Don't return here - task was created successfully
        }
      }

      // Reset form
      title = "";
      description = "";
      priority = "Medium";
      category = "";
      start_date = "";
      due_date = "";
      estimated_hours = "";
      done_ratio = 0;
      attachmentFiles = [];
      error = "";
      dateError = "";
      doneRatioError = "";
      estimatedHoursError = "";
      attachmentError = "";

      dispatch("created");
    } catch (err) {
      error = err.message || "ایجاد وظیفه با خطا مواجه شد";
    }
  }

  // Attachment event handlers
  function handleFilesAdded(event) {
    const newFiles = event.detail.files;
    attachmentFiles = [...attachmentFiles, ...newFiles];
    attachmentError = "";
  }

  function handleFileRemoved(event) {
    const index = event.detail.index;
    attachmentFiles = attachmentFiles.filter((_, i) => i !== index);
  }

  function handleAttachmentError(event) {
    attachmentError = event.detail.message;
  }
</script>

<form
  onsubmit={(e) => { e.preventDefault(); handleSubmit(); }}
  class="space-y-4 p-4 border rounded-lg bg-white"
>
  <h3 class="text-lg font-semibold text-gray-800">ایجاد وظیفه جدید</h3>

  {#if error}
    <div class="p-3 bg-red-100 text-red-700 rounded-lg text-sm">
      {error}
    </div>
  {/if}

  {#if attachmentError}
    <div class="p-3 bg-amber-100 text-amber-700 rounded-lg text-sm">
      {attachmentError}
    </div>
  {/if}

  <div>
    <label for="task-title" class="block text-sm font-medium text-gray-700 mb-1"
      >عنوان</label
    >
    <input
      type="text"
      id="task-title"
      bind:value={title}
      class="w-full px-3 py-3 min-h-[44px] border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
      placeholder="عنوان وظیفه"
      required
    />
  </div>

  <div>
    <label
      for="description"
      class="block text-sm font-medium text-gray-700 mb-1">توضیحات</label
    >
    <textarea
      id="description"
      bind:value={description}
      rows="3"
      class="w-full px-3 py-3 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500 resize-none"
      placeholder="توضیحات وظیفه (اختیاری)"
    ></textarea>
  </div>

  <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
    <div>
      <label for="priority" class="block text-sm font-medium text-gray-700 mb-1"
        >اولویت</label
      >
      <select
        id="priority"
        bind:value={priority}
        class="w-full px-3 py-3 min-h-[44px] border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
      >
        <option value="Low">پایین</option>
        <option value="Medium">متوسط</option>
        <option value="High">بالا</option>
      </select>
    </div>

    <div>
      <label for="category" class="block text-sm font-medium text-gray-700 mb-1"
        >دسته‌بندی</label
      >
      <input
        type="text"
        id="category"
        bind:value={category}
        class="w-full px-3 py-3 min-h-[44px] border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
        placeholder="بک‌اند، فرانت‌اند، ..."
      />
    </div>
  </div>

  <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
    <div>
      <label
        for="start_date"
        class="block text-sm font-medium text-gray-700 mb-1">تاریخ شروع</label
      >
      <JalaliDatePicker
        bind:value={start_date}
        onchange={validateDates}
        placeholder="1403/10/10"
        error={dateError}
      />
    </div>

    <div>
      <label for="due_date" class="block text-sm font-medium text-gray-700 mb-1"
        >تاریخ مهلت</label
      >
      <JalaliDatePicker
        bind:value={due_date}
        onchange={validateDates}
        placeholder="1403/10/20"
        error={dateError}
      />
    </div>
  </div>
  {#if dateError}
    <p class="text-red-500 text-xs -mt-2">{dateError}</p>
  {/if}

  <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
    <div>
      <label
        for="estimated_hours"
        class="block text-sm font-medium text-gray-700 mb-1"
        >ساعات تخمینی</label
      >
      <input
        type="number"
        id="estimated_hours"
        bind:value={estimated_hours}
        onblur={validateEstimatedHours}
        min="0"
        step="0.5"
        class="w-full px-3 py-3 min-h-[44px] border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
        class:border-red-500={estimatedHoursError}
        placeholder="8.5"
      />
      {#if estimatedHoursError}
        <p class="text-red-500 text-xs mt-1">{estimatedHoursError}</p>
      {/if}
    </div>

    <div>
      <label
        for="done_ratio"
        class="block text-sm font-medium text-gray-700 mb-1"
      >
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
        class="w-full h-3 bg-gray-200 rounded-lg appearance-none cursor-pointer"
        class:border-red-500={doneRatioError}
      />
      {#if doneRatioError}
        <p class="text-red-500 text-xs mt-1">{doneRatioError}</p>
      {/if}
    </div>
  </div>

  <!-- Attachments Section -->
  <div>
    <label class="block text-sm font-medium text-gray-700 mb-2">فایل‌های پیوست (اختیاری)</label>
    <AttachmentFormUploader
      bind:files={attachmentFiles}
      maxFiles={5}
      on:filesAdded={handleFilesAdded}
      on:fileRemoved={handleFileRemoved}
      on:error={handleAttachmentError}
    />
  </div>

  <button
    type="submit"
    disabled={!title.trim()}
    class="w-full min-h-[44px] bg-blue-500 hover:bg-blue-600 disabled:bg-gray-300 text-white px-4 py-3 rounded-lg transition-colors font-medium"
  >
    ایجاد وظیفه
  </button>
</form>
