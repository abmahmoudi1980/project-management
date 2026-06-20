<script>
  let {
    value = $bindable(""),
    label = "",
    id = undefined,
    name = undefined,
    placeholder = "",
    required = false,
    disabled = false,
    readonly = false,
    rows = 3,
    maxlength = undefined,
    error = "",
    helper = "",
    class: extraClass = "",
    oninput = undefined,
    onchange = undefined,
    onblur = undefined,
  } = $props();

  let fieldId = $state(id);
  let errorId = $state();
  let helperId = $state();

  $effect(() => {
    if (!fieldId) {
      fieldId = `txt-${Math.random().toString(36).slice(2, 9)}`;
    }
    errorId = `${fieldId}-err`;
    helperId = `${fieldId}-help`;
  });

  let describedBy = $derived(
    [error ? errorId : null, helper ? helperId : null].filter(Boolean).join(" ") || undefined
  );
</script>

<div class="space-y-1.5 {extraClass}">
  {#if label}
    <label for={fieldId} class="block text-sm font-medium text-slate-700">
      {label}{#if required}<span class="text-danger-500 ms-1">*</span>{/if}
    </label>
  {/if}

  <textarea
    bind:value
    id={fieldId}
    {name}
    {placeholder}
    {required}
    {disabled}
    {readonly}
    {rows}
    {maxlength}
    aria-invalid={error ? "true" : undefined}
    aria-describedby={describedBy}
    class="w-full px-3 py-2.5 text-sm border rounded-lg transition-colors placeholder-slate-400 resize-none
      focus:outline-none focus:ring-2 focus:ring-brand-500 focus:border-transparent
      disabled:bg-slate-100 disabled:text-slate-500 disabled:cursor-not-allowed
      {error ? 'border-danger-500 focus:ring-danger-500' : 'border-slate-300'}"
    {oninput}
    {onchange}
    {onblur}
  ></textarea>

  {#if error}
    <p id={errorId} class="text-xs text-danger-600" role="alert">{error}</p>
  {:else if helper}
    <p id={helperId} class="text-xs text-slate-500">{helper}</p>
  {/if}
</div>
