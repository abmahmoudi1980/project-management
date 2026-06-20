<script>
  let { size = "md", color = "brand", label = "در حال بارگذاری" } = $props();

  const sizeClasses = {
    sm: "h-4 w-4 border-2",
    md: "h-6 w-6 border-2",
    lg: "h-10 w-10 border-[3px]",
  };

  const colorClasses = {
    brand: "border-brand-600",
    white: "border-white",
    slate: "border-slate-400",
    danger: "border-danger-600",
    success: "border-success-600",
  };
</script>

<span class="inline-flex items-center" role="status" aria-label={label}>
  <svg
    class="animate-spin rounded-full {sizeClasses[size]} border-t-transparent {colorClasses[color]}"
    viewBox="0 0 24 24"
    fill="none"
    aria-hidden="true"
  >
    <circle cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4" opacity="0.25" />
    <path
      d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"
      fill="currentColor"
      opacity="0.75"
    />
  </svg>
</span>
