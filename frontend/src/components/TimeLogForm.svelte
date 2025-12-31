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
  <form onsubmit={(e) => { e.preventDefault(); handleSubmit(); }} class="flex flex-col md:flex-row gap-3">
    <div class="w-full md:w-40">
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
      class="w-full md:w-24 px-3 py-3 min-h-[44px] border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
    />
    <input
      type="text"
      bind:value={note}
      placeholder="یادداشت (اختیاری)"
      class="w-full md:flex-1 px-3 py-3 min-h-[44px] border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
    />
    <button
      type="submit"
      disabled={!durationMinutes || durationMinutes <= 0}
      class="w-full md:w-auto min-h-[44px] bg-green-500 hover:bg-green-600 disabled:bg-gray-300 text-white px-4 py-3 rounded-lg transition-colors font-medium"
    >
      افزودن
    </button>
  </form>

  {#if ($timeLogs || []).length > 0}
    <div class="space-y-2">
      <h4 class="text-sm font-medium text-gray-700">زمان‌های ثبت شده</h4>
      {#each $timeLogs || [] as log (log.id)}
        <div
          class="flex flex-col md:flex-row md:justify-between md:items-start gap-2 p-3 bg-gray-50 rounded-lg text-sm"
        >
          <div class="flex-1">
            <span class="font-medium block"
              >{formatJalaliDate(log.date)}</span
            >
            <div class="flex items-center gap-2 mt-1">
              <span class="text-blue-600 font-semibold"
                >{formatMinutes(log.duration_minutes)}</span
              >
              {#if log.note}
                <span class="text-gray-600">- {log.note}</span>
              {/if}
            </div>
          </div>
          <button
            onclick={() => handleDelete(log.id)}
            class="self-start md:self-auto px-3 py-2 min-h-[44px] text-red-500 hover:text-red-700 hover:bg-red-50 rounded transition-colors text-xs font-medium"
          >
            حذف
          </button>
        </div>
      {/each}
    </div>
  {/if}
</div>
