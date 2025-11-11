import { describe, it, expect } from 'vitest';
import { base64Encode } from './base64';

describe('base64Encode', () => {
	it('should encode empty ArrayBuffer', () => {
		const buffer = new ArrayBuffer(0);
		const result = base64Encode(buffer);
		expect(result).toBe('');
	});

	it('should encode single byte', () => {
		const buffer = new Uint8Array([65]).buffer; // 'A' in ASCII
		const result = base64Encode(buffer);
		expect(result).toBe('QQ==');
	});

	it('should encode simple string data', () => {
		const text = 'Hello';
		const buffer = new TextEncoder().encode(text).buffer;
		const result = base64Encode(buffer);
		expect(result).toBe('SGVsbG8=');
	});

	it('should encode "Hello, World!"', () => {
		const text = 'Hello, World!';
		const buffer = new TextEncoder().encode(text).buffer;
		const result = base64Encode(buffer);
		expect(result).toBe('SGVsbG8sIFdvcmxkIQ==');
	});

	it('should encode binary data correctly', () => {
		// Test with various byte values
		const buffer = new Uint8Array([0, 1, 2, 255, 128, 64]).buffer;
		const result = base64Encode(buffer);
		expect(result).toBe('AAEC/4BA');
	});

	it('should encode all ASCII printable characters', () => {
		const text = 'ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789';
		const buffer = new TextEncoder().encode(text).buffer;
		const result = base64Encode(buffer);
		// This should produce a valid base64 string
		expect(result).toBeTruthy();
		expect(result.length).toBeGreaterThan(0);
		// Verify it can be decoded back
		const decoded = atob(result);
		expect(decoded).toBe(text);
	});

	it('should encode Unicode characters', () => {
		const text = 'Hello 世界 🌍';
		const buffer = new TextEncoder().encode(text).buffer;
		const result = base64Encode(buffer);
		expect(result).toBeTruthy();
		expect(result.length).toBeGreaterThan(0);
		// Verify it can be decoded back
		const decoded = new TextDecoder().decode(new Uint8Array(Array.from(atob(result), (c) => c.charCodeAt(0))));
		expect(decoded).toBe(text);
	});

	it('should handle large buffers without stack overflow', () => {
		// Create a buffer with 100KB of data
		const size = 100 * 1024;
		const data = new Uint8Array(size);
		// Fill with pattern
		for (let i = 0; i < size; i++) {
			data[i] = i % 256;
		}
		const buffer = data.buffer;
		
		// Should not throw
		expect(() => {
			const result = base64Encode(buffer);
			expect(result).toBeTruthy();
			expect(result.length).toBeGreaterThan(0);
		}).not.toThrow();
	});

	it('should handle very large buffers', () => {
		// Create a buffer with 1MB of data
		const size = 1024 * 1024;
		const data = new Uint8Array(size);
		// Fill with zeros for speed
		data.fill(0);
		const buffer = data.buffer;
		
		// Should not throw
		expect(() => {
			const result = base64Encode(buffer);
			expect(result).toBeTruthy();
			expect(result.length).toBeGreaterThan(0);
		}).not.toThrow();
	});

	it('should produce correct base64 encoding for known values', () => {
		// Test cases with known base64 encodings
		const testCases = [
			{ input: 'Man', expected: 'TWFu' },
			{ input: 'Ma', expected: 'TWE=' },
			{ input: 'M', expected: 'TQ==' },
			{ input: 'foobar', expected: 'Zm9vYmFy' }
		];

		for (const testCase of testCases) {
			const buffer = new TextEncoder().encode(testCase.input).buffer;
			const result = base64Encode(buffer);
			expect(result).toBe(testCase.expected);
		}
	});

	it('should handle buffer with all zero bytes', () => {
		const buffer = new Uint8Array(10).buffer; // 10 zero bytes
		const result = base64Encode(buffer);
		expect(result).toBe('AAAAAAAAAAAAAA==');
	});

	it('should handle buffer with all 255 bytes', () => {
		const buffer = new Uint8Array([255, 255, 255]).buffer;
		const result = base64Encode(buffer);
		expect(result).toBe('////');
	});

	it('should produce output that can be decoded back', () => {
		const originalText = 'Test data with special chars: !@#$%^&*()';
		const buffer = new TextEncoder().encode(originalText).buffer;
		const encoded = base64Encode(buffer);
		
		// Decode and verify
		const decodedBytes = new Uint8Array(
			Array.from(atob(encoded), (c) => c.charCodeAt(0))
		);
		const decoded = new TextDecoder().decode(decodedBytes);
		expect(decoded).toBe(originalText);
	});

	it('should handle single byte values correctly', () => {
		for (let i = 0; i < 256; i++) {
			const buffer = new Uint8Array([i]).buffer;
			const result = base64Encode(buffer);
			expect(result).toBeTruthy();
			expect(result.length).toBeGreaterThan(0);
		}
	});

	it('should match browser btoa for string data', () => {
		const text = 'Hello, World!';
		const buffer = new TextEncoder().encode(text).buffer;
		const result = base64Encode(buffer);
		const browserResult = btoa(text);
		expect(result).toBe(browserResult);
	});
});

