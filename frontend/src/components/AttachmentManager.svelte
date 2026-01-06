<script>
  import { onMount } from "svelte";
  import { api } from "../lib/api.js";
  import AttachmentUploader from "./AttachmentUploader.svelte";
  import AttachmentList from "./AttachmentList.svelte";

  let { taskId, currentUser = null, canUpload = true, canDelete = false } = $props();

  let attachments = $state([]);
  let isLoading = $state(false);
  let error = $state("");

  // API functions
  async function loadAttachments() {
    if (!taskId) return;
    
    isLoading = true;
    error = "";
    
    try {
      const data = await api.attachments.getByTask(taskId);
      attachments = data.attachments || [];
    } catch (err) {
      console.error('Load attachments error:', err);
      error = err.message || 'خطا در بارگذاری فایل‌های پیوست';
    } finally {
      isLoading = false;
    }
  }

  async function deleteAttachment(attachment) {
    try {
      await api.attachments.delete(attachment.id);
      
      // Remove from local state
      attachments = attachments.filter(att => att.id !== attachment.id);
      
    } catch (err) {
      console.error('Delete attachment error:', err);
      error = err.message || 'خطا در حذف فایل';
    }
  }

  // Event handlers
  function handleUploaded(event) {
    const response = event.detail;
    
    // Add new attachments to the list
    if (response.attachments) {
      attachments = [...attachments, ...response.attachments];
    }
    
    // Clear any previous errors
    error = "";
  }

  function handleDelete(event) {
    const attachment = event.detail;
    deleteAttachment(attachment);
  }

  function handleError(event) {
    error = event.detail.message;
  }

  // Load attachments on mount and when taskId changes
  onMount(() => {
    loadAttachments();
  });

  // Watch for taskId changes
  $effect(() => {
    if (taskId) {
      loadAttachments();
    }
  });
</script>

<div class="space-y-6">
  <!-- Header -->
  <div class="flex items-center justify-between">
    <h3 class="text-lg font-medium text-slate-900">فایل‌های پیوست</h3>
    {#if attachments.length > 0}
      <span class="text-sm text-slate-500">
        {attachments.length} فایل
      </span>
    {/if}
  </div>

  <!-- Error Display -->
  {#if error}
    <div class="bg-rose-50 border border-rose-200 rounded-lg p-4">
      <div class="flex items-start">
        <div class="flex-shrink-0">
          <svg class="w-5 h-5 text-rose-400" fill="currentColor" viewBox="0 0 20 20">
            <path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zM8.707 7.293a1 1 0 00-1.414 1.414L8.586 10l-1.293 1.293a1 1 0 101.414 1.414L10 11.414l1.293 1.293a1 1 0 001.414-1.414L11.414 10l1.293-1.293a1 1 0 00-1.414-1.414L10 8.586 8.707 7.293z" clip-rule="evenodd" />
          </svg>
        </div>
        <div class="mr-3">
          <p class="text-sm text-rose-700">{error}</p>
        </div>
        <div class="mr-auto">
          <button
            onclick={() => error = ""}
            class="text-rose-400 hover:text-rose-600"
          >
            <svg class="w-4 h-4" fill="currentColor" viewBox="0 0 20 20">
              <path fill-rule="evenodd" d="M4.293 4.293a1 1 0 011.414 0L10 8.586l4.293-4.293a1 1 0 111.414 1.414L11.414 10l4.293 4.293a1 1 0 01-1.414 1.414L10 11.414l-4.293 4.293a1 1 0 01-1.414-1.414L8.586 10 4.293 5.707a1 1 0 010-1.414z" clip-rule="evenodd" />
            </svg>
          </button>
        </div>
      </div>
    </div>
  {/if}

  <!-- Upload Section -->
  {#if canUpload}
    <div class="border border-slate-200 rounded-lg p-4">
      <h4 class="text-sm font-medium text-slate-700 mb-3">آپلود فایل جدید</h4>
      <AttachmentUploader
        {taskId}
        disabled={!canUpload}
        on:uploaded={handleUploaded}
        on:error={handleError}
      />
    </div>
  {/if}

  <!-- Attachments List -->
  <div class="border border-slate-200 rounded-lg p-4">
    <h4 class="text-sm font-medium text-slate-700 mb-3">فایل‌های موجود</h4>
    
    {#if isLoading}
      <div class="flex items-center justify-center py-8">
        <div class="flex items-center gap-3 text-slate-500">
          <div class="animate-spin rounded-full h-5 w-5 border-b-2 border-indigo-600"></div>
          <span class="text-sm">در حال بارگذاری...</span>
        </div>
      </div>
    {:else}
      <AttachmentList
        {attachments}
        {currentUser}
        {canDelete}
        on:delete={handleDelete}
        on:error={handleError}
      />
    {/if}
  </div>

  <!-- Refresh Button -->
  <div class="flex justify-center">
    <button
      onclick={loadAttachments}
      disabled={isLoading}
      class="px-4 py-2 text-sm text-slate-600 hover:text-slate-900 hover:bg-slate-100 rounded-lg transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
    >
      <div class="flex items-center gap-2">
        <svg class="w-4 h-4 {isLoading ? 'animate-spin' : ''}" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
        </svg>
        <span>به‌روزرسانی</span>
      </div>
    </button>
  </div>
</div>