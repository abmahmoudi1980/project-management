<script>
  let {
    width = "w-full",
    height = "h-4",
    rounded = "rounded",
    class: extraClass = "",
  } = $props();
</script>

<div
  class="animate-pulse bg-slate-200 {width} {height} {rounded} {extraClass}"
  aria-hidden="true"
></div>
