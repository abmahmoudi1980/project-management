<script>
  import { timeLogs } from "../stores/timeLogStore";
  import JalaliDatePicker from "./JalaliDatePicker.svelte";
  import moment from "jalali-moment";

  let { task } = $props();

  let date = $state(new Date().toISOString().split("T")[0]);
  let durationMinutes = $state(30);
  let note = $state("");

  function formatJalaliDate(dateString) {
    if (!dateString) return "";
    return moment(dateString).locale("fa").format("YYYY/MM/DD");
  }

  async function handleSubmit() {
    if (!durationMinutes || durationMinutes <= 0) return;

    await timeLogs.create(task.id, {
      date: new Date(date).toISOString(),
      duration_minutes: durationMinutes,
      note: note.trim() || null,
    });

    durationMinutes = 30;
    note = "";
  }

  async function handleDelete(timeLogId) {
    if (confirm("آیا مطمئن هستید که می‌خواهید این زمان ثبت شده را حذف کنید؟")) {
      await timeLogs.delete(timeLogId);
    }
  }

  function formatMinutes(minutes) {
    const hours = Math.floor(minutes / 60);
    const mins = minutes % 60;
    if (hours > 0) {
      return `${hours}h ${mins}m`;
    }
    return `${mins}m`;
  }
</script>

<div class="space-y-4">
  <form on:submit|preventDefault={handleSubmit} class="flex gap-3">
    <div class="w-40">
      <JalaliDatePicker
        bind:value={date}
        placeholder="1403/10/10"
      />
    </div>
    <input
      type="number"
      bind:value={durationMinutes}
      min="1"
      placeholder="دقیقه"
      class="px-3 py-2 border border-gray-300 rounded-lg w-24 focus:outline-none focus:ring-2 focus:ring-blue-500"
    />
    <input
      type="text"
      bind:value={note}
      placeholder="یادداشت (اختیاری)"
      class="flex-1 px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
    />
    <button
      type="submit"
      disabled={!durationMinutes || durationMinutes <= 0}
      class="bg-green-500 hover:bg-green-600 disabled:bg-gray-300 text-white px-4 py-2 rounded-lg transition-colors"
    >
      افزودن
    </button>
  </form>

  {#if ($timeLogs || []).length > 0}
    <div class="space-y-2">
      <h4 class="text-sm font-medium text-gray-700">زمان‌های ثبت شده</h4>
      {#each $timeLogs || [] as log (log.id)}
        <div
          class="flex justify-between items-center p-3 bg-gray-50 rounded-lg text-sm"
        >
          <div>
            <span class="font-medium"
              >{formatJalaliDate(log.date)}</span
            >
            <span class="mx-2">•</span>
            <span class="text-blue-600 font-semibold"
              >{formatMinutes(log.duration_minutes)}</span
            >
            {#if log.note}
              <span class="ml-2 text-gray-600">- {log.note}</span>
            {/if}
          </div>
          <button
            on:click={() => handleDelete(log.id)}
            class="text-red-500 hover:text-red-700 text-xs"
          >
            حذف
          </button>
        </div>
      {/each}
    </div>
  {/if}
</div>
