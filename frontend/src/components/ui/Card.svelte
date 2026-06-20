<script>
  let {
    padding = "md",
    shadow = "sm",
    hover = false,
    bordered = true,
    class: extraClass = "",
    children,
  } = $props();

  const paddingClasses = {
    none: "",
    sm: "p-3",
    md: "p-5",
    lg: "p-6",
    xl: "p-8",
  };

  const shadowClasses = {
    none: "",
    sm: "shadow-sm",
    md: "shadow-md",
    lg: "shadow-lg",
  };
</script>

<div
  class="bg-white rounded-xl
    {bordered ? 'border border-slate-200' : ''}
    {paddingClasses[padding]}
    {shadowClasses[shadow]}
    {hover ? 'transition-shadow hover:shadow-md' : ''}
    {extraClass}"
>
  {@render children?.()}
</div>
