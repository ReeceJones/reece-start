/**
 * Encodes an ArrayBuffer to a base64 string. It is necessary to use this function for large buffers to avoid stack overflows.
 * @param data
 * @returns
 */
export function base64Encode(data: ArrayBuffer): string {
	return btoa(
		new Uint8Array(data).reduce(function (data, byte) {
			return data + String.fromCharCode(byte);
		}, '')
	);
}
