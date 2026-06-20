<script>
  import Spinner from "./Spinner.svelte";

  let {
    variant = "primary",
    size = "md",
    type = "button",
    loading = false,
    disabled = false,
    fullWidth = false,
    onclick = undefined,
    children,
    class: extraClass = "",
    ...rest
  } = $props();

  const sizeClasses = {
    sm: "px-3 py-1.5 text-xs min-h-[36px]",
    md: "px-4 py-2.5 text-sm min-h-[44px]",
    lg: "px-5 py-3 text-base min-h-[48px]",
  };

  const variantClasses = {
    primary:
      "bg-brand-600 text-white hover:bg-brand-700 active:bg-brand-800 disabled:bg-slate-300 disabled:text-slate-500",
    secondary:
      "bg-white text-slate-700 border border-slate-300 hover:bg-slate-50 active:bg-slate-100 disabled:bg-slate-100 disabled:text-slate-400",
    danger:
      "bg-danger-600 text-white hover:bg-danger-700 active:bg-danger-700 disabled:bg-slate-300 disabled:text-slate-500",
    ghost:
      "bg-transparent text-slate-700 hover:bg-slate-100 active:bg-slate-200 disabled:text-slate-400 disabled:hover:bg-transparent",
    "ghost-primary":
      "bg-transparent text-brand-600 hover:bg-brand-50 active:bg-brand-100 disabled:text-slate-400 disabled:hover:bg-transparent",
  };

  const spinnerColor = {
    primary: "white",
    secondary: "slate",
    danger: "white",
    ghost: "slate",
    "ghost-primary": "brand",
  };

  const isDisabled = $derived(disabled || loading);
</script>

<button
  {type}
  onclick={onclick}
  disabled={isDisabled}
  aria-busy={loading || undefined}
  class="inline-flex items-center justify-center gap-2 font-medium rounded-lg transition-colors disabled:cursor-not-allowed focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-brand-500
    {sizeClasses[size]}
    {variantClasses[variant]}
    {fullWidth ? 'w-full' : ''}
    {extraClass}"
  {...rest}
>
  {#if loading}
    <Spinner size="sm" color={spinnerColor[variant]} label="" />
  {/if}
  {@render children?.()}
</button>
