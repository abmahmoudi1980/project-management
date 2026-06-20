<script>
  import { createEventDispatcher } from "svelte";
  import Modal from "./Modal.svelte";

  let {
    show = false,
    title = "تأیید عملیات",
    message = "آیا مطمئن هستید؟",
    confirmText = "تأیید",
    cancelText = "انصراف",
    variant = "danger",
  } = $props();

  const dispatch = createEventDispatcher();

  function handleConfirm() {
    dispatch("confirm");
  }

  function handleCancel() {
    dispatch("cancel");
  }
</script>

<Modal show={show} {title} maxWidth="md" fullScreen={false} on:close={handleCancel}>
  <div class="p-4 sm:p-6">
    <p class="text-slate-600 mb-6 leading-relaxed">{message}</p>
    <div class="flex flex-col sm:flex-row gap-3 sm:justify-end">
      <button
        type="button"
        onclick={handleCancel}
        class="w-full sm:w-auto px-4 py-3 min-h-[44px] sm:min-h-0 bg-slate-200 text-slate-700 rounded-lg hover:bg-slate-300 transition-colors font-medium"
      >
        {cancelText}
      </button>
      <button
        type="button"
        onclick={handleConfirm}
        class="w-full sm:w-auto px-4 py-3 min-h-[44px] sm:min-h-0 text-white rounded-lg font-medium transition-colors
          {variant === 'danger'
            ? 'bg-danger-600 hover:bg-danger-700'
            : 'bg-brand-600 hover:bg-brand-700'}"
      >
        {confirmText}
      </button>
    </div>
  </div>
</Modal>
