<script>
  import { createEventDispatcher } from "svelte";
  import moment from "jalali-moment";

  let { value = $bindable(""), placeholder = "مثال: 1403/10/10", error = false } = $props();

  const dispatch = createEventDispatcher();

  let jalaliDate = $state("");
  let inputElement = $state();
  let showCalendar = $state(false);
  let currentMonth = $state(moment());
  let calendarContainer = $state();

  // Convert Gregorian to Jalali when value changes
  $effect(() => {
    if (value) {
      jalaliDate = moment(value, "YYYY-MM-DD").locale("fa").format("YYYY/MM/DD");
      currentMonth = moment(value, "YYYY-MM-DD");
    } else {
      jalaliDate = "";
      currentMonth = moment();
    }
  });

  // Generate calendar days
  let calendarDays = $derived(generateCalendar(currentMonth));

  function handleInput(event) {
    let input = event.target.value;
    
    // Remove non-numeric and non-slash characters
    input = input.replace(/[^\d/]/g, '');
    
    // Auto-add slashes
    if (input.length === 4 && !input.includes('/')) {
      input = input + '/';
    } else if (input.length === 7 && input.split('/').length === 2) {
      input = input + '/';
    }
    
    jalaliDate = input;

    // Try to parse the Jalali date and convert to Gregorian
    if (input && input.length >= 10) {
      try {
        // Parse as Jalali date
        const parsed = moment(input, "jYYYY/jMM/jDD");
        
        // Validate the date
        if (parsed.isValid()) {
          const gregorianDate = parsed.locale("en").format("YYYY-MM-DD");
          value = gregorianDate;
          dispatch("change", gregorianDate);
        }
      } catch (e) {
        // Invalid date format
      }
    } else if (!input) {
      value = "";
      dispatch("change", "");
    }
  }

  function handleBlur() {
    // Format the date properly on blur
    if (jalaliDate && jalaliDate.length >= 8) {
      try {
        // Try different formats
        const formats = ["jYYYY/jMM/jDD", "jYYYY/jM/jD", "jYYYY/jMM/jD", "jYYYY/jM/jDD"];
        let parsed = null;
        
        for (const format of formats) {
          const attempt = moment(jalaliDate, format, true);
          if (attempt.isValid()) {
            parsed = attempt;
            break;
          }
        }
        
        if (parsed && parsed.isValid()) {
          jalaliDate = parsed.locale("fa").format("YYYY/MM/DD");
          value = parsed.locale("en").format("YYYY-MM-DD");
          dispatch("change", value);
        }
      } catch (e) {
        // Keep as is
      }
    }
  }

  function handleKeydown(event) {
    // Close calendar on Escape
    if (event.keyCode === 27) {
      showCalendar = false;
      return;
    }
    // Allow: backspace, delete, tab, enter
    if ([46, 8, 9, 13].indexOf(event.keyCode) !== -1 ||
        // Allow: Ctrl/cmd+A
        (event.keyCode === 65 && (event.ctrlKey === true || event.metaKey === true)) ||
        // Allow: Ctrl/cmd+C
        (event.keyCode === 67 && (event.ctrlKey === true || event.metaKey === true)) ||
        // Allow: Ctrl/cmd+V
        (event.keyCode === 86 && (event.ctrlKey === true || event.metaKey === true)) ||
        // Allow: Ctrl/cmd+X
        (event.keyCode === 88 && (event.ctrlKey === true || event.metaKey === true)) ||
        // Allow: home, end, left, right
        (event.keyCode >= 35 && event.keyCode <= 39)) {
      return;
    }
    // Allow: forward slash
    if (event.key === '/') {
      return;
    }
    // Ensure that it is a number
    if ((event.shiftKey || (event.keyCode < 48 || event.keyCode > 57)) && (event.keyCode < 96 || event.keyCode > 105)) {
      event.preventDefault();
    }
  }

  function generateCalendar(date) {
    const jalaliDate = moment(date).locale('fa');
    const year = jalaliDate.jYear();
    const month = jalaliDate.jMonth();
    
    const firstDay = moment.from(`${year}/${month + 1}/1`, 'fa', 'jYYYY/jM/jD');
    const daysInMonth = jalaliDate.jDaysInMonth();
    const startDayOfWeek = firstDay.day();
    
    const days = [];
    
    // Add empty cells for days before the month starts
    for (let i = 0; i < (startDayOfWeek + 1) % 7; i++) {
      days.push(null);
    }
    
    // Add days of the month
    for (let day = 1; day <= daysInMonth; day++) {
      const dayMoment = moment.from(`${year}/${month + 1}/${day}`, 'fa', 'jYYYY/jM/jD');
      days.push({
        day,
        gregorian: dayMoment.format('YYYY-MM-DD'),
        isToday: dayMoment.isSame(moment(), 'day'),
        isSelected: value && dayMoment.isSame(moment(value, 'YYYY-MM-DD'), 'day')
      });
    }
    
    return days;
  }

  function selectDate(day) {
    if (day) {
      value = day.gregorian;
      jalaliDate = moment(day.gregorian, "YYYY-MM-DD").locale("fa").format("YYYY/MM/DD");
      dispatch("change", day.gregorian);
      showCalendar = false;
    }
  }

  function previousMonth() {
    currentMonth = moment(currentMonth).subtract(1, 'jMonth');
  }

  function nextMonth() {
    currentMonth = moment(currentMonth).add(1, 'jMonth');
  }

  function toggleCalendar() {
    showCalendar = !showCalendar;
  }

  function handleClickOutside(event) {
    if (calendarContainer && !calendarContainer.contains(event.target) && !inputElement.contains(event.target)) {
      showCalendar = false;
    }
  }
</script>

<svelte:window on:click={handleClickOutside} />

<div class="relative">
  <input
    type="text"
    bind:this={inputElement}
    bind:value={jalaliDate}
    on:input={handleInput}
    on:blur={handleBlur}
    on:keydown={handleKeydown}
    on:click={toggleCalendar}
    {placeholder}
    maxlength="10"
    class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500 cursor-pointer"
    class:border-red-500={error}
    dir="ltr"
    style="text-align: right; font-family: 'Vazirmatn', sans-serif;"
  />
  
  {#if showCalendar}
    <div
      bind:this={calendarContainer}
      class="absolute z-50 mt-2 bg-white border border-gray-300 rounded-lg shadow-lg p-4"
      style="min-width: 280px; left: 0;"
    >
      <!-- Calendar Header -->
      <div class="flex items-center justify-between mb-4">
        <button
          type="button"
          on:click={nextMonth}
          class="p-1 hover:bg-gray-100 rounded"
        >
          <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7" />
          </svg>
        </button>
        
        <div class="font-semibold">
          {moment(currentMonth).locale('fa').format('jMMMM jYYYY')}
        </div>
        
        <button
          type="button"
          on:click={previousMonth}
          class="p-1 hover:bg-gray-100 rounded"
        >
          <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7" />
          </svg>
        </button>
      </div>
      
      <!-- Week days header -->
      <div class="grid grid-cols-7 gap-1 mb-2">
        {#each ['ش', 'ی', 'د', 'س', 'چ', 'پ', 'ج'] as day}
          <div class="text-center text-xs font-semibold text-gray-600 py-1">
            {day}
          </div>
        {/each}
      </div>
      
      <!-- Calendar days -->
      <div class="grid grid-cols-7 gap-1">
        {#each calendarDays as day}
          {#if day}
            <button
              type="button"
              on:click={() => selectDate(day)}
              class="p-2 text-sm rounded hover:bg-blue-100 transition-colors"
              class:bg-blue-500={day.isSelected}
              class:text-white={day.isSelected}
              class:font-bold={day.isToday}
              class:ring-2={day.isToday}
              class:ring-blue-300={day.isToday}
            >
              {day.day}
            </button>
          {:else}
            <div class="p-2"></div>
          {/if}
        {/each}
      </div>
      
      <!-- Today button -->
      <div class="mt-3 pt-3 border-t border-gray-200">
        <button
          type="button"
          on:click={() => selectDate({
            day: moment().locale('fa').jDate(),
            gregorian: moment().format('YYYY-MM-DD'),
            isToday: true,
            isSelected: false
          })}
          class="w-full px-3 py-2 text-sm bg-gray-100 hover:bg-gray-200 rounded transition-colors"
        >
          امروز
        </button>
      </div>
    </div>
  {/if}
</div>
