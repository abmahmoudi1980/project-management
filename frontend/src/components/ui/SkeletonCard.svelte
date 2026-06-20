<script>
  import Skeleton from "./Skeleton.svelte";
  import Card from "./Card.svelte";

  let {
    lines = 3,
    showHeader = true,
    showFooter = false,
    class: extraClass = "",
  } = $props();
</script>

<Card padding="md" class="space-y-3 {extraClass}">
  {#if showHeader}
    <div class="flex items-center justify-between">
      <Skeleton width="w-1/3" height="h-5" />
      <Skeleton width="w-6" height="h-6" rounded="rounded-full" />
    </div>
  {/if}

  <div class="space-y-2">
    {#each Array(lines) as _, i}
      <Skeleton width={i === lines - 1 ? 'w-2/3' : 'w-full'} height="h-3" />
    {/each}
  </div>

  {#if showFooter}
    <div class="pt-2 flex justify-between">
      <Skeleton width="w-1/4" height="h-3" />
      <Skeleton width="w-1/4" height="h-3" />
    </div>
  {/if}
</Card>
