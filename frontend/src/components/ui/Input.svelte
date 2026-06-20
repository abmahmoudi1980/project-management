<script>
  let {
    value = $bindable(""),
    label = "",
    id = undefined,
    name = undefined,
    type = "text",
    placeholder = "",
    required = false,
    disabled = false,
    readonly = false,
    autocomplete = undefined,
    inputmode = undefined,
    min = undefined,
    max = undefined,
    maxlength = undefined,
    minlength = undefined,
    step = undefined,
    pattern = undefined,
    error = "",
    helper = "",
    size = "md",
    class: extraClass = "",
    oninput = undefined,
    onchange = undefined,
    onblur = undefined,
    onfocus = undefined,
  } = $props();

  let fieldId = $state(id);
  let errorId = $state();
  let helperId = $state();

  $effect(() => {
    if (!fieldId) {
      fieldId = `inp-${Math.random().toString(36).slice(2, 9)}`;
    }
    errorId = `${fieldId}-err`;
    helperId = `${fieldId}-help`;
  });

  const sizeClasses = {
    sm: "px-2.5 py-1.5 text-sm min-h-[36px]",
    md: "px-3 py-2.5 text-sm min-h-[44px]",
    lg: "px-4 py-3 text-base min-h-[48px]",
  };

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

  <input
    bind:value
    id={fieldId}
    {name}
    {type}
    {placeholder}
    {required}
    {disabled}
    {readonly}
    {autocomplete}
    {inputmode}
    {min}
    {max}
    {maxlength}
    {minlength}
    {step}
    {pattern}
    aria-invalid={error ? "true" : undefined}
    aria-describedby={describedBy}
    class="w-full border rounded-lg transition-colors placeholder-slate-400
      focus:outline-none focus:ring-2 focus:ring-brand-500 focus:border-transparent
      disabled:bg-slate-100 disabled:text-slate-500 disabled:cursor-not-allowed
      {sizeClasses[size]}
      {error ? 'border-danger-500 focus:ring-danger-500' : 'border-slate-300'}"
    {oninput}
    {onchange}
    {onblur}
    {onfocus}
  />

  {#if error}
    <p id={errorId} class="text-xs text-danger-600" role="alert">{error}</p>
  {:else if helper}
    <p id={helperId} class="text-xs text-slate-500">{helper}</p>
  {/if}
</div>
