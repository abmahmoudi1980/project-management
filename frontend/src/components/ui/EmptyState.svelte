<script>
  let {
    title = "",
    description = "",
    icon = undefined,
    children,
    action,
  } = $props();
</script>

<div class="text-center py-8 md:py-12 px-4">
  {#if icon}
    <div class="w-16 h-16 mx-auto mb-4 rounded-full bg-slate-100 flex items-center justify-center text-slate-400">
      {@render icon()}
    </div>
  {/if}
  {#if title}
    <h2 class="text-xl font-semibold text-slate-700 mb-2">{title}</h2>
  {/if}
  {#if description}
    <p class="text-slate-500 max-w-md mx-auto">{description}</p>
  {/if}
  {#if children}
    <div class="mt-2">{@render children()}</div>
  {/if}
  {#if action}
    <div class="mt-6 flex justify-center">{@render action()}</div>
  {/if}
</div>
