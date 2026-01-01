/**
 * Filter utility functions for task search and filtering
 * Used by TaskSearch and TaskList components
 */

/**
 * Evaluate if a task matches the text search filter
 * @param {Object} task - Task object with title and description properties
 * @param {string} searchTerm - Search term to match (case-insensitive)
 * @returns {boolean} True if task matches search term or no search term provided
 */
export function evaluateTextFilter(task, searchTerm) {
  // If no search term, include all tasks
  if (!searchTerm || typeof searchTerm !== 'string') {
    return true;
  }

  // Normalize search term to lowercase
  const normalizedSearchTerm = searchTerm.toLowerCase().trim();

  // If search term is empty, include all tasks
  if (normalizedSearchTerm.length === 0) {
    return true;
  }

  // Check if title matches (case-insensitive)
  if (task.title && task.title.toLowerCase().includes(normalizedSearchTerm)) {
    return true;
  }

  // Check if description matches (case-insensitive)
  if (task.description && task.description.toLowerCase().includes(normalizedSearchTerm)) {
    return true;
  }

  // No match found
  return false;
}

/**
 * Evaluate if a task's date falls within a specified range
 * @param {Object} task - Task object with a date field (start_date or due_date)
 * @param {string} dateField - Name of the date field to check (e.g., 'start_date', 'due_date')
 * @param {Date|null} fromDate - Range start (inclusive), null means no lower bound
 * @param {Date|null} toDate - Range end (inclusive), null means no upper bound
 * @returns {boolean} True if task date is within range or no filter applied
 */
export function evaluateDateFilter(task, dateField, fromDate, toDate) {
  // If both boundaries are null, no filter applied - include task
  if (!fromDate && !toDate) {
    return true;
  }

  // Get task date and convert to Date object if it's a string
  const taskDateValue = task[dateField];
  if (!taskDateValue) {
    // Task has no date value - exclude it when filter is active
    return false;
  }

  // Convert ISO string to Date if needed
  let taskDate;
  if (typeof taskDateValue === 'string') {
    try {
      taskDate = new Date(taskDateValue);
      if (isNaN(taskDate.getTime())) {
        return false; // Invalid date
      }
    } catch (error) {
      return false; // Error parsing date
    }
  } else if (taskDateValue instanceof Date) {
    taskDate = taskDateValue;
  } else {
    return false; // Unsupported date type
  }

  // Normalize task date to start of day for comparison
  const normalizedTaskDate = new Date(
    taskDate.getFullYear(),
    taskDate.getMonth(),
    taskDate.getDate()
  );

  // Check lower bound (from date)
  if (fromDate) {
    const normalizedFromDate = new Date(
      fromDate.getFullYear(),
      fromDate.getMonth(),
      fromDate.getDate()
    );
    if (normalizedTaskDate < normalizedFromDate) {
      return false; // Task date is before range start
    }
  }

  // Check upper bound (to date)
  if (toDate) {
    const normalizedToDate = new Date(
      toDate.getFullYear(),
      toDate.getMonth(),
      toDate.getDate()
    );
    if (normalizedTaskDate > normalizedToDate) {
      return false; // Task date is after range end
    }
  }

  // Date is within range
  return true;
}

/**
 * Evaluate if a task matches all active filters (AND logic)
 * @param {Object} task - Task object to evaluate
 * @param {Object} filters - Filter state object with:
 *   - text: string (search term)
 *   - start_date_from: Date|null
 *   - start_date_to: Date|null
 *   - due_date_from: Date|null
 *   - due_date_to: Date|null
 * @returns {boolean} True if task matches all active filters
 */
export function evaluateAllFilters(task, filters) {
  if (!task || !filters) {
    return false;
  }

  // Evaluate text filter
  const textMatch = evaluateTextFilter(task, filters.text);
  if (!textMatch) {
    return false;
  }

  // Evaluate start date filter
  const startDateMatch = evaluateDateFilter(
    task,
    'start_date',
    filters.start_date_from,
    filters.start_date_to
  );
  if (!startDateMatch) {
    return false;
  }

  // Evaluate due date filter
  const dueDateMatch = evaluateDateFilter(
    task,
    'due_date',
    filters.due_date_from,
    filters.due_date_to
  );
  if (!dueDateMatch) {
    return false;
  }

  // All filters passed
  return true;
}

/**
 * Check if any filter is currently active
 * @param {Object} filters - Filter state object
 * @returns {boolean} True if any filter has a non-null/non-empty value
 */
export function hasActiveFilters(filters) {
  if (!filters) {
    return false;
  }

  return (
    (filters.text && filters.text.trim().length > 0) ||
    filters.start_date_from !== null ||
    filters.start_date_to !== null ||
    filters.due_date_from !== null ||
    filters.due_date_to !== null
  );
}

/**
 * Get count of active filters
 * @param {Object} filters - Filter state object
 * @returns {number} Number of active filters
 */
export function getActiveFilterCount(filters) {
  if (!filters) {
    return 0;
  }

  let count = 0;

  if (filters.text && filters.text.trim().length > 0) {
    count++;
  }
  if (filters.start_date_from !== null || filters.start_date_to !== null) {
    count++; // Count as one filter pair
  }
  if (filters.due_date_from !== null || filters.due_date_to !== null) {
    count++; // Count as one filter pair
  }

  return count;
}
