<script>
  import { comments } from "../stores/commentStore.js";
  import { authStore } from "../stores/authStore.js";
  import moment from "jalali-moment";
  import Modal from "./Modal.svelte";

  let { task, authUser } = $props();

  let newComment = $state("");
  let showDeleteModal = $state(false);
  let commentToDelete = $state(null);
  let editingCommentId = $state(null);
  let editingContent = $state("");

  function formatDate(dateString) {
    if (!dateString) return "";
    return moment(dateString).locale("fa").fromNow();
  }

  async function handleAddComment() {
    if (!newComment.trim()) return;

    try {
      await comments.create(task.id, { content: newComment });
      newComment = "";
    } catch (error) {
      alert(error.message);
    }
  }

  function startEdit(comment) {
    editingCommentId = comment.id;
    editingContent = comment.content;
  }

  async function handleSaveEdit(commentId) {
    if (!editingContent.trim()) return;

    try {
      await comments.update(commentId, { content: editingContent });
      editingCommentId = null;
      editingContent = "";
    } catch (error) {
      alert(error.message);
    }
  }

  function cancelEdit() {
    editingCommentId = null;
    editingContent = "";
  }

  function confirmDelete(comment) {
    showDeleteModal = true;
    commentToDelete = comment;
  }

  async function handleDelete() {
    if (!commentToDelete) return;

    try {
      await comments.delete(commentToDelete.id);
      showDeleteModal = false;
      commentToDelete = null;
    } catch (error) {
      alert(error.message);
    }
  }

  function isOwnComment(comment) {
    return authUser?.id === comment.user_id;
  }
</script>

<div class="space-y-4">
  <h3 class="text-lg font-semibold text-slate-900 mb-4">کامنت‌ها</h3>

  {#if $comments.length === 0}
    <div class="text-center py-8 text-slate-500">
      هنوز کامنتی ثبت نشده است
    </div>
  {/if}

  {#each $comments as comment (comment.id)}
    <div class="bg-white rounded-xl border border-slate-200 p-4">
      <div class="flex items-start justify-between gap-4">
        <div class="flex-1">
          <div class="flex items-center gap-2 mb-2">
            <span class="font-medium text-slate-900">{comment.username || "کاربر"}</span>
            <span class="text-xs text-slate-500">{formatDate(comment.created_at)}</span>
          </div>

          {#if editingCommentId === comment.id}
            <div class="space-y-2">
              <textarea
                bind:value={editingContent}
                class="w-full px-3 py-2 border border-slate-200 rounded-lg resize-none"
                rows="3"
              ></textarea>
              <div class="flex gap-2">
                <button
                  onclick={() => handleSaveEdit(comment.id)}
                  class="px-3 py-1.5 text-sm bg-emerald-500 text-white rounded-lg hover:bg-emerald-600"
                >
                  ذخیره
                </button>
                <button
                  onclick={cancelEdit}
                  class="px-3 py-1.5 text-sm bg-slate-200 text-slate-700 rounded-lg hover:bg-slate-300"
                >
                  لغو
                </button>
              </div>
            </div>
          {:else}
            <p class="text-slate-700 whitespace-pre-wrap">{comment.content}</p>
          {/if}
        </div>

        {#if isOwnComment(comment)}
          <div class="flex gap-2">
            {#if editingCommentId !== comment.id}
              <button
                onclick={() => startEdit(comment)}
                class="p-1.5 hover:bg-slate-100 rounded-lg text-slate-600"
                title="ویرایش"
              >
                <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path
                    stroke-linecap="round"
                    stroke-linejoin="round"
                    stroke-width="2"
                    d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z"
                  />
                </svg>
              </button>
            {/if}
            <button
              onclick={() => confirmDelete(comment)}
              class="p-1.5 hover:bg-rose-50 rounded-lg text-slate-600 hover:text-rose-600"
              title="حذف"
            >
              <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path
                  stroke-linecap="round"
                  stroke-linejoin="round"
                  stroke-width="2"
                  d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"
                />
              </svg>
            </button>
          </div>
        {/if}
      </div>
    </div>
  {/each}

  <div class="mt-6">
    <textarea
      bind:value={newComment}
      class="w-full px-4 py-3 border border-slate-200 rounded-lg resize-none"
      placeholder="نظر خود را بنویسید..."
      rows="3"
    ></textarea>
    <button
      onclick={handleAddComment}
      class="mt-2 px-4 py-2 bg-slate-900 text-white rounded-lg hover:bg-slate-800 font-medium"
    >
      ارسال کامنت
    </button>
  </div>
</div>

{#if showDeleteModal && commentToDelete}
  <Modal on:close={() => { showDeleteModal = false; commentToDelete = null; }}>
    <div class="p-6">
      <h3 class="text-lg font-semibold text-slate-900 mb-2">
        حذف کامنت
      </h3>
      <p class="text-slate-600 mb-4">
        آیا مطمئن هستید که می‌خواهید این کامنت را حذف کنید؟
      </p>
      <div class="flex gap-3 justify-end">
        <button
          onclick={() => { showDeleteModal = false; commentToDelete = null; }}
          class="px-4 py-2 bg-slate-200 text-slate-700 rounded-lg hover:bg-slate-300"
        >
          لغو
        </button>
        <button
          onclick={handleDelete}
          class="px-4 py-2 bg-rose-600 text-white rounded-lg hover:bg-rose-700"
        >
          حذف
        </button>
      </div>
    </div>
  </Modal>
{/if}
