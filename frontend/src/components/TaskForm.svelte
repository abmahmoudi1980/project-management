<script>
  import { tasks } from "../stores/taskStore";
  import { api } from "../lib/api.js";
  import { toasts } from "../lib/toastStore.js";
  import { createEventDispatcher } from "svelte";
  import JalaliDatePicker from "./JalaliDatePicker.svelte";
  import AttachmentFormUploader from "./AttachmentFormUploader.svelte";
  import Button from "./ui/Button.svelte";
  import Input from "./ui/Input.svelte";
  import Textarea from "./ui/Textarea.svelte";

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
  let isSubmitting = $state(false);
  let error = $state("");
  let dateError = $state("");
  let doneRatioError = $state("");
  let estimatedHoursError = $state("");
  let attachmentError = $state("");

  function validateDates() {
    if (start_date && due_date && new Date(due_date) < new Date(start_date)) {
      dateError = "تاریخ مهلت باید بعد از تاریخ شروع یا برابر با آن باشد";
      return false;
    }
    dateError = "";
    return true;
  }

  function validateDoneRatio() {
    const ratio = parseInt(done_ratio);
    if (isNaN(ratio) || ratio < 0 || ratio > 100) {
      doneRatioError = "درصد پیشرفت باید بین 0 تا 100 باشد";
      return false;
    }
    doneRatioError = "";
    return true;
  }

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

  async function handleSubmit() {
    error = "";
    attachmentError = "";

    if (!title.trim()) {
      error = "عنوان الزامی است";
      return;
    }

    if (!validateDates() || !validateDoneRatio() || !validateEstimatedHours()) {
      return;
    }

    isSubmitting = true;

    try {
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

      if (attachmentFiles.length > 0) {
        try {
          await api.attachments.upload(newTask.id, attachmentFiles);
        } catch (attachmentErr) {
          console.error('Attachment upload error:', attachmentErr);
          attachmentError = "وظیفه ایجاد شد اما آپلود فایل‌ها با خطا مواجه شد: " + (attachmentErr.message || "خطای نامشخص");
        }
      }

      toasts.success('وظیفه با موفقیت ایجاد شد');

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
    } finally {
      isSubmitting = false;
    }
  }

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
  class="space-y-4"
>
  <h3 class="text-lg font-semibold text-slate-800">ایجاد وظیفه جدید</h3>

  {#if error}
    <div class="p-3 bg-danger-50 text-danger-700 rounded-lg text-sm" role="alert">
      {error}
    </div>
  {/if}

  {#if attachmentError}
    <div class="p-3 bg-warning-50 text-warning-700 rounded-lg text-sm" role="alert">
      {attachmentError}
    </div>
  {/if}

  <Input
    id="task-title"
    label="عنوان"
    bind:value={title}
    placeholder="عنوان وظیفه"
    required
  />

  <Textarea
    id="description"
    label="توضیحات"
    bind:value={description}
    rows={3}
    placeholder="توضیحات وظیفه (اختیاری)"
  />

  <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
    <div>
      <label for="priority" class="block text-sm font-medium text-slate-700 mb-1.5">اولویت</label>
      <select
        id="priority"
        bind:value={priority}
        class="w-full px-3 py-2.5 text-sm min-h-[44px] border border-slate-300 rounded-lg transition-colors focus:outline-none focus:ring-2 focus:ring-brand-500 focus:border-transparent"
      >
        <option value="Low">پایین</option>
        <option value="Medium">متوسط</option>
        <option value="High">بالا</option>
      </select>
    </div>

    <Input
      id="category"
      label="دسته‌بندی"
      bind:value={category}
      placeholder="بک‌اند، فرانت‌اند، ..."
    />
  </div>

  <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
    <div>
      <label for="start_date" class="block text-sm font-medium text-slate-700 mb-1.5">تاریخ شروع</label>
      <JalaliDatePicker
        bind:value={start_date}
        onchange={validateDates}
        placeholder="1403/10/10"
        error={dateError}
      />
    </div>

    <div>
      <label for="due_date" class="block text-sm font-medium text-slate-700 mb-1.5">تاریخ مهلت</label>
      <JalaliDatePicker
        bind:value={due_date}
        onchange={validateDates}
        placeholder="1403/10/20"
        error={dateError}
      />
    </div>
  </div>

  <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
    <Input
      id="estimated_hours"
      type="number"
      label="ساعات تخمینی"
      bind:value={estimated_hours}
      onblur={validateEstimatedHours}
      min="0"
      step="0.5"
      placeholder="8.5"
      error={estimatedHoursError}
    />

    <div>
      <label for="done_ratio" class="block text-sm font-medium text-slate-700 mb-1.5">
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
        class="w-full h-3 bg-slate-200 rounded-lg appearance-none cursor-pointer
          {doneRatioError ? 'ring-2 ring-danger-500' : ''}"
      />
      {#if doneRatioError}
        <p class="text-xs text-danger-600 mt-1" role="alert">{doneRatioError}</p>
      {/if}
    </div>
  </div>

  <div>
    <label class="block text-sm font-medium text-slate-700 mb-2">فایل‌های پیوست (اختیاری)</label>
    <AttachmentFormUploader
      bind:files={attachmentFiles}
      maxFiles={5}
      on:filesAdded={handleFilesAdded}
      on:fileRemoved={handleFileRemoved}
      on:error={handleAttachmentError}
    />
  </div>

  <Button type="submit" disabled={!title.trim()} loading={isSubmitting} fullWidth>
    ایجاد وظیفه
  </Button>
</form>
