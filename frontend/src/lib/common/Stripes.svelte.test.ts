// @vitest-environment jsdom
import { render } from '@testing-library/svelte';
import { expect, it } from 'vitest';
import Stripes from './Stripes.svelte';
import { hour } from '../../test/fixtures';

const time = Date.parse('2026-11-10T00:00:00Z') / 1000;

it.each([0, 1])('renders an empty state for %i hourly entries', (count) => {
	const { getByText } = render(Stripes, { hours: count ? [hour(time)] : [], timezone: 'UTC' });
	expect(getByText('No hourly timeline is available.')).toBeTruthy();
});

it('can transition between unavailable and populated timelines', async () => {
	const { queryByText, getByText, rerender } = render(Stripes, { hours: [], timezone: 'UTC' });
	await rerender({ hours: [hour(time), hour(time + 3600), hour(time + 7200)] });
	expect(queryByText('No hourly timeline is available.')).toBeNull();
	expect(getByText('Clear')).toBeTruthy();
	await rerender({ hours: [] });
	expect(getByText('No hourly timeline is available.')).toBeTruthy();
});

it('renders weather periods without treating the closing tick as another hour', () => {
	const hours = [
		hour(time),
		hour(time + 3600),
		hour(time + 7200, { precipType: 'rain', precipIntensity: 0.1 }),
		hour(time + 10800, { precipType: 'rain', precipIntensity: 0.1 }),
		hour(time + 14400, { precipType: 'snow', precipIntensity: 0.1 })
	];
	const { getByText, queryByText } = render(Stripes, { hours, timezone: 'UTC' });
	expect(getByText('Clear')).toBeTruthy();
	expect(getByText('Rain')).toBeTruthy();
	expect(queryByText('Snow')).toBeNull();
});
