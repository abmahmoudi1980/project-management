/**
 * Date utility functions for Jalali/Gregorian conversion
 * Used by TaskSearch component for date filtering
 */

import moment from 'jalali-moment';

/**
 * Convert a Jalali date string (1403/10/15) to JavaScript Date (Gregorian)
 * @param {string} jalaliStr - Date string in format "YYYY/MM/DD" (Jalali)
 * @returns {Date|null} JavaScript Date object in Gregorian calendar, or null if invalid
 */
export function jalaliStringToDate(jalaliStr) {
  if (!jalaliStr || typeof jalaliStr !== 'string') {
    return null;
  }

  try {
    // Parse Jalali date with multiple format attempts
    const formats = ['YYYY/MM/DD', 'YYYY/M/DD', 'YYYY/MM/D', 'YYYY/M/D'];
    let parsed = null;

    for (const format of formats) {
      const attempt = moment(jalaliStr, format, true);
      if (attempt.isValid()) {
        parsed = attempt;
        break;
      }
    }

    if (parsed && parsed.isValid()) {
      // Convert to Gregorian and return JavaScript Date
      return parsed.toDate();
    }

    return null;
  } catch (error) {
    console.error('Error parsing Jalali date:', jalaliStr, error);
    return null;
  }
}

/**
 * Convert a JavaScript Date (Gregorian) to Jalali date string (1403/10/15)
 * @param {Date} date - JavaScript Date object in Gregorian calendar
 * @returns {string} Date string in format "YYYY/MM/DD" (Jalali), or empty string if invalid
 */
export function dateToJalaliString(date) {
  if (!date || !(date instanceof Date)) {
    return '';
  }

  try {
    // Convert Gregorian to Jalali using moment
    const jalaliMoment = moment(date).locale('fa');
    return jalaliMoment.format('YYYY/MM/DD');
  } catch (error) {
    console.error('Error converting to Jalali date:', date, error);
    return '';
  }
}

/**
 * Format a date into a Jalali string with various formats
 * @param {Date|string} date - Date object or ISO string
 * @param {string} formatType - 'full', 'short', 'time', or 'relative'
 * @returns {string} Formatted Jalali date string
 */
export function formatJalaliDate(date, formatType = 'full') {
  if (!date) return '';
  const m = moment(date);
  if (!m.isValid()) return '';

  m.locale('fa');

  switch (formatType) {
    case 'full':
      return m.format('jYYYY/jMM/jDD');
    case 'short':
      return m.format('jMM/jDD');
    case 'time':
      return m.format('HH:mm');
    case 'relative':
      const now = moment();
      if (m.isSame(now, 'day')) return 'امروز';
      if (m.isSame(now.clone().add(1, 'day'), 'day')) return 'فردا';
      if (m.isSame(now.clone().subtract(1, 'day'), 'day')) return 'دیروز';
      return m.format('jDD jMMMM');
    default:
      return m.format('jYYYY/jMM/jDD');
  }
}

/**
 * Normalize a date to start of day (00:00:00) for accurate range comparison
 * @param {Date} date - JavaScript Date object
 * @returns {Date} Date at start of day (midnight)
 */
export function normalizeDate(date) {
  if (!date || !(date instanceof Date)) {
    return null;
  }

  const normalized = new Date(date.getFullYear(), date.getMonth(), date.getDate());
  return normalized;
}

/**
 * Check if a date falls within a range (inclusive on both ends)
 * @param {Date} date - Date to check
 * @param {Date|null} fromDate - Range start (inclusive), null means no lower bound
 * @param {Date|null} toDate - Range end (inclusive), null means no upper bound
 * @returns {boolean} True if date is within range, false otherwise
 */
export function isDateInRange(date, fromDate, toDate) {
  if (!date || !(date instanceof Date)) {
    return false;
  }

  // Normalize all dates to start of day for accurate comparison
  const normalizedDate = normalizeDate(date);
  const normalizedFrom = fromDate ? normalizeDate(fromDate) : null;
  const normalizedTo = toDate ? normalizeDate(toDate) : null;

  // If both boundaries are null, date is in range (no filter applied)
  if (!normalizedFrom && !normalizedTo) {
    return true;
  }

  // Check lower bound
  if (normalizedFrom && normalizedDate < normalizedFrom) {
    return false;
  }

  // Check upper bound
  if (normalizedTo && normalizedDate > normalizedTo) {
    return false;
  }

  return true;
}

/**
 * Parse an ISO 8601 date string to JavaScript Date
 * @param {string} isoString - ISO 8601 date string (e.g., "2024-03-15T00:00:00Z")
 * @returns {Date|null} JavaScript Date object or null if invalid
 */
export function isoStringToDate(isoString) {
  if (!isoString || typeof isoString !== 'string') {
    return null;
  }

  try {
    const date = new Date(isoString);
    // Check if date is valid
    if (isNaN(date.getTime())) {
      return null;
    }
    return date;
  } catch (error) {
    console.error('Error parsing ISO date string:', isoString, error);
    return null;
  }
}
