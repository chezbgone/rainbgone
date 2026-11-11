// @vitest-environment jsdom
import { render, within } from '@testing-library/svelte';
import { expect, it } from 'vitest';
import Daily from './Daily.svelte';
import { day, hoursBetween } from '../../test/fixtures';

it('keeps each daily row on its local date when morning backfill is unavailable', () => {
	const days = ['2026-11-10', '2026-11-11'].map((date) =>
		day(Date.parse(`${date}T00:00:00Z`) / 1000)
	);
	const { container } = render(Daily, {
		daily: { summary: 'Fixture', icon: 'clear-day', data: days },
		hourlyFromMidnight: hoursBetween('2026-11-10T15:00:00Z', '2026-11-13T00:00:00Z'),
		lat: 0,
		lng: 0,
		timezone: 'UTC'
	});
	const rows = container.querySelectorAll('details');
	expect(within(rows[0]).queryByText('1am')).toBeNull();
	const firstTimelineLabel = rows[1].querySelector('.top-0');
	expect(firstTimelineLabel?.textContent).toBe('12am');
});

it('keeps daily rows usable when a later day has no hourly coverage', () => {
	const days = ['2026-11-10', '2026-11-11'].map((date) =>
		day(Date.parse(`${date}T00:00:00Z`) / 1000)
	);
	const { getByText, getAllByRole } = render(Daily, {
		daily: { summary: 'Fixture', icon: 'clear-day', data: days },
		hourlyFromMidnight: hoursBetween('2026-11-10T15:00:00Z', '2026-11-10T23:00:00Z'),
		lat: 0,
		lng: 0,
		timezone: 'UTC'
	});
	expect(getByText('No hourly timeline is available.')).toBeTruthy();
	expect(getAllByRole('link', { name: 'more details', hidden: true })).toHaveLength(2);
});
