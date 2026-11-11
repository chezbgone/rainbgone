import { describe, expect, it } from 'vitest';
import { hoursForDate } from './hourly';
import { formatDateKey } from './time';
import { hour, hoursBetween } from '../../test/fixtures';

describe('daily hourly selection', () => {
	it('selects local dates rather than the UTC date or a fixed offset', () => {
		const series = hoursBetween('2026-11-10T23:00:00Z', '2026-11-13T08:00:00Z');
		const today = hoursForDate(series, '2026-11-10', 'America/Los_Angeles');
		const tomorrow = hoursForDate(series, '2026-11-11', 'America/Los_Angeles');
		expect(today).toHaveLength(10); // 3pm to midnight, with the closing tick.
		expect(tomorrow).toHaveLength(25);
		expect(tomorrow[0].time).toBe(Date.parse('2026-11-11T08:00:00Z') / 1000);
		expect(
			today.slice(0, -1).every((h) => formatDateKey(h.time, 'America/Los_Angeles') === '2026-11-10')
		).toBe(true);
		expect(today.at(-1)).toBe(tomorrow[0]);
	});

	it.each([
		['2026-03-08', '2026-03-08T08:00:00Z', '2026-03-09T07:00:00Z', 23],
		['2026-11-01', '2026-11-01T07:00:00Z', '2026-11-02T08:00:00Z', 25]
	])('preserves all hours on the DST day %s', (date, start, end, count) => {
		const series = hoursBetween(start, end);
		const hours = hoursForDate(series, date, 'America/Los_Angeles');
		expect(hours).toHaveLength(Number(count) + 1);
		expect(hours).toEqual(series);
	});

	it('closes a partial final day without mutating the source', () => {
		const series = hoursBetween('2026-11-10T15:00:00Z', '2026-11-10T17:00:00Z');
		const hours = hoursForDate(series, '2026-11-10', 'UTC');
		expect(hours).toHaveLength(4);
		expect(hours.at(-1)?.time).toBe(series.at(-1)!.time + 3600);
		expect(series).toHaveLength(3);
	});

	it('does not stretch the final hour across a gap in the series', () => {
		const last = hour(Date.parse('2026-11-10T23:00:00Z') / 1000);
		const later = hour(last.time + 3 * 3600);
		expect(hoursForDate([last, later], '2026-11-10', 'UTC').at(-1)?.time).toBe(last.time + 3600);
	});

	it('returns no entries for an unavailable day', () => {
		expect(hoursForDate([], '2026-11-10', 'UTC')).toEqual([]);
	});
});
