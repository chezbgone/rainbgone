// @vitest-environment jsdom
import { fireEvent, render } from '@testing-library/svelte';
import { goto } from '$app/navigation';
import { beforeEach, expect, it, vi } from 'vitest';
import SearchBar from './SearchBar.svelte';

vi.mock('$app/navigation', () => ({ goto: vi.fn(), replaceState: vi.fn() }));
beforeEach(() => vi.mocked(goto).mockClear());

it.each([
	['123 Main St #4, Seattle', '/123%20Main%20St%20%234%2C%20Seattle'],
	['Paris/France', '/Paris%2FFrance'],
	['Where? A&B=50%', '/Where%3F%20A%26B%3D50%25'],
	['  Seattle  ', '/Seattle']
])('preserves the complete search address: %s', async (address, expectedPath) => {
	const { container, getByRole } = render(SearchBar, { defaultAddress: '' });
	await fireEvent.input(getByRole('textbox'), { target: { value: address } });
	await fireEvent.submit(container.querySelector('form')!);
	expect(goto).toHaveBeenCalledExactlyOnceWith(expectedPath);
	const url = new URL(vi.mocked(goto).mock.calls[0][0], 'http://app.example');
	expect(url.hash).toBe('');
	expect(url.search).toBe('');
	expect(decodeURIComponent(url.pathname.slice(1))).toBe(address.trim());
});

it('does not navigate for a blank address', async () => {
	const { container, getByRole } = render(SearchBar, { defaultAddress: '' });
	await fireEvent.input(getByRole('textbox'), { target: { value: '   ' } });
	await fireEvent.submit(container.querySelector('form')!);
	expect(goto).not.toHaveBeenCalled();
});
